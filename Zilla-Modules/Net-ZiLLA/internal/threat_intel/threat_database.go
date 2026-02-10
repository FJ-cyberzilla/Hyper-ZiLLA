package threat_intel

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"net-zilla/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type ThreatDatabase struct {
	db      *sql.DB
	mu      sync.RWMutex
	cache   map[string]*cacheEntry
	metrics *DBMetrics
	logger  *log.Logger
}

type cacheEntry struct {
	indicator *models.Indicator
	timestamp time.Time
	hits      int64
}

type DBMetrics struct {
	TotalQueries      int64
	CacheHits         int64
	CacheMisses       int64
	AverageQueryTime  time.Duration
	TotalIndicators   int64
	mu                sync.RWMutex
}

func NewThreatDatabase(path string, logger *log.Logger) (*ThreatDatabase, error) {
	if path == "" {
		path = "netzilla_threats.db"
	}

	if logger == nil {
		logger = log.New(log.Writer(), "[ThreatDB] ", log.LstdFlags)
	}

	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=5000&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed to open threat db: %w", err)
	}

	// Set connection limits
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	td := &ThreatDatabase{
		db:      db,
		cache:   make(map[string]*cacheEntry),
		metrics: &DBMetrics{},
		logger:  logger,
	}

	if err := td.initThreatSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	// Start cache cleanup goroutine
	go td.startCacheCleanup()

	td.logger.Printf("ThreatDatabase initialized at %s", path)
	return td, nil
}

func (td *ThreatDatabase) initThreatSchema() error {
	queries := []string{
		// Main threat indicators table
		`CREATE TABLE IF NOT EXISTS threat_indicators (
			value TEXT PRIMARY KEY,
			type TEXT NOT NULL CHECK(type IN ('ip', 'domain', 'url', 'hash', 'email', 'asn', 'cidr')),
			source TEXT NOT NULL,
			confidence REAL NOT NULL CHECK(confidence >= 0 AND confidence <= 1),
			severity TEXT NOT NULL CHECK(severity IN ('low', 'medium', 'high', 'critical')),
			last_seen DATETIME NOT NULL,
			first_seen DATETIME NOT NULL,
			description TEXT,
			tags TEXT,
			"references" TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Indexes for faster lookups
		`CREATE INDEX IF NOT EXISTS idx_type ON threat_indicators(type)`,
		`CREATE INDEX IF NOT EXISTS idx_severity ON threat_indicators(severity)`,
		`CREATE INDEX IF NOT EXISTS idx_last_seen ON threat_indicators(last_seen)`,
		`CREATE INDEX IF NOT EXISTS idx_confidence ON threat_indicators(confidence)`,

		// Threat intelligence feed sources
		`CREATE TABLE IF NOT EXISTS threat_feeds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			last_updated DATETIME,
			update_interval INTEGER DEFAULT 3600,
			enabled BOOLEAN DEFAULT 1,
			api_key TEXT,
			trust_level INTEGER DEFAULT 3 CHECK(trust_level >= 1 AND trust_level <= 5),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Indicator relationships
		`CREATE TABLE IF NOT EXISTS indicator_relationships (
			source_indicator TEXT NOT NULL,
			related_indicator TEXT NOT NULL,
			relationship_type TEXT NOT NULL,
			confidence REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (source_indicator, related_indicator),
			FOREIGN KEY (source_indicator) REFERENCES threat_indicators(value) ON DELETE CASCADE,
			FOREIGN KEY (related_indicator) REFERENCES threat_indicators(value) ON DELETE CASCADE
		)`,

		// Query performance table
		`CREATE TABLE IF NOT EXISTS query_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			query_value TEXT,
			query_type TEXT,
			found BOOLEAN,
			response_time_ms INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	tx, err := td.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (td *ThreatDatabase) AddIndicator(ctx context.Context, i models.Indicator) error {
	// Validate indicator
	if i.Value == "" {
		return fmt.Errorf("indicator value cannot be empty")
	}
	if i.Confidence < 0 || i.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}
	if i.Severity == "" {
		i.Severity = "medium"
	}
	if i.Source == "" {
		i.Source = "manual"
	}
	if i.LastSeen.IsZero() {
		i.LastSeen = time.Now()
	}
	if i.FirstSeen.IsZero() {
		i.FirstSeen = i.LastSeen
	}

	// Prepare tags as JSON array
	tagsJSON := "[]"
	if len(i.Tags) > 0 {
		tagsJSON = fmt.Sprintf(`["%s"]`, strings.Join(i.Tags, `","`))
	}

	query := `INSERT OR REPLACE INTO threat_indicators 
			  (value, type, source, confidence, severity, last_seen, first_seen, description, tags, "references", updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := td.db.ExecContext(ctx, query,
		i.Value,
		strings.ToLower(string(i.Type)),
		i.Source,
		i.Confidence,
		i.Severity,
		i.LastSeen,
		i.FirstSeen,
		i.Description,
		tagsJSON,
		strings.Join(i.References, ";"),
	)

	if err != nil {
		return fmt.Errorf("failed to add indicator: %w", err)
	}

	// Update cache
	td.mu.Lock()
	td.cache[i.Value] = &cacheEntry{
		indicator: &i,
		timestamp: time.Now(),
	}
	td.metrics.mu.Lock()
	td.metrics.TotalIndicators++
	td.metrics.mu.Unlock()
	td.mu.Unlock()

	td.logger.Printf("Added indicator: %s (type: %s, severity: %s)", i.Value, i.Type, i.Severity)
	return nil
}

func (td *ThreatDatabase) Lookup(ctx context.Context, value string) (*models.Indicator, error) {
	start := time.Now()
	
	// Check cache first
	td.mu.RLock()
	if entry, found := td.cache[value]; found {
		if time.Since(entry.timestamp) < 5*time.Minute {
			td.mu.RUnlock()
			td.metrics.mu.Lock()
			td.metrics.CacheHits++
			entry.hits++
			td.metrics.mu.Unlock()
			
			duration := time.Since(start)
			td.logger.Printf("Cache hit for %s (took %v)", value, duration)
			return entry.indicator, nil
		}
		// Cache expired
		delete(td.cache, value)
	}
	td.mu.RUnlock()

	td.metrics.mu.Lock()
	td.metrics.CacheMisses++
	td.metrics.mu.Unlock()

	// Query database
	query := `SELECT type, source, confidence, severity, last_seen, first_seen, description, tags, "references" 
			  FROM threat_indicators WHERE value = ?`

	row := td.db.QueryRowContext(ctx, query, value)

	var i models.Indicator
	var tagsJSON, refsStr string
	i.Value = value

	err := row.Scan(&i.Type, &i.Source, &i.Confidence, &i.Severity, &i.LastSeen, &i.FirstSeen,
		&i.Description, &tagsJSON, &refsStr)

	if err == sql.ErrNoRows {
		duration := time.Since(start)
		td.logger.Printf("Lookup miss for %s (took %v)", value, duration)
		return nil, nil
	}
	if err != nil {
		duration := time.Since(start)
		td.logger.Printf("Lookup error for %s: %v (took %v)", value, err, duration)
		return nil, fmt.Errorf("database lookup failed: %w", err)
	}

	// Parse tags and references
	i.Tags = parseJSONArray(tagsJSON)
	i.References = strings.Split(refsStr, ";")

	// Update cache
	td.mu.Lock()
	td.cache[value] = &cacheEntry{
		indicator: &i,
		timestamp: time.Now(),
		hits:      1,
	}
	td.mu.Unlock()

	duration := time.Since(start)
	td.logger.Printf("Database lookup for %s (took %v)", value, duration)
	
	// Update query metrics
	td.metrics.mu.Lock()
	td.metrics.TotalQueries++
	if td.metrics.TotalQueries == 1 {
		td.metrics.AverageQueryTime = duration
	} else {
		alpha := 0.1
		td.metrics.AverageQueryTime = time.Duration(
			float64(td.metrics.AverageQueryTime)*(1-alpha) + float64(duration)*alpha,
		)
	}
	td.metrics.mu.Unlock()

	return &i, nil
}

func (td *ThreatDatabase) BulkLookup(ctx context.Context, values []string) (map[string]*models.Indicator, error) {
	if len(values) == 0 {
		return make(map[string]*models.Indicator), nil
	}

	if len(values) > 1000 {
		return nil, fmt.Errorf("too many values for bulk lookup: %d (max 1000)", len(values))
	}

	start := time.Now()
	results := make(map[string]*models.Indicator)
	var uncached []string

	// Check cache first
	td.mu.RLock()
	for _, value := range values {
		if entry, found := td.cache[value]; found {
			if time.Since(entry.timestamp) < 5*time.Minute {
				results[value] = entry.indicator
				td.metrics.mu.Lock()
				td.metrics.CacheHits++
				entry.hits++
				td.metrics.mu.Unlock()
				continue
			}
			delete(td.cache, value)
		}
		uncached = append(uncached, value)
	}
	td.mu.RUnlock()

	if len(uncached) == 0 {
		duration := time.Since(start)
		td.logger.Printf("Bulk lookup all cache hits: %d values (took %v)", len(values), duration)
		return results, nil
	}

	td.metrics.mu.Lock()
	td.metrics.CacheMisses += int64(len(uncached))
	td.metrics.mu.Unlock()

	// Prepare batch query
	var placeholders []string
	var args []interface{}
	for _, value := range uncached {
		placeholders = append(placeholders, "?")
		args = append(args, value)
	}

	query := fmt.Sprintf(`SELECT value, type, source, confidence, severity, last_seen, first_seen, description, tags, "references" 
						   FROM threat_indicators WHERE value IN (%s)`, 
						   strings.Join(placeholders, ","))

	rows, err := td.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("bulk query failed: %w", err)
	}
	defer rows.Close()

	td.mu.Lock()
	for rows.Next() {
		var i models.Indicator
		var tagsJSON, refsStr string

		err := rows.Scan(&i.Value, &i.Type, &i.Source, &i.Confidence, &i.Severity, &i.LastSeen,
			&i.FirstSeen, &i.Description, &tagsJSON, &refsStr)
		if err != nil {
			continue
		}

		i.Tags = parseJSONArray(tagsJSON)
		i.References = strings.Split(refsStr, ";")
		
		results[i.Value] = &i
		td.cache[i.Value] = &cacheEntry{
			indicator: &i,
			timestamp: time.Now(),
			hits:      1,
		}
	}
	td.mu.Unlock()

	// Update metrics
	duration := time.Since(start)
	td.metrics.mu.Lock()
	td.metrics.TotalQueries++
	if td.metrics.TotalQueries == 1 {
		td.metrics.AverageQueryTime = duration
	} else {
		alpha := 0.1
		td.metrics.AverageQueryTime = time.Duration(
			float64(td.metrics.AverageQueryTime)*(1-alpha) + float64(duration)*alpha,
		)
	}
	td.metrics.mu.Unlock()

	td.logger.Printf("Bulk lookup: %d total, %d from cache, %d from DB (took %v)", 
		len(values), len(values)-len(uncached), len(uncached), duration)
	
	return results, nil
}

func (td *ThreatDatabase) GetStats(ctx context.Context) (*models.ThreatDBStats, error) {
	stats := &models.ThreatDBStats{
		Timestamp: time.Now(),
	}

	// Get count by type
	typeCountQuery := `SELECT type, COUNT(*) FROM threat_indicators GROUP BY type`
	rows, err := td.db.QueryContext(ctx, typeCountQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.CountByType = make(map[string]int64)
	for rows.Next() {
		var typ string
		var count int64
		rows.Scan(&typ, &count)
		stats.CountByType[typ] = count
		stats.TotalIndicators += count
	}

	// Get count by severity
	severityQuery := `SELECT severity, COUNT(*) FROM threat_indicators GROUP BY severity`
	rows2, err := td.db.QueryContext(ctx, severityQuery)
	if err != nil {
		return stats, nil
	}
	defer rows2.Close()

	stats.CountBySeverity = make(map[string]int64)
	for rows2.Next() {
		var severity string
		var count int64
		rows2.Scan(&severity, &count)
		stats.CountBySeverity[severity] = count
	}

	// Get recent activity
	recentQuery := `SELECT COUNT(*) FROM threat_indicators WHERE last_seen > datetime('now', '-7 days')`
	row := td.db.QueryRowContext(ctx, recentQuery)
	row.Scan(&stats.RecentActivity7d)

	// Get cache metrics
	td.metrics.mu.RLock()
	stats.CacheHits = td.metrics.CacheHits
	stats.CacheMisses = td.metrics.CacheMisses
	stats.TotalQueries = td.metrics.TotalQueries
	td.metrics.mu.RUnlock()

	td.logger.Printf("Stats collected: %d total indicators, %d recent (7d)", 
		stats.TotalIndicators, stats.RecentActivity7d)
	
	return stats, nil
}

func (td *ThreatDatabase) Cleanup(ctx context.Context, olderThanDays int) (int64, error) {
	if olderThanDays < 30 {
		olderThanDays = 90
	}

	query := `DELETE FROM threat_indicators WHERE last_seen < datetime('now', '-? days')`
	result, err := td.db.ExecContext(ctx, query, olderThanDays)
	if err != nil {
		return 0, fmt.Errorf("cleanup failed: %w", err)
	}

	rows, _ := result.RowsAffected()
	
	// Clear cache after cleanup
	td.mu.Lock()
	td.cache = make(map[string]*cacheEntry)
	td.mu.Unlock()

	td.logger.Printf("Cleanup removed %d indicators older than %d days", rows, olderThanDays)
	return rows, nil
}

func (td *ThreatDatabase) Close() error {
	// Log final metrics
	td.metrics.mu.RLock()
	td.logger.Printf("Closing ThreatDatabase - Final metrics: Queries=%d, CacheHits=%d, CacheMisses=%d, AvgQueryTime=%v",
		td.metrics.TotalQueries, td.metrics.CacheHits, td.metrics.CacheMisses, td.metrics.AverageQueryTime)
	td.metrics.mu.RUnlock()
	
	return td.db.Close()
}

func (td *ThreatDatabase) startCacheCleanup() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			td.cleanupCache()
		}
	}
}

func (td *ThreatDatabase) cleanupCache() {
	td.mu.Lock()
	defer td.mu.Unlock()

	now := time.Now()
	removed := 0
	for key, entry := range td.cache {
		if now.Sub(entry.timestamp) > 30*time.Minute || entry.hits == 0 {
			delete(td.cache, key)
			removed++
		}
	}
	
	if removed > 0 {
		td.logger.Printf("Cache cleanup removed %d entries", removed)
	}
}

// Helper function to parse JSON array
func parseJSONArray(jsonStr string) []string {
	if jsonStr == "" || jsonStr == "[]" {
		return []string{}
	}
	jsonStr = strings.Trim(jsonStr, "[]")
	if jsonStr == "" {
		return []string{}
	}
	return strings.Split(strings.ReplaceAll(jsonStr, `"`, ""), ",")
}
