#!/bin/bash
# Build script for syswatch-daemon (Linux/Mac)

echo "Temizleniyor..."
rm -f syswatch-daemon

echo "Go modülleri indiriliyor..."
go mod tidy

echo "Daemon derleniyor..."
go build -o syswatch-daemon ./cmd/syswatch-daemon

if [ $? -eq 0 ]; then
    echo "✅ Build başarılı!"
    echo "Daemon'u çalıştırmak için: sudo ./syswatch-daemon"
else
    echo "❌ Build başarısız!"
    exit 1
fi