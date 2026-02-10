#pragma once

#include "base_analyzer.h"
#include <regex>
#include <fstream>
#include <sstream>

namespace CodezillA {

    class CppAnalyzer : public BaseAnalyzer {
    private:
        std::vector<std::regex> error_patterns_;
        std::vector<std::regex> warning_patterns_;
        std::vector<std::regex> security_patterns_;
        
    public:
        CppAnalyzer(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
            : BaseAnalyzer(std::move(error_handler), std::move(logger)) {
            initializePatterns();
        }
        
        std::string getLanguageName() const override {
            return "C++";
        }
        
        std::vector<std::string> getSupportedExtensions() const override {
            return {".cpp", ".cc", ".cxx", ".c", ".h", ".hpp", ".hh", ".hxx"};
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
                error_handler_->handleError("CppAnalyzer", 
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
            
            // Basic syntax and pattern checks
            while (std::getline(stream, line)) {
                line_number++;
                checkLine(line, line_number, file_name, results);
            }
            
            // Advanced checks
            checkAdvancedPatterns(code, file_name, results);
            checkSecurityPatterns(code, file_name, results); // New: Check for security patterns
            
            return results;
        }
        
        bool canAutoFix(const AnalysisResult& result) const override {
            // Simple auto-fixable patterns and any AI-fixable security vulnerabilities
            return (result.rule_id == "MISSING_SEMICOLON" ||
                    result.rule_id == "BRACE_STYLE" ||
                    result.rule_id == "INCLUDE_GUARD_MISSING" ||
                    result.rule_id == "SECURITY_VULNERABILITY"); // Allow AI to fix security vulns
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
                if (result.rule_id == "MISSING_SEMICOLON" && result.line_number > 0) {
                    int idx = result.line_number - 1;
                    if (idx < lines.size() && !lines[idx].empty() && 
                        lines[idx].back() != ';' && lines[idx].back() != '{') {
                        lines[idx] += ';';
                    }
                }
                // ... other rule-based fixes ...
                
                // Write back
                std::ofstream out(file_path);
                for (size_t i = 0; i < lines.size(); ++i) {
                    out << lines[i];
                    if (i != lines.size() - 1) out << "\n";
                }
                
                logger_->info("Applied rule-based fix for " + result.rule_id + " in " + file_path);
                return true;
                
            } catch (const std::exception& e) {
                error_handler_->handleError("CppAnalyzer", 
                                          "Rule-based fix application failed: " + std::string(e.what()));
                return false;
            }
        }

        bool applyFix(const std::string& file_path, const AnalysisResult& result, const std::string& ai_suggested_fix) override {
            if (result.rule_id == "SECURITY_VULNERABILITY" && !ai_suggested_fix.empty()) {
                logger_->warn("AI suggested fix for security vulnerability: " + result.rule_id + " in " + file_path + " at line " + std::to_string(result.line_number) + ". Suggestion: " + ai_suggested_fix + ". Actual application of AI fix is not fully implemented yet.");
                // Placeholder: In a real scenario, this would involve parsing ai_suggested_fix
                // and applying it to the file. For now, we'll just log and return false.
                return false; 
            }
            // Fallback to rule-based fix if no AI suggestion or not a security vulnerability
            return applyFix(file_path, result);
        }
        
    private:
        void initializePatterns() {
            // Common C++ error patterns
            error_patterns_ = {
                std::regex(R"(undefined reference to)", std::regex::icase),
                std::regex(R"(expected ';' after)"),
                std::regex(R"(use of undeclared identifier)"),
                std::regex(R"(no matching function for call)")
            };
            
            warning_patterns_ = {
                std::regex(R"(unused variable)", std::regex::icase),
                std::regex(R"(comparison between signed and unsigned)"),
                std::regex(R"(deprecated declaration)")
            };

            // Security patterns
            security_patterns_ = {
                // Buffer overflow related (e.g., unsafe functions)
                std::regex(R"(strcpy\()"),
                std::regex(R"(strcat\()"),
                std::regex(R"(sprintf\()"),
                std::regex(R"(vsprintf\()"),
                std::regex(R"(gets\()"),
                // Format string vulnerability (simplified)
                std::regex(R"(printf\()"),
                std::regex(R"(fprintf\()"),
                // Command injection
                std::regex(R"(system\()"),
                std::regex(R"(exec\()"),
                std::regex(R"(popen\()"),
                // SQL injection (simplified examples)
                std::regex(R"(SELECT.*FROM)"), // Basic SELECT detection
                std::regex(R"(INSERT INTO)"),  // Basic INSERT detection
                std::regex(R"(UPDATE.*SET)"),   // Basic UPDATE detection
                std::regex(R"(DELETE FROM)")    // Basic DELETE detection
            };
        }
        
        void checkLine(const std::string& line, int line_number, 
                      const std::string& file_name, std::vector<AnalysisResult>& results) {
            
            // Check for missing semicolon (basic heuristic)
            if (!line.empty() && line.find('{') == std::string::npos &&
                line.find('}') == std::string::npos &&
                line.find('#') != 0 && // Not a preprocessor directive
                line.find("//") != 0 && // Not a comment
                line.back() != ';' && line.back() != '{' && 
                !line.empty() && !std::all_of(line.begin(), line.end(), isspace)) {
                
                // Exclude certain patterns that don't need semicolons
                if (line.find("if ") != 0 && line.find("for ") != 0 && 
                    line.find("while ") != 0 && line.find("switch ") != 0 &&
                    line.find("namespace ") != 0 && line.find("class ") != 0 &&
                    line.find("struct ") != 0 && line.find("enum ") != 0) {
                    
                    results.emplace_back(file_name, "MISSING_SEMICOLON",
                                       "Possible missing semicolon", "WARNING", 
                                       line_number);
                }
            }
            
            // Check for common patterns
            for (const auto& pattern : error_patterns_) {
                if (std::regex_search(line, pattern)) {
                    results.emplace_back(file_name, "SYNTAX_ERROR", 
                                       "Potential syntax issue detected", "ERROR",
                                       line_number);
                }
            }
            
            for (const auto& pattern : warning_patterns_) {
                if (std::regex_search(line, pattern)) {
                    results.emplace_back(file_name, "CODE_SMELL", 
                                       "Code quality issue", "WARNING",
                                       line_number);
                }
            }
        }
        
        void checkAdvancedPatterns(const std::string& code, const std::string& file_name,
                                 std::vector<AnalysisResult>& results) {
            // Check for include guards in header files
            if (file_name.find(".h") != std::string::npos) {
                checkIncludeGuard(code, file_name, results);
            }
            
            // Check for modern C++ features
            checkModernCpp(code, file_name, results);
        }
        
        void checkIncludeGuard(const std::string& code, const std::string& file_name,
                             std::vector<AnalysisResult>& results) {
            if (code.find("#ifndef") == std::string::npos || 
                code.find("#define") == std::string::npos) {
                results.emplace_back(file_name, "INCLUDE_GUARD_MISSING",
                                   "Header file missing include guard", "WARNING");
            }
        }
        
        void checkModernCpp(const std::string& code, const std::string& file_name,
                          std::vector<AnalysisResult>& results) {
            // Suggest modern C++ alternatives
            if (code.find("malloc(") != std::string::npos || 
                code.find("free(") != std::string::npos) {
                results.emplace_back(file_name, "USE_MODERN_MEMORY",
                                   "Consider using new/delete or smart pointers instead of malloc/free", 
                                   "INFO");
            }
            
            if (code.find("printf(") != std::string::npos) {
                results.emplace_back(file_name, "USE_IOSTREAMS",
                                   "Consider using iostreams instead of printf", 
                                   "INFO");
            }
        }

        void checkSecurityPatterns(const std::string& code, const std::string& file_name,
                                 std::vector<AnalysisResult>& results) {
            std::istringstream stream(code);
            std::string line;
            int line_number = 0;

            while (std::getline(stream, line)) {
                line_number++;
                for (const auto& pattern : security_patterns_) {
                    if (std::regex_search(line, pattern)) {
                        results.emplace_back(file_name, "SECURITY_VULNERABILITY",
                                           "Potential security vulnerability detected", "CRITICAL",
                                           line_number);
                    }
                }
            }
        }
    };
}
