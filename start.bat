@echo off
go run main.go
set COLS=800
set LINES = 600
mode con: cols=%COLS% lines=%LINES%
go run main.go
pause