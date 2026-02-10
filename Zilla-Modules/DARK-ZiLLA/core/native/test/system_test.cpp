#include <iostream>
#include <string>
#include <cstdio> // For popen, pclose
#include <cassert> // For assert
#include <stdexcept> // For std::runtime_error
#include <vector>
#include <sqlite3.h> // For SQLite operations
#include <filesystem> // For std::filesystem::remove

// Function to execute a command and capture its output
std::string exec(const char* cmd) {
    char buffer[128];
    std::string result = "";
    FILE* pipe = popen(cmd, "r");
    if (!pipe) throw std::runtime_error("popen() failed!");
    try {
        while (fgets(buffer, sizeof buffer, pipe) != NULL) {
            result += buffer;
        }
    } catch (...) {
        pclose(pipe);
        throw;
    }
    pclose(pipe);
    return result;
}

// Callback for SQLite queries to collect results
static int collect_results_callback(void* data, int argc, char** argv, char** azColName) {
    std::vector<std::vector<std::string>>* results = 
        static_cast<std::vector<std::vector<std::string>>*>(data);
    std::vector<std::string> row;
    for (int i = 0; i < argc; i++) {
        row.push_back(argv[i] ? argv[i] : "NULL");
    }
    results->push_back(row);
    return 0;
}

    int main() {
        std::cout << "Running comprehensive system tests for CodezillA...\n";
        int rc; // Declare rc here for consistent use
    // Ensure test_files/vulnerable.cpp exists for analysis
    std::string create_vulnerable_file_cmd = R"(
        mkdir -p test_files && \
        echo '#include <iostream>' > test_files/vulnerable.cpp && \
        echo '#include <string.h> // For strcpy' >> test_files/vulnerable.cpp && \
        echo '' >> test_files/vulnerable.cpp && \
        echo 'void vulnerable_function(char* input) {' >> test_files/vulnerable.cpp && \
        echo '    char buffer[10];' >> test_files/vulnerable.cpp && \
        echo '    strcpy(buffer, input); // Buffer overflow vulnerability' >> test_files/vulnerable.cpp && \
        echo '    std::cout << "Buffer: " << buffer << std::endl;' >> test_files/vulnerable.cpp && \
        echo '}' >> test_files/vulnerable.cpp && \
        echo '' >> test_files/vulnerable.cpp && \
        echo 'int main(int argc, char** argv) {' >> test_files/vulnerable.cpp && \
        echo '    if (argc < 2) {' >> test_files/vulnerable.cpp && \
        echo '        std::cout << "Usage: " << argv[0] << " <input_string>" << std::endl;' >> test_files/vulnerable.cpp && \
        echo '        return 1;' >> test_files/vulnerable.cpp && \
        echo '    }' >> test_files/vulnerable.cpp && \
        echo '    vulnerable_function(argv[1]);' >> test_files/vulnerable.cpp && \
        echo '    return 0;' >> test_files/vulnerable.cpp && \
        echo '}' >> test_files/vulnerable.cpp
)";
system(create_vulnerable_file_cmd.c_str());


    // Test 1: Graceful exit and SCC reporting
    std::cout << "\n--- Running Test 1: Analysis and SCC Reporting ---\n";
    // Input sequence: 0 (Analyze Current Directory), then 7 (Exit)
    std::string command_analyze_and_exit = "echo -e \"0\n7\" | ../build/codezilla";
    std::string output_analyze_and_exit = exec(command_analyze_and_exit.c_str());

    std::cout << "--- Application Output (Analyze & Exit) ---\n";
    std::cout << output_analyze_and_exit;
    std::cout << "-------------------------------------------\n";

    assert(output_analyze_and_exit.find("Exiting CodezillA. Goodbye!") != std::string::npos && "Exit message not found!");
    assert(output_analyze_and_exit.find("CodezillA Shutdown") != std::string::npos && "Shutdown message not found!");
    
    // Verify SCC output (check for some non-zero lines of code)
    // This is a basic check. A more robust test would parse the SCC section.
    assert(output_analyze_and_exit.find("Total Code Lines: ") != std::string::npos && "SCC Code Lines not found!");
    assert(output_analyze_and_exit.find("Total Files: ") != std::string::npos && "SCC Total Files not found!");
    // You might want to add more specific checks for actual numbers if predictable
    
    std::cout << "System Test 1 Passed: Application launched, analyzed, and exited gracefully with SCC report.\n";

    // Test 2: AI Auto-Fix simulation and database entry
    std::cout << "\n--- Running Test 2: AI Auto-Fix and Database Entry ---\n";
    // Clean up previous database file if exists
    if (std::filesystem::exists("codezilla.db")) {
        std::filesystem::remove("codezilla.db");
        std::cout << "Cleaned up old codezilla.db\n";
    }
    
    // Command to analyze the current directory (which includes vulnerable.cpp) and then run AI Auto-Fix
    // Input sequence: 0 (Analyze Current Directory), then 3 (Run AI Auto-Fix), then 7 (Exit)
    std::string command_ai_fix = "echo -e \"0\n3\n7\" | ../build/codezilla";
    std::string output_ai_fix = exec(command_ai_fix.c_str());

    std::cout << "--- Application Output (AI Fix) ---\n";
    std::cout << output_ai_fix;
    std::cout << "-----------------------------------\n";

    // Verify AI suggestion was logged/attempted
    // The CppAnalyzer logs a specific message for strcpy
    assert(output_ai_fix.find("AI suggests reviewing this critical security vulnerability manually.") != std::string::npos ||
           output_ai_fix.find("Consider replacing `strcpy` with `strncpy`") != std::string::npos);
    assert(output_ai_fix.find("AI suggestion saved: SECURITY_VULNERABILITY for test_files/vulnerable.cpp") != std::string::npos && "AI suggestion not logged/saved.");

    // Now, query the database
    sqlite3* db;
    int rc_open = sqlite3_open("codezilla.db", &db);
    assert(rc_open == SQLITE_OK && "Failed to open codezilla.db");

    std::vector<std::vector<std::string>> saved_suggestions;
    char* err_msg = nullptr;
    std::string sql_query = "SELECT rule_id, file_path, line_number, suggested_fix FROM ai_suggestions WHERE rule_id = 'SECURITY_VULNERABILITY' AND file_path LIKE '%vulnerable.cpp%';";
    
    rc = sqlite3_exec(db, sql_query.c_str(), collect_results_callback, &saved_suggestions, &err_msg);
    if (rc != SQLITE_OK) {
        std::cerr << "SQLite error: " << err_msg << std::endl;
        sqlite3_free(err_msg);
    }
    sqlite3_close(db);

    assert(rc == SQLITE_OK && "Failed to query database for AI suggestions.");
    assert(!saved_suggestions.empty() && "No AI suggestions found for vulnerable.cpp in database!");
    assert(saved_suggestions[0][0] == "SECURITY_VULNERABILITY" && "Incorrect rule_id in saved suggestion.");
    assert(saved_suggestions[0][1].find("vulnerable.cpp") != std::string::npos && "Incorrect file_path in saved suggestion.");
    assert(saved_suggestions[0][3].find("Consider replacing `strcpy` with `strncpy`") != std::string::npos && "Incorrect suggested_fix content.");
    
    std::cout << "System Test 2 Passed: AI Auto-Fix initiated and suggestion saved to database.\n";
    
    std::cout << "All comprehensive system tests passed!\n";

    // Clean up test files and database
    std::filesystem::remove("test_files/vulnerable.cpp");
    std::filesystem::remove("test_files");
    std::filesystem::remove("codezilla.db");
    std::cout << "Cleaned up test environment.\n";

    return 0;
}