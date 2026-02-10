#pragma once

#include <string>
#include <vector>
#include <map>

namespace CodezillA {
namespace SCC {

struct FileStats {
    std::string path;
    long code_lines;
    long comment_lines;
    long blank_lines;
};

struct LanguageStats {
    std::string name;
    long total_files;
    long total_code_lines;
    long total_comment_lines;
    long total_blank_lines;
    std::vector<FileStats> files; // Optional: detailed stats per file
};

struct OverallStats {
    long total_files;
    long total_code_lines;
    long total_comment_lines;
    long total_blank_lines;
    std::vector<LanguageStats> languages;
};

} // namespace SCC
} // namespace CodezillA
