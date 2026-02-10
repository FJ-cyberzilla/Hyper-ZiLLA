cmake_minimum_required(VERSION 3.15)
project(CodezillA VERSION 3.0.0 LANGUAGES CXX)

# C++ Standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Enable colored output
if(CMAKE_VERSION VERSION_GREATER_EQUAL "3.24")
    set(CMAKE_COLOR_DIAGNOSTICS ON)
endif()

# Build type
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

# Compiler flags
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra -Wpedantic")
set(CMAKE_CXX_FLAGS_DEBUG "-g -O0 -DDEBUG")
set(CMAKE_CXX_FLAGS_RELEASE "-O3 -DNDEBUG")

# Output directories
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR})
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

# ============================================================================
# Find Required Packages
# ============================================================================

# Core dependencies
find_package(Threads REQUIRED)

# SQLite3 - Required for database functionality
find_package(SQLite3 QUIET)
if(NOT SQLite3_FOUND)
    message(STATUS "SQLite3 not found via find_package, trying pkg-config...")
    find_package(PkgConfig QUIET)
    if(PkgConfig_FOUND)
        pkg_check_modules(SQLite3 QUIET sqlite3)
    endif()
endif()

if(SQLite3_FOUND OR SQLite3_LIBRARIES)
    message(STATUS "SQLite3 found: ${SQLite3_LIBRARIES}")
    set(HAVE_SQLITE3 TRUE)
else()
    message(STATUS "SQLite3 not found - will try to link with -lsqlite3")
    set(HAVE_SQLITE3 FALSE)
endif()

# OpenSSL - Required for SHA256 hashing in AI engine
find_package(OpenSSL QUIET)
if(OpenSSL_FOUND)
    message(STATUS "OpenSSL found: ${OPENSSL_VERSION}")
    set(HAVE_OPENSSL TRUE)
else()
    message(STATUS "OpenSSL not found - will try to link with -lssl -lcrypto")
    set(HAVE_OPENSSL FALSE)
endif()

# Python3 - Optional for AI service
find_package(Python3 COMPONENTS Interpreter QUIET)
if(Python3_FOUND)
    message(STATUS "Python3 found: ${Python3_EXECUTABLE}")
    add_definitions(-DPYTHON3_FOUND)
    add_definitions(-DPYTHON3_EXECUTABLE="${Python3_EXECUTABLE}")
else()
    message(WARNING "Python3 not found - AI features will be limited")
endif()

# Optional AI dependencies (your original ONNX support)
find_package(ONNXRuntime QUIET)
if(ONNXRuntime_FOUND)
    message(STATUS "ONNXRuntime found - advanced AI features enabled")
    add_definitions(-DHAVE_ONNXRUNTIME)
endif()

# ============================================================================
# Source Files Organization
# ============================================================================

# Core sources
set(CORE_SOURCES
    src/core/error_handler.cpp
    src/core/database_manager.cpp
    src/core/configuration_manager.cpp
)

# Analysis sources
set(ANALYSIS_SOURCES
    src/analysis/analysis_manager.cpp
    src/analysis/ai/ai_engine.cpp
    src/analysis/languages/cpp_analyzer.cpp
)

# UI sources
set(UI_SOURCES
    src/ui/menu_system.cpp
)

# Utility sources
set(UTILS_SOURCES
    src/utils/logger.cpp
    src/utils/scc_parser.cpp
)

# All sources combined
set(ALL_SOURCES
    ${CORE_SOURCES}
    ${ANALYSIS_SOURCES}
    ${UI_SOURCES}
    ${UTILS_SOURCES}
)

# ============================================================================
# Main Executable
# ============================================================================

add_executable(codezilla
    main.cpp
    ${ALL_SOURCES}
)

# Global include directories for all headers
target_include_directories(codezilla PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/core
    ${CMAKE_CURRENT_SOURCE_DIR}/include/db
    ${CMAKE_CURRENT_SOURCE_DIR}/include/ui
    ${CMAKE_CURRENT_SOURCE_DIR}/include/utils
    ${CMAKE_CURRENT_SOURCE_DIR}/include/analysis
    ${CMAKE_CURRENT_SOURCE_DIR}/include/nlohmann
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/core
    ${CMAKE_CURRENT_SOURCE_DIR}/src/ui
    ${CMAKE_CURRENT_SOURCE_DIR}/src/analysis
    ${CMAKE_CURRENT_SOURCE_DIR}/src/analysis/ai
    ${CMAKE_CURRENT_SOURCE_DIR}/src/analysis/languages
    ${CMAKE_CURRENT_SOURCE_DIR}/src/utils
)

# Add SQLite3 includes if found
if(HAVE_SQLITE3 AND SQLite3_INCLUDE_DIRS)
    target_include_directories(codezilla PUBLIC ${SQLite3_INCLUDE_DIRS})
endif()

# Add OpenSSL includes if found
if(HAVE_OPENSSL)
    target_include_directories(codezilla PUBLIC ${OPENSSL_INCLUDE_DIR})
endif()

# ============================================================================
# Linking Libraries
# ============================================================================

target_link_libraries(codezilla
    PRIVATE
    Threads::Threads
    ${CMAKE_DL_LIBS}
)

# Link SQLite3
if(HAVE_SQLITE3 AND SQLite3_LIBRARIES)
    target_link_libraries(codezilla PRIVATE ${SQLite3_LIBRARIES})
else()
    # Try direct linking
    target_link_libraries(codezilla PRIVATE sqlite3)
endif()

# Link OpenSSL
if(HAVE_OPENSSL)
    target_link_libraries(codezilla PRIVATE OpenSSL::SSL OpenSSL::Crypto)
else()
    # Try direct linking
    target_link_libraries(codezilla PRIVATE ssl crypto)
endif()

# Link ONNXRuntime if available (your original support)
if(ONNXRuntime_FOUND)
    target_link_libraries(codezilla PRIVATE ONNXRuntime::ONNXRuntime)
endif()

# ============================================================================
# Test Executable
# ============================================================================

add_executable(system_test
    test/system_test.cpp
    ${ALL_SOURCES}
)

# Include directories for test
target_include_directories(system_test PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/core
    ${CMAKE_CURRENT_SOURCE_DIR}/include/db
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/core
    ${CMAKE_CURRENT_SOURCE_DIR}/src/analysis
    ${CMAKE_CURRENT_SOURCE_DIR}/src/analysis/ai
    ${CMAKE_CURRENT_SOURCE_DIR}/src/utils
)

# Add SQLite3 includes for test if found
if(HAVE_SQLITE3 AND SQLite3_INCLUDE_DIRS)
    target_include_directories(system_test PUBLIC ${SQLite3_INCLUDE_DIRS})
endif()

# Add OpenSSL includes for test if found
if(HAVE_OPENSSL)
    target_include_directories(system_test PUBLIC ${OPENSSL_INCLUDE_DIR})
endif()

# Link libraries for test
target_link_libraries(system_test
    PRIVATE
    Threads::Threads
    ${CMAKE_DL_LIBS}
)

# Link SQLite3 for test
if(HAVE_SQLITE3 AND SQLite3_LIBRARIES)
    target_link_libraries(system_test PRIVATE ${SQLite3_LIBRARIES})
else()
    target_link_libraries(system_test PRIVATE sqlite3)
endif()

# Link OpenSSL for test
if(HAVE_OPENSSL)
    target_link_libraries(system_test PRIVATE OpenSSL::SSL OpenSSL::Crypto)
else()
    target_link_libraries(system_test PRIVATE ssl crypto)
endif()

# ============================================================================
# Custom Targets
# ============================================================================

# Add custom command to run the test
add_custom_target(test_codezilla
    COMMAND system_test
    WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
    COMMENT "Running CodezillA system tests..."
)

# Format target (if clang-format is available)
find_program(CLANG_FORMAT clang-format)
if(CLANG_FORMAT)
    add_custom_target(format
        COMMAND find ${CMAKE_SOURCE_DIR}/src ${CMAKE_SOURCE_DIR}/include ${CMAKE_SOURCE_DIR}/main.cpp
                -name "*.cpp" -o -name "*.h" | xargs ${CLANG_FORMAT} -i
        COMMENT "Formatting source code with clang-format"
    )
endif()

# Lint target (if cppcheck is available)
find_program(CPPCHECK cppcheck)
if(CPPCHECK)
    add_custom_target(lint
        COMMAND ${CPPCHECK} --enable=all --inconclusive --std=c++17
                ${CMAKE_SOURCE_DIR}/src ${CMAKE_SOURCE_DIR}/main.cpp
        COMMENT "Running cppcheck static analysis"
    )
endif()

# Run target for convenience
add_custom_target(run
    COMMAND ${CMAKE_BINARY_DIR}/codezilla
    DEPENDS codezilla
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Running CodezillA"
)

# ============================================================================
# Installation Rules
# ============================================================================

# Install main executable
install(TARGETS codezilla
    RUNTIME DESTINATION bin
    COMPONENT runtime
)

# Install test executable (optional)
install(TARGETS system_test
    RUNTIME DESTINATION bin
    COMPONENT testing
    OPTIONAL
)

# Install AI service Python script
install(FILES
    ${CMAKE_SOURCE_DIR}/src/analysis/ai/ai_service.py
    DESTINATION share/codezilla
    PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE
    COMPONENT runtime
)

# Install headers (for development) - FIX APPLIED HERE
install(DIRECTORY
    ${CMAKE_SOURCE_DIR}/include/
    DESTINATION include/codezilla
    FILES_MATCHING PATTERN "*.h" PATTERN "*.hpp"
    # FIX: COMPONENT MUST be the last argument in install(DIRECTORY) when using patterns.
    COMPONENT development
)

# Install configuration files (if directory exists)
if(EXISTS ${CMAKE_SOURCE_DIR}/config)
    install(DIRECTORY config/
        DESTINATION etc/codezilla
        COMPONENT configuration
    )
endif()

# Install models (if directory exists - for your ONNX support)
if(EXISTS ${CMAKE_SOURCE_DIR}/models)
    install(DIRECTORY models/
        DESTINATION share/codezilla/models
        COMPONENT models
    )
endif()

# Create default config.json template during installation
install(CODE "
    set(CONFIG_FILE \"\${CMAKE_INSTALL_PREFIX}/share/codezilla/config.json.template\")
    file(WRITE \"\${CONFIG_FILE}\" \"{
      \\\"application\\\": {
        \\\"name\\\": \\\"CodeZilla\\\",
        \\\"version\\\": \\\"3.0.0\\\",
        \\\"environment\\\": \\\"production\\\"
      },
      \\\"database\\\": {
        \\\"path\\\": \\\"codezilla.db\\\",
        \\\"connection_timeout\\\": 30,
        \\\"max_connections\\\": 10
      },
      \\\"logging\\\": {
        \\\"level\\\": \\\"INFO\\\",
        \\\"file\\\": \\\"codezilla.log\\\"
      },
      \\\"ai_engine\\\": {
        \\\"python_executable\\\": \\\"python3\\\",
        \\\"service_path\\\": \\\"share/codezilla/ai_service.py\\\",
        \\\"model_type\\\": \\\"advanced\\\",
        \\\"timeout_seconds\\\": 30,
        \\\"max_retries\\\": 3,
        \\\"enable_caching\\\": true,
        \\\"enable_learning\\\": true,
        \\\"cache_max_size\\\": 1000
      }
    }\")
    message(STATUS \"Installed config template to: \${CONFIG_FILE}\")
")

# ============================================================================
# Testing Support
# ============================================================================

enable_testing()
add_test(NAME SystemTest COMMAND system_test)
set_tests_properties(SystemTest PROPERTIES
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    TIMEOUT 60
)

# ============================================================================
# Build Information Display
# ============================================================================

message(STATUS "")
message(STATUS "╔════════════════════════════════════════════════════════════════╗")
message(STATUS "║        CodezillA v${PROJECT_VERSION} Build Configuration              ║")
message(STATUS "╚════════════════════════════════════════════════════════════════╝")
message(STATUS "")
message(STATUS "Build Configuration:")
message(STATUS "  Build type:           ${CMAKE_BUILD_TYPE}")
message(STATUS "  C++ Compiler:         ${CMAKE_CXX_COMPILER}")
message(STATUS "  C++ Standard:         C++${CMAKE_CXX_STANDARD}")
message(STATUS "  Install prefix:       ${CMAKE_INSTALL_PREFIX}")
message(STATUS "")
message(STATUS "Dependencies:")
message(STATUS "  Threads:              ✓ Found")
if(HAVE_SQLITE3)
    message(STATUS "  SQLite3:              ✓ Found (${SQLite3_LIBRARIES})")
else()
    message(STATUS "  SQLite3:              ⚠ Not found via find_package (will try -lsqlite3)")
endif()
if(HAVE_OPENSSL)
    message(STATUS "  OpenSSL:              ✓ Found (${OPENSSL_VERSION})")
else()
    message(STATUS "  OpenSSL:              ⚠ Not found via find_package (will try -lssl -lcrypto)")
endif()
if(Python3_FOUND)
    message(STATUS "  Python3:              ✓ Found (${Python3_EXECUTABLE})")
else()
    message(STATUS "  Python3:              ✗ Not found (AI features limited)")
endif()
if(ONNXRuntime_FOUND)
    message(STATUS "  ONNXRuntime:          ✓ Found (advanced AI enabled)")
else()
    message(STATUS "  ONNXRuntime:          ✗ Not found (optional)")
endif()
message(STATUS "")
message(STATUS "Build Targets:")
message(STATUS "  codezilla             Main application")
message(STATUS "  system_test           Test suite")
message(STATUS "  test_codezilla        Run tests")
if(CLANG_FORMAT)
    message(STATUS "  format                Format code with clang-format")
endif()
if(CPPCHECK)
    message(STATUS "  lint                  Static analysis with cppcheck")
endif()
message(STATUS "  run                   Build and run application")
message(STATUS "")
message(STATUS "Compiler Flags:")
message(STATUS "  CXX FLAGS:            ${CMAKE_CXX_FLAGS}")
if(CMAKE_BUILD_TYPE STREQUAL "Debug")
    message(STATUS "  Debug FLAGS:          ${CMAKE_CXX_FLAGS_DEBUG}")
elseif(CMAKE_BUILD_TYPE STREQUAL "Release")
    message(STATUS "  Release FLAGS:        ${CMAKE_CXX_FLAGS_RELEASE}")
endif()
message(STATUS "")
message(STATUS "╚════════════════════════════════════════════════════════════════╝")
message(STATUS "")

# ============================================================================
# Post-build Commands
# ============================================================================

# Make Python script executable after build
add_custom_command(TARGET codezilla POST_BUILD
    COMMAND chmod +x ${CMAKE_SOURCE_DIR}/src/analysis/ai/ai_service.py
    COMMENT "Making AI service script executable"
    VERBATIM
)

# Copy config.json if it exists
if(EXISTS ${CMAKE_SOURCE_DIR}/config.json)
    add_custom_command(TARGET codezilla POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E copy_if_different
                ${CMAKE_SOURCE_DIR}/config.json
                ${CMAKE_BINARY_DIR}/config.json
        COMMENT "Copying config.json to build directory"
        VERBATIM
    )
endif()
