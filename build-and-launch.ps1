$Env:GOOS="js"
$Env:GOARCH="wasm"
go build -o .\main.wasm
python3.10.exe .\server.py