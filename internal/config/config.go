package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/karsterr/syswatch-daemon/internal/logger"
)

// Config uygulama konfigürasyonu
type Config struct {
	// Daemon ayarları
	Daemon DaemonConfig `json:"daemon"`
	
	// Dashboard ayarları
	Dashboard DashboardConfig `json:"dashboard"`
	
	// Logging ayarları
	Logging LoggingConfig `json:"logging"`
	
	// Metrics ayarları  
	Metrics MetricsConfig `json:"metrics"`
}

// DaemonConfig daemon ayarları
type DaemonConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// DashboardConfig dashboard ayarları
type DashboardConfig struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Host    string `json:"host"`
}

// LoggingConfig logging ayarları
type LoggingConfig struct {
	Level      string `json:"level"`       // debug, info, warn, error
	Format     string `json:"format"`      // text, json
	Output     string `json:"output"`      // stdout, file
	Filename   string `json:"filename"`    // log dosyası adı (output=file ise)
}

// MetricsConfig metrics ayarları
type MetricsConfig struct {
	Interval     int  `json:"interval"`      // Metrik toplama aralığı (saniye)
	EnableCPU    bool `json:"enable_cpu"`
	EnableMemory bool `json:"enable_memory"`
	EnableDisk   bool `json:"enable_disk"`
	EnableNet    bool `json:"enable_network"`
}

// Default varsayılan konfigürasyon
func Default() *Config {
	return &Config{
		Daemon: DaemonConfig{
			Name:    "syswatch-daemon",
			Version: "0.1.0",
		},
		Dashboard: DashboardConfig{
			Enabled: true,
			Port:    8080,
			Host:    "localhost",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Metrics: MetricsConfig{
			Interval:     5,
			EnableCPU:    true,
			EnableMemory: true,
			EnableDisk:   true,
			EnableNet:    true,
		},
	}
}

// Load konfigürasyon dosyasından ayarları yükler
func Load(configPath string) (*Config, error) {
	log := logger.GetLogger()
	
	// Varsayılan konfigürasyon ile başla
	config := Default()
	
	// Dosya var mı kontrol et
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Infof("Konfigürasyon dosyası bulunamadı: %s, varsayılan ayarlar kullanılıyor", configPath)
		return config, nil
	}
	
	// Dosyayı oku
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("konfigürasyon dosyası okunamadı: %w", err)
	}
	
	// JSON'dan parse et
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("konfigürasyon dosyası parse edilemedi: %w", err)
	}
	
	log.Infof("Konfigürasyon dosyası başarıyla yüklendi: %s", configPath)
	return config, nil
}

// Save konfigürasyonu dosyaya kaydeder
func (c *Config) Save(configPath string) error {
	log := logger.GetLogger()
	
	// Dizini oluştur
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("konfigürasyon dizini oluşturulamadı: %w", err)
	}
	
	// JSON formatında serialize et
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("konfigürasyon serialize edilemedi: %w", err)
	}
	
	// Dosyaya yaz
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("konfigürasyon dosyası yazılamadı: %w", err)
	}
	
	log.Infof("Konfigürasyon dosyası başarıyla kaydedildi: %s", configPath)
	return nil
}

// Validate konfigürasyonu doğrular
func (c *Config) Validate() error {
	// Dashboard port kontrolü
	if c.Dashboard.Port < 1024 || c.Dashboard.Port > 65535 {
		return fmt.Errorf("dashboard port geçersiz: %d (1024-65535 arası olmalı)", c.Dashboard.Port)
	}
	
	// Logging level kontrolü
	validLevels := []string{"debug", "info", "warn", "error"}
	isValid := false
	for _, level := range validLevels {
		if c.Logging.Level == level {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("geçersiz log level: %s (debug, info, warn, error olmalı)", c.Logging.Level)
	}
	
	// Metrics interval kontrolü
	if c.Metrics.Interval < 1 || c.Metrics.Interval > 3600 {
		return fmt.Errorf("metrics interval geçersiz: %d (1-3600 saniye arası olmalı)", c.Metrics.Interval)
	}
	
	return nil
}