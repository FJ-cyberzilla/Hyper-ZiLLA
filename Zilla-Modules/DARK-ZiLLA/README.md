# Cyberzilla Codex - Enterprise Code Quality Analyzer

## Overview
The Cyberzilla Codex is a powerful and flexible static analysis framework designed to integrate various external code quality tools (linters, formatters, etc.) into a unified, concurrent workflow. It helps maintain high code quality by automating checks and fixes across multiple programming languages, providing detailed reports and historical trend analysis.

## Features
-   **Multi-language Support**: Easily configure analysis tools for different programming languages.
-   **Concurrent Analysis**: Processes files in parallel for efficient scanning of large codebases.
-   **Auto-Fix Mode**: Automatically applies fixes suggested by formatter tools with built-in backup and rollback safety.
-   **Customizable Configuration**: Define tools, analysis rules, and project-specific settings via a `.codexrc.json` file.
-   **Detailed Reporting**: Generates comprehensive reports of analysis results.
-   **History and Trend Analysis**: Tracks code quality over time to identify trends and regressions.

## Installation

### Prerequisites
-   Python 3.8 or higher.
-   Ensure you have the required external code quality tools installed for the languages you intend to analyze (e.g., `pylint`, `black` for Python; `eslint`, `prettier` for JavaScript).

To install `black` (example):
```bash
pip install black
```

### Install from PyPI (Recommended)
(Once published, you can install directly via pip)
```bash
pip install codex-enterprise
```

### Install from Source (for development)

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/codex-enterprise.git
    cd codex-enterprise
    ```
2.  **Install dependencies**:
    ```bash
    pip install -r requirements.txt
    ```
3.  **Install the package in editable mode**:
    ```bash
    pip install -e .
    ```
    This makes the `codex` command available globally and reflects any changes you make in the `src/codex_enterprise` directory.

## Usage

The Cyberzilla Codex can be run in two modes: command-line or interactive.

### Configuration (`.codexrc.json`)
The tool looks for a `.codexrc.json` file in the current directory or its parent directories to load configuration. This file defines which tools to use for which languages, their commands, and application settings. An example configuration is implicitly provided by the `src/codex_enterprise/config.py` file.

```json
{
  "app_settings": {
    "max_workers": 4,
    "history_file": "codex_history.json",
    "output_dir": "reports",
    "default_timeout": 30,
    "skip_dirs": ["node_modules", ".git", ".venv"]
  },
  "language_tools": {
    "python": [
      {"tool": "pylint", "command": ["pylint", "--output-format=json"], "check": true, "fix": false, "timeout": 60},
      {"tool": "black", "command": ["black", "-q"], "check": false, "fix": true, "timeout": 30}
    ],
    "javascript": [
      {"tool": "eslint", "command": ["eslint", "--format=json"], "check": true, "fix": false, "timeout": 60},
      {"tool": "prettier", "command": ["prettier", "--write"], "check": false, "fix": true, "timeout": 30}
    ]
  }
}
```

### Command-Line Mode
Run analysis directly from your terminal.

```bash
codex <path_to_directory_to_analyze> [--fix] [--verbose]
```

-   `<path_to_directory_to_analyze>`: The path to the directory containing code you want to analyze (e.g., `.`, `my_project/`). Default is the current directory.
-   `--fix`: Enable auto-fix mode. If enabled, the tool will attempt to fix issues using configured 'fixer' tools. A backup of the original file is created (`.codex.bak`) before any changes, and rolled back if the fixer fails.
-   `--verbose`: Enable verbose output, showing progress for each file.

**Examples:**

-   Analyze the current directory in verbose mode:
    ```bash
    codex . --verbose
    ```
-   Analyze and fix issues in the `my_python_project` directory:
    ```bash
    codex my_python_project --fix
    ```

### Interactive Mode
Launch an interactive menu for guided usage.

```bash
codex --interactive
```

Follow the on-screen prompts to perform analysis, fix issues, view history, or inspect configuration.

### Direct Execution (without installation)
If you prefer not to install the package, you can run it directly from the project root:

```bash
python -m src.codex_enterprise.cli <path_to_directory_to_analyze> [--fix] [--verbose]
```
Example: `python -m src.codex_enterprise.cli . --fix --verbose`

## Contributing
We welcome contributions! Please feel free to open issues or submit pull requests.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details (if applicable).
