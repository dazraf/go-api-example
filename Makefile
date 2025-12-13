.PHONY: docs run build clean

docs:
	swag init

run: docs
	go run main.go

build: docs
	go build -o api-server main.go

clean:
	rm -rf docs/ api-server

test:
	go test ./...