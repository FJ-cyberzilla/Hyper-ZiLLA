#ifndef BASE_ANALYZER_H
#define BASE_ANALYZER_H

#include <analysis/analysis_result.h>
#include "utils/logger.h"    // Include for Logger
#include "core/error_handler.h" // Include for ErrorHandler
#include <string>
#include <vector>
#include <unordered_set>
#include <memory> // For std::shared_ptr

namespace CodezillA {

class BaseAnalyzer {
protected:
    std::shared_ptr<Logger> logger_;
    std::shared_ptr<ErrorHandler> error_handler_;

public:
    BaseAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
        : error_handler_(std::move(error_handler)), logger_(std::move(logger)) {}
    virtual ~BaseAnalyzer() = default;
    virtual bool isSupportedFile(const std::string& file_path) const = 0;
    virtual std::vector<AnalysisResult> analyze(const std::string& file_path) = 0;
    virtual std::vector<AnalysisResult> analyzeCode(const std::string& code, const std::string& file_name) = 0;
    virtual bool canAutoFix(const AnalysisResult& result) const { return false; }
    virtual bool applyFix(const std::string& file_path, const AnalysisResult& result) { return false; }
    // Overload for AI-driven fixes
    virtual bool applyFix(const std::string& file_path, const AnalysisResult& result, const std::string& ai_suggested_fix) {
        // Default implementation: if AI suggested fix is empty, fallback to rule-based fix
        if (ai_suggested_fix.empty()) {
            return applyFix(file_path, result);
        }
        // Otherwise, a derived class should handle the AI suggested fix
        return false; 
    }
    
    // Add virtual methods that were missing in the base class
    virtual std::string getLanguageName() const { return "Unknown"; }
    virtual std::vector<std::string> getSupportedExtensions() const { return {}; }

    void setLogger(std::shared_ptr<Logger> logger) {
        logger_ = std::move(logger);
    }

    void setErrorHandler(std::shared_ptr<ErrorHandler> error_handler) {
        error_handler_ = std::move(error_handler);
    }
};

} // namespace CodezillA

#endif // BASE_ANALYZER_H