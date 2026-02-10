#pragma once

#include <string>
#include <vector>

namespace CodezillA {

struct AnalysisResult {
    std::string file_path;
    std::string rule_id;
    std::string message;
    std::string severity;
    int line_number;

    AnalysisResult(std::string fp, std::string rid, std::string msg, std::string sev, int line = 0)
        : file_path(std::move(fp)), rule_id(std::move(rid)), message(std::move(msg)),
          severity(std::move(sev)), line_number(line) {}
};

} // namespace CodezillA
