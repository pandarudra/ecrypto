# 🚀 Ecrypto GUI - Quick Start Guide

## What You Have Now

A complete **Electron.js desktop application** with:

- Beautiful modern UI (HTML/CSS/JS)
- Your existing Go backend as HTTP API
- All CLI features exposed in GUI
- AI-powered suggestions integrated
- Real-time progress tracking

## Project Structure

```
ecrypto/
├── main.go                    # Go CLI (updated with --serve flag)
├── gui/
│   └── server.go             # HTTP API server for Electron
├── electron/
│   ├── main.js               # Electron main process
│   ├── preload.js            # Security bridge
│   ├── package.json          # Node dependencies
│   └── renderer/
│       ├── index.html        # GUI interface
│       ├── styles.css        # Modern styling
│       └── app.js            # Frontend logic
```

## Getting Started

### 1. Install Node.js Dependencies

```bash
cd electron
npm install
```

This installs:

- Electron (desktop app framework)
- Axios (HTTP client)
- electron-builder (packaging tool)

### 2. Build Go Binary

From the root directory:

```bash
# Windows
go build -o ecrypto.exe .

# macOS/Linux
go build -o ecrypto .
```

### 3. Run Development Mode

```bash
cd electron
npm start
```

Or from root:

```bash
make gui-dev
```

This will:

1. Start Go HTTP API server on port 8765
2. Launch Electron app
3. Open DevTools for debugging

### 4. Test the GUI

Try these features:

- **Encrypt Tab**: Select folder/file → Enter password → Encrypt
- **Decrypt Tab**: Select .ecrypt file → Enter password → Decrypt
- **Key Gen**: Generate and copy encryption keys
- **Info**: View container metadata without decrypting
- **History**: See recent operations

## Building Distributable Apps

### Windows Installer (.exe)

```bash
make gui-build-win
```

Output: `electron/dist/Ecrypto Setup.exe`

### macOS App (.dmg)

```bash
make gui-build-mac
```

Output: `electron/dist/Ecrypto.dmg`

### Linux Package

```bash
make gui-build-linux
```

Output: `electron/dist/Ecrypto.AppImage` and `.deb`

## How It Works

### Architecture

```
┌─────────────────┐
│  Electron UI    │ (HTML/CSS/JS)
│  (Renderer)     │
└────────┬────────┘
         │ IPC
┌────────▼────────┐
│  Electron Main  │ (spawns Go server, handles dialogs)
│  (main.js)      │
└────────┬────────┘
         │ HTTP (localhost:8765)
┌────────▼────────┐
│  Go API Server  │ (gui/server.go)
│  --serve mode   │
└────────┬────────┘
         │
┌────────▼────────┐
│  Your Existing  │ (crypto/, cmd/, ai/)
│  Go Code        │
└─────────────────┘
```

### Communication Flow

1. **User clicks "Encrypt" button** in HTML
2. **JavaScript** calls `window.electronAPI.encrypt(data)`
3. **Preload.js** receives call via contextBridge (security)
4. **Main.js** forwards to Go HTTP API: `POST http://localhost:8765/encrypt`
5. **Go server** (`gui/server.go`) calls your existing `cmd.EncryptWithPassphrase()`
6. **Result** flows back through the chain
7. **UI updates** with success/error message

## CLI Still Works!

Your CLI is **completely unchanged**:

```bash
# Interactive TUI (no args)
./ecrypto

# Direct CLI commands
./ecrypto encrypt --in folder --out file.ecrypt --pass mypass
./ecrypto decrypt --in file.ecrypt --out folder --pass mypass
./ecrypto keygen --out key.key
./ecrypto info --file container.ecrypt

# New: API server mode for GUI
./ecrypto --serve --port=8765
```

## Customization

### Change UI Colors

Edit `electron/renderer/styles.css`:

```css
:root {
  --primary: #6366f1; /* Change to your color */
  --primary-hover: #4f46e5;
}
```

### Add New Features

1. **Go backend**: Add endpoint in `gui/server.go`
2. **Electron main**: Add IPC handler in `main.js`
3. **Preload**: Expose API in `preload.js`
4. **Frontend**: Call from `renderer/app.js`

### Change Window Size

Edit `electron/main.js`:

```javascript
mainWindow = new BrowserWindow({
  width: 1200, // Change width
  height: 800, // Change height
  // ...
});
```

## Debugging

### Go API Logs

The Go server logs to console when started with Electron:

```
[Go Server]: Server started on http://localhost:8765
[Go Server]: POST /encrypt - 200 OK
```

### Electron DevTools

Press `F12` or `Ctrl+Shift+I` to open Chrome DevTools for debugging JavaScript.

### Common Issues

**Problem**: "Go binary not found"

- **Solution**: Build Go binary first: `go build -o ecrypto.exe .`

**Problem**: "Port 8765 already in use"

- **Solution**: Change port in `main.js` and `gui/server.go`

**Problem**: "Electron not starting"

- **Solution**: Run `cd electron && npm install` first

## Next Steps

### Icons

Add app icons to `electron/assets/`:

- `icon.png` (512x512 for Linux)
- `icon.ico` (for Windows)
- `icon.icns` (for macOS)

### Code Signing

For production releases, sign your apps:

- **Windows**: Get code signing certificate
- **macOS**: Apple Developer account + notarization
- **Linux**: No signing required

### Auto-Updates

Add `electron-updater` for automatic app updates:

```bash
cd electron
npm install electron-updater
```

## Resources

- **Electron Docs**: https://www.electronjs.org/docs/latest/
- **electron-builder**: https://www.electron.build/
- **Your Go API**: http://localhost:8765/health (when running)

## Distributing

After building:

1. **Windows**: Share `.exe` installer
2. **macOS**: Share `.dmg` disk image
3. **Linux**: Share `.AppImage` or `.deb`

Users just double-click to install - no Go, no Node.js required!

---

**Enjoy your new GUI! 🎉**

Your CLI encryption tool now has a beautiful desktop interface while keeping all existing functionality intact.
