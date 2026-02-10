#include "db/database_manager.h"
#include "utils/logger.h"
#include "core/error_handler.h"
#include <sqlite3.h> // SQLite C API
#include <iostream>

namespace CodezillA {

// Callback function for sqlite3_exec (not used for data retrieval here, but useful for errors)
static int callback(void* data, int argc, char** argv, char** azColName) {
    return 0; // Don't do anything with data
}

DatabaseManager::DatabaseManager(const std::string& db_path,
                                 std::shared_ptr<ErrorHandler> error_handler,
                                 std::shared_ptr<Logger> logger)
    : db_path_(db_path), logger_(std::move(logger)), error_handler_(std::move(error_handler)) {
}

DatabaseManager::~DatabaseManager() {
    if (db_) {
        sqlite3_close(db_);
        logger_->info("Database connection closed for: " + db_path_);
    }
}

void DatabaseManager::log_error(const std::string& message) {
    if (logger_) {
        logger_->error("Database Error: " + message);
    }
    if (error_handler_) {
        error_handler_->handleError("DatabaseManager", message);
    }
}

bool DatabaseManager::execute_sql(const std::string& sql) {
    char* err_msg = nullptr;
    int rc = sqlite3_exec(db_, sql.c_str(), callback, 0, &err_msg);
    if (rc != SQLITE_OK) {
        log_error("SQL error: " + std::string(err_msg) + " while executing: " + sql);
        sqlite3_free(err_msg);
        return false;
    }
    return true;
}

bool DatabaseManager::table_exists(const std::string& table_name) {
    std::string sql = "SELECT name FROM sqlite_master WHERE type='table' AND name='" + table_name + "';";
    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db_, sql.c_str(), -1, &stmt, nullptr);
    if (rc != SQLITE_OK) {
        log_error("Failed to prepare table_exists statement: " + std::string(sqlite3_errmsg(db_)));
        return false;
    }
    bool exists = (sqlite3_step(stmt) == SQLITE_ROW);
    sqlite3_finalize(stmt);
    return exists;
}


bool DatabaseManager::connect() {
    int rc = sqlite3_open(db_path_.c_str(), &db_);
    if (rc) {
        log_error("Can't open database: " + std::string(sqlite3_errmsg(db_)));
        return false;
    } else {
        logger_->info("Opened database successfully: " + db_path_);
    }

    // Create table if it doesn't exist
    if (!table_exists("ai_suggestions")) {
        std::string create_table_sql = R"(
            CREATE TABLE ai_suggestions (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                rule_id TEXT NOT NULL,
                file_path TEXT NOT NULL,
                line_number INTEGER,
                original_code TEXT,
                suggested_fix TEXT,
                timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
            );
        )";
        if (!execute_sql(create_table_sql)) {
            log_error("Failed to create 'ai_suggestions' table.");
            return false;
        }
    }
    
    return true;
}

bool DatabaseManager::saveAISuggestion(const std::string& rule_id,
                                      const std::string& file_path,
                                      int line_number,
                                      const std::string& original_code,
                                      const std::string& suggested_fix) {
    if (!db_) {
        log_error("Attempted to save suggestion to an unconnected database.");
        return false;
    }

    std::string sql = "INSERT INTO ai_suggestions (rule_id, file_path, line_number, original_code, suggested_fix) VALUES (?, ?, ?, ?, ?);";
    sqlite3_stmt* stmt;
    int rc = sqlite3_prepare_v2(db_, sql.c_str(), -1, &stmt, nullptr);
    if (rc != SQLITE_OK) {
        log_error("Failed to prepare insert statement: " + std::string(sqlite3_errmsg(db_)));
        return false;
    }

    sqlite3_bind_text(stmt, 1, rule_id.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, file_path.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_int(stmt, 3, line_number);
    sqlite3_bind_text(stmt, 4, original_code.c_str(), -1, SQLITE_STATIC);
    sqlite3_bind_text(stmt, 5, suggested_fix.c_str(), -1, SQLITE_STATIC);

    rc = sqlite3_step(stmt);
    if (rc != SQLITE_DONE) {
        log_error("Failed to execute insert statement: " + std::string(sqlite3_errmsg(db_)));
        sqlite3_finalize(stmt);
        return false;
    }

    sqlite3_finalize(stmt);
    logger_->info("AI suggestion saved: " + rule_id + " for " + file_path + " at line " + std::to_string(line_number));
    return true;
}

} // namespace CodezillA
