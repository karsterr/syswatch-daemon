package metrics

import (
	"fmt"
	"time"

	"github.com/karsterr/syswatch-daemon/internal/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemMetrics sistem metriklerini içerir
type SystemMetrics struct {
	Timestamp time.Time    `json:"timestamp"`
	CPU       CPUMetrics   `json:"cpu"`
	Memory    MemMetrics   `json:"memory"`
	Disk      DiskMetrics  `json:"disk"`
	Network   NetMetrics   `json:"network"`
}

// CPUMetrics CPU ile ilgili metrikleri içerir
type CPUMetrics struct {
	Usage   float64 `json:"usage"`   // CPU kullanım yüzdesi
	Count   int     `json:"count"`   // CPU çekirdek sayısı
}

// MemMetrics bellek ile ilgili metrikleri içerir
type MemMetrics struct {
	Usage       float64 `json:"usage"`        // Bellek kullanım yüzdesi
	Total       uint64  `json:"total"`        // Toplam bellek (bytes)
	Available   uint64  `json:"available"`    // Kullanılabilir bellek (bytes)
	Used        uint64  `json:"used"`         // Kullanılan bellek (bytes)
}

// DiskMetrics disk ile ilgili metrikleri içerir
type DiskMetrics struct {
	Usage       float64 `json:"usage"`        // Disk kullanım yüzdesi
	Total       uint64  `json:"total"`        // Toplam disk alanı (bytes)
	Used        uint64  `json:"used"`         // Kullanılan disk alanı (bytes)
	Free        uint64  `json:"free"`         // Boş disk alanı (bytes)
}

// NetMetrics ağ ile ilgili metrikleri içerir
type NetMetrics struct {
	BytesRecv   uint64 `json:"bytes_recv"`   // Alınan bytes
	BytesSent   uint64 `json:"bytes_sent"`   // Gönderilen bytes
	PacketsRecv uint64 `json:"packets_recv"` // Alınan paket sayısı
	PacketsSent uint64 `json:"packets_sent"` // Gönderilen paket sayısı
}

// Collector sistem metriklerini toplayan yapı
type Collector struct {
	// Ağ istatistikleri için önceki değerleri sakla
	prevNetStats map[string]net.IOCountersStat
}

// NewCollector yeni collector oluşturur
func NewCollector() *Collector {
	return &Collector{
		prevNetStats: make(map[string]net.IOCountersStat),
	}
}

// Start collector'ı başlatır
func (c *Collector) Start() error {
	log := logger.GetLogger()
	log.Info("Metrics collector başlatıldı")
	
	// İlk ağ istatistiklerini al
	netStats, err := net.IOCounters(true)
	if err == nil {
		for _, stat := range netStats {
			c.prevNetStats[stat.Name] = stat
		}
	}
	
	return nil
}

// Stop collector'ı durdurur
func (c *Collector) Stop() {
	log := logger.GetLogger()
	log.Info("Metrics collector durduruldu")
}

// CollectAll tüm sistem metriklerini toplar
func (c *Collector) CollectAll() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// CPU metrikleri
	cpuMetrics, err := c.collectCPU()
	if err != nil {
		return nil, fmt.Errorf("CPU metrikleri toplanamadı: %w", err)
	}
	metrics.CPU = *cpuMetrics

	// Bellek metrikleri
	memMetrics, err := c.collectMemory()
	if err != nil {
		return nil, fmt.Errorf("bellek metrikleri toplanamadı: %w", err)
	}
	metrics.Memory = *memMetrics

	// Disk metrikleri
	diskMetrics, err := c.collectDisk()
	if err != nil {
		return nil, fmt.Errorf("disk metrikleri toplanamadı: %w", err)
	}
	metrics.Disk = *diskMetrics

	// Ağ metrikleri
	netMetrics, err := c.collectNetwork()
	if err != nil {
		return nil, fmt.Errorf("ağ metrikleri toplanamadı: %w", err)
	}
	metrics.Network = *netMetrics

	return metrics, nil
}

// collectCPU CPU metriklerini toplar
func (c *Collector) collectCPU() (*CPUMetrics, error) {
	// CPU kullanım yüzdesi (1 saniye bekleme ile)
	usage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// CPU çekirdek sayısı
	count, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	var cpuUsage float64
	if len(usage) > 0 {
		cpuUsage = usage[0]
	}

	return &CPUMetrics{
		Usage: cpuUsage,
		Count: count,
	}, nil
}

// collectMemory bellek metriklerini toplar
func (c *Collector) collectMemory() (*MemMetrics, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemMetrics{
		Usage:     memInfo.UsedPercent,
		Total:     memInfo.Total,
		Available: memInfo.Available,
		Used:      memInfo.Used,
	}, nil
}

// collectDisk disk metriklerini toplar  
func (c *Collector) collectDisk() (*DiskMetrics, error) {
	// Ana disk partition'ını al (genellikle "/" veya "C:")
	var path string
	if partitions, err := disk.Partitions(false); err == nil && len(partitions) > 0 {
		path = partitions[0].Mountpoint
	} else {
		path = "/" // Linux default
	}

	diskInfo, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}

	return &DiskMetrics{
		Usage: diskInfo.UsedPercent,
		Total: diskInfo.Total,
		Used:  diskInfo.Used,
		Free:  diskInfo.Free,
	}, nil
}

// collectNetwork ağ metriklerini toplar
func (c *Collector) collectNetwork() (*NetMetrics, error) {
	netStats, err := net.IOCounters(false) // false = toplam tüm interface'ler
	if err != nil {
		return nil, err
	}

	if len(netStats) == 0 {
		return &NetMetrics{}, nil
	}

	// İlk (toplam) istatistikleri al
	stat := netStats[0]

	return &NetMetrics{
		BytesRecv:   stat.BytesRecv,
		BytesSent:   stat.BytesSent,
		PacketsRecv: stat.PacketsRecv,
		PacketsSent: stat.PacketsSent,
	}, nil
}