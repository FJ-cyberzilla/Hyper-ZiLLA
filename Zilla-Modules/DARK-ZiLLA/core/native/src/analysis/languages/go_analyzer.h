#pragma once

#include "base_analyzer.h"
#include <regex>
#include <fstream>
#include <sstream>
#include <unordered_set>

namespace CodezillA {

    class GoAnalyzer : public BaseAnalyzer {
    private:
        std::vector<std::regex> error_patterns_;
        std::vector<std::regex> warning_patterns_;
        std::unordered_set<std::string> go_keywords_;
        
    public:
        GoAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
            : BaseAnalyzer(std::move(error_handler), std::move(logger)) {
            initializePatterns();
            initializeKeywords();
        }
        
        std::string getLanguageName() const override {
            return "Go";
        }
        
        std::vector<std::string> getSupportedExtensions() const override {
            return {".go"};
        }

        bool isSupportedFile(const std::string& file_path) const override {
            for (const auto& ext : getSupportedExtensions()) {
                if (file_path.size() >= ext.size() && 
                    file_path.substr(file_path.size() - ext.size()) == ext) {
                    return true;
                }
            }
            return false;
        }
        
        std::vector<AnalysisResult> analyze(const std::string& file_path) override {
            std::vector<AnalysisResult> results;
            
            try {
                std::ifstream file(file_path);
                if (!file.is_open()) {
                    results.emplace_back(file_path, "FILE_ERROR", 
                                       "Cannot open file: " + file_path, "ERROR");
                    return results;
                }
                
                std::stringstream buffer;
                buffer << file.rdbuf();
                std::string code = buffer.str();
                
                return analyzeCode(code, file_path);
                
            } catch (const std::exception& e) {
                error_handler_->handleError("GoAnalyzer", 
                                          "Error analyzing file: " + std::string(e.what()));
                results.emplace_back(file_path, "ANALYSIS_ERROR", 
                                   "Analysis failed: " + std::string(e.what()), "ERROR");
            }
            
            return results;
        }
        
        std::vector<AnalysisResult> analyzeCode(const std::string& code, const std::string& file_name) override {
            std::vector<AnalysisResult> results;
            std::istringstream stream(code);
            std::string line;
            int line_number = 0;
            bool in_block_comment = false;
            
            while (std::getline(stream, line)) {
                line_number++;
                std::string clean_line = trim(line);
                
                // Handle block comments
                if (in_block_comment) {
                    if (clean_line.find("*/") != std::string::npos) {
                        in_block_comment = false;
                    }
                    continue;
                }
                
                if (clean_line.find("/*") != std::string::npos) {
                    in_block_comment = true;
                    continue;
                }
                
                // Skip empty lines and single line comments
                if (clean_line.empty() || clean_line.find("//") == 0) {
                    continue;
                }
                
                // Analyze Go-specific patterns
                checkGoSyntax(line, line_number, file_name, results);
                checkGoConventions(line, line_number, file_name, results);
                checkCommonErrors(line, line_number, file_name, results);
            }
            
            // File-level checks
            checkFileLevelIssues(code, file_name, results);
            checkGoModuleStructure(code, file_name, results);
            
            return results;
        }
        
        bool canAutoFix(const AnalysisResult& result) const override {
            return (result.rule_id == "MISSING_PACKAGE" ||
                    result.rule_id == "UNUSED_IMPORT" ||
                    result.rule_id == "MISSING_IMPORT" ||
                    result.rule_id == "INCORRECT_FORMATTING");
        }
        
        bool applyFix(const std::string& file_path, const AnalysisResult& result) override {
            if (!canAutoFix(result)) {
                return false;
            }
            
            try {
                std::ifstream file(file_path);
                if (!file.is_open()) return false;
                
                std::vector<std::string> lines;
                std::string line;
                while (std::getline(file, line)) {
                    lines.push_back(line);
                }
                file.close();
                
                bool modified = false;
                
                if (result.rule_id == "MISSING_PACKAGE" && result.line_number == 1) {
                    // Add package declaration
                    std::string package_name = extractPackageName(file_path);
                    if (!package_name.empty()) {
                        lines.insert(lines.begin(), "package " + package_name);
                        modified = true;
                    }
                }
                else if (result.rule_id == "UNUSED_IMPORT" && result.line_number > 0) {
                    int idx = result.line_number - 1;
                    if (idx < lines.size() && lines[idx].find("import") != std::string::npos) {
                        // Comment out unused import
                        lines[idx] = "// " + lines[idx] + "  // Auto-removed: unused import";
                        modified = true;
                    }
                }
                
                if (modified) {
                    // Write back
                    std::ofstream out(file_path);
                    for (size_t i = 0; i < lines.size(); ++i) {
                        out << lines[i];
                        if (i != lines.size() - 1) out << "\n";
                    }
                    
                    logger_->info("Applied Go fix for " + result.rule_id + " in " + file_path);
                    return true;
                }
                
            } catch (const std::exception& e) {
                error_handler_->handleError("GoAnalyzer", 
                                          "Fix application failed: " + std::string(e.what()));
            }
            
            return false;
        }
        
    private:
        void initializePatterns() {
            // Common Go error patterns
            error_patterns_ = {
                std::regex(R"(fmt\.Print)"), // Using fmt.Print instead of log
                std::regex(R"(panic\([^)]+)"), // Use of panic
                std::regex(R"(\.Close\(\))"), // Potential unhandled Close
                std::regex(R"(go\s+func\([^)]*\)\s*\{)"), // Goroutine without recovery
            };
            
            warning_patterns_ = {
                std::regex(R"(var\s+\w+\s+int)"), // Uninitialized variable
                std::regex(R"(_\s*:?=)"), // Blank identifier assignment
                std::regex(R"(interface\{\})"), // Empty interface
                std::regex(R"(make\(\[\]\.+)"), // Potential slice capacity issue
            };
        }
        
        void initializeKeywords() {
            go_keywords_ = {
                "break", "case", "chan", "const", "continue", "default", "defer",
                "else", "fallthrough", "for", "func", "go", "goto", "if", "import",
                "interface", "map", "package", "range", "return", "select", "struct",
                "switch", "type", "var"
            };
        }
        
        void checkGoSyntax(const std::string& line, int line_number,
                         const std::string& file_name, std::vector<AnalysisResult>& results) {
            
            // Check for package declaration (must be first line)
            if (line_number == 1 && !std::regex_search(line, std::regex(R"(^package\s+\w+)"))) {
                results.emplace_back(file_name, "MISSING_PACKAGE",
                                   "Go file must start with package declaration", 
                                   "ERROR", line_number);
            }
            
            // Check for unused imports (basic detection)
            if (std::regex_search(line, std::regex(R"(import\s+\([^)]*\"[^\"]+\"[^)]*\))")) ||
                std::regex_search(line, std::regex(R"(import\s+\"[^\"]+\")"))) {
                // This is a simple check - real unused import detection requires deeper analysis
                results.emplace_back(file_name, "POTENTIAL_UNUSED_IMPORT",
                                   "Verify all imports are used", "INFO", line_number);
            }
            
            // Check for incorrect error handling
            if (std::regex_search(line, std::regex(R"(err\s*:?=\s*\w+\([^)]*\))")) &&
                !std::regex_search(line, std::regex(R"(err\s*!= nil)")) &&
                !std::regex_search(line, std::regex(R"(if\s+err)"))) {
                results.emplace_back(file_name, "UNCHECKED_ERROR",
                                   "Error return value not checked", "WARNING", line_number);
            }
        }
        
        void checkGoConventions(const std::string& line, int line_number,
                              const std::string& file_name, std::vector<AnalysisResult>& results) {
            
            // Check for mixed caps (exported identifiers)
            if (std::regex_search(line, std::regex(R"(func\s+[a-z]\w*\s*\()"))) {
                // Check if it should be exported (if it's in a test file or has documentation)
                if (file_name.find("_test.go") == std::string::npos) {
                    results.emplace_back(file_name, "UNEXPORTED_FUNCTION",
                                       "Consider exporting function if it needs external access", 
                                       "INFO", line_number);
                }
            }
            
            // Check for receiver names
            if (std::regex_search(line, std::regex(R"(func\s+\(\w+\s+\*?\w+\))"))) {
                auto match = std::smatch{};
                if (std::regex_search(line, match, std::regex(R"(func\s+\((\w+)\s+\*?\w+\))"))) {
                    std::string receiver = match[1];
                    if (receiver.length() > 1 || !std::isalpha(receiver[0])) {
                        results.emplace_back(file_name, "RECEIVER_NAME",
                                           "Receiver name should be 1-2 letters", "INFO", line_number);
                    }
                }
            }
            
            // Check line length (Go convention)
            if (line.length() > 100) {
                results.emplace_back(file_name, "LINE_TOO_LONG",
                                   "Line exceeds 100 characters (Go convention)", 
                                   "WARNING", line_number);
            }
        }
        
        void checkCommonErrors(const std::string& line, int line_number,
                             const std::string& file_name, std::vector<AnalysisResult>& results) {
            
            // Check for use of panic
            if (std::regex_search(line, std::regex(R"(panic\([^)]*\))"))) {
                results.emplace_back(file_name, "USE_OF_PANIC",
                                   "Avoid using panic for normal error handling", 
                                   "WARNING", line_number);
            }
            
            // Check for naked returns
            if (std::regex_search(line, std::regex(R"(return\s*$)")) &&
                std::regex_search(line, std::regex(R"(func\s+\w+\([^)]*\)\s*\([^)]*\)\s*\{)"))) {
                results.emplace_back(file_name, "NAKED_RETURN",
                                   "Naked returns can reduce code clarity", 
                                   "INFO", line_number);
            }
            
            // Check for potential data races
            if (std::regex_search(line, std::regex(R"(go\s+)")) &&
                std::regex_search(line, std::regex(R"(\w+\s*:?=\s*&?\w+\[[^]]+\])"))) {
                results.emplace_back(file_name, "POTENTIAL_DATA_RACE",
                                   "Potential data race with slice/map in goroutine", 
                                   "WARNING", line_number);
            }
        }
        
        void checkFileLevelIssues(const std::string& code, const std::string& file_name,
                                std::vector<AnalysisResult>& results) {
            // Check for proper main function in main packages
            if (code.find("package main") != std::string::npos &&
                code.find("func main()") == std::string::npos) {
                results.emplace_back(file_name, "MISSING_MAIN_FUNCTION",
                                   "Main package should contain func main()", "ERROR");
            }
            
            // Check for init functions
            if (std::regex_search(code, std::regex(R"(func init\(\))"))) {
                results.emplace_back(file_name, "INIT_FUNCTION",
                                   "Be cautious with init() functions - they can make code harder to test", 
                                   "INFO");
            }
        }
        
        void checkGoModuleStructure(const std::string& code, const std::string& file_name,
                                  std::vector<AnalysisResult>& results) {
            // Check for module compatibility
            if (code.find("// +build") != std::string::npos) {
                results.emplace_back(file_name, "BUILD_CONSTRAINTS",
                                   "Consider using go:build constraints instead of +build", 
                                   "INFO");
            }
            
            // Check for proper error wrapping
            if (code.find("fmt.Errorf") != std::string::npos && 
                code.find("%w") == std::string::npos) {
                results.emplace_back(file_name, "ERROR_WRAPPING",
                                   "Consider using %w with fmt.Errorf for error wrapping", 
                                   "INFO");
            }
        }
        
        std::string extractPackageName(const std::string& file_path) {
            // Extract directory name as package name
            size_t last_slash = file_path.find_last_of("/\\");
            if (last_slash != std::string::npos) {
                std::string dir = file_path.substr(0, last_slash);
                size_t prev_slash = dir.find_last_of("/\\");
                if (prev_slash != std::string::npos) {
                    return dir.substr(prev_slash + 1);
                }
                return dir;
            }
            return "main"; // Default fallback
        }
        
        std::string trim(const std::string& str) {
            size_t start = str.find_first_not_of(" \t\n\r");
            size_t end = str.find_last_not_of(" \t\n\r");
            return (start == std::string::npos) ? "" : str.substr(start, end - start + 1);
        }
    };
}

