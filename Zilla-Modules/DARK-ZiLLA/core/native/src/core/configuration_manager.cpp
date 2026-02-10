#include "core/configuration_manager.h"
#include <fstream>
#include <iostream>

namespace CodezillA {

ConfigurationManager::ConfigurationManager(const std::string& config_file_path, std::shared_ptr<Logger> logger)
    : config_file_path_(config_file_path), logger_(logger) {
    // Attempt to load existing configuration, or initialize with defaults
    if (!loadConfiguration()) {
        logger_->warn("Configuration file not found or invalid. Initializing with default settings.");
        resetToDefaults();
        saveConfiguration(); // Save defaults if loaded failed
    }
}

bool ConfigurationManager::loadConfiguration() {
    std::ifstream file(config_file_path_);
    if (!file.is_open()) {
        logger_->info("Configuration file not found: " + config_file_path_);
        return false;
    }

    try {
        file >> config_data_;
        logger_->info("Configuration loaded from: " + config_file_path_);
        return true;
    } catch (const nlohmann::json::parse_error& e) {
        logger_->error("Error parsing configuration file " + config_file_path_ + ": " + e.what());
        return false;
    }
}

bool ConfigurationManager::saveConfiguration() {
    std::ofstream file(config_file_path_);
    if (!file.is_open()) {
        logger_->error("Failed to open configuration file for writing: " + config_file_path_);
        return false;
    }

    try {
        file << std::setw(4) << config_data_ << std::endl; // Pretty print with 4-space indent
        logger_->info("Configuration saved to: " + config_file_path_);
        return true;
    } catch (const std::exception& e) {
        logger_->error("Error writing configuration file " + config_file_path_ + ": " + e.what());
        return false;
    }
}

void ConfigurationManager::resetToDefaults() {
    config_data_.clear();
    initializeDefaultSettings();
    logger_->info("Configuration reset to default settings.");
}

void ConfigurationManager::initializeDefaultSettings() {
    // Set default values for all known settings
    config_data_["python_ai_service_path"] = "./src/analysis/ai/ai_service.py";
    config_data_["default_analysis_directory"] = ".";
    // Add other default settings here
}

std::string ConfigurationManager::getPythonAIServicePath() const {
    return get("python_ai_service_path", std::string("./src/analysis/ai/ai_service.py"));
}

std::string ConfigurationManager::getDefaultAnalysisDirectory() const {
    return get("default_analysis_directory", std::string("."));
}

void ConfigurationManager::setPythonAIServicePath(const std::string& path) {
    set("python_ai_service_path", path);
}

void ConfigurationManager::setDefaultAnalysisDirectory(const std::string& path) {
    set("default_analysis_directory", path);
}

} // namespace CodezillA
