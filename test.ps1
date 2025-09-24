# Test runner for syswatch-daemon

Write-Host "🧪 Syswatch Daemon Test Suite" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan

# Test all packages
Write-Host "`n📦 Tüm paketleri test ediliyor..." -ForegroundColor Yellow
$testResult = go test ./... -v

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Tüm testler başarılı!" -ForegroundColor Green
} else {
    Write-Host "`n❌ Bazı testler başarısız!" -ForegroundColor Red
    exit 1
}

# Run benchmarks
Write-Host "`n⚡ Benchmark testleri çalıştırılıyor..." -ForegroundColor Yellow
go test ./internal/metrics -bench=. -benchmem

# Test coverage
Write-Host "`n📊 Test coverage hesaplanıyor..." -ForegroundColor Yellow
go test ./... -cover

Write-Host "`n🎉 Test suite tamamlandı!" -ForegroundColor Cyan