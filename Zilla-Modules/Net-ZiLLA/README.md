# Net-ZiLLA

**Net-ZiLLA** is a high-performance, production-ready security analysis engine and network intelligence toolset. Designed for rapid URL and IP assessment, it combines multi-vector analysis, concurrent processing, and automated risk scoring to provide deep visibility into potential cyber threats.

---

## 🚀 Key Features

*   **AI-Powered Intelligence**:
    *   **ML Link Analysis**: Integrated `MLAgent` utilizing specialized Go-based models for phishing and malware detection.
    *   **SMS Security**: Specialized analysis for SMS-based phishing (smishing) patterns.
    *   **AI Orchestration**: Correlates AI findings with traditional security signals for high-confidence verdicts.
*   **Hybrid Analysis Sandbox**:
    *   Supports multiple analysis modes: **Docker** (local isolation), **ANY.RUN** (remote analysis), or **Hybrid** (simultaneous local and remote execution).
*   **Enhanced URL Enrichment**:
    *   Deep inspection of URL structural integrity, including **Entropy calculation**, **TLD risk assessment**, and **Homograph attack detection** (Punycode/IDN).
*   **Concurrent Analysis Pipeline**: Utilizes Go's concurrency primitives to run multiple security checks (DNS, WHOIS, Threat Intel, SSL) in parallel.
*   **Advanced Network Toolkit**:
    *   **Redirect Tracer**: Deep inspection of HTTP redirect chains with hop-by-hop threat scoring.
    *   **SSL/TLS Analyzer**: Evaluation of certificate validity, key strength, and expiration monitoring.
    *   **Infrastructure Intel**: Automated DNSSEC, MX, and NS record verification.
*   **Dual-Interface Support**:
    *   **Interactive CLI**: Professional ASCII-based menu system with detailed enrichment data display.
    *   **REST API**: Structured JSON API for integration into SOC workflows or CI/CD pipelines.

---

## 🏗️ Architecture

Net-ZiLLA follows a modular, service-oriented architecture:

*   **`internal/ai/`**: AI/ML analysis bridge and Go-based threat models.
*   **`internal/analyzer/`**: Core orchestration logic, domain analysis, and safety screening.
*   **`internal/network/`**: Low-level network utilities (HTTP, DNS, SSL, etc.).
*   **`internal/services/`**: Business logic layer bridging the API/CLI and the analysis engine.
*   **`internal/storage/`**: Persistence layer using SQLite for historical analysis data.
*   **`internal/utils/`**: Common utilities including exponential **backoff** implementations.
*   **`pkg/`**: Core library components for logging, metrics tracking, and distributed tracing.

---

## ⚙️ Configuration

The application can be configured via `config.yaml` or Environment Variables.

| Environment Variable | Description |
| :--- | :--- |
| `SERVER_PORT` | Port for the REST API (Default: 8080) |
| `LOG_FORMAT` | Set to `json` for structured logging |
| `VT_API_KEY` | VirusTotal API Key |
| `ABUSEIPDB_API_KEY` | AbuseIPDB API Key |

---

## 🛠️ Installation & Setup

### Prerequisites
*   Go 1.24+
*   SQLite3

### Quick Start
1.  **Build the project**:
    ```bash
    make build
    ```
2.  **Run Tests**:
    ```bash
    make test
    ```
3.  **Start the Application**:
    *   **CLI Mode**: `./netzilla`
    *   **API Mode**: Set `server.enable_api: true` in `config.yaml` and run `./netzilla`.

---

## 📖 Usage

### Interactive CLI
Launch the tool without flags to enter the secure menu. The CLI now displays **URL Enrichment** details such as entropy and TLD risk.

### REST API
**Endpoint**: `POST /api/v1/analyze`
**Request**:
```json
{
  "target": "https://suspicious-target.com"
}
```

---

## 🧪 Development & Quality

The project maintains a rigorous testing standard with a focus on reliability:

*   **Test Coverage**: Currently **~60%** and climbing, covering all core orchestration and analysis logic.
*   **Resiliency**: Integrated **Circuit Breakers** and **Exponential Backoff** for external network requests.
*   **Observability**: Integrated metrics and span-based tracing for performance monitoring.

*   **`make fmt`**: Auto-format all Go source files.
*   **`make test`**: Execute the full test suite.
*   **`make coverage`**: Generate an HTML coverage report.

---

## ⚖️ License
This project is licensed under the MIT License - see the `LICENSE` file for details.
