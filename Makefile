assets/setaspect.wasm: main_wasm.go setaspect.go
	GOOS=js GOARCH=wasm  go build -o $@ main_wasm.go setaspect.go

setaspect: main.go
	go build -o $@ main.go

.PHONY: server
server:
	cd assets && python3 -m http.server
