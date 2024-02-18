$Env:GOOS="js"
#$Env:GOOS="wasip1"
$Env:GOARCH="wasm"
go build -o .\main.wasm
python3.10.exe .\server.py