package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Service     *ServiceConfig    `mapstructure:"service"`
	Security    SecurityConfig    `mapstructure:"security"`
	AI          AIConfig          `mapstructure:"ai"`
	Network     NetworkConfig     `mapstructure:"network"`
	ThreatIntel ThreatIntelConfig `mapstructure:"threat_intel"`
	Sandbox     SandboxConfig     `mapstructure:"sandbox"`
	Analysis    *AnalysisConfig   `mapstructure:"analysis"`
	Output      OutputConfig      `mapstructure:"output"`
}

type ServerConfig struct {
	EnableAPI bool   `mapstructure:"enable_api"`
	EnableCLI bool   `mapstructure:"enable_cli"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
}

type ServiceConfig struct {
	WaitForDuplicateAnalysis bool `mapstructure:"wait_for_duplicate_analysis"`
}

type SecurityConfig struct {
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	RateLimit      int           `mapstructure:"rate_limit"`
}

type AIConfig struct {
	EnableAI            bool    `mapstructure:"enable_ai"`
	ConfidenceThreshold float64 `mapstructure:"confidence_threshold"`
}

type NetworkConfig struct {
	ProxyEnabled   bool   `mapstructure:"proxy_enabled"`
	ProxyURL       string `mapstructure:"proxy_url"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
	MaxRedirects   int    `mapstructure:"max_redirects"`
	UserAgent      string `mapstructure:"user_agent"`
}

type ThreatIntelConfig struct {
	EnabledProviders []string `mapstructure:"enabled_providers"`
	CacheTTLHours    int      `mapstructure:"cache_ttl_hours"`
}

type SandboxConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Type           string `mapstructure:"type"`
	DockerImage    string `mapstructure:"docker_image"`
	TimeoutMinutes int    `mapstructure:"timeout_minutes"`
	AutoDestroy    bool   `mapstructure:"auto_destroy"`
}

type AnalysisConfig struct {
	DeepScan          bool              `mapstructure:"deep_scan"`
	ScoringThresholds ScoringThresholds `mapstructure:"scoring_thresholds"`
	TimeoutSeconds    int               `mapstructure:"timeout_seconds"`
}

type ScoringThresholds struct {
	HighRisk   int `mapstructure:"high_risk"`
	MediumRisk int `mapstructure:"medium_risk"`
	LowRisk    int `mapstructure:"low_risk"`
}

type OutputConfig struct {
	SaveReports  bool   `mapstructure:"save_reports"`
	ReportFormat string `mapstructure:"report_format"`
	ReportPath   string `mapstructure:"report_path"`
	EnableColors bool   `mapstructure:"enable_colors"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("security.request_timeout", 30*time.Second)

	viper.AutomaticEnv()
	viper.BindEnv("threat_intel.vt_key", "VT_API_KEY")
	viper.BindEnv("threat_intel.abuse_key", "ABUSEIPDB_API_KEY")
	viper.BindEnv("threat_intel.av_key", "ALIENVAULT_API_KEY")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
