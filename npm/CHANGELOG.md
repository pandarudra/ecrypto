# Changelog - ecrypto-cli

All notable changes to the npm package will be documented here.

## [1.0.7] - 2026-01-12

### 🎉 New Features

- **Interactive File Browser**: Navigate your file system with a visual interface
  - Browse drives (C:, D:, etc.) on Windows
  - Navigate folders with numbered menu
  - See file sizes before selecting
  - Quick access to common folders (Documents, Downloads, Pictures, Desktop)
- **Enhanced File Selection**: Choose between interactive browsing or direct path paste
- **Smart Path Detection**: Auto-detects whether input is a file or folder
- **Single File Encryption**: Now supports encrypting individual files (not just folders)

### ✨ Improvements

- Better UX with dual-mode input (browse or paste)
- Quote-aware path input (handles paths with spaces)
- Visual feedback with icons (📁, 📄, 💾, ⬆️)
- Cleaner error messages
- Enhanced decryption (auto-detects file vs folder encryption)

### 🔧 Technical

- Updated to Go 1.24
- Improved error handling
- Better cross-platform support
- Optimized binary size

---

## [1.0.6] - 2025-XX-XX

### ✨ Improvements

- Updated binary download mechanism
- Better Windows compatibility
- Enhanced postinstall script

---

## [1.0.5] - 2025-XX-XX

### 🎉 Initial npm Release

- Interactive menu-driven interface
- Folder encryption with XChaCha20-Poly1305
- Argon2id key derivation
- Passphrase or key file authentication
- Cross-platform support (Windows, macOS, Linux)
- Undo functionality
- Container info viewer
- Key generation tool

### 🔒 Security Features

- XChaCha20-Poly1305 AEAD cipher
- 256-bit keys
- Argon2id KDF (256MB memory, 3 iterations)
- Random nonces and salts
- Tamper detection

---

## Upcoming Features

- [ ] Bookmark/favorites for frequent paths
- [ ] Recent operations history
- [ ] Search/filter in file browser
- [ ] Batch operations (multi-select)
- [ ] Compression level options
- [ ] Password strength meter
- [ ] Hardware key support

---

For detailed release notes, see: https://github.com/pandarudra/ecrypto/releases
