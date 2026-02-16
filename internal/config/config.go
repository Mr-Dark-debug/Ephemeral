package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nick       string         `yaml:"nick"`
	Port       int            `yaml:"port"`
	Discovery  DiscoveryConfig `yaml:"discovery"`
	Rooms      []RoomConfig   `yaml:"rooms"`
	Security   SecurityConfig `yaml:"security"`
	Logging    LoggingConfig  `yaml:"logging"`
}

type DiscoveryConfig struct {
	MDNS        bool `yaml:"mdns"`
	UDPFallback bool `yaml:"udp_fallback"`
}

type RoomConfig struct {
	Name      string `yaml:"name"`
	Encrypted bool   `yaml:"encrypted"`
}

type SecurityConfig struct {
	PersistKeys     bool `yaml:"persist_keys"`
	KeyRotationDays int  `yaml:"key_rotation_days"`
}

type LoggingConfig struct {
	Level         string `yaml:"level"`
	EphemeralLogs bool   `yaml:"ephemeral_logs"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		Nick: "guest",
		Port: 9999,
		Discovery: DiscoveryConfig{
			MDNS:        true,
			UDPFallback: true,
		},
		Rooms: []RoomConfig{
			{Name: "global", Encrypted: false},
		},
		Security: SecurityConfig{
			PersistKeys: false,
		},
		Logging: LoggingConfig{
			Level:         "info",
			EphemeralLogs: true,
		},
	}
}
