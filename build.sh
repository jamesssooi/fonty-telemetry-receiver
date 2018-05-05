# Build `telemetry-listen` binary
GOOS=linux GOARCH=amd64 go build -o dist/telemetry-listen cmd/telemetry-listen/telemetry-listen.go
chmod +x dist/telemetry-listen

# Build `telemetry-receive` binary
GOOS=linux GOARCH=amd64 go build -o dist/telemetry-receive cmd/telemetry-receive/telemetry-receive.go
chmod +x dist/telemetry-receive