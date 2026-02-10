#include "scc_parser.h"
#include "utils/json.hpp" // New include for nlohmann/json
#include <iostream>
#include <cstdio> // For popen, pclose
#include <array>   // For std::array
#include <stdexcept> // For std::runtime_error
#include <algorithm> // For std::remove_if

namespace CodezillA {

SccParser::SccParser(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger)
    : logger_(std::move(logger)), error_handler_(std::move(error_handler)) {
}

std::string SccParser::executeSccCommand(const std::string& directory_path) const {
    std::string command = "scc --by-file --format json " + directory_path; // Use --format json for structured output
    std::array<char, 128> buffer;
    std::string result;
    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(command.c_str(), "r"), pclose);
    if (!pipe) {
        error_handler_->handleError("SccParser", "Failed to run scc command: " + command);
        return "";
    }
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }
    return result;
}

SCC::OverallStats SccParser::parseSccJsonOutput(const std::string& json_output) const {
    SCC::OverallStats overall_stats = {};
    overall_stats.total_files = 0;
    overall_stats.total_code_lines = 0;
    overall_stats.total_comment_lines = 0;
    overall_stats.total_blank_lines = 0;

    try {
        nlohmann::json j = nlohmann::json::parse(json_output);

        if (j.is_array()) {
            for (const auto& lang_entry : j) {
                SCC::LanguageStats lang_stats = {};
                lang_stats.name = lang_entry.value("Name", "Unknown");
                lang_stats.total_files = lang_entry.value("Count", 0);
                lang_stats.total_code_lines = lang_entry.value("Code", 0);
                lang_stats.total_comment_lines = lang_entry.value("Comment", 0);
                lang_stats.total_blank_lines = lang_entry.value("Blank", 0);

                overall_stats.total_files += lang_stats.total_files;
                overall_stats.total_code_lines += lang_stats.total_code_lines;
                overall_stats.total_comment_lines += lang_stats.total_comment_lines;
                overall_stats.total_blank_lines += lang_stats.total_blank_lines;
                
                // Detailed file stats are not directly in this top-level json for scc by default
                // scc --by-file --format json would require deeper parsing if needed.
                // For now, we focus on overall language stats.

                overall_stats.languages.push_back(lang_stats);
            }
        }
    } catch (const nlohmann::json::exception& e) {
        error_handler_->handleError("SccParser", "JSON parsing error: " + std::string(e.what()));
    } catch (const std::exception& e) {
        error_handler_->handleError("SccParser", "General parsing error: " + std::string(e.what()));
    }

    return overall_stats;
}

std::optional<SCC::OverallStats> SccParser::analyzeDirectory(const std::string& directory_path) const {
    std::string json_output = executeSccCommand(directory_path);
    if (json_output.empty()) {
        logger_->warn("scc command returned empty output or failed.");
        return std::nullopt;
    }

    try {
        // SCC's JSON output for `--by-file --format json` can be an array of objects for each file.
        // For overall stats, it's typically a single object or an array of language summaries.
        // Let's assume the top-level output is an array of language summaries.
        // If `--by-file` is used, it often returns an array where each element is a file.
        // For simplicity, I'll first try to parse it as an array of language summaries,
        // which `scc --format json` produces. The `--by-file` option with `--format json`
        // produces a more complex structure (array of file objects), but the prompt
        // implies overall stats, so I'll prioritize that.
        // If scc --by-file --format json gives an array of files, a different parsing logic is needed.
        // For now, I will use `scc --format json` which gives language summaries.

        // Re-execute scc command without --by-file if only overall stats are needed.
        // For this task, let's keep --by-file and adapt parsing.
        // No, the original command was "scc --by-file --json ". This yields an array of file objects.
        // I'll adjust the parsing to sum up stats from --by-file output.

        nlohmann::json j = nlohmann::json::parse(json_output);
        SCC::OverallStats overall_stats = {};

        if (j.is_array()) {
            std::map<std::string, SCC::LanguageStats> language_map;

            for (const auto& file_entry : j) {
                std::string lang_name = file_entry.value("Language", "Unknown");
                long code = file_entry.value("Code", 0);
                long comment = file_entry.value("Comment", 0);
                long blank = file_entry.value("Blank", 0);

                // Aggregate stats by language
                language_map[lang_name].name = lang_name;
                language_map[lang_name].total_files++;
                language_map[lang_name].total_code_lines += code;
                language_map[lang_name].total_comment_lines += comment;
                language_map[lang_name].total_blank_lines += blank;

                overall_stats.total_files++;
                overall_stats.total_code_lines += code;
                overall_stats.total_comment_lines += comment;
                overall_stats.total_blank_lines += blank;
            }

            for (const auto& pair : language_map) {
                overall_stats.languages.push_back(pair.second);
            }
        }
        return overall_stats;
    } catch (const std::exception& e) {
        error_handler_->handleError("SccParser", "Failed to parse scc output with nlohmann/json: " + std::string(e.what()));
        return std::nullopt;
    }
}

} // namespace CodezillA
