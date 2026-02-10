#pragma once

#include <memory>
#include <map>
#include <vector>
#include <string>
#include <thread>
#include <future>
#include <atomic>
#include <filesystem>
#include <mutex>
#include <stdexcept> // For std::runtime_error

#include <analysis/analysis_result.h>
#include "languages/base_analyzer.h"
#include "languages/cpp_analyzer.h"
#include "languages/python_analyzer.h"
#include "languages/go_analyzer.h"
#include "languages/java_analyzer.h"
#include "languages/javascript_analyzer.h"
#include "ai/ai_engine.h"
#include "core/error_handler.h"
#include "utils/logger.h"
#include "utils/scc_parser.h" // New include for SccParser
#include "utils/scc_types.h"  // New include for SCC types
#include "db/database_manager.h" // New include for DatabaseManager
#include "core/configuration_manager.h" // New: Include ConfigurationManager

namespace CodezillA {

class AnalysisManager {
private:
    std::map<std::string, std::shared_ptr<BaseAnalyzer>> analyzers_;
    std::shared_ptr<AIEngine> ai_engine_;
    std::shared_ptr<ErrorHandler> error_handler_;
    std::shared_ptr<Logger> logger_;
    std::shared_ptr<SccParser> scc_parser_; // New: SccParser instance
    std::optional<SCC::OverallStats> scc_results_; // New: To store scc results
    std::shared_ptr<DatabaseManager> db_manager_; // New: Database Manager
    std::atomic<bool> analysis_cancelled_{false};
    
public:
    AnalysisManager();
    
    void initializeAnalyzers();
    void initializeAIEngine();
    void initializeSccParser(); // New: Initialize SccParser
    void initializeDatabaseManager(); // New: Initialize DatabaseManager
    
    std::vector<AnalysisResult> analyzeFile(const std::string& file_path);
    
    std::map<std::string, std::vector<AnalysisResult>> 
    analyzeDirectory(const std::string& directory_path);
    
    std::map<std::string, std::vector<AnalysisResult>>
    analyzeDirectoryParallel(const std::string& directory_path, int max_threads = 4);
    
    bool applyAutoFixes(const std::vector<AnalysisResult>& results);

    std::optional<SCC::OverallStats> runSccAnalysis(const std::string& directory_path); // New: Run SCC analysis
    std::optional<SCC::OverallStats> getSccResults() const { return scc_results_; } // New: Get SCC results
    
    void cancelAnalysis();
    void resetCancellation();
    
    std::shared_ptr<ErrorHandler> getErrorHandler() const; // Declared here, defined in .cpp
    
private:
    std::shared_ptr<BaseAnalyzer> getAnalyzerForFile(const std::string& file_path);
};

} // namespace CodezillA
