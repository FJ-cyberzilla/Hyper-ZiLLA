#ifndef CONFIGURATION_MANAGER_H
#define CONFIGURATION_MANAGER_H

#include <string>
#include <map>
#include <memory> // For std::shared_ptr
#include <utils/json.hpp> // For nlohmann::json
#include "utils/logger.h" // For Logger

namespace CodezillA {

class ConfigurationManager {
public:
    // Define a structure for application settings
    struct AppSettings {
        std::string python_ai_service_path;
        std::string default_analysis_directory;
        // Add more settings as needed
    };

private:
    std::string config_file_path_;
    nlohmann::json config_data_;
    std::shared_ptr<Logger> logger_;

public:
    ConfigurationManager(const std::string& config_file_path, std::shared_ptr<Logger> logger);

    bool loadConfiguration();
    bool saveConfiguration();
    void resetToDefaults();

    // Getters for specific settings
    std::string getPythonAIServicePath() const;
    std::string getDefaultAnalysisDirectory() const;

    // Setter for a specific setting
    void setPythonAIServicePath(const std::string& path);
    void setDefaultAnalysisDirectory(const std::string& path);

    // Generic getter/setter for flexibility
    template<typename T>
    T get(const std::string& key, const T& default_value) const {
        if (config_data_.contains(key)) {
            return config_data_[key].get<T>();
        }
        return default_value;
    }

    template<typename T>
    void set(const std::string& key, const T& value) {
        config_data_[key] = value;
    }

private:
    // Helper to initialize with default values if config is empty or specific keys are missing
    void initializeDefaultSettings();
};

} // namespace CodezillA

#endif // CONFIGURATION_MANAGER_H
