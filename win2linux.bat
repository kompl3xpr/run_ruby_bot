set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o run_ruby_bot -ldflags="-s -w" main.go