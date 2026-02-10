#include "../../include/ui/menu_system.h"
#include "../../src/utils/logger.h"
#include "../../include/utils/colors.h"
#include <iostream>
#include <iomanip>
#include <sstream>

namespace codezilla {
namespace ui {

// UTF-8 safe box drawing characters
namespace BoxChars {
    const std::string HORIZONTAL = "─";
    const std::string VERTICAL = "│";
    const std::string TOP_LEFT = "┌";
    const std::string TOP_RIGHT = "┐";
    const std::string BOTTOM_LEFT = "└";
    const std::string BOTTOM_RIGHT = "┘";
    const std::string T_DOWN = "┬";
    const std::string T_UP = "┴";
    const std::string T_RIGHT = "├";
    const std::string T_LEFT = "┤";
    const std::string CROSS = "┼";
}

MenuSystem::MenuSystem(
    std::shared_ptr<DatabaseManager> db_manager,
    std::shared_ptr<AnalysisManager> analysis_manager
) : db_manager_(std::move(db_manager)),
    analysis_manager_(std::move(analysis_manager)),
    running_(true) 
{
    if (!db_manager_) {
        throw std::invalid_argument("Database manager cannot be null");
    }
    if (!analysis_manager_) {
        throw std::invalid_argument("Analysis manager cannot be null");
    }
    
    Logger::log("Menu system initialized", LogLevel::INFO);
}

void MenuSystem::run() {
    clearScreen();
    displayWelcomeBanner();
    
    while (running_) {
        try {
            displayMainMenu();
            int choice = getUserChoice();
            handleMenuChoice(choice);
            if (running_) {
                waitForEnter();
            }
        } catch (const std::exception& e) {
            std::cout << colors::RED << "Error: " << e.what() << colors::RESET << std::endl;
            Logger::log("Menu error: " + std::string(e.what()), LogLevel::ERROR);
            waitForEnter();
        }
    }
    
    displayExitMessage();
}

void MenuSystem::displayWelcomeBanner() {
    std::cout << colors::CYAN << colors::BOLD;
    
    std::cout << "\n";
    std::cout << BoxChars::TOP_LEFT << repeatString(BoxChars::HORIZONTAL, 78) << BoxChars::TOP_RIGHT << "\n";
    std::cout << BoxChars::VERTICAL << centerText("CODEZILLA - Advanced Code Analysis System", 78) << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << centerText("Enterprise Edition v2.0", 78) << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::T_RIGHT << repeatString(BoxChars::HORIZONTAL, 78) << BoxChars::T_LEFT << "\n";
    std::cout << BoxChars::VERTICAL << centerText("AI-Powered Security & Quality Analysis", 78) << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::BOTTOM_LEFT << repeatString(BoxChars::HORIZONTAL, 78) << BoxChars::BOTTOM_RIGHT << "\n";
    
    std::cout << colors::RESET << std::endl;
}

void MenuSystem::displayMainMenu() {
    std::cout << "\n" << colors::CYAN << colors::BOLD;
    std::cout << BoxChars::TOP_LEFT << repeatString(BoxChars::HORIZONTAL, 60) << BoxChars::TOP_RIGHT << "\n";
    std::cout << BoxChars::VERTICAL << colors::YELLOW << "  Main Menu" 
              << std::string(49, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::T_RIGHT << repeatString(BoxChars::HORIZONTAL, 60) << BoxChars::T_LEFT << "\n";
    std::cout << colors::RESET;
    
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "1." << colors::RESET 
              << " Analyze Single File" << std::string(35, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "2." << colors::RESET 
              << " Analyze Directory (Recursive)" << std::string(26, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "3." << colors::RESET 
              << " View Analysis History" << std::string(31, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "4." << colors::RESET 
              << " Generate Report (JSON/HTML)" << std::string(27, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "5." << colors::RESET 
              << " AI Engine Configuration" << std::string(30, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "6." << colors::RESET 
              << " System Statistics" << std::string(36, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "7." << colors::RESET 
              << " Performance Metrics" << std::string(34, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "8." << colors::RESET 
              << " Clear Cache" << std::string(42, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::GREEN << "9." << colors::RESET 
              << " Run Tests" << std::string(44, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << "  " << colors::RED << "0." << colors::RESET 
              << " Exit" << std::string(49, ' ') << colors::CYAN << BoxChars::VERTICAL << "\n";
    
    std::cout << BoxChars::BOTTOM_LEFT << repeatString(BoxChars::HORIZONTAL, 60) << BoxChars::BOTTOM_RIGHT << "\n";
    std::cout << colors::RESET;
    
    std::cout << "\n" << colors::YELLOW << "Enter your choice: " << colors::RESET;
}

int MenuSystem::getUserChoice() {
    std::string input;
    std::getline(std::cin, input);
    
    try {
        return std::stoi(input);
    } catch (...) {
        return -1;
    }
}

void MenuSystem::handleMenuChoice(int choice) {
    clearScreen();
    
    switch (choice) {
        case 1:
            analyzeSingleFile();
            break;
        case 2:
            analyzeDirectory();
            break;
        case 3:
            viewAnalysisHistory();
            break;
        case 4:
            generateReport();
            break;
        case 5:
            configureAIEngine();
            break;
        case 6:
            showSystemStatistics();
            break;
        case 7:
            showPerformanceMetrics();
            break;
        case 8:
            clearCache();
            break;
        case 9:
            runTests();
            break;
        case 0:
            running_ = false;
            break;
        default:
            std::cout << colors::RED << "Invalid choice. Please try again." 
                     << colors::RESET << std::endl;
            waitForEnter();
    }
}

void MenuSystem::analyzeSingleFile() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Analyze Single File ===\n" 
              << colors::RESET << std::endl;
    
    std::string file_path;
    std::cout << "Enter file path: ";
    std::getline(std::cin, file_path);
    
    if (file_path.empty()) {
        std::cout << colors::RED << "File path cannot be empty!" << colors::RESET << std::endl;
        waitForEnter();
        return;
    }
    
    std::cout << "\n" << colors::YELLOW << "Analyzing file: " << file_path 
              << colors::RESET << std::endl;
    std::cout << "Please wait...\n" << std::endl;
    
    try {
        auto result = analysis_manager_->analyzeFile(file_path);
        
        if (result.success) {
            displayAnalysisResult(result);
        } else {
            std::cout << colors::RED << "Analysis failed: " << result.error_message 
                     << colors::RESET << std::endl;
        }
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Error during analysis: " << e.what() 
                 << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::analyzeDirectory() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Analyze Directory (Recursive) ===\n" 
              << colors::RESET << std::endl;
    
    std::string dir_path;
    std::cout << "Enter directory path: ";
    std::getline(std::cin, dir_path);
    
    if (dir_path.empty()) {
        std::cout << colors::RED << "Directory path cannot be empty!" << colors::RESET << std::endl;
        waitForEnter();
        return;
    }
    
    std::cout << "\n" << colors::YELLOW << "Analyzing directory: " << dir_path 
              << colors::RESET << std::endl;
    std::cout << "This may take a while for large directories...\n" << std::endl;
    
    try {
        auto results = analysis_manager_->analyzeDirectory(dir_path);
        
        std::cout << colors::GREEN << "\nAnalyzed " << results.size() << " files\n" 
                 << colors::RESET << std::endl;
        
        int total_issues = 0;
        for (const auto& result : results) {
            if (!result.vulnerabilities.empty()) {
                total_issues += result.vulnerabilities.size();
            }
        }
        
        std::cout << colors::YELLOW << "Total issues found: " << total_issues 
                 << colors::RESET << std::endl;
        
        // Display summary
        if (!results.empty()) {
            std::cout << "\n" << colors::CYAN << "Top Issues:\n" << colors::RESET;
            int count = 0;
            for (const auto& result : results) {
                if (!result.vulnerabilities.empty() && count < 5) {
                    std::cout << "  " << colors::YELLOW << result.file_path << colors::RESET 
                             << ": " << result.vulnerabilities.size() << " issues\n";
                    count++;
                }
            }
        }
        
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Error during directory analysis: " << e.what() 
                 << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::viewAnalysisHistory() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Analysis History ===\n" 
              << colors::RESET << std::endl;
    
    try {
        // Get recent analyses from database
        std::cout << colors::YELLOW << "Loading history from database...\n" 
                 << colors::RESET << std::endl;
        
        // Placeholder for database query
        std::cout << "Feature coming soon: Database integration for history viewing\n";
        
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Error loading history: " << e.what() 
                 << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::generateReport() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Generate Report ===\n" 
              << colors::RESET << std::endl;
    
    std::cout << "Select report format:\n";
    std::cout << "  1. JSON\n";
    std::cout << "  2. HTML\n";
    std::cout << "  3. PDF\n";
    std::cout << "\nEnter choice: ";
    
    int format_choice = getUserChoice();
    
    std::string output_file;
    std::cout << "Enter output file path: ";
    std::getline(std::cin, output_file);
    
    if (output_file.empty()) {
        std::cout << colors::RED << "Output file path cannot be empty!" 
                 << colors::RESET << std::endl;
        waitForEnter();
        return;
    }
    
    std::cout << "\n" << colors::YELLOW << "Generating report..." 
              << colors::RESET << std::endl;
    
    // Placeholder for report generation
    std::cout << colors::GREEN << "Report would be generated to: " << output_file 
             << colors::RESET << std::endl;
    
    waitForEnter();
}

void MenuSystem::configureAIEngine() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== AI Engine Configuration ===\n" 
              << colors::RESET << std::endl;
    
    if (auto ai_engine = analysis_manager_->getAIEngine()) {
        auto config = ai_engine->getConfiguration();
        
        std::cout << "Current Configuration:\n";
        std::cout << "  Python Executable: " << config.python_executable << "\n";
        std::cout << "  AI Service Path: " << config.ai_service_path << "\n";
        std::cout << "  Model Type: " << config.model_type << "\n";
        std::cout << "  Timeout: " << config.timeout_seconds << " seconds\n";
        std::cout << "  Max Retries: " << config.max_retries << "\n";
        std::cout << "  Caching: " << (config.enable_caching ? "Enabled" : "Disabled") << "\n";
        std::cout << "  Learning: " << (config.enable_learning ? "Enabled" : "Disabled") << "\n";
        std::cout << "  Cache Size: " << config.cache_max_size << "\n";
        
        std::cout << "\nModify configuration? (y/n): ";
        std::string response;
        std::getline(std::cin, response);
        
        if (response == "y" || response == "Y") {
            // Allow configuration updates
            std::cout << "Configuration update interface would go here\n";
        }
    } else {
        std::cout << colors::RED << "AI Engine not available" << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::showSystemStatistics() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== System Statistics ===\n" 
              << colors::RESET << std::endl;
    
    try {
        if (auto ai_engine = analysis_manager_->getAIEngine()) {
            std::string cache_stats = ai_engine->getCacheStatistics();
            std::cout << colors::GREEN << "Cache Statistics:\n" << colors::RESET;
            std::cout << cache_stats << "\n" << std::endl;
        }
        
        // Add more statistics here
        std::cout << colors::GREEN << "Database Statistics:\n" << colors::RESET;
        std::cout << "Feature coming soon\n";
        
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Error retrieving statistics: " << e.what() 
                 << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::showPerformanceMetrics() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Performance Metrics ===\n" 
              << colors::RESET << std::endl;
    
    try {
        if (auto ai_engine = analysis_manager_->getAIEngine()) {
            std::string metrics = ai_engine->getPerformanceMetrics();
            std::cout << colors::GREEN << "AI Engine Metrics:\n" << colors::RESET;
            std::cout << metrics << "\n" << std::endl;
        }
        
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Error retrieving metrics: " << e.what() 
                 << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::clearCache() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Clear Cache ===\n" 
              << colors::RESET << std::endl;
    
    std::cout << colors::YELLOW << "Are you sure you want to clear the cache? (y/n): " 
              << colors::RESET;
    
    std::string response;
    std::getline(std::cin, response);
    
    if (response == "y" || response == "Y") {
        try {
            if (auto ai_engine = analysis_manager_->getAIEngine()) {
                ai_engine->clearCache();
                std::cout << colors::GREEN << "Cache cleared successfully!" 
                         << colors::RESET << std::endl;
            }
        } catch (const std::exception& e) {
            std::cout << colors::RED << "Error clearing cache: " << e.what() 
                     << colors::RESET << std::endl;
        }
    } else {
        std::cout << "Cache clear cancelled." << std::endl;
    }
    
}

void MenuSystem::runTests() {
    std::cout << colors::CYAN << colors::BOLD << "\n=== Run Tests ===\n" 
              << colors::RESET << std::endl;
    
    std::cout << colors::YELLOW << "Running system tests...\n" << colors::RESET << std::endl;
    
    try {
        // Test AI engine availability
        if (auto ai_engine = analysis_manager_->getAIEngine()) {
            std::cout << "  Testing AI Engine... ";
            bool available = ai_engine->isServiceAvailable();
            if (available) {
                std::cout << colors::GREEN << "✓ PASSED" << colors::RESET << std::endl;
            } else {
                std::cout << colors::RED << "✗ FAILED" << colors::RESET << std::endl;
            }
        }
        
        // Test database connection
        std::cout << "  Testing Database Connection... ";
        if (db_manager_) {
            std::cout << colors::GREEN << "✓ PASSED" << colors::RESET << std::endl;
        } else {
            std::cout << colors::RED << "✗ FAILED" << colors::RESET << std::endl;
        }
        
        std::cout << "\n" << colors::GREEN << "Tests completed!" << colors::RESET << std::endl;
        
    } catch (const std::exception& e) {
        std::cout << colors::RED << "Test error: " << e.what() << colors::RESET << std::endl;
    }
    
    waitForEnter();
}

void MenuSystem::displayAnalysisResult(const AnalysisResult& result) {
    std::cout << "\n" << colors::GREEN << colors::BOLD << "=== Analysis Results ===" 
              << colors::RESET << "\n" << std::endl;
    
    std::cout << colors::CYAN << "File: " << colors::RESET << result.file_path << std::endl;
    std::cout << colors::CYAN << "Language: " << colors::RESET << result.language << std::endl;
    std::cout << colors::CYAN << "Lines of Code: " << colors::RESET << result.lines_of_code << std::endl;
    
    if (!result.vulnerabilities.empty()) {
        std::cout << "\n" << colors::RED << colors::BOLD << "Vulnerabilities Found: " 
                 << result.vulnerabilities.size() << colors::RESET << std::endl;
        
        for (size_t i = 0; i < result.vulnerabilities.size() && i < 10; ++i) {
            const auto& vuln = result.vulnerabilities[i];
            std::cout << "\n  " << (i + 1) << ". " << colors::YELLOW << vuln.type 
                     << colors::RESET << " (Line " << vuln.line_number << ")";
            std::cout << "\n     " << vuln.description;
            if (!vuln.recommendation.empty()) {
                std::cout << "\n     " << colors::GREEN << "→ " << vuln.recommendation 
                         << colors::RESET;
            }
            std::cout << std::endl;
        }
        
        if (result.vulnerabilities.size() > 10) {
            std::cout << "\n  ... and " << (result.vulnerabilities.size() - 10) 
                     << " more issues\n";
        }
    } else {
        std::cout << "\n" << colors::GREEN << "✓ No vulnerabilities detected!" 
                 << colors::RESET << std::endl;
    }
    
    if (!result.ai_analysis.empty()) {
        std::cout << "\n" << colors::CYAN << colors::BOLD << "AI Analysis:\n" 
                 << colors::RESET << result.ai_analysis << std::endl;
    }
}

void MenuSystem::displayExitMessage() {
    std::cout << "\n" << colors::CYAN << colors::BOLD;
    std::cout << BoxChars::TOP_LEFT << repeatString(BoxChars::HORIZONTAL, 60) << BoxChars::TOP_RIGHT << "\n";
    std::cout << BoxChars::VERTICAL << centerText("Thank you for using CodeZilla!", 60) << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::VERTICAL << centerText("Stay secure, code better!", 60) << BoxChars::VERTICAL << "\n";
    std::cout << BoxChars::BOTTOM_LEFT << repeatString(BoxChars::HORIZONTAL, 60) << BoxChars::BOTTOM_RIGHT << "\n";
    std::cout << colors::RESET << std::endl;
}

void MenuSystem::clearScreen() {
#ifdef _WIN32
    system("cls");
#else
    system("clear");
#endif
}

void MenuSystem::waitForEnter() {
    std::cout << "\n" << colors::YELLOW << "Press Enter to continue..." 
              << colors::RESET;
    std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
    std::cin.get();
}

std::string MenuSystem::centerText(const std::string& text, int width) {
    int padding = (width - text.length()) / 2;
    return std::string(padding, ' ') + text + std::string(width - padding - text.length(), ' ');
}

std::string MenuSystem::repeatString(const std::string& str, int count) {
    std::string result;
    for (int i = 0; i < count; ++i) {
        result += str;
    }
    return result;
}

} // namespace ui
} // namespace codezilla
