package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init logger'ı başlatır
func Init() {
	log = logrus.New()

	// Log formatını ayarla
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors: true,
	})

	// Log level'ı ayarla (development için debug)
	log.SetLevel(logrus.InfoLevel)

	// Output'u ayarla
	log.SetOutput(os.Stdout)

	log.Info("Logger başarıyla başlatıldı")
}

// GetLogger global logger instance'ını döndürür
func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}