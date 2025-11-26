# build.ps1 - Build, proto-gen, test, tidy for Windows PowerShell

# Variables
$ProtoDir = "proto"
$ServerDir = "server"
$ClientDir = "client"
$BinDir = "bin"
$ServerBin = "$BinDir\server.exe"
$ClientBin = "$BinDir\client.exe"

# Create bin directory
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

Write-Host "`n==> Running go mod tidy"
go mod tidy

Write-Host "`n==> Generating gRPC code from proto files"
protoc --go_out=. --go-grpc_out=. "$ProtoDir\*.proto"

Write-Host "`n==> Building server"
go build -o $ServerBin .\$ServerDir

Write-Host "`n==> Building client"
go build -o $ClientBin .\$ClientDir

Write-Host "`n==> Running tests"
go test -v ./...

Write-Host "`n==> Build complete!"