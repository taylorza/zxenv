@echo off
go build .
if %errorlevel% neq 0 goto exit

copy /y zxenv.exe c:\nextdev

:exit