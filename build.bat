@echo off
chcp 65001
echo 📦 Start building...

echo 🔧   1. Set up environment

set GOOS=js
set GOARCH=wasm
set WASMOPT=D:\apps\binaryen-version_122\bin\wasm-opt.exe

echo 🔧   2. Build (tinygo)

:: tinygo build -o visual-lyric-core.wasm -target wasm
go build -o visual-lyric-core.wasm

echo 🚤   3. Copy file to D:\xiaowumin\projects\test\klok

copy visual-lyric-core.wasm D:\xiaowumin\projects\test\klok

echo ✅   4. Done