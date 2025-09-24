package dashboard

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/karsterr/syswatch-daemon/internal/logger"
	"github.com/karsterr/syswatch-daemon/internal/metrics"
)

// Server web dashboard HTTP sunucusu
type Server struct {
	server     *http.Server
	router     *gin.Engine
	collector  *metrics.Collector
	port       int
}

// NewServer yeni dashboard server oluşturur
func NewServer(collector *metrics.Collector, port int) *Server {
	// Production modda gin loglarını kapat
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(gin.Recovery())
	
	return &Server{
		router:    router,
		collector: collector,
		port:      port,
	}
}

// Start dashboard sunucusunu başlatır
func (s *Server) Start() error {
	log := logger.GetLogger()
	
	// Routes'ları tanımla
	s.setupRoutes()
	
	// HTTP server'ı oluştur
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}
	
	log.Infof("Web dashboard başlatılıyor: http://localhost:%d", s.port)
	
	// Server'ı background'da başlat
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Dashboard server hatası: %v", err)
		}
	}()
	
	return nil
}

// Stop dashboard sunucusunu durdurur
func (s *Server) Stop(ctx context.Context) error {
	log := logger.GetLogger()
	
	if s.server == nil {
		return nil
	}
	
	log.Info("Web dashboard kapatılıyor...")
	
	// Graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("Dashboard shutdown hatası: %v", err)
		return err
	}
	
	log.Info("Web dashboard başarıyla kapatıldı")
	return nil
}

// setupRoutes HTTP endpoint'lerini tanımlar
func (s *Server) setupRoutes() {
	// Ana sayfa
	s.router.GET("/", s.handleHome)
	
	// API endpoints
	api := s.router.Group("/api")
	{
		api.GET("/metrics", s.handleMetrics)
		api.GET("/health", s.handleHealth)
	}
	
	// Static files (CSS, JS)
	s.router.Static("/static", "./web/static")
}

// handleHome ana sayfa handler'ı
func (s *Server) handleHome(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Syswatch Dashboard</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
        }
        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .metric-card {
            background: rgba(255, 255, 255, 0.1);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.2);
            border-radius: 15px;
            padding: 20px;
            text-align: center;
        }
        .metric-title {
            font-size: 1.2em;
            margin-bottom: 10px;
            opacity: 0.9;
        }
        .metric-value {
            font-size: 2.5em;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .metric-unit {
            font-size: 0.9em;
            opacity: 0.7;
        }
        .last-update {
            text-align: center;
            opacity: 0.6;
            font-size: 0.9em;
        }
        .status-indicator {
            display: inline-block;
            width: 10px;
            height: 10px;
            background: #4CAF50;
            border-radius: 50%;
            margin-right: 8px;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.5; }
            100% { opacity: 1; }
        }
        .cpu { color: #FF6B6B; }
        .memory { color: #4ECDC4; }
        .disk { color: #45B7D1; }
        .network { color: #FFA726; }
    </style>
    <script>
        function updateMetrics() {
            fetch('/api/metrics')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('cpu-value').textContent = data.cpu.usage.toFixed(1);
                    document.getElementById('memory-value').textContent = data.memory.usage.toFixed(1);
                    document.getElementById('disk-value').textContent = data.disk.usage.toFixed(1);
                    document.getElementById('network-recv').textContent = (data.network.bytes_recv / (1024*1024)).toFixed(2);
                    document.getElementById('network-sent').textContent = (data.network.bytes_sent / (1024*1024)).toFixed(2);
                    document.getElementById('last-update').textContent = 'Son güncelleme: ' + new Date().toLocaleTimeString();
                })
                .catch(error => {
                    console.error('Metrics yüklenemedi:', error);
                });
        }
        
        // Sayfa yüklendiğinde ve her 5 saniyede bir güncelle
        document.addEventListener('DOMContentLoaded', function() {
            updateMetrics();
            setInterval(updateMetrics, 5000);
        });
    </script>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🖥️ Syswatch Dashboard</h1>
            <p>Gerçek Zamanlı Sistem İzleme</p>
        </div>
        
        <div class="metrics-grid">
            <div class="metric-card cpu">
                <div class="metric-title">🔥 CPU Kullanımı</div>
                <div class="metric-value"><span id="cpu-value">--</span></div>
                <div class="metric-unit">%</div>
            </div>
            
            <div class="metric-card memory">
                <div class="metric-title">🧠 RAM Kullanımı</div>
                <div class="metric-value"><span id="memory-value">--</span></div>
                <div class="metric-unit">%</div>
            </div>
            
            <div class="metric-card disk">
                <div class="metric-title">💾 Disk Kullanımı</div>
                <div class="metric-value"><span id="disk-value">--</span></div>
                <div class="metric-unit">%</div>
            </div>
            
            <div class="metric-card network">
                <div class="metric-title">🌐 Ağ Trafiği</div>
                <div class="metric-value">
                    ⬇️ <span id="network-recv">--</span> MB/s<br>
                    ⬆️ <span id="network-sent">--</span> MB/s
                </div>
                <div class="metric-unit">İndirme / Yükleme</div>
            </div>
        </div>
        
        <div class="last-update">
            <span class="status-indicator"></span>
            <span id="last-update">Bağlanıyor...</span>
        </div>
    </div>
</body>
</html>`
	
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// handleMetrics API endpoint for metrics
func (s *Server) handleMetrics(c *gin.Context) {
	metrics, err := s.collector.CollectAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Metrikler alınamadı",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, metrics)
}

// handleHealth health check endpoint
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service": "syswatch-daemon",
		"version": "0.1.0",
	})
}