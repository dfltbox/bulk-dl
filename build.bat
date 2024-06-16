@echo off

echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o ./Linux/bulk-dl
if %errorlevel% neq 0 (
    echo Failed to build for Linux
    exit /b 1
)

echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o ./Windows/bulk-dl.exe
if %errorlevel% neq 0 (
    echo Failed to build for Windows
    exit /b 1
)

echo Builds completed successfully
