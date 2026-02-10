# CodezillA Project Wiki

## 1. Project Overview

CodezillA is a C++ Command Line Interface (CLI) application designed for static code analysis. Its primary goal is to assist developers in identifying potential security vulnerabilities, performance bottlenecks, and code quality issues. The application features an interactive menu system with a distinct visual theme.

## 2. Key Features

*   **Interactive Menu System:** A user-friendly CLI interface with a sleek, modern yet retro DOS-style design, featuring a new ASCII art banner and themed colors.
*   **Code Analysis:** Supports analysis for C++, Python, Go, Java, and JavaScript.
*   **AI Simulation:** Integrates a simulated AI engine for:
    *   **Code Smell Detection & Suggestions:** Identifies and suggests fixes for God objects/functions, deep nesting, high cyclomatic complexity, and dead code.
    *   **Security Analysis:** Simulates detection of buffer overflow risks, SQL injection, XSS vulnerabilities, and general warnings for use-after-free/TOCTOU.
    *   **Performance Analysis:** Simulates detection of unnecessary allocations, N+1 query patterns, and inefficient algorithms.
    *   The AI logic is simulated via IPC with a Python script (`ai_service.py`), providing dynamic, contextualized placeholder responses.
*   **Configuration Management:** A `ConfigurationManager` allows users to view and modify application settings (e.g., default analysis directory).
*   **Database Integration:** Stores AI suggestions in `ai_database.db`.
*   **Build System:** Uses CMake for building the C++ application.

## 3. Architecture

CodezillA follows a modular architecture:

*   **`ui/`**: Handles the user interface and interactive menu.
*   **`analysis/`**: Contains the core analysis logic, including `AnalysisManager`, AI engine (`ai/ai_engine.cpp`), and language-specific analyzers (`languages/`).
*   **`core/`**: Provides fundamental services like error handling (`error_handler.cpp`), database management (`database_manager.cpp`), and configuration (`configuration_manager.cpp`).
*   **`utils/`**: Offers utility functions like logging (`logger.cpp`) and code statistics parsing (`scc_parser.cpp`).
*   **`ai_service.py`**: A Python script acting as the backend for simulated AI analysis, communicating via IPC.

## 4. Technologies Used

*   **Language:** C++ (modern standards with CMake)
*   **JSON Handling:** `nlohmann/json` library
*   **Database:** SQLite (for AI suggestions)
*   **AI Simulation:** Python script (`ai_service.py`) for simulated AI logic.

## 5. Build and Run Instructions

### Prerequisites
*   C++ Compiler (GCC/Clang)
*   CMake (version 3.10 or higher)
*   Python 3 (for AI service simulation)

### Building the Project

Navigate to the project root and run:
```bash
mkdir build
cd build
cmake ..
make
```
Alternatively:
```bash
cmake -B build
cmake --build build
```

### Running the Application

After a successful build, execute from the project root:
```bash
./build/codezilla
```
This launches the interactive CLI menu.

## 6. AI Simulation Limitations

*   **Simulated Logic:** The AI engine's analysis and suggestions are currently simulated using heuristics and pattern matching in Python. It does not integrate with actual ML models or external AI APIs.
*   **AI Service Path:** Due to persistent compilation issues in the build environment related to `shared_ptr` construction, the path to the Python AI service (`ai_service.py`) is hardcoded within the C++ code. It cannot be configured via the UI at this time.

## 7. Contribution and Future Work

*   **Contributions:** Open to contributions for enhancing AI detection logic, improving error handling, and expanding language support.
*   **Future Enhancements:**
    *   Integration of actual AI models (if environment permits).
    *   More sophisticated static analysis for deeper code smell detection.
    *   Error handling improvements.
    *   Potentially enabling configurable AI service path if build environment issues are resolved.

This wiki provides a detailed overview of the project, its features, architecture, and current limitations.