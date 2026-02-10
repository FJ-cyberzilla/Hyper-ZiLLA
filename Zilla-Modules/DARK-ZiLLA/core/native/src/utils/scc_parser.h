#pragma once

#include <string>
#include <vector>
#include <memory> // For std::shared_ptr
#include <optional> // For std::optional

#include "utils/scc_types.h"
#include "utils/logger.h"
#include "core/error_handler.h"

namespace CodezillA {

class SccParser {
private:
    std::shared_ptr<Logger> logger_;
    std::shared_ptr<ErrorHandler> error_handler_;

    std::string executeSccCommand(const std::string& directory_path) const;
    SCC::OverallStats parseSccJsonOutput(const std::string& json_output) const;

public:
    SccParser(std::shared_ptr<ErrorHandler> error_handler, std::shared_ptr<Logger> logger);
    
    std::optional<SCC::OverallStats> analyzeDirectory(const std::string& directory_path) const;
};

} // namespace CodezillA
