package config

import "time"

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Strategy    string            `yaml:"strategy"`
	HealthCheck HealthCheckConfig `yaml:"health_check"`
	TLS         TLSConfig         `yaml:"tls"`
	BackendTLS  BackendTLSConfig  `yaml:"backend_tls"`
	Upstreams   []UpstreamConfig  `yaml:"upstreams"`
	Logging     LoggingConfig     `yaml:"logging"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type HealthCheckConfig struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type BackendTLSConfig struct {
	Enabled            bool `yaml:"enabled"`
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
}

type UpstreamConfig struct {
	URL               string `yaml:"url"`
	Weight            int    `yaml:"weight"`
	HealthCheckPath   string `yaml:"health_check_path"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}
