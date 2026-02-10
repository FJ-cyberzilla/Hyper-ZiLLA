#pragma once

#include "base_analyzer.h"
#include <regex>
#include <fstream>
#include <sstream>
#include <unordered_set>

namespace CodezillA {

    class PythonAnalyzer : public BaseAnalyzer {
    private:
        std::vector<std::regex> error_patterns_;
        std::vector<std::regex> warning_patterns_;
        std::unordered_set<std::string> python_keywords_;
        
    public:
        PythonAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
            : BaseAnalyzer(std::move(error_handler), std::move(logger)) {
            initializePatterns();
            initializeKeywords();
        }
        
        std::string getLanguageName() const override {
            return "Python";
        }
        
        std::vector<std::string> getSupportedExtensions() const override {
            return {".py", ".pyw", ".pyi"};
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
            int indent_level = 0;
            std::vector<int> indent_stack = {0};
            
            // Track imports for unused import detection
            std::vector<std::string> imports;
            
            while (std::getline(stream, line)) {
                line_number++;
                
                // Clean the line for analysis
                std::string clean_line = trim(line);
                
                // Skip empty lines and comments
                if (clean_line.empty() || clean_line[0] == '#') {
                    continue;
                }
                
                // Check indentation
                checkIndentation(line, line_number, file_name, results);
                
                // Analyze line content
                checkLineContent(clean_line, line_number, file_name, results);
                
                // Check for common Python issues
                checkPythonSpecifics(clean_line, line_number, file_name, results);
                
                // Track imports
                trackImports(clean_line, line_number, imports);
            }
            
            // Advanced Python-specific checks
            checkAdvancedPythonPatterns(code, file_name, results);
            checkPEP8Violations(code, file_name, results);
            
            return results;
        }
        
        bool canAutoFix(const AnalysisResult& result) const override {
            return (result.rule_id == "MISSING_IMPORT" || // Corrected error_code to rule_id
                    result.rule_id == "UNUSED_IMPORT" ||
                    result.rule_id == "MISSING_WHITESPACE" ||
                    result.rule_id == "EXTRA_WHITESPACE" ||
                    result.rule_id == "TRAILING_WHITESPACE");
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
                
                // Apply fixes based on error type
                if (result.line_number > 0 && result.line_number <= lines.size()) {
                    int idx = result.line_number - 1;
                    
                    if (result.rule_id == "TRAILING_WHITESPACE") {
                        // Remove trailing whitespace
                        lines[idx] = trim_right(lines[idx]);
                    }
                    else if (result.rule_id == "MISSING_WHITESPACE") {
                        // Add spaces around operators (basic implementation)
                        std::string& current_line = lines[idx];
                        // Simple operator spacing
                        current_line = std::regex_replace(current_line, std::regex(R"((\\w)([+\\-\\/=!<>]=?))" ), "$1 $2");
                        current_line = std::regex_replace(current_line, std::regex(R"(([+\\-\\/=!<>]=?)(\\w))" ), "$1 $2");
                    }
                    else if (result.rule_id == "UNUSED_IMPORT" && result.line_number > 0) {
                        // Comment out unused import
                        if (lines[idx].find("import") != std::string::npos) {
                            lines[idx] = "# " + lines[idx] + "  # Auto-removed: unused import";
                        }
                    }
                }
                
                // Write back
                std::ofstream out(file_path);
                for (size_t i = 0; i < lines.size(); ++i) {
                    out << lines[i];
                    if (i != lines.size() - 1) out << "\n";
                }
                
                logger_->info("Applied Python fix for " + result.rule_id + " in " + file_path);
                return true;
                
            } catch (const std::exception& e) {
                error_handler_->handleError("PythonAnalyzer", 
                                          "Fix application failed: " + std::string(e.what()));
                return false;
            }
        }
        
    private:
        void initializePatterns() {
            // Common Python error patterns
            error_patterns_ = {
                std::regex(R"(except\s*:)", std::regex::icase), // Bare except
                std::regex(R"(except\s+Exception\s*:)", std::regex::icase), // Broad exception
                std::regex(R"(from\s+\w+\s+import\s*\*)"), // Wildcard import
                std::regex(R"(print\s+[^(])"), // Python 2 style print
                std::regex(R"(\.iterkeys\(\)|\.itervalues\(\)|\.iteritems\(\))"), // Python 2 iter methods
            };
            
            warning_patterns_ = {
                std::regex(R"(import\s+os\s*$)", std::regex::icase), // Unused import pattern
                std::regex(R"(from\s+\w+\s+import\s+[^)]+$)"), // Potential unused import
                std::regex(R"(def\s+\w+\(\)\s*:)", std::regex::icase), // Function with no parameters
                std::regex(R"(class\s+\w+\(\)\s*:)", std::regex::icase), // Class with no inheritance
            };
        }
        
        void initializeKeywords() {
            python_keywords_ = {
                "False", "None", "True", "and", "as", "assert", "async", "await",
                "break", "class", "continue", "def", "del", "elif", "else", "except",
                "finally", "for", "from", "global", "if", "import", "in", "is",
                "lambda", "nonlocal", "not", "or", "pass", "raise", "return",
                "try", "while", "with", "yield"
            };
        }
        
        void checkIndentation(const std::string& line, int line_number, 
                            const std::string& file_name, std::vector<AnalysisResult>& results) {
            if (line.empty()) return;
            
            // Count leading spaces
            int spaces = 0;
            for (char c : line) {
                if (c == ' ') spaces++;
                else if (c == '\t') {
                    results.emplace_back(file_name, "TABS_USED", 
                                       "Use spaces instead of tabs for indentation", 
                                       "WARNING", line_number);
                    return;
                }
                else break;
            }
            
            // Check for inconsistent indentation
            if (spaces % 4 != 0 && spaces > 0) {
                results.emplace_back(file_name, "INDENTATION_ERROR", 
                                   "Indentation should be multiple of 4 spaces", 
                                   "ERROR", line_number);
            }
        }
        
        void checkLineContent(const std::string& line, int line_number,
                            const std::string& file_name, std::vector<AnalysisResult>& results) {
            
            // Check for missing whitespace around operators
            if (std::regex_search(line, std::regex(R"(\w[=+\-*/<>!]+\w)")) &&
                !std::regex_search(line, std::regex(R"(:\s*$)"))) { // Not in slice
                results.emplace_back(file_name, "MISSING_WHITESPACE",
                                   "Missing whitespace around operator", "WARNING", line_number);
            }
            
            // Check for trailing whitespace
            if (!line.empty() && std::isspace(line.back())) {
                results.emplace_back(file_name, "TRAILING_WHITESPACE",
                                   "Trailing whitespace detected", "INFO", line_number);
            }
            
            // Check line length (PEP8)
            if (line.length() > 79 && !line.substr(0, 1).empty()) { // Allow long strings/comments
                if (line.find("#") == std::string::npos && 
                    line.find("\"\"\"") == std::string::npos &&
                    line.find("'''") == std::string::npos) {
                    results.emplace_back(file_name, "LINE_TOO_LONG",
                                       "Line exceeds 79 characters (PEP8)", "WARNING", line_number);
                }
            }
            
            // Check Python patterns
            for (const auto& pattern : error_patterns_) {
                if (std::regex_search(line, pattern)) {
                    std::string issue = "Code style issue detected";
                    if (std::regex_search(line, std::regex(R"(except\s*:)"))) {
                        issue = "Avoid bare except clause";
                    }
                    results.emplace_back(file_name, "CODE_STYLE_ISSUE", 
                                       issue, "WARNING", line_number);
                }
            }
        }
        
        void checkPythonSpecifics(const std::string& line, int line_number,
                                const std::string& file_name, std::vector<AnalysisResult>& results) {
            // Check for mutable default arguments
            if (std::regex_search(line, std::regex(R"(def\s+\w+\([^)]*=\s*\[)")) ||
                std::regex_search(line, std::regex(R"(def\s+\w+\([^)]*=\s*\{)"))) {
                results.emplace_back(file_name, "MUTABLE_DEFAULT_ARG",
                                   "Mutable default argument detected - can lead to unexpected behavior", 
                                   "WARNING", line_number);
            }
            
            // Check for == None instead of is None
            if (std::regex_search(line, std::regex(R"(==\s*None)"))) {
                results.emplace_back(file_name, "USE_IS_NONE",
                                   "Use 'is None' instead of '== None' for identity check", 
                                   "INFO", line_number);
            }
            
            // Check for potential NameError
            if (std::regex_search(line, std::regex(R"(\bprint\b[^(])"))) {
                results.emplace_back(file_name, "PYTHON2_PRINT",
                                   "Python 2 style print statement detected", 
                                   "ERROR", line_number);
            }
        }
        
        void trackImports(const std::string& line, int line_number, 
                         std::vector<std::string>& imports) {
            if (line.find("import ") == 0 || line.find("from ") == 0) {
                imports.push_back(line);
            }
        }
        
        void checkAdvancedPythonPatterns(const std::string& code, const std::string& file_name,
                                       std::vector<AnalysisResult>& results) {
            // Check for missing __init__.py in package (if this is __init__.py)
            if (file_name.find("__init__.py") != std::string::npos) {
                if (code.empty() || trim(code).empty()) {
                    results.emplace_back(file_name, "EMPTY_INIT",
                                       "__init__.py file is empty", "INFO");
                }
            }
            
            // Check for shebang in executable scripts
            if (file_name.find(".py") != std::string::npos && 
                code.substr(0, 3) != "#!/") {
                // This is just informational, not an error
                results.emplace_back(file_name, "MISSING_SHEBANG",
                                   "Consider adding shebang for executable scripts", 
                                   "INFO", 1);
            }
        }
        
        void checkPEP8Violations(const std::string& code, const std::string& file_name,
                               std::vector<AnalysisResult>& results) {
            // Check for multiple imports on one line
            if (std::regex_search(code, std::regex(R"(import\s+\w+\s*,\s*\w+)"))) {
                results.emplace_back(file_name, "MULTIPLE_IMPORTS",
                                   "Import each module on separate line (PEP8)", "INFO");
            }
            
            // Check for wildcard imports
            if (std::regex_search(code, std::regex(R"(from\s+\w+\s+import\s*\*)"))) {
                results.emplace_back(file_name, "WILDCARD_IMPORT",
                                   "Avoid wildcard imports (from module import *)", "WARNING");
            }
        }
        
        // Utility functions
        std::string trim(const std::string& str) {
            size_t start = str.find_first_not_of(" \t\n\r");
            size_t end = str.find_last_not_of(" \t\n\r");
            return (start == std::string::npos) ? "" : str.substr(start, end - start + 1);
        }
        
        std::string trim_right(const std::string& str) {
            size_t end = str.find_last_not_of(" \t\n\r");
            return (end == std::string::npos) ? "" : str.substr(0, end + 1);
        }
    };
}
