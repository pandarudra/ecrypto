.PHONY: dev start build test clean gui gui-dev gui-build gui-build-win gui-build-mac gui-build-linux





Thumbs.db.DS_Store*.logdev:
	go run main.go

start:
	./ecrypto

build:
	go build -o ecrypto

test:
	go test ./...

clean:
	rm -f ecrypto
	rm -rf electron/dist
	rm -rf electron/node_modules

# GUI Development
gui-dev:
	@echo "Building Go binary for GUI..."
	go build -o ecrypto.exe .
	@echo "Starting Electron app..."
	cd electron && npm start

# GUI Installation
gui-install:
	@echo "Installing Electron dependencies..."
	cd electron && npm install

# Build GUI for all platforms
gui-build: gui-build-win gui-build-mac gui-build-linux

# Build GUI for Windows
gui-build-win:
	@echo "Building Go binary for Windows..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ecrypto.exe .
	@echo "Building Electron app for Windows..."
	cd electron && npm run build:win

# Build GUI for macOS
gui-build-mac:
	@echo "Building Go binary for macOS..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ecrypto .
	@echo "Building Electron app for macOS..."
	cd electron && npm run build:mac

# Build GUI for Linux
gui-build-linux:
	@echo "Building Go binary for Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ecrypto .
	@echo "Building Electron app for Linux..."
	cd electron && npm run build:linux