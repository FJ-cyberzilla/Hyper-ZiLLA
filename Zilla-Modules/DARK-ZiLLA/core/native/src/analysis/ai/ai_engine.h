#ifndef AI_ENGINE_H
#define AI_ENGINE_H

#include <string>
#include <memory>
#include <vector>
#include <mutex>
#include <optional>
#include <functional>
#include "../../include/db/database_manager.h"
#include "../../include/analysis/analysis_result.h"

namespace codezilla {
namespace analysis {

/**
 * @brief Configuration structure for AI Engine
 */
struct AIEngineConfig {
    std::string python_executable = "python3";
    std::string ai_service_path = "src/analysis/ai/ai_service.py";
    std::string model_type = "advanced";
    int timeout_seconds = 30;
    int max_retries = 3;
    bool enable_caching = true;
    bool enable_learning = true;
    size_t cache_max_size = 1000;
    
    bool validate() const {
        return !python_executable.empty() && 
               !ai_service_path.empty() && 
               timeout_seconds > 0 && 
               max_retries >= 0;
    }
};

/**
 * @brief Result structure for AI analysis operations
 */
struct AIAnalysisResult {
    bool success = false;
    std::string analysis;
    std::string error_message;
    double confidence_score = 0.0;
    int severity_level = 0;
    std::vector<std::string> recommendations;
    std::chrono::milliseconds processing_time{0};
    bool from_cache = false;
    
    explicit operator bool() const { return success; }
};

/**
 * @brief Cache entry for AI analysis results
 */
struct CacheEntry {
    std::string key;
    AIAnalysisResult result;
    std::chrono::system_clock::time_point timestamp;
    size_t access_count = 0;
    
    bool isExpired(std::chrono::seconds ttl) const {
        auto now = std::chrono::system_clock::now();
        return (now - timestamp) > ttl;
    }
};

/**
 * @brief Enterprise-grade AI Engine for code analysis
 * 
 * This class provides thread-safe AI-powered code analysis with:
 * - Configurable execution environments
 * - Intelligent caching with LRU eviction
 * - Automatic retry with exponential backoff
 * - Performance monitoring and metrics
 * - Database integration for persistent learning
 */
class AIEngine {
public:
    /**
     * @brief Factory method to create AIEngine with configuration
     * @param db_manager Database manager for persistent storage
     * @param config AI engine configuration
     * @return Optional shared_ptr to AIEngine, nullopt on failure
     */
    static std::optional<std::shared_ptr<AIEngine>> create(
        std::shared_ptr<DatabaseManager> db_manager,
        const AIEngineConfig& config = AIEngineConfig{}
    );
    
    /**
     * @brief Constructor with explicit configuration
     * @param db_manager Database manager shared pointer
     * @param config AI engine configuration
     * @throws std::invalid_argument if db_manager is null or config invalid
     */
    explicit AIEngine(
        std::shared_ptr<DatabaseManager> db_manager,
        AIEngineConfig config = AIEngineConfig{}
    );
    
    /**
     * @brief Destructor - ensures cleanup of resources
     */
    ~AIEngine();
    
    // Disable copy operations for safety
    AIEngine(const AIEngine&) = delete;
    AIEngine& operator=(const AIEngine&) = delete;
    
    // Enable move operations
    AIEngine(AIEngine&&) noexcept = default;
    AIEngine& operator=(AIEngine&&) noexcept = default;
    
    /**
     * @brief Analyze code with AI capabilities
     * @param code Source code to analyze
     * @param language Programming language
     * @param file_path Optional file path for context
     * @return AIAnalysisResult containing analysis or error details
     */
    AIAnalysisResult analyzeCode(
        const std::string& code,
        const std::string& language,
        const std::string& file_path = ""
    );
    
    /**
     * @brief Analyze code asynchronously
     * @param code Source code to analyze
     * @param language Programming language
     * @param callback Callback function for result
     * @param file_path Optional file path for context
     */
    void analyzeCodeAsync(
        const std::string& code,
        const std::string& language,
        std::function<void(AIAnalysisResult)> callback,
        const std::string& file_path = ""
    );
    
    /**
     * @brief Get AI-powered recommendations for code improvement
     * @param code Source code
     * @param language Programming language
     * @param analysis_type Type of analysis ("security", "performance", "quality")
     * @return Vector of recommendations
     */
    std::vector<std::string> getRecommendations(
        const std::string& code,
        const std::string& language,
        const std::string& analysis_type = "general"
    );
    
    /**
     * @brief Update configuration at runtime
     * @param config New configuration
     * @return true if configuration was successfully updated
     */
    bool updateConfiguration(const AIEngineConfig& config);
    
    /**
     * @brief Get current configuration
     * @return Current AIEngineConfig
     */
    AIEngineConfig getConfiguration() const;
    
    /**
     * @brief Clear analysis cache
     */
    void clearCache();
    
    /**
     * @brief Get cache statistics
     * @return JSON string with cache statistics
     */
    std::string getCacheStatistics() const;
    
    /**
     * @brief Get engine performance metrics
     * @return JSON string with performance metrics
     */
    std::string getPerformanceMetrics() const;
    
    /**
     * @brief Check if AI service is available and responsive
     * @return true if AI service is operational
     */
    bool isServiceAvailable() const;
    
    /**
     * @brief Warm up the AI service (pre-load models, etc.)
     * @return true if warmup successful
     */
    bool warmup();

private:
    /**
     * @brief Execute Python AI service with retry logic
     * @param code Source code
     * @param language Programming language
     * @param file_path File path for context
     * @return AIAnalysisResult
     */
    AIAnalysisResult executePythonService(
        const std::string& code,
        const std::string& language,
        const std::string& file_path
    );
    
    /**
     * @brief Execute with exponential backoff retry
     * @param operation Operation to execute
     * @param max_retries Maximum number of retries
     * @return AIAnalysisResult
     */
    AIAnalysisResult executeWithRetry(
        std::function<AIAnalysisResult()> operation,
        int max_retries
    );
    
    /**
     * @brief Generate cache key for analysis request
     * @param code Source code
     * @param language Programming language
     * @param file_path File path
     * @return Cache key string
     */
    std::string generateCacheKey(
        const std::string& code,
        const std::string& language,
        const std::string& file_path
    ) const;
    
    /**
     * @brief Get result from cache if available
     * @param key Cache key
     * @return Optional AIAnalysisResult
     */
    std::optional<AIAnalysisResult> getFromCache(const std::string& key);
    
    /**
     * @brief Store result in cache
     * @param key Cache key
     * @param result Analysis result
     */
    void storeInCache(const std::string& key, const AIAnalysisResult& result);
    
    /**
     * @brief Evict oldest entries from cache using LRU
     */
    void evictCacheEntries();
    
    /**
     * @brief Parse JSON response from Python service
     * @param json_response JSON string
     * @return AIAnalysisResult
     */
    AIAnalysisResult parseServiceResponse(const std::string& json_response);
    
    /**
     * @brief Store analysis result in database for learning
     * @param result Analysis result
     * @param code Source code
     * @param language Programming language
     */
    void storeAnalysisResult(
        const AIAnalysisResult& result,
        const std::string& code,
        const std::string& language
    );
    
    /**
     * @brief Sanitize input to prevent injection attacks
     * @param input Input string
     * @return Sanitized string
     */
    std::string sanitizeInput(const std::string& input) const;
    
    /**
     * @brief Update performance metrics
     * @param duration Processing duration
     * @param success Whether operation was successful
     */
    void updateMetrics(std::chrono::milliseconds duration, bool success);

    // Member variables
    std::shared_ptr<DatabaseManager> db_manager_;
    AIEngineConfig config_;
    
    // Thread-safe cache
    mutable std::mutex cache_mutex_;
    std::unordered_map<std::string, CacheEntry> analysis_cache_;
    
    // Performance metrics
    mutable std::mutex metrics_mutex_;
    struct Metrics {
        size_t total_requests = 0;
        size_t successful_requests = 0;
        size_t failed_requests = 0;
        size_t cache_hits = 0;
        size_t cache_misses = 0;
        std::chrono::milliseconds total_processing_time{0};
        std::chrono::milliseconds avg_processing_time{0};
    } metrics_;
    
    // Service state
    mutable std::mutex state_mutex_;
    bool service_available_;
    std::chrono::system_clock::time_point last_health_check_;
    
    static constexpr std::chrono::seconds CACHE_TTL{3600}; // 1 hour
    static constexpr std::chrono::seconds HEALTH_CHECK_INTERVAL{300}; // 5 minutes
};

} // namespace analysis
} // namespace codezilla

#endif // AI_ENGINE_H
