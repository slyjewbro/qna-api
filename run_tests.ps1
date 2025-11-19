Write-Host "=== Q&A API TEST RUNNER ===" -ForegroundColor Cyan

# Шаг 1: Останавливаем все контейнеры
Write-Host "`n1. Stopping all containers..." -ForegroundColor Yellow
docker compose down --remove-orphans 2>$null
Write-Host "   ✅ All containers stopped" -ForegroundColor Green

# Шаг 2: Запускаем тестовую базу данных
Write-Host "`n2. Starting test database..." -ForegroundColor Yellow
docker compose -f docker-compose.test.yml up test-db -d

# Шаг 3: Ждем пока база будет готова
Write-Host "`n3. Waiting for database to be ready..." -ForegroundColor Yellow
$attempts = 0
do {
    Start-Sleep -Seconds 2
    $attempts++
    $status = docker compose -f docker-compose.test.yml exec test-db pg_isready 2>$null
    Write-Host "   Attempt $attempts: $status" -ForegroundColor Gray
} while ($status -notlike "*accepting connections*" -and $attempts -lt 10)

if ($attempts -ge 10) {
    Write-Host "   ❌ Database failed to start" -ForegroundColor Red
    exit 1
}

Write-Host "   ✅ Database is ready" -ForegroundColor Green

# Шаг 4: Запускаем миграции для тестовой БД
Write-Host "`n4. Running migrations..." -ForegroundColor Yellow
docker compose -f docker-compose.test.yml run --rm migration

# Шаг 5: Запускаем тесты
Write-Host "`n5. Running tests..." -ForegroundColor Yellow
Write-Host "   Handler tests:" -ForegroundColor White
go test -v ./internal/handler

Write-Host "`n   Service tests:" -ForegroundColor White  
go test -v ./internal/service

Write-Host "`n   Repository tests:" -ForegroundColor White
go test -v ./internal/repository

Write-Host "`n   All tests:" -ForegroundColor White
go test -v ./...

# Шаг 6: Останавливаем тестовую базу
Write-Host "`n6. Cleaning up..." -ForegroundColor Yellow
docker compose -f docker-compose.test.yml down
Write-Host "   ✅ Test environment cleaned up" -ForegroundColor Green

Write-Host "`n=== TESTS COMPLETED ===" -ForegroundColor Cyan