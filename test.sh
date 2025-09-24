#!/bin/bash
# Test runner for syswatch-daemon (Linux/Mac)

echo "🧪 Syswatch Daemon Test Suite"
echo "================================"

# Test all packages
echo ""
echo "📦 Tüm paketleri test ediliyor..."
if go test ./... -v; then
    echo ""
    echo "✅ Tüm testler başarılı!"
else
    echo ""
    echo "❌ Bazı testler başarısız!"
    exit 1
fi

# Run benchmarks
echo ""
echo "⚡ Benchmark testleri çalıştırılıyor..."
go test ./internal/metrics -bench=. -benchmem

# Test coverage
echo ""
echo "📊 Test coverage hesaplanıyor..."
go test ./... -cover

echo ""
echo "🎉 Test suite tamamlandı!"