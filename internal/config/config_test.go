package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	
	// Varsayılan değerleri kontrol et
	if cfg.Daemon.Name != "syswatch-daemon" {
		t.Errorf("Expected daemon name 'syswatch-daemon', got '%s'", cfg.Daemon.Name)
	}
	
	if cfg.Dashboard.Port != 8080 {
		t.Errorf("Expected dashboard port 8080, got %d", cfg.Dashboard.Port)
	}
	
	if cfg.Metrics.Interval != 5 {
		t.Errorf("Expected metrics interval 5, got %d", cfg.Metrics.Interval)
	}
}

func TestLoad(t *testing.T) {
	// Test için geçici dosya oluştur
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")
	
	testConfig := &Config{
		Daemon: DaemonConfig{
			Name:    "test-daemon",
			Version: "1.0.0",
		},
		Dashboard: DashboardConfig{
			Enabled: true,
			Port:    9090,
			Host:    "test-host",
		},
		Logging: LoggingConfig{
			Level:  "debug",
			Format: "json",
			Output: "stdout",
		},
		Metrics: MetricsConfig{
			Interval:     10,
			EnableCPU:    true,
			EnableMemory: false,
			EnableDisk:   true,
			EnableNet:    false,
		},
	}
	
	// Konfigürasyonu dosyaya yaz
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	
	// Konfigürasyonu yükle
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Değerleri kontrol et
	if loadedConfig.Daemon.Name != testConfig.Daemon.Name {
		t.Errorf("Expected daemon name '%s', got '%s'", testConfig.Daemon.Name, loadedConfig.Daemon.Name)
	}
	
	if loadedConfig.Dashboard.Port != testConfig.Dashboard.Port {
		t.Errorf("Expected dashboard port %d, got %d", testConfig.Dashboard.Port, loadedConfig.Dashboard.Port)
	}
	
	if loadedConfig.Metrics.Interval != testConfig.Metrics.Interval {
		t.Errorf("Expected metrics interval %d, got %d", testConfig.Metrics.Interval, loadedConfig.Metrics.Interval)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	// Var olmayan dosya yükle - varsayılan konfigürasyon dönmeli
	cfg, err := Load("nonexistent.json")
	if err != nil {
		t.Fatalf("Expected no error for nonexistent file, got: %v", err)
	}
	
	// Varsayılan değerler olmalı
	if cfg.Daemon.Name != "syswatch-daemon" {
		t.Errorf("Expected default daemon name, got '%s'", cfg.Daemon.Name)
	}
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "save_test.json")
	
	cfg := Default()
	cfg.Dashboard.Port = 9999
	
	// Konfigürasyonu kaydet
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Dosyanın var olduğunu kontrol et
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
	
	// Dosyayı tekrar yükle ve kontrol et
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}
	
	if loadedConfig.Dashboard.Port != 9999 {
		t.Errorf("Expected port 9999, got %d", loadedConfig.Dashboard.Port)
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name:      "valid config",
			config:    Default(),
			expectErr: false,
		},
		{
			name: "invalid port - too low",
			config: &Config{
				Dashboard: DashboardConfig{Port: 100},
				Logging:   LoggingConfig{Level: "info"},
				Metrics:   MetricsConfig{Interval: 5},
			},
			expectErr: true,
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Dashboard: DashboardConfig{Port: 70000},
				Logging:   LoggingConfig{Level: "info"},
				Metrics:   MetricsConfig{Interval: 5},
			},
			expectErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Dashboard: DashboardConfig{Port: 8080},
				Logging:   LoggingConfig{Level: "invalid"},
				Metrics:   MetricsConfig{Interval: 5},
			},
			expectErr: true,
		},
		{
			name: "invalid metrics interval - too low",
			config: &Config{
				Dashboard: DashboardConfig{Port: 8080},
				Logging:   LoggingConfig{Level: "info"},
				Metrics:   MetricsConfig{Interval: 0},
			},
			expectErr: true,
		},
		{
			name: "invalid metrics interval - too high",
			config: &Config{
				Dashboard: DashboardConfig{Port: 8080},
				Logging:   LoggingConfig{Level: "info"},
				Metrics:   MetricsConfig{Interval: 5000},
			},
			expectErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectErr && err == nil {
				t.Error("Expected validation error, but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}