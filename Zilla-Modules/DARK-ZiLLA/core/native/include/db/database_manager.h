#pragma once

#include <string>
#include <memory>
#include <vector>

// Forward declaration of sqlite3 struct
struct sqlite3;

namespace CodezillA {

// Assuming Logger and ErrorHandler are defined in core/ or utils/
// Add forward declarations or includes if needed
// #include "utils/logger.h"
// #include "core/error_handler.h"

// Forward declarations to avoid circular includes for shared_ptr
class Logger;
class ErrorHandler;

class DatabaseManager {
private:
    std::string db_path_;
    sqlite3* db_ = nullptr; // SQLite database connection handle
    std::shared_ptr<Logger> logger_; // Assuming Logger is available
    std::shared_ptr<ErrorHandler> error_handler_; // Assuming ErrorHandler is available

    void log_error(const std::string& message);
    bool execute_sql(const std::string& sql);
    bool table_exists(const std::string& table_name);


public:
    // Constructor
    DatabaseManager(const std::string& db_path, 
                    std::shared_ptr<ErrorHandler> error_handler, 
                    std::shared_ptr<Logger> logger);
    // Destructor
    ~DatabaseManager();

    // Connect to the database and create tables
    bool connect();

    // Save an AI suggestion to the database
    bool saveAISuggestion(const std::string& rule_id,
                          const std::string& file_path,
                          int line_number,
                          const std::string& original_code,
                          const std::string& suggested_fix);
};

} // namespace CodezillA
