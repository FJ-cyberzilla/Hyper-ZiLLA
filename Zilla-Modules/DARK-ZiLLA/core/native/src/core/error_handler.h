#pragma once

#include <exception>
#include <string>
#include <memory>
#include <vector>
#include <unordered_map>
#include <sstream>
#include <iomanip>
#include "../utils/logger.h"

namespace CodezillA {

    // Exception Hierarchy
    class AnalyzerException : public std::exception {
    protected:
        std::string message_;
        std::string component_;
        int error_code_;
        
    public:
        AnalyzerException(const std::string& message, const std::string& component = "", int error_code = 0)
            : message_(message), component_(component), error_code_(error_code) {}
        
        const char* what() const noexcept override { return message_.c_str(); }
        const std::string& getComponent() const { return component_; }
        int getErrorCode() const { return error_code_; }
        virtual std::string getFormattedMessage() const {
            return "[" + component_ + "] " + message_;
        }
    };

    class FileSystemException : public AnalyzerException {
    public:
        FileSystemException(const std::string& message, const std::string& file_path = "")
            : AnalyzerException(message, "FileSystem", 1000) {
            if (!file_path.empty()) {
                message_ += " - File: " + file_path;
            }
        }
    };

    class AnalysisException : public AnalyzerException {
    public:
        AnalysisException(const std::string& message, const std::string& language = "")
            : AnalyzerException(message, "Analysis", 2000) {
            if (!language.empty()) {
                message_ += " - Language: " + language;
            }
        }
    };

    class AIEngineException : public AnalyzerException {
    public:
        AIEngineException(const std::string& message, const std::string& model = "")
            : AnalyzerException(message, "AIEngine", 3000) {
            if (!model.empty()) {
                message_ += " - Model: " + model;
            }
        }
    };

    class ConfigurationException : public AnalyzerException {
    public:
        ConfigurationException(const std::string& message, const std::string& config_key = "")
            : AnalyzerException(message, "Configuration", 4000) {
            if (!config_key.empty()) {
                message_ += " - Key: " + config_key;
            }
        }
    };

    class PluginException : public AnalyzerException {
    public:
        PluginException(const std::string& message, const std::string& plugin_name = "")
            : AnalyzerException(message, "Plugin", 5000) {
            if (!plugin_name.empty()) {
                message_ += " - Plugin: " + plugin_name;
            }
        }
    };

    class MemoryException : public AnalyzerException {
    public:
        MemoryException(const std::string& message, size_t memory_usage = 0)
            : AnalyzerException(message, "Memory", 6000) {
            if (memory_usage > 0) {
                std::ostringstream oss;
                oss << " - Memory: " << (memory_usage / 1024 / 1024) << "MB";
                message_ += oss.str();
            }
        }
    };

    class TimeoutException : public AnalyzerException {
    public:
        TimeoutException(const std::string& message, int timeout_seconds = 0)
            : AnalyzerException(message, "Timeout", 7000) {
            if (timeout_seconds > 0) {
                message_ += " - Timeout: " + std::to_string(timeout_seconds) + "s";
            }
        }
    };

    // Error Handler Class
    class ErrorHandler {
    private:
        std::shared_ptr<Logger> logger_;
        std::unordered_map<int, int> error_counts_;
        bool recovery_enabled_;
        size_t max_memory_mb_;
        int operation_timeout_seconds_;
        
    public:
        ErrorHandler() 
            : logger_(std::make_shared<Logger>("ErrorHandler"))
            , recovery_enabled_(true)
            , max_memory_mb_(512)
            , operation_timeout_seconds_(30) {}
        
        void handleError(const std::string& component, const std::string& message, int error_code = 0) {
            // Log the error
            logger_->error("[" + component + "] " + message);
            
            // Track error counts
            error_counts_[error_code]++;
            
            // Check for critical error patterns
            checkForCriticalPatterns(component, message, error_code);
        }
        
        void handleException(const AnalyzerException& e) {
            handleError(e.getComponent(), e.what(), e.getErrorCode());
            
            // Attempt recovery for certain exception types
            if (recovery_enabled_) {
                attemptRecovery(e);
            }
        }
        
        template<typename Func, typename... Args>
        auto executeWithRecovery(Func&& func, Args&&... args) {
            try {
                return func(std::forward<Args>(args)...);
            } catch (const AnalyzerException& e) {
                handleException(e);
                throw; // Re-throw after handling
            } catch (const std::exception& e) {
                handleError("Unknown", std::string("Standard exception: ") + e.what());
                throw AnalyzerException(e.what(), "Unknown", 9999);
            } catch (...) {
                handleError("Unknown", "Unknown exception occurred");
                throw AnalyzerException("Unknown exception", "Unknown", 9999);
            }
        }
        
        template<typename Func, typename... Args>
        auto executeWithTimeout(Func&& func, Args&&... args) {
            // For now, simple execution - could be enhanced with actual timeout
            return executeWithRecovery(std::forward<Func>(func), std::forward<Args>(args)...);
        }
        
        void enableRecovery(bool enable) { recovery_enabled_ = enable; }
        void setMemoryLimit(size_t mb) { max_memory_mb_ = mb; }
        void setTimeout(int seconds) { operation_timeout_seconds_ = seconds; }
        
        std::unordered_map<int, int> getErrorStatistics() const { return error_counts_; }
        void resetErrorCounts() { error_counts_.clear(); }
        
    private:
        void checkForCriticalPatterns(const std::string& component, const std::string& message, int error_code) {
            // Check for memory-related errors
            if (message.find("memory") != std::string::npos || 
                message.find("alloc") != std::string::npos) {
                logger_->warn("Memory-related error detected - consider increasing memory limits");
            }
            
            // Check for file system errors
            if (message.find("permission") != std::string::npos ||
                message.find("access") != std::string::npos) {
                logger_->warn("File permission issue detected");
            }
            
            // Check for network-related errors (for future cloud features)
            if (message.find("network") != std::string::npos ||
                message.find("connection") != std::string::npos) {
                logger_->warn("Network-related error detected");
            }
        }
        
        void attemptRecovery(const AnalyzerException& e) {
            switch (e.getErrorCode()) {
                case 1000: // FileSystem
                    logger_->info("Attempting filesystem error recovery...");
                    // Could implement file handle cleanup, retry logic, etc.
                    break;
                    
                case 6000: // Memory
                    logger_->info("Attempting memory error recovery...");
                    // Could implement garbage collection, cache clearing, etc.
                    break;
                    
                case 7000: // Timeout
                    logger_->info("Attempting timeout recovery...");
                    // Could implement operation retry with backoff
                    break;
                    
                default:
                    logger_->debug("No specific recovery strategy for error code: " + 
                                 std::to_string(e.getErrorCode()));
                    break;
            }
        }
    };
}
