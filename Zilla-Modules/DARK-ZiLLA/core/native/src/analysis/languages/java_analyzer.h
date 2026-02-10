#pragma once
#include "base_analyzer.h"
#include <string>
#include <vector>

namespace CodezillA {

class JavaAnalyzer : public BaseAnalyzer {
public:
    JavaAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
        : BaseAnalyzer(std::move(error_handler), std::move(logger)) {}
    
    std::string getLanguageName() const override {
        return "Java";
    }
    
    std::vector<std::string> getSupportedExtensions() const override {
        return {".java"};
    }
    
    bool isSupportedFile(const std::string& file_path) const override {
        return file_path.size() >= 5 && file_path.substr(file_path.size() - 5) == ".java";
    }
    
    std::vector<AnalysisResult> analyze(const std::string& file_path) override {
        std::vector<AnalysisResult> results;
        results.emplace_back(file_path, "JAVA_PLACEHOLDER", "Java analysis not yet implemented", "INFO");
        return results;
    }

    std::vector<AnalysisResult> analyzeCode(const std::string& code, const std::string& file_name) override {
        std::vector<AnalysisResult> results;
        results.emplace_back(file_name, "JAVA_PLACEHOLDER", "Java code analysis not yet implemented", "INFO");
        return results;
    }
};

} // namespace CodezillA
