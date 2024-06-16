@echo off

echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o myprogram-linux
if %errorlevel% neq 0 (
    echo Failed to build for Linux
    exit /b 1
)

echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o myprogram-windows.exe
if %errorlevel% neq 0 (
    echo Failed to build for Windows
    exit /b 1
)

echo Building for macOS (Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -o myprogram-macos
if %errorlevel% neq 0 (
    echo Failed to build for macOS (Intel)
    exit /b 1
)

echo Building for macOS (ARM)...
set GOOS=darwin
set GOARCH=arm64
go build -o myprogram-macos-arm64
if %errorlevel% neq 0 (
    echo Failed to build for macOS (ARM)
    exit /b 1
)

echo Builds completed successfully
