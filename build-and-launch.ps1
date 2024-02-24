$Env:GOOS="js"
$Env:GOARCH="wasm"
cd src
go build -o ..\bin\main.wasm
cd ..
python3.10.exe .\server.py