package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/karsterr/syswatch-daemon/internal/logger"
	"github.com/karsterr/syswatch-daemon/internal/metrics"
)

// Daemon ana daemon yapısı
type Daemon struct {
	mu            sync.RWMutex
	running       bool
	metricsCol    *metrics.Collector
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

// New yeni daemon instance oluşturur
func New() *Daemon {
	return &Daemon{
		metricsCol: metrics.NewCollector(),
		stopChan:   make(chan struct{}),
	}
}

// Start daemon'u başlatır
func (d *Daemon) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.running {
		return nil
	}

	log := logger.GetLogger()
	log.Info("Daemon başlatılıyor...")

	// Metrics collector'ı başlat
	if err := d.metricsCol.Start(); err != nil {
		return err
	}

	d.running = true

	// Ana iş döngüsünü başlat
	d.wg.Add(1)
	go d.mainLoop(ctx)

	log.Info("Daemon başarıyla başlatıldı")
	return nil
}

// Stop daemon'u durdurur
func (d *Daemon) Stop(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return nil
	}

	log := logger.GetLogger()
	log.Info("Daemon durduruluyor...")

	// Stop sinyali gönder
	close(d.stopChan)

	// Tüm goroutine'lerin bitmesini bekle
	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	// Timeout veya tamamlanma
	select {
	case <-done:
		log.Info("Tüm işlemler temiz şekilde durduruldu")
	case <-ctx.Done():
		log.Warn("Shutdown timeout, zorla çıkılıyor")
	}

	// Metrics collector'ı durdur
	d.metricsCol.Stop()

	d.running = false
	log.Info("Daemon başarıyla durduruldu")
	return nil
}

// IsRunning daemon'un çalışıp çalışmadığını kontrol eder
func (d *Daemon) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}

// mainLoop ana iş döngüsü
func (d *Daemon) mainLoop(ctx context.Context) {
	defer d.wg.Done()

	log := logger.GetLogger()
	ticker := time.NewTicker(5 * time.Second) // Her 5 saniyede bir çalış
	defer ticker.Stop()

	log.Info("Ana iş döngüsü başlatıldı")

	for {
		select {
		case <-ticker.C:
			// Metrikleri topla ve işle
			d.collectAndProcessMetrics()
		case <-d.stopChan:
			log.Info("Stop sinyali alındı, ana döngü sonlandırılıyor")
			return
		case <-ctx.Done():
			log.Info("Context iptal edildi, ana döngü sonlandırılıyor")
			return
		}
	}
}

// collectAndProcessMetrics metrikleri toplar ve işler
func (d *Daemon) collectAndProcessMetrics() {
	log := logger.GetLogger()
	
	// Sistem metriklerini topla
	metrics, err := d.metricsCol.CollectAll()
	if err != nil {
		log.Errorf("Metrikler toplanırken hata: %v", err)
		return
	}

	// Şu an için sadece log'la (ilerleyen aşamalarda dashboard'a gönderilecek)
	log.Infof("CPU: %.1f%%, RAM: %.1f%%, Disk: %.1f%%, Network: R:%.2fMB/s W:%.2fMB/s", 
		metrics.CPU.Usage,
		metrics.Memory.Usage,
		metrics.Disk.Usage,
		float64(metrics.Network.BytesRecv)/(1024*1024),
		float64(metrics.Network.BytesSent)/(1024*1024),
	)
}