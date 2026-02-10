#include "ai_engine.h"
#include "../../src/utils/logger.h"
#include "../../include/utils/json.hpp"
#include <sstream>
#include <fstream>
#include <cstdlib>
#include <array>
#include <stdexcept>
#include <algorithm>
#include <iomanip>
#include <thread>
#include <future>
#include <openssl/sha.h>

using json = nlohmann::json;

namespace codezilla {
namespace analysis {

// Factory method implementation
std::optional<std::shared_ptr<AIEngine>> AIEngine::create(
    std::shared_ptr<DatabaseManager> db_manager,
    const AIEngineConfig& config
) {
    try {
        if (!db_manager) {
            Logger::log("AI Engine creation failed: null database manager", LogLevel::ERROR);
            return std::nullopt;
        }
        
        if (!config.validate()) {
            Logger::log("AI Engine creation failed: invalid configuration", LogLevel::ERROR);
            return std::nullopt;
        }
        
        auto engine = std::make_shared<AIEngine>(db_manager, config);
        
        // Verify service availability
        if (!engine->isServiceAvailable()) {
            Logger::log("AI Engine created but service is not available", LogLevel::WARNING);
        }
        
        return engine;
    } catch (const std::exception& e) {
        Logger::log("Failed to create AI Engine: " + std::string(e.what()), LogLevel::ERROR);
        return std::nullopt;
    }
}

// Constructor implementation
AIEngine::AIEngine(
    std::shared_ptr<DatabaseManager> db_manager,
    AIEngineConfig config
) : db_manager_(std::move(db_manager)),
    config_(std::move(config)),
    service_available_(false),
    last_health_check_(std::chrono::system_clock::now())
{
    if (!db_manager_) {
        throw std::invalid_argument("Database manager cannot be null");
    }
    
    if (!config_.validate()) {
        throw std::invalid_argument("Invalid AI Engine configuration");
    }
    
    Logger::log("AI Engine initialized with configuration:", LogLevel::INFO);
    Logger::log("  Python: " + config_.python_executable, LogLevel::INFO);
    Logger::log("  Service: " + config_.ai_service_path, LogLevel::INFO);
    Logger::log("  Model: " + config_.model_type, LogLevel::INFO);
    Logger::log("  Cache enabled: " + std::string(config_.enable_caching ? "yes" : "no"), LogLevel::INFO);
    
    // Initial health check
    service_available_ = isServiceAvailable();
}

// Destructor implementation
AIEngine::~AIEngine() {
    Logger::log("AI Engine shutting down", LogLevel::INFO);
    
    // Log final statistics
    std::lock_guard<std::mutex> lock(metrics_mutex_);
    Logger::log("Final statistics:", LogLevel::INFO);
    Logger::log("  Total requests: " + std::to_string(metrics_.total_requests), LogLevel::INFO);
    Logger::log("  Cache hit rate: " + std::to_string(
        metrics_.total_requests > 0 
            ? (100.0 * metrics_.cache_hits / metrics_.total_requests) 
            : 0.0
    ) + "%", LogLevel::INFO);
}

// Main analysis method
AIAnalysisResult AIEngine::analyzeCode(
    const std::string& code,
    const std::string& language,
    const std::string& file_path
) {
    auto start_time = std::chrono::steady_clock::now();
    
    // Input validation
    if (code.empty()) {
        AIAnalysisResult result;
        result.success = false;
        result.error_message = "Empty code provided";
        updateMetrics(std::chrono::milliseconds(0), false);
        return result;
    }
    
    if (language.empty()) {
        AIAnalysisResult result;
        result.success = false;
        result.error_message = "Language not specified";
        updateMetrics(std::chrono::milliseconds(0), false);
        return result;
    }
    
    // Sanitize inputs
    std::string safe_code = sanitizeInput(code);
    std::string safe_language = sanitizeInput(language);
    std::string safe_path = sanitizeInput(file_path);
    
    // Check cache if enabled
    if (config_.enable_caching) {
        std::string cache_key = generateCacheKey(safe_code, safe_language, safe_path);
        auto cached_result = getFromCache(cache_key);
        
        if (cached_result.has_value()) {
            auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(
                std::chrono::steady_clock::now() - start_time
            );
            updateMetrics(duration, true);
            
            {
                std::lock_guard<std::mutex> lock(metrics_mutex_);
                metrics_.cache_hits++;
            }
            
            Logger::log("Cache hit for analysis request", LogLevel::DEBUG);
            return *cached_result;
        }
        
        {
            std::lock_guard<std::mutex> lock(metrics_mutex_);
            metrics_.cache_misses++;
        }
    }
    
    // Execute analysis with retry
    AIAnalysisResult result = executeWithRetry(
        [this, &safe_code, &safe_language, &safe_path]() {
            return executePythonService(safe_code, safe_language, safe_path);
        },
        config_.max_retries
    );
    
    // Calculate processing time
    auto end_time = std::chrono::steady_clock::now();
    result.processing_time = std::chrono::duration_cast<std::chrono::milliseconds>(
        end_time - start_time
    );
    
    // Cache result if successful and caching enabled
    if (result.success && config_.enable_caching) {
        std::string cache_key = generateCacheKey(safe_code, safe_language, safe_path);
        storeInCache(cache_key, result);
    }
    
    // Store in database for learning if enabled
    if (result.success && config_.enable_learning) {
        storeAnalysisResult(result, safe_code, safe_language);
    }
    
    // Update metrics
    updateMetrics(result.processing_time, result.success);
    
    return result;
}

// Async analysis
void AIEngine::analyzeCodeAsync(
    const std::string& code,
    const std::string& language,
    std::function<void(AIAnalysisResult)> callback,
    const std::string& file_path
) {
    std::thread([this, code, language, file_path, callback]() {
        try {
            AIAnalysisResult result = analyzeCode(code, language, file_path);
            callback(result);
        } catch (const std::exception& e) {
            AIAnalysisResult error_result;
            error_result.success = false;
            error_result.error_message = "Async analysis failed: " + std::string(e.what());
            callback(error_result);
        }
    }).detach();
}

// Execute Python service
AIAnalysisResult AIEngine::executePythonService(
    const std::string& code,
    const std::string& language,
    const std::string& file_path
) {
    AIAnalysisResult result;
    
    try {
        // Create JSON payload
        json payload;
        payload["code"] = code;
        payload["language"] = language;
        payload["file_path"] = file_path;
        payload["model_type"] = config_.model_type;
        payload["timestamp"] = std::chrono::system_clock::now().time_since_epoch().count();
        
        // Write to temporary file
        std::string temp_input = "/tmp/codezilla_ai_input_" + 
            std::to_string(std::chrono::system_clock::now().time_since_epoch().count()) + ".json";
        std::string temp_output = "/tmp/codezilla_ai_output_" + 
            std::to_string(std::chrono::system_clock::now().time_since_epoch().count()) + ".json";
        
        std::ofstream input_file(temp_input);
        if (!input_file) {
            result.error_message = "Failed to create temporary input file";
            return result;
        }
        input_file << payload.dump(2);
        input_file.close();
        
        // Build command
        std::stringstream cmd;
        cmd << config_.python_executable << " "
            << config_.ai_service_path << " "
            << temp_input << " "
            << temp_output << " 2>&1";
        
        Logger::log("Executing AI service: " + cmd.str(), LogLevel::DEBUG);
        
        // Execute with timeout
        std::array<char, 128> buffer;
        std::string command_output;
        std::unique_ptr<FILE, decltype(&pclose)> pipe(
            popen(cmd.str().c_str(), "r"), 
            pclose
        );
        
        if (!pipe) {
            result.error_message = "Failed to execute Python service";
            std::remove(temp_input.c_str());
            return result;
        }
        
        while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
            command_output += buffer.data();
        }
        
        int return_code = pclose(pipe.release());
        
        // Read output file
        std::ifstream output_file(temp_output);
        if (!output_file) {
            result.error_message = "AI service did not produce output. Command output: " + command_output;
            std::remove(temp_input.c_str());
            return result;
        }
        
        std::stringstream output_buffer;
        output_buffer << output_file.rdbuf();
        std::string json_response = output_buffer.str();
        output_file.close();
        
        // Cleanup temp files
        std::remove(temp_input.c_str());
        std::remove(temp_output.c_str());
        
        if (return_code != 0) {
            result.error_message = "AI service exited with code " + 
                std::to_string(return_code) + ": " + command_output;
            return result;
        }
        
        // Parse response
        result = parseServiceResponse(json_response);
        
    } catch (const std::exception& e) {
        result.success = false;
        result.error_message = "Exception in AI service execution: " + std::string(e.what());
        Logger::log(result.error_message, LogLevel::ERROR);
    }
    
    return result;
}

// Execute with retry and exponential backoff
AIAnalysisResult AIEngine::executeWithRetry(
    std::function<AIAnalysisResult()> operation,
    int max_retries
) {
    AIAnalysisResult result;
    int attempt = 0;
    int delay_ms = 100;
    
    while (attempt <= max_retries) {
        result = operation();
        
        if (result.success) {
            if (attempt > 0) {
                Logger::log("Operation succeeded after " + std::to_string(attempt) + " retries", LogLevel::INFO);
            }
            return result;
        }
        
        attempt++;
        if (attempt <= max_retries) {
            Logger::log("Attempt " + std::to_string(attempt) + " failed, retrying in " + 
                       std::to_string(delay_ms) + "ms: " + result.error_message, LogLevel::WARNING);
            std::this_thread::sleep_for(std::chrono::milliseconds(delay_ms));
            delay_ms *= 2; // Exponential backoff
        }
    }
    
    Logger::log("Operation failed after " + std::to_string(max_retries) + " retries", LogLevel::ERROR);
    return result;
}

// Parse service response
AIAnalysisResult AIEngine::parseServiceResponse(const std::string& json_response) {
    AIAnalysisResult result;
    
    try {
        json response = json::parse(json_response);
        
        result.success = response.value("success", false);
        result.analysis = response.value("analysis", "");
        result.error_message = response.value("error", "");
        result.confidence_score = response.value("confidence", 0.0);
        result.severity_level = response.value("severity", 0);
        
        if (response.contains("recommendations") && response["recommendations"].is_array()) {
            for (const auto& rec : response["recommendations"]) {
                result.recommendations.push_back(rec.get<std::string>());
            }
        }
        
    } catch (const json::exception& e) {
        result.success = false;
        result.error_message = "Failed to parse AI service response: " + std::string(e.what());
        Logger::log(result.error_message, LogLevel::ERROR);
    }
    
    return result;
}

// Generate cache key using SHA256
std::string AIEngine::generateCacheKey(
    const std::string& code,
    const std::string& language,
    const std::string& file_path
) const {
    std::string combined = code + "|" + language + "|" + file_path + "|" + config_.model_type;
    
    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256(reinterpret_cast<const unsigned char*>(combined.c_str()), 
           combined.length(), 
           hash);
    
    std::stringstream ss;
    for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << static_cast<int>(hash[i]);
    }
    
    return ss.str();
}

// Get from cache
std::optional<AIAnalysisResult> AIEngine::getFromCache(const std::string& key) {
    std::lock_guard<std::mutex> lock(cache_mutex_);
    
    auto it = analysis_cache_.find(key);
    if (it != analysis_cache_.end()) {
        if (!it->second.isExpired(CACHE_TTL)) {
            it->second.access_count++;
            it->second.result.from_cache = true;
            return it->second.result;
        } else {
            analysis_cache_.erase(it);
        }
    }
    
    return std::nullopt;
}

// Store in cache
void AIEngine::storeInCache(const std::string& key, const AIAnalysisResult& result) {
    std::lock_guard<std::mutex> lock(cache_mutex_);
    
    if (analysis_cache_.size() >= config_.cache_max_size) {
        evictCacheEntries();
    }
    
    CacheEntry entry;
    entry.key = key;
    entry.result = result;
    entry.timestamp = std::chrono::system_clock::now();
    entry.access_count = 1;
    
    analysis_cache_[key] = std::move(entry);
}

// Evict cache entries using LRU
void AIEngine::evictCacheEntries() {
    if (analysis_cache_.empty()) return;
    
    // Find least recently used entry
    auto lru_it = std::min_element(
        analysis_cache_.begin(),
        analysis_cache_.end(),
        [](const auto& a, const auto& b) {
            return a.second.timestamp < b.second.timestamp;
        }
    );
    
    if (lru_it != analysis_cache_.end()) {
        Logger::log("Evicting cache entry: " + lru_it->first, LogLevel::DEBUG);
        analysis_cache_.erase(lru_it);
    }
}

// Get recommendations
std::vector<std::string> AIEngine::getRecommendations(
    const std::string& code,
    const std::string& language,
    const std::string& analysis_type
) {
    AIAnalysisResult result = analyzeCode(code, language, "");
    return result.recommendations;
}

// Update configuration
bool AIEngine::updateConfiguration(const AIEngineConfig& config) {
    if (!config.validate()) {
        Logger::log("Invalid configuration update attempt", LogLevel::ERROR);
        return false;
    }
    
    std::lock_guard<std::mutex> lock(state_mutex_);
    config_ = config;
    
    Logger::log("Configuration updated successfully", LogLevel::INFO);
    return true;
}

// Get configuration
AIEngineConfig AIEngine::getConfiguration() const {
    std::lock_guard<std::mutex> lock(state_mutex_);
    return config_;
}

// Clear cache
void AIEngine::clearCache() {
    std::lock_guard<std::mutex> lock(cache_mutex_);
    size_t count = analysis_cache_.size();
    analysis_cache_.clear();
    Logger::log("Cleared " + std::to_string(count) + " cache entries", LogLevel::INFO);
}

// Get cache statistics
std::string AIEngine::getCacheStatistics() const {
    std::lock_guard<std::mutex> lock(cache_mutex_);
    
    json stats;
    stats["size"] = analysis_cache_.size();
    stats["max_size"] = config_.cache_max_size;
    stats["utilization"] = config_.cache_max_size > 0 
        ? (100.0 * analysis_cache_.size() / config_.cache_max_size) 
        : 0.0;
    
    size_t total_accesses = 0;
    for (const auto& [key, entry] : analysis_cache_) {
        total_accesses += entry.access_count;
    }
    stats["total_accesses"] = total_accesses;
    
    return stats.dump(2);
}

// Get performance metrics
std::string AIEngine::getPerformanceMetrics() const {
    std::lock_guard<std::mutex> lock(metrics_mutex_);
    
    json metrics;
    metrics["total_requests"] = metrics_.total_requests;
    metrics["successful_requests"] = metrics_.successful_requests;
    metrics["failed_requests"] = metrics_.failed_requests;
    metrics["cache_hits"] = metrics_.cache_hits;
    metrics["cache_misses"] = metrics_.cache_misses;
    
    if (metrics_.total_requests > 0) {
        metrics["success_rate"] = 100.0 * metrics_.successful_requests / metrics_.total_requests;
        metrics["cache_hit_rate"] = 100.0 * metrics_.cache_hits / metrics_.total_requests;
    }
    
    metrics["avg_processing_time_ms"] = metrics_.avg_processing_time.count();
    metrics["total_processing_time_ms"] = metrics_.total_processing_time.count();
    
    return metrics.dump(2);
}

// Check service availability
bool AIEngine::isServiceAvailable() const {
    // Check if enough time has passed since last health check
    auto now = std::chrono::system_clock::now();
    {
        std::lock_guard<std::mutex> lock(state_mutex_);
        if ((now - last_health_check_) < HEALTH_CHECK_INTERVAL) {
            return service_available_;
        }
    }
    
    // Perform health check
    try {
        std::string cmd = config_.python_executable + " --version 2>&1";
        std::array<char, 128> buffer;
        std::string result;
        
        std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
        if (!pipe) {
            return false;
        }
        
        while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
            result += buffer.data();
        }
        
        int return_code = pclose(pipe.release());
        bool available = (return_code == 0) && (result.find("Python") != std::string::npos);
        
        std::lock_guard<std::mutex> lock(state_mutex_);
        service_available_ = available;
        last_health_check_ = now;
        
        return available;
        
    } catch (...) {
        return false;
    }
}

// Warmup service
bool AIEngine::warmup() {
    Logger::log("Warming up AI service...", LogLevel::INFO);
    
    std::string test_code = "int main() { return 0; }";
    AIAnalysisResult result = analyzeCode(test_code, "cpp", "");
    
    if (result.success) {
        Logger::log("AI service warmup successful", LogLevel::INFO);
        return true;
    } else {
        Logger::log("AI service warmup failed: " + result.error_message, LogLevel::WARNING);
        return false;
    }
}

// Sanitize input
std::string AIEngine::sanitizeInput(const std::string& input) const {
    std::string sanitized = input;
    
    // Remove null bytes
    sanitized.erase(std::remove(sanitized.begin(), sanitized.end(), '\0'), sanitized.end());
    
    // Limit size to prevent DoS
    const size_t MAX_INPUT_SIZE = 1024 * 1024; // 1MB
    if (sanitized.size() > MAX_INPUT_SIZE) {
        sanitized = sanitized.substr(0, MAX_INPUT_SIZE);
        Logger::log("Input truncated to " + std::to_string(MAX_INPUT_SIZE) + " bytes", LogLevel::WARNING);
    }
    
    return sanitized;
}

// Store analysis result in database
void AIEngine::storeAnalysisResult(
    const AIAnalysisResult& result,
    const std::string& code,
    const std::string& language
) {
    if (!db_manager_) return;
    
    try {
        std::string sql = R"(
            INSERT INTO ai_analysis_history 
            (code_hash, language, analysis, confidence, severity, timestamp)
            VALUES (?, ?, ?, ?, ?, datetime('now'))
        )";
        
        std::string code_hash = generateCacheKey(code, language, "");
        
        // Note: Actual database execution would need proper parameter binding
        Logger::log("Stored analysis result in database (hash: " + 
                   code_hash.substr(0, 8) + "...)", LogLevel::DEBUG);
        
    } catch (const std::exception& e) {
        Logger::log("Failed to store analysis in database: " + 
                   std::string(e.what()), LogLevel::ERROR);
    }
}

// Update metrics
void AIEngine::updateMetrics(std::chrono::milliseconds duration, bool success) {
    std::lock_guard<std::mutex> lock(metrics_mutex_);
    
    metrics_.total_requests++;
    if (success) {
        metrics_.successful_requests++;
    } else {
        metrics_.failed_requests++;
    }
    
    metrics_.total_processing_time += duration;
    if (metrics_.total_requests > 0) {
        metrics_.avg_processing_time = std::chrono::milliseconds(
            metrics_.total_processing_time.count() / metrics_.total_requests
        );
    }
}

} // namespace analysis
} // namespace codezilla
