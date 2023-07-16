build-wasm:
	GOOS=js GOARCH=wasm go build -o assets/bf.wasm .
run: build-wasm
	go run server.go
build: build-wasm
	go build server.go
