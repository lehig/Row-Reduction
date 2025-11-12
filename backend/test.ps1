# Test script for Lambda function
# Clears cross-compilation environment variables before running tests

# Clear any cross-compilation settings
$env:GOOS = $null
$env:GOARCH = $null

# Run tests with lambda build tag
Write-Host "Running tests with lambda build tag..." -ForegroundColor Green
go test -tags lambda -v

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nAll tests passed!" -ForegroundColor Green
} else {
    Write-Host "`nTests failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

