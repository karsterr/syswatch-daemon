# Build script for syswatch-daemon

# Clean previous builds
Write-Host "Temizleniyor..."
Remove-Item -Path "syswatch-daemon.exe" -Force -ErrorAction SilentlyContinue

# Initialize Go modules
Write-Host "Go modulleri indiriliyor..."
go mod tidy

# Build the daemon
Write-Host "Daemon derleniyor..."
go build -o syswatch-daemon.exe ./cmd/syswatch-daemon

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build başarılı!"
    Write-Host "Daemon'u çalıştırmak için: .\syswatch-daemon.exe"
} else {
    Write-Host "❌ Build başarısız!"
    exit 1
}