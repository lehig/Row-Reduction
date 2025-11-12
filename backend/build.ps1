# Clean old files
if (Test-Path .\bootstrap) {
    Remove-Item .\bootstrap -Force
}
if (Test-Path .\rref-api.zip) {
    Remove-Item .\rref-api.zip -Force
}

# Optional: clear Go build cache
go clean -cache

# Save current environment variables (if any)
$originalGOOS = $env:GOOS
$originalGOARCH = $env:GOARCH

try {
    # Cross-compile for Lambda on Amazon Linux 2023 custom runtime
    $env:GOOS  = "linux"
    $env:GOARCH = "arm64"   # <-- make sure the Lambda is also set to ARM64 architecture

    # Build Lambda bootstrap binary with lambda build tag
    # This includes main.go (lambda handler) and types.go (shared types)
    go build -tags lambda -o bootstrap .

    # Zip it
    Compress-Archive -Path .\bootstrap -DestinationPath .\rref-api.zip -Force

    Write-Host "Build and packaging complete."
    Write-Host "Upload rref-api.zip to AWS Lambda"
    Write-Host "Handler: bootstrap"
    Write-Host "Runtime: Custom runtime on Amazon Linux 2023"
    Write-Host "Architecture: arm64"
} finally {
    # Restore original environment variables (or clear if they weren't set)
    if ($originalGOOS) {
        $env:GOOS = $originalGOOS
    } else {
        $env:GOOS = $null
    }
    if ($originalGOARCH) {
        $env:GOARCH = $originalGOARCH
    } else {
        $env:GOARCH = $null
    }
}
