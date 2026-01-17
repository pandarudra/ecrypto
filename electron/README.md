# Ecrypto GUI

Beautiful desktop application for Ecrypto encryption tool built with Electron.js.

## Architecture

- **Frontend**: HTML, CSS, JavaScript (Electron Renderer)
- **Backend**: Go HTTP API server
- **IPC**: Electron Main Process bridges frontend and Go backend

## Development

### Prerequisites

- Node.js 18+ and npm
- Go 1.24+
- Built ecrypto binary in parent directory

### Install Dependencies

```bash
cd electron
npm install
```

### Run Development Mode

1. **Build Go backend** (from root directory):

```bash
go build -o ecrypto.exe .
```

2. **Start Electron app**:

```bash
cd electron
npm start
```

The app will automatically spawn the Go API server on port 8765.

## Building Distributable

### Windows

```bash
npm run build:win
```

Output: `dist/Ecrypto Setup.exe`

### macOS

```bash
npm run build:mac
```

Output: `dist/Ecrypto.dmg`

### Linux

```bash
npm run build:linux
```

Output: `dist/Ecrypto.AppImage` and `.deb`

## Features

- ✅ Drag-and-drop file/folder selection
- ✅ Native OS file dialogs
- ✅ AI-powered path suggestions
- ✅ Real-time password strength meter
- ✅ Operation history with undo
- ✅ Progress tracking
- ✅ Container metadata viewer
- ✅ Key generation with clipboard copy
- ✅ Beautiful modern UI

## Project Structure

```
electron/
├── main.js              # Electron main process (spawns Go server)
├── preload.js           # Security bridge (contextBridge API)
├── package.json         # Dependencies and build config
└── renderer/
    ├── index.html       # Main window UI
    ├── styles.css       # Modern styling
    └── app.js           # Frontend logic
```

## API Endpoints

The Go backend exposes these REST endpoints:

- `POST /encrypt` - Encrypt file/folder
- `POST /decrypt` - Decrypt container
- `POST /keygen` - Generate encryption key
- `POST /info` - Get container metadata
- `GET /history` - List operations
- `POST /undo` - Undo operation
- `POST /suggest-path` - AI path suggestions
- `POST /check-password` - Password strength check
- `GET /progress` - SSE progress updates
- `GET /health` - Server health check

## Security

- Context isolation enabled
- Node integration disabled
- Secure IPC via contextBridge
- No eval() or remote code execution
- Passphrases never logged

## License

Same as parent project (MIT)
