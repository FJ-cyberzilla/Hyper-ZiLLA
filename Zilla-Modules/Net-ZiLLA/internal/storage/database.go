package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"net-zilla/internal/models"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dataSource string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &Database{db: db}, nil
}

func initSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS analyses (
			id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			threat_level TEXT NOT NULL,
			threat_score INTEGER NOT NULL,
			analysis_data TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_analyses_created_at ON analyses(created_at)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) SaveAnalysis(ctx context.Context, analysis *models.ThreatAnalysis) error {
	query := `INSERT INTO analyses (id, url, threat_level, threat_score, analysis_data) 
	          VALUES (?, ?, ?, ?, ?)`

	analysisData, err := json.Marshal(analysis)
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, query,
		analysis.AnalysisID,
		analysis.URL,
		string(analysis.ThreatLevel),
		analysis.ThreatScore,
		string(analysisData),
	)
	return err
}

func (d *Database) GetAnalysisByID(ctx context.Context, id string) (*models.ThreatAnalysis, error) {
	query := `SELECT analysis_data FROM analyses WHERE id = ?`
	var data string
	err := d.db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		return nil, err
	}
	var analysis models.ThreatAnalysis
	if err := json.Unmarshal([]byte(data), &analysis); err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *Database) GetAnalysisHistory(ctx context.Context, limit int) ([]*models.ThreatAnalysis, error) {
	query := `SELECT analysis_data FROM analyses ORDER BY created_at DESC LIMIT ?`
	rows, err := d.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.ThreatAnalysis
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var analysis models.ThreatAnalysis
		if err := json.Unmarshal([]byte(data), &analysis); err == nil {
			history = append(history, &analysis)
		}
	}
	return history, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
