$Env:GOOS="js"
#$Env:GOOS="wasip1"
$Env:GOARCH="wasm"
cd src
go build -o ..\bin\main.wasm
cd ..
python3.10.exe .\server.py