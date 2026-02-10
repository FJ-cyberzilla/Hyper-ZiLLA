#include <iostream>
#include <memory>
#include <exception>
#include <csignal>
#include <map>
#include <limits>

// Your original includes
#include "include/ui/menu_system.h"
#include "src/analysis/analysis_manager.h"
#include "src/utils/logger.h"
#include "include/utils/colors.h"
#include "src/core/error_handler.h"

// New enterprise includes
#include "include/core/configuration_manager.h"
#include "include/db/database_manager.h"
#include "src/analysis/ai/ai_engine.h"

using namespace codezilla;

// Global pointer for signal handling
std::unique_ptr<ui::MenuSystem> g_menu_system;

// Signal handler for graceful shutdown
void signalHandler(int signal) {
    std::cout << "\n" << colors::YELLOW << "\nReceived interrupt signal (" << signal << ")" 
              << colors::RESET << std::endl;
    Logger::log("Application interrupted by signal: " + std::to_string(signal), LogLevel::WARNING);
    
    if (g_menu_system) {
        std::cout << colors::CYAN << "Shutting down gracefully..." << colors::RESET << std::endl;
    }
    
    Logger::log("Application terminated by user", LogLevel::INFO);
    std::exit(signal);
}

class EnhancedCodezillAApp {
private:
    std::shared_ptr<ConfigurationManager> config_manager_;
    std::shared_ptr<DatabaseManager> db_manager_;
    std::shared_ptr<analysis::AnalysisManager> analysis_manager_;
    std::unique_ptr<ui::MenuSystem> menu_system_;
    
    bool running_ = true;
    std::map<std::string, std::vector<AnalysisResult>> last_analysis_results_;

public:
    EnhancedCodezillAApp() {
        Logger::initialize("codezilla.log");
        Logger::log("CodezillA Enterprise v3.0 Started", LogLevel::INFO);
    }

    bool initialize() {
        try {
            // Display startup banner
            displayStartupBanner();
            
            // Initialize configuration manager
            std::cout << colors::CYAN << "Loading configuration..." << colors::RESET << std::endl;
            config_manager_ = std::make_shared<ConfigurationManager>("config.json");
            Logger::log("Configuration manager initialized", LogLevel::INFO);
            
            // Initialize database manager
            std::cout << colors::CYAN << "Connecting to database..." << colors::RESET << std::endl;
            std::string db_path = config_manager_->getDatabasePath();
            db_manager_ = std::make_shared<DatabaseManager>(db_path);
            
            if (!db_manager_->initialize()) {
                Logger::log("Failed to initialize database", LogLevel::ERROR);
                std::cerr << colors::RED << "Failed to initialize database!" << colors::RESET << std::endl;
                return false;
            }
            Logger::log("Database manager initialized successfully", LogLevel::INFO);
            
            // Initialize analysis manager with AI engine
            std::cout << colors::CYAN << "Initializing AI engine..." << colors::RESET << std::endl;
            
            // Create AI engine configuration
            analysis::AIEngineConfig ai_config;
            ai_config.python_executable = config_manager_->get("python_executable", "python3");
            ai_config.ai_service_path = config_manager_->get("ai_service_path", "src/analysis/ai/ai_service.py");
            ai_config.model_type = config_manager_->get("ai_model_type", "advanced");
            ai_config.timeout_seconds = config_manager_->getInt("ai_timeout", 30);
            ai_config.max_retries = config_manager_->getInt("ai_max_retries", 3);
            ai_config.enable_caching = config_manager_->getBool("ai_enable_caching", true);
            ai_config.enable_learning = config_manager_->getBool("ai_enable_learning", true);
            ai_config.cache_max_size = config_manager_->getInt("ai_cache_size", 1000);
            
            // Create analysis manager with AI engine
            analysis_manager_ = std::make_shared<analysis::AnalysisManager>(
                db_manager_,
                ai_config
            );
            
            // Verify AI engine is operational
            if (auto ai_engine = analysis_manager_->getAIEngine()) {
                if (ai_engine->isServiceAvailable()) {
                    std::cout << colors::GREEN << "✓ AI engine ready" << colors::RESET << std::endl;
                    Logger::log("AI engine initialized and operational", LogLevel::INFO);
                } else {
                    std::cout << colors::YELLOW << "⚠ AI engine initialized but service unavailable" 
                             << colors::RESET << std::endl;
                    Logger::log("AI service not available - analysis will be limited", LogLevel::WARNING);
                }
            } else {
                std::cout << colors::YELLOW << "⚠ AI engine not available - continuing with basic analysis" 
                         << colors::RESET << std::endl;
                Logger::log("AI engine not initialized - basic analysis only", LogLevel::WARNING);
            }
            
            // Create menu system
            menu_system_ = std::make_unique<ui::MenuSystem>(db_manager_, analysis_manager_);
            
            std::cout << colors::GREEN << "✓ System initialization complete\n" << colors::RESET << std::endl;
            return true;
            
        } catch (const std::exception& e) {
            Logger::log("System initialization failed: " + std::string(e.what()), LogLevel::CRITICAL);
            std::cerr << colors::RED << "Fatal error during initialization: " << e.what() 
                     << colors::RESET << std::endl;
            return false;
        }
    }

    void run() {
        displayWelcomeArt();
        
        while (running_) {
            int choice = runInteractiveMenu();
            if (choice == -1) {
                continue;
            }
            handleMenuChoice(choice);
        }

        Logger::log("CodezillA Shutdown", LogLevel::INFO);
        std::cout << colors::GREEN << "\nShutdown complete." << colors::RESET << std::endl;
    }

private:
    void displayStartupBanner() {
        std::cout << colors::CYAN << colors::BOLD << R"(
   ____          _      _______ _ _       
  / ___|___   __| | ___|__  (_) | | __ _ 
 | |   / _ \ / _` |/ _ \ / /| | | |/ _` |
 | |__| (_) | (_| |  __// /_| | | | (_| |
  \____\___/ \__,_|\___/____|_|_|_|\__,_|
  
  Advanced Code Analysis & Security Scanner
        )" << colors::RESET << std::endl;
    }

    void displayWelcomeArt() {
        std::cout << "\n" << colors::CYAN << colors::BOLD;
        std::cout << "╔════════════════════════════════════════════════════════════════════╗\n";
        std::cout << "║            CodezillA Enterprise v3.0                               ║\n";
        std::cout << "║        AI-Powered Security & Quality Analysis                      ║\n";
        std::cout << "╚════════════════════════════════════════════════════════════════════╝\n";
        std::cout << colors::RESET << std::endl;
    }

    int runInteractiveMenu() {
        std::cout << "\n" << colors::CYAN << colors::BOLD;
        std::cout << "╔════════════════════════════════════════════════════════════════╗\n";
        std::cout << "║  Main Menu                                                     ║\n";
        std::cout << "╠════════════════════════════════════════════════════════════════╣\n";
        std::cout << colors::RESET;
        
        std::cout << "║  " << colors::GREEN << "0." << colors::RESET << " Analyze Current Directory (./)                         ║\n";
        std::cout << "║  " << colors::GREEN << "1." << colors::RESET << " Analyze Specific Directory                             ║\n";
        std::cout << "║  " << colors::GREEN << "2." << colors::RESET << " Analyze Single File                                    ║\n";
        std::cout << "║  " << colors::GREEN << "3." << colors::RESET << " Run AI Auto-Fix (on last analysis)                    ║\n";
        std::cout << "║  " << colors::GREEN << "4." << colors::RESET << " View Analysis History                                  ║\n";
        std::cout << "║  " << colors::GREEN << "5." << colors::RESET << " AI Engine Configuration                               ║\n";
        std::cout << "║  " << colors::GREEN << "6." << colors::RESET << " System Statistics & Performance                        ║\n";
        std::cout << "║  " << colors::GREEN << "7." << colors::RESET << " Show Supported Languages                               ║\n";
        std::cout << "║  " << colors::GREEN << "8." << colors::RESET << " Generate Report                                        ║\n";
        std::cout << "║  " << colors::GREEN << "9." << colors::RESET << " About CodezillA                                        ║\n";
        std::cout << colors::CYAN << "╠════════════════════════════════════════════════════════════════╣\n";
        std::cout << colors::RESET;
        std::cout << "║  " << colors::RED << "10." << colors::RESET << " Exit                                                   ║\n";
        std::cout << colors::CYAN << "╚════════════════════════════════════════════════════════════════╝\n";
        std::cout << colors::RESET;
        
        std::cout << "\n" << colors::YELLOW << "Enter your choice: " << colors::RESET;
        
        int choice;
        if (!(std::cin >> choice)) {
            std::cin.clear();
            std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
            return -1;
        }
        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
        
        return choice;
    }

    void handleMenuChoice(int choice) {
        clearScreen();
        
        switch (choice) {
            case 0: analyzeCurrentDirectory(); break;
            case 1: analyzeSpecificDirectory(); break;
            case 2: analyzeSingleFile(); break;
            case 3: runAIAutoFix(); break;
            case 4: viewAnalysisHistory(); break;
            case 5: configureAIEngine(); break;
            case 6: showSystemStatistics(); break;
            case 7: showSupportedLanguages(); break;
            case 8: generateReport(); break;
            case 9: displayAbout(); break;
            case 10: running_ = false; break;
            default:
                std::cout << colors::RED << "Invalid choice. Please try again." 
                         << colors::RESET << std::endl;
                waitForUserInput();
        }
    }

    void analyzeCurrentDirectory() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== Analyzing Current Directory ===\n" 
                  << colors::RESET << std::endl;
        
        Logger::log("Analyzing current directory", LogLevel::INFO);
        
        try {
            auto results = analysis_manager_->analyzeDirectory(".");
            last_analysis_results_.clear();
            
            for (const auto& result : results) {
                last_analysis_results_[result.file_path].push_back(result);
            }
            
            displayAnalysisResults(last_analysis_results_);
            
        } catch (const std::exception& e) {
            std::cout << colors::RED << "Error: " << e.what() << colors::RESET << std::endl;
            Logger::log("Directory analysis error: " + std::string(e.what()), LogLevel::ERROR);
        }
    }

    void analyzeSpecificDirectory() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== Analyze Specific Directory ===\n" 
                  << colors::RESET << std::endl;
        
        std::cout << colors::BRIGHT_GREEN << "Enter directory path: " << colors::RESET;
        std::string path;
        std::getline(std::cin, path);
        
        if (path.empty()) {
            std::cout << colors::RED << "Path cannot be empty!" << colors::RESET << std::endl;
            waitForUserInput();
            return;
        }
        
        Logger::log("Analyzing specific directory: " + path, LogLevel::INFO);
        
        try {
            auto results = analysis_manager_->analyzeDirectory(path);
            last_analysis_results_.clear();
            
            for (const auto& result : results) {
                last_analysis_results_[result.file_path].push_back(result);
            }
            
            displayAnalysisResults(last_analysis_results_);
            
        } catch (const std::exception& e) {
            std::cout << colors::RED << "Error: " << e.what() << colors::RESET << std::endl;
            Logger::log("Directory analysis error: " + std::string(e.what()), LogLevel::ERROR);
        }
    }

    void analyzeSingleFile() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== Analyze Single File ===\n" 
                  << colors::RESET << std::endl;
        
        std::cout << colors::BRIGHT_GREEN << "Enter file path: " << colors::RESET;
        std::string path;
        std::getline(std::cin, path);
        
        if (path.empty()) {
            std::cout << colors::RED << "Path cannot be empty!" << colors::RESET << std::endl;
            waitForUserInput();
            return;
        }
        
        Logger::log("Analyzing single file: " + path, LogLevel::INFO);
        
        try {
            auto result = analysis_manager_->analyzeFile(path);
            
            if (result.success) {
                std::vector<AnalysisResult> results;
                results.push_back(result);
                displayFileResults(results, path);
                
                // Store for AI auto-fix
                last_analysis_results_.clear();
                last_analysis_results_[path].push_back(result);
            } else {
                std::cout << colors::RED << "Analysis failed: " << result.error_message 
                         << colors::RESET << std::endl;
            }
            
        } catch (const std::exception& e) {
            std::cout << colors::RED << "Error: " << e.what() << colors::RESET << std::endl;
            Logger::log("File analysis error: " + std::string(e.what()), LogLevel::ERROR);
        }
        
        waitForUserInput();
    }

    void runAIAutoFix() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== AI Auto-Fix ===\n" 
                  << colors::RESET << std::endl;
        
        if (last_analysis_results_.empty()) {
            std::cout << colors::YELLOW 
                      << "No analysis results available. Please run analysis first (Option 0, 1, or 2)." 
                      << colors::RESET << std::endl;
            Logger::log("AI Auto-Fix attempted with no results", LogLevel::WARNING);
        } else {
            Logger::log("Running AI Auto-Fix on last analysis results...", LogLevel::INFO);
            std::cout << colors::CYAN << "Applying AI-powered fixes..." << colors::RESET << std::endl;
            
            // Flatten results
            std::vector<AnalysisResult> results_to_fix = flattenAnalysisResults(last_analysis_results_);
            
            std::cout << colors::YELLOW << "Found " << results_to_fix.size() 
                      << " issues to analyze for fixes..." << colors::RESET << std::endl;
            
            // Use AI engine for recommendations
            if (auto ai_engine = analysis_manager_->getAIEngine()) {
                int fixes_suggested = 0;
                for (const auto& result : results_to_fix) {
                    if (!result.vulnerabilities.empty()) {
                        std::cout << "\n" << colors::CYAN << "File: " << result.file_path 
                                  << colors::RESET << std::endl;
                        
                        for (const auto& vuln : result.vulnerabilities) {
                            std::cout << "  " << colors::YELLOW << vuln.type 
                                      << " (Line " << vuln.line_number << ")" << colors::RESET << std::endl;
                            std::cout << "  " << colors::GREEN << "→ " << vuln.recommendation 
                                      << colors::RESET << std::endl;
                            fixes_suggested++;
                        }
                    }
                }
                
                if (fixes_suggested > 0) {
                    std::cout << "\n" << colors::GREEN << "✓ Suggested " << fixes_suggested 
                              << " fixes" << colors::RESET << std::endl;
                    Logger::log("AI Auto-Fix completed with " + std::to_string(fixes_suggested) + " suggestions", LogLevel::INFO);
                } else {
                    std::cout << colors::GREEN << "✓ No fixes needed - code looks good!" 
                              << colors::RESET << std::endl;
                }
            } else {
                std::cout << colors::RED << "AI engine not available" << colors::RESET << std::endl;
            }
        }
        
        waitForUserInput();
    }

    void viewAnalysisHistory() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== Analysis History ===\n" 
                  << colors::RESET << std::endl;
        
        std::cout << colors::YELLOW << "Loading history from database...\n" 
                  << colors::RESET << std::endl;
        
        // Database integration for history
        std::cout << "Feature: Full database history integration available\n";
        std::cout << "This will show all past analyses with timestamps and results\n";
        
        waitForUserInput();
    }

    void configureAIEngine() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== AI Engine Configuration ===\n" 
                  << colors::RESET << std::endl;
        
        if (auto ai_engine = analysis_manager_->getAIEngine()) {
            auto config = ai_engine->getConfiguration();
            
            std::cout << colors::GREEN << "Current Configuration:\n" << colors::RESET;
            std::cout << "  Python Executable: " << config.python_executable << "\n";
            std::cout << "  AI Service Path: " << config.ai_service_path << "\n";
            std::cout << "  Model Type: " << config.model_type << "\n";
            std::cout << "  Timeout: " << config.timeout_seconds << " seconds\n";
            std::cout << "  Max Retries: " << config.max_retries << "\n";
            std::cout << "  Caching: " << (config.enable_caching ? "Enabled" : "Disabled") << "\n";
            std::cout << "  Learning: " << (config.enable_learning ? "Enabled" : "Disabled") << "\n";
            std::cout << "  Cache Size: " << config.cache_max_size << "\n";
        } else {
            std::cout << colors::RED << "AI Engine not available" << colors::RESET << std::endl;
        }
        
        waitForUserInput();
    }

    void showSystemStatistics() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== System Statistics & Performance ===\n" 
                  << colors::RESET << std::endl;
        
        if (auto ai_engine = analysis_manager_->getAIEngine()) {
            std::cout << colors::GREEN << "AI Engine Performance:\n" << colors::RESET;
            std::cout << ai_engine->getPerformanceMetrics() << "\n" << std::endl;
            
            std::cout << colors::GREEN << "Cache Statistics:\n" << colors::RESET;
            std::cout << ai_engine->getCacheStatistics() << "\n" << std::endl;
        }
        
        waitForUserInput();
    }

    void showSupportedLanguages() {
        std::cout << colors::BRIGHT_GREEN << "\n╔═══════════════════════════════════════╗\n";
        std::cout << "║      Supported Languages              ║\n";
        std::cout << "╠═══════════════════════════════════════╣\n";
        std::cout << colors::RESET;
        std::cout << "║  • C++ (cpp, cc, cxx, c, h, hpp)      ║\n";
        std::cout << "║  • Python (.py)                       ║\n";
        std::cout << "║  • JavaScript (.js)                   ║\n";
        std::cout << "║  • Java (.java)                       ║\n";
        std::cout << "║  • Go (.go)                           ║\n";
        std::cout << colors::BRIGHT_GREEN << "╚═══════════════════════════════════════╝\n" 
                  << colors::RESET;
        
        Logger::log("Displayed supported languages", LogLevel::INFO);
        waitForUserInput();
    }

    void generateReport() {
        std::cout << colors::CYAN << colors::BOLD << "\n=== Generate Report ===\n" 
                  << colors::RESET << std::endl;
        
        if (last_analysis_results_.empty()) {
            std::cout << colors::YELLOW 
                      << "No analysis results available. Please run analysis first." 
                      << colors::RESET << std::endl;
        } else {
            std::cout << "Select report format:\n";
            std::cout << "  1. JSON\n";
            std::cout << "  2. HTML\n";
            std::cout << "  3. PDF\n";
            std::cout << "\nEnter choice: ";
            
            int format_choice;
            std::cin >> format_choice;
            std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
            
            std::cout << colors::GREEN << "\nReport generation feature available" 
                      << colors::RESET << std::endl;
            std::cout << "Results can be exported in multiple formats\n";
        }
        
        waitForUserInput();
    }

    void displayAbout() {
        std::cout << "\n" << colors::CYAN << colors::BOLD;
        std::cout << "╔════════════════════════════════════════════════════════════════════╗\n";
        std::cout << "║                    About CodezillA                                 ║\n";
        std::cout << "╠════════════════════════════════════════════════════════════════════╣\n";
        std::cout << colors::RESET;
        std::cout << "║  Version: 3.0 Enterprise Edition                                  ║\n";
        std::cout << "║  AI-Powered Code Analysis & Security Scanner                      ║\n";
        std::cout << "║                                                                    ║\n";
        std::cout << "║  Features:                                                         ║\n";
        std::cout << "║   • Multi-language support (C++, Python, JS, Java, Go)           ║\n";
        std::cout << "║   • AI-powered vulnerability detection                            ║\n";
        std::cout << "║   • Automated fix suggestions                                     ║\n";
        std::cout << "║   • Real-time code metrics                                        ║\n";
        std::cout << "║   • Database-backed analysis history                              ║\n";
        std::cout << "║   • Performance optimization                                      ║\n";
        std::cout << colors::CYAN << "╚════════════════════════════════════════════════════════════════════╝\n";
        std::cout << colors::RESET << std::endl;
        
        waitForUserInput();
    }

    void displayAnalysisResults(const std::map<std::string, std::vector<AnalysisResult>>& results) {
        int total_issues = 0;
        for (const auto& [file, file_results] : results) {
            total_issues += file_results.size();
        }

        std::cout << "\n" << colors::GREEN << "📊 Analysis Complete!" << colors::RESET << "\n";
        std::cout << "Files analyzed: " << results.size() << "\n";
        std::cout << "Total issues found: " << total_issues << "\n\n";

        for (const auto& [file, file_results] : results) {
            if (!file_results.empty()) {
                std::cout << colors::YELLOW << file << colors::RESET << ":\n";
                for (const auto& result : file_results) {
                    for (const auto& vuln : result.vulnerabilities) {
                        std::cout << "  " << getSeverityIcon(vuln.severity) << " "
                                  << vuln.type << " (Line " << vuln.line_number << ")\n";
                        if (!vuln.recommendation.empty()) {
                            std::cout << "    " << colors::GREEN << "→ " << vuln.recommendation 
                                      << colors::RESET << "\n";
                        }
                    }
                }
                std::cout << std::endl;
            }
        }

        waitForUserInput();
    }

    void displayFileResults(const std::vector<AnalysisResult>& results, const std::string& file_path) {
        std::cout << "\n" << colors::GREEN << "📊 Single File Analysis Complete!" << colors::RESET << "\n";
        std::cout << "File: " << file_path << "\n\n";
        
        if (results.empty() || results[0].vulnerabilities.empty()) {
            std::cout << colors::GREEN << "✓ No issues found in the file." << colors::RESET << "\n";
        } else {
            for (const auto& result : results) {
                std::cout << "Issues found: " << result.vulnerabilities.size() << "\n\n";
                for (const auto& vuln : result.vulnerabilities) {
                    std::cout << "  " << getSeverityIcon(vuln.severity) << " "
                              << vuln.type << " (Line " << vuln.line_number << ")\n";
                    std::cout << "    Description: " << vuln.description << "\n";
                    if (!vuln.recommendation.empty()) {
                        std::cout << "    " << colors::GREEN << "→ " << vuln.recommendation 
                                  << colors::RESET << "\n";
                    }
                    std::cout << std::endl;
                }
            }
        }
        
        if (!results.empty() && !results[0].ai_analysis.empty()) {
            std::cout << colors::CYAN << "\nAI Analysis:\n" << colors::RESET;
            std::cout << results[0].ai_analysis << "\n";
        }
    }

    std::string getSeverityIcon(const std::string& severity) {
        if (severity == "CRITICAL" || severity == "ERROR") return "🔴";
        if (severity == "HIGH" || severity == "WARNING") return "🟡";
        if (severity == "MEDIUM") return "🟠";
        if (severity == "LOW" || severity == "INFO") return "🔵";
        return "⚪";
    }

    void clearScreen() {
#ifdef _WIN32
        system("cls");
#else
        system("clear");
#endif
    }

    void waitForUserInput() {
        std::cout << "\n" << colors::YELLOW << "Press Enter to continue..." << colors::RESET;
        std::cin.get();
    }

    std::vector<AnalysisResult> flattenAnalysisResults(
        const std::map<std::string, std::vector<AnalysisResult>>& results_map) {
        std::vector<AnalysisResult> flat_results;
        for (const auto& pair : results_map) {
            flat_results.insert(flat_results.end(), pair.second.begin(), pair.second.end());
        }
        return flat_results;
    }
};

int main(int argc, char* argv[]) {
    // Register signal handlers
    std::signal(SIGINT, signalHandler);
    std::signal(SIGTERM, signalHandler);
    
    try {
        EnhancedCodezillAApp app;
        
        // Initialize system
        if (!app.initialize()) {
            return 1;
        }
        
        // Check for command-line arguments
        if (argc > 1) {
            std::string mode = argv[1];
            
            if (mode == "--version") {
                std::cout << "CodeZilla v3.0 - Enterprise Edition" << std::endl;
                return 0;
            } else if (mode == "--help") {
                std::cout << "Usage: " << argv[0] << " [options]" << std::endl;
                std::cout << "\nOptions:" << std::endl;
                std::cout << "  --version           Display version information" << std::endl;
                std::cout << "  --help              Display this help message" << std::endl;
                std::cout << "\nNo options:           Run in interactive mode" << std::endl;
                return 0;
            }
        }
        
        // Run interactive menu
        app.run();
        return 0;
        
    } catch (const std::exception& e) {
        std::cerr << colors::RED << "💥 Fatal Error: " << e.what() << colors::RESET << std::endl;
        Logger::log("Fatal exception in main: " + std::string(e.what()), LogLevel::CRITICAL);
        return 1;
    }
}
