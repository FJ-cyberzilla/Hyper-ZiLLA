#ifndef LOGGER_H
#define LOGGER_H

#include <string>
#include <iostream>

namespace CodezillA {

class Logger {
public:
    explicit Logger(const std::string& name) : name_(name) {}

    void info(const std::string& message) {
        std::cout << "[INFO] [" << name_ << "] " << message << std::endl;
    }

    void warn(const std::string& message) {
        std::cout << "[WARN] [" << name_ << "] " << message << std::endl;
    }

    void error(const std::string& message) {
        std::cerr << "[ERROR] [" << name_ << "] " << message << std::endl;
    }

    void debug(const std::string& message) {
        // In a production system, this would typically be conditional on a debug flag
        std::cout << "[DEBUG] [" << name_ << "] " << message << std::endl;
    }

private:
    std::string name_;
};

} // namespace CodezillA

#endif // LOGGER_H