@echo off

echo Windows x86-64 build
set GOOS=windows
set GOARCH=amd64
go build -o bin/windows-amd64/zxenv.exe .
powershell Compress-Archive -LiteralPath 'bin/windows-amd64/zxenv.exe' -DestinationPath "bin/windows-amd64.zip" -Force

echo Linux x86-64 Build
set GOOS=linux
set GOARCH=amd64
go build -o bin/linux-amd64/zxenv .
powershell Compress-Archive -LiteralPath 'bin/linux-amd64/zxenv' -DestinationPath "bin/linux-amd64.zip" -Force

echo Linux ARM Build
set GOOS=linux
set GOARCH=arm
go build -o bin/linux-arm/zxenv .
powershell Compress-Archive -LiteralPath 'bin/linux-arm/zxenv' -DestinationPath "bin/linux-arm.zip" -Force

echo Linux ARM64 Build
set GOOS=linux
set GOARCH=arm64
go build -o bin/linux-arm64/zxenv .
powershell Compress-Archive -LiteralPath 'bin/linux-arm64/zxenv' -DestinationPath "bin/linux-arm64.zip" -Force

echo MacOS x86-64 Build
set GOOS=darwin
set GOARCH=amd64
go build -o bin/mac-amd64/zxenv .
powershell Compress-Archive -LiteralPath 'bin/mac-m1/zxenv' -DestinationPath "bin/mac-m1.zip" -Force

echo MacOS ARM64 Build
set GOOS=darwin
set GOARCH=arm64
go build -o bin/mac-m1/zxenv .
powershell Compress-Archive -LiteralPath 'bin/mac-m1/zxenv' -DestinationPath "bin/mac-m1.zip" -Force

REM Reset
set GOOS=
set GOARCH=