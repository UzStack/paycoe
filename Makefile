

run:
	@air --build.cmd "go build -o ./tmp/main ./cmd/main.go" --build.bin "./tmp/main"

build:
	go build -o ./bin/paycue ./cmd
