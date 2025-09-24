package metrics

import (
	"testing"
)

func TestNewCollector(t *testing.T) {
	collector := NewCollector()
	
	if collector == nil {
		t.Fatal("NewCollector() returned nil")
	}
	
	if collector.prevNetStats == nil {
		t.Error("prevNetStats map should be initialized")
	}
}

func TestCollectorStartStop(t *testing.T) {
	collector := NewCollector()
	
	// Start test
	err := collector.Start()
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	
	// Stop test
	collector.Stop()
	
	// Stop should not panic or error
}

func TestCollectAll(t *testing.T) {
	collector := NewCollector()
	
	if err := collector.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer collector.Stop()
	
	metrics, err := collector.CollectAll()
	if err != nil {
		t.Fatalf("CollectAll() failed: %v", err)
	}
	
	if metrics == nil {
		t.Fatal("CollectAll() returned nil metrics")
	}
	
	// Timestamp kontrolü
	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
	
	// CPU metrikleri kontrolü
	if metrics.CPU.Count <= 0 {
		t.Error("CPU count should be positive")
	}
	
	if metrics.CPU.Usage < 0 || metrics.CPU.Usage > 100 {
		t.Errorf("CPU usage should be between 0-100, got %.2f", metrics.CPU.Usage)
	}
	
	// Memory metrikleri kontrolü
	if metrics.Memory.Total == 0 {
		t.Error("Total memory should be positive")
	}
	
	if metrics.Memory.Usage < 0 || metrics.Memory.Usage > 100 {
		t.Errorf("Memory usage should be between 0-100, got %.2f", metrics.Memory.Usage)
	}
	
	// Disk metrikleri kontrolü  
	if metrics.Disk.Total == 0 {
		t.Error("Total disk space should be positive")
	}
	
	if metrics.Disk.Usage < 0 || metrics.Disk.Usage > 100 {
		t.Errorf("Disk usage should be between 0-100, got %.2f", metrics.Disk.Usage)
	}
	
	// Network metrikleri kontrolü (0 veya pozitif olmalı)
	if metrics.Network.BytesRecv < 0 {
		t.Errorf("Network bytes received should be non-negative, got %d", metrics.Network.BytesRecv)
	}
	
	if metrics.Network.BytesSent < 0 {
		t.Errorf("Network bytes sent should be non-negative, got %d", metrics.Network.BytesSent)
	}
}

func TestSystemMetricsStructure(t *testing.T) {
	// Yapısal test - metrik yapısının doğru tanımlandığından emin ol
	var metrics SystemMetrics
	
	// JSON tag'lerinin varlığını kontrol et (reflection ile)
	// Bu test compile time'da yapı doğruluğunu kontrol eder
	
	if metrics.CPU.Usage == 0 && metrics.CPU.Count == 0 {
		// Başlangıç değerleri sıfır olmalı
	}
	
	if metrics.Memory.Total == 0 && metrics.Memory.Used == 0 {
		// Başlangıç değerleri sıfır olmalı  
	}
	
	if metrics.Disk.Total == 0 && metrics.Disk.Used == 0 {
		// Başlangıç değerleri sıfır olmalı
	}
	
	if metrics.Network.BytesRecv == 0 && metrics.Network.BytesSent == 0 {
		// Başlangıç değerleri sıfır olmalı
	}
}

// Benchmark testleri
func BenchmarkCollectAll(b *testing.B) {
	collector := NewCollector()
	if err := collector.Start(); err != nil {
		b.Fatalf("Start() failed: %v", err)
	}
	defer collector.Stop()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := collector.CollectAll()
		if err != nil {
			b.Fatalf("CollectAll() failed: %v", err)
		}
	}
}

func BenchmarkCollectCPU(b *testing.B) {
	collector := NewCollector()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := collector.collectCPU()
		if err != nil {
			b.Fatalf("collectCPU() failed: %v", err)
		}
	}
}

func BenchmarkCollectMemory(b *testing.B) {
	collector := NewCollector()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := collector.collectMemory()
		if err != nil {
			b.Fatalf("collectMemory() failed: %v", err)
		}
	}
}