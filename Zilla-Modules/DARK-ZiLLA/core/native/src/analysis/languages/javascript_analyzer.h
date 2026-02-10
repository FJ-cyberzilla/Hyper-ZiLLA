#pragma once
#include "base_analyzer.h"
#include <string>
#include <vector>

namespace CodezillA {

class JavaScriptAnalyzer : public BaseAnalyzer {
public:
    JavaScriptAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
        : BaseAnalyzer(std::move(error_handler), std::move(logger)) {}
    
    std::string getLanguageName() const override {
        return "JavaScript";
    }
    
    std::vector<std::string> getSupportedExtensions() const override {
        return {".js", ".jsx", ".ts", ".tsx"};
    }
    
    bool isSupportedFile(const std::string& file_path) const override {
        return file_path.size() >= 3 && 
               (file_path.substr(file_path.size() - 3) == ".js" ||
                file_path.substr(file_path.size() - 4) == ".jsx" ||
                file_path.substr(file_path.size() - 3) == ".ts" ||
                file_path.substr(file_path.size() - 4) == ".tsx");
    }
    
    std::vector<AnalysisResult> analyze(const std::string& file_path) override {
        std::vector<AnalysisResult> results;
        results.emplace_back(file_path, "JS_PLACEHOLDER", "JavaScript analysis not yet implemented", "INFO");
        return results;
    }

    std::vector<AnalysisResult> analyzeCode(const std::string& code, const std::string& file_name) override {
        std::vector<AnalysisResult> results;
        results.emplace_back(file_name, "JS_PLACEHOLDER", "JavaScript code analysis not yet implemented", "INFO");
        return results;
    }
};

} // namespace CodezillA
