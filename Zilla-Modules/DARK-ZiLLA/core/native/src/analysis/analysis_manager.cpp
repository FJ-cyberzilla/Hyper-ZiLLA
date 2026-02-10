#include "analysis_manager.h"
#include "utils/logger.h" // For Logger
#include "core/error_handler.h" // For ErrorHandler
#include "ai/ai_engine.h" // For AIEngine
#include "languages/base_analyzer.h" // For BaseAnalyzer
#include "languages/cpp_analyzer.h" // For CppAnalyzer
#include <iostream>
#include <filesystem> // For std::filesystem

namespace CodezillA {

AnalysisManager::AnalysisManager() {
    logger_ = std::make_shared<Logger>("AnalysisManager");
    error_handler_ = std::make_shared<ErrorHandler>();
    initializeDatabaseManager(); // New: Initialize DatabaseManager first
    initializeAnalyzers();
    initializeAIEngine();
    initializeSccParser();
}

void AnalysisManager::initializeAnalyzers() {
    try {
        analyzers_["cpp"] = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<CppAnalyzer>(error_handler_, logger_); });

        analyzers_["python"] = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<PythonAnalyzer>(error_handler_, logger_); });
        analyzers_["go"] = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<GoAnalyzer>(error_handler_, logger_); });
        analyzers_["java"] = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<JavaAnalyzer>(error_handler_, logger_); });
        analyzers_["javascript"] = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<JavaScriptAnalyzer>(error_handler_, logger_); });

        logger_->info("Successfully initialized " + 
                     std::to_string(analyzers_.size()) + " language analyzers");

    } catch (const AnalyzerException& e) {
        error_handler_->handleException(e);
        logger_->error("Failed to initialize some analyzers");
    }
}

void AnalysisManager::initializeAIEngine() {
    try {
        ai_engine_ = std::shared_ptr<AIEngine>(new AIEngine(db_manager_)); // Direct construction
        
        if (ai_engine_ && ai_engine_->isAvailable()) {
            logger_->info("AI Engine initialized successfully");
        }
        else {
            logger_->warn("AI Engine not available - running in basic mode");
        }
        
    } catch (const std::exception& e) { // Catch std::exception for general AI init errors
        error_handler_->handleError("AnalysisManager", "Error during AI Engine initialization: " + std::string(e.what()));
        ai_engine_ = nullptr; // Ensure it's null if creation fails
    }
}

void AnalysisManager::initializeSccParser() {
    try {
        scc_parser_ = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<SccParser>(error_handler_, logger_); });
        logger_->info("SCC Parser initialized successfully");
    } catch (const std::exception& e) { // Catch all exceptions for robustness
        error_handler_->handleError("AnalysisManager", "Failed to initialize SCC Parser: " + std::string(e.what()));
        scc_parser_ = nullptr; // Indicate failure
    }
}

void AnalysisManager::initializeDatabaseManager() {
    try {
        db_manager_ = error_handler_->executeWithRecovery(
            [this]() { return std::make_shared<DatabaseManager>("codezilla.db", error_handler_, logger_); });
        if (db_manager_ && db_manager_->connect()) {
            logger_->info("Database Manager initialized and connected successfully.");
        } else {
            logger_->error("Failed to initialize or connect Database Manager.");
            db_manager_ = nullptr; // Ensure it's null if connection failed
        }
    } catch (const std::exception& e) {
        error_handler_->handleError("AnalysisManager", "Error during Database Manager initialization: " + std::string(e.what()));
        db_manager_ = nullptr;
    }
}

std::optional<SCC::OverallStats> AnalysisManager::runSccAnalysis(const std::string& directory_path) {
    if (!scc_parser_) {
        logger_->error("SCC Parser not initialized. Cannot run SCC analysis.");
        return std::nullopt;
    }
    scc_results_ = scc_parser_->analyzeDirectory(directory_path);
    if (scc_results_) {
        logger_->info("SCC analysis completed for " + directory_path);
    } else {
        logger_->warn("SCC analysis failed for " + directory_path);
    }
    return scc_results_;
}

std::vector<AnalysisResult> AnalysisManager::analyzeFile(const std::string& file_path) {
    if (analysis_cancelled_) {
        throw AnalysisException("Analysis cancelled by user", "All");
    }
    
    return error_handler_->executeWithRecovery([this, &file_path]() {
        auto analyzer = getAnalyzerForFile(file_path);
        if (!analyzer) {
            return std::vector<AnalysisResult>{
                AnalysisResult(file_path, "UNSUPPORTED_LANGUAGE", 
                             "File type not supported", "ERROR")
            };
        }
        
        logger_->info("Analyzing: " + file_path);
        // Read file content to provide context for AI enhancement
        std::ifstream file(file_path);
        std::string code_context((std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>());
        file.close();

        auto results = analyzer->analyze(file_path);
        
        // Enhance with AI insights if available
        if (ai_engine_ && ai_engine_->isAvailable()) {
            results = ai_engine_->enhanceAnalysis(results, file_path, code_context); // Pass code_context
        }
        
        return results;
    });
}

std::map<std::string, std::vector<AnalysisResult>> 
AnalysisManager::analyzeDirectory(const std::string& directory_path) {
    std::map<std::string, std::vector<AnalysisResult>> all_results;
    
    try {
        error_handler_->executeWithRecovery([this, &directory_path, &all_results]() {
            for (const auto& entry : std::filesystem::recursive_directory_iterator(directory_path)) {
                if (analysis_cancelled_) {
                    throw AnalysisException("Analysis cancelled", "Directory");
                }
                
                if (entry.is_regular_file()) {
                    try {
                        auto file_results = analyzeFile(entry.path().string());
                        if (!file_results.empty()) {
                            all_results[entry.path().string()] = file_results;
                        }
                    } catch (const FileSystemException& e) {
                        // Log but continue with other files
                        error_handler_->handleException(e);
                    } catch (const AnalysisException& e) {
                        // Log but continue with other files  
                        error_handler_->handleException(e);
                    }
                }
            }
        });
        
    } catch (const AnalyzerException& e) {
        error_handler_->handleException(e);
        // Return partial results even if some files failed
    }
    
    return all_results;
}

std::map<std::string, std::vector<AnalysisResult>>
AnalysisManager::analyzeDirectoryParallel(const std::string& directory_path, int max_threads) {
    std::map<std::string, std::vector<AnalysisResult>> all_results;
    std::mutex results_mutex;
    std::vector<std::future<void>> futures;
    
    try {
        std::vector<std::string> files_to_analyze;
        
        // Collect all files first
        for (const auto& entry : std::filesystem::recursive_directory_iterator(directory_path)) {
            if (entry.is_regular_file() && getAnalyzerForFile(entry.path().string())) {
                files_to_analyze.push_back(entry.path().string());
            }
        }
        
        // Process files in parallel with limited threads
        std::atomic<size_t> current_index{0};
        std::atomic<bool> has_errors{false};
        
        auto worker = [&]() {
            while (!analysis_cancelled_ && !has_errors) {
                size_t index = current_index++;
                if (index >= files_to_analyze.size()) break;
                
                try {
                    auto file_results = analyzeFile(files_to_analyze[index]);
                    if (!file_results.empty()) {
                        std::lock_guard<std::mutex> lock(results_mutex);
                        all_results[files_to_analyze[index]] = file_results;
                    }
                } catch (const std::exception& e) {
                    has_errors = true;
                    error_handler_->handleError("ParallelAnalysis", 
                                              "Failed to analyze: " + files_to_analyze[index]);
                }
            }
        };
        
        // Launch worker threads
        int num_threads = std::min(max_threads, static_cast<int>(files_to_analyze.size()));
        for (int i = 0; i < num_threads; ++i) {
            futures.push_back(std::async(std::launch::async, worker));
        }
        
        // Wait for all threads to complete
        for (auto& future : futures) {
            future.wait();
        }
        
    } catch (const std::exception& e) {
        error_handler_->handleError("ParallelAnalysis", 
                                  "Parallel analysis failed: " + std::string(e.what()));
    }
    
    return all_results;
}

bool AnalysisManager::applyAutoFixes(const std::vector<AnalysisResult>& results) {
    int fixes_applied = 0;
    int total_fixable = 0;
    
    try {
        error_handler_->executeWithRecovery([&]() {
            for (const auto& result : results) {
                auto analyzer = getAnalyzerForFile(result.file_path);
                if (analyzer && analyzer->canAutoFix(result)) {
                    total_fixable++;
                    
                    std::string ai_suggested_fix = "";
                    if (ai_engine_ && ai_engine_->isAvailable() && result.rule_id == "SECURITY_VULNERABILITY") {
                        // Read the entire file content to pass as code_context to AI
                        std::ifstream file(result.file_path);
                        std::string code_context((std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>());
                        file.close();
                        
                        ai_suggested_fix = ai_engine_->suggestFixes(result, code_context);
                    }

                    if (analyzer->applyFix(result.file_path, result, ai_suggested_fix)) {
                        fixes_applied++;
                    }
                }
            }
        });
        
        logger_->info("Applied " + std::to_string(fixes_applied) + 
                     " out of " + std::to_string(total_fixable) + " auto-fixes");
        return fixes_applied > 0;
        
    } catch (const AnalyzerException& e) {
        error_handler_->handleException(e);
        return false;
    }
}

void AnalysisManager::cancelAnalysis() {
    analysis_cancelled_ = true;
    logger_->info("Analysis cancellation requested");
}

void AnalysisManager::resetCancellation() {
    analysis_cancelled_ = false;
}

std::shared_ptr<ErrorHandler> AnalysisManager::getErrorHandler() const { return error_handler_; }

std::shared_ptr<BaseAnalyzer> AnalysisManager::getAnalyzerForFile(const std::string& file_path) {
    for (const auto& [lang, analyzer] : analyzers_) {
        if (analyzer->isSupportedFile(file_path)) {
            return analyzer;
        }
    }
    return nullptr;
}

} // namespace CodezillA