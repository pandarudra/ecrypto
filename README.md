# 🔐 ECRYPTO

**Military-Grade Folder Encryption Tool for Windows**

Encrypt entire folders into a single secure container using **XChaCha20-Poly1305** (AEAD) and **Argon2id** (KDF). Protects filenames, metadata, and contents with cutting-edge cryptography.

---

## ✨ Features

- 🔒 **Strong Encryption**: XChaCha20-Poly1305 AEAD cipher
- 🔑 **Flexible Key Management**: Passphrase or raw 32-byte keys
- 📦 **Single Container**: Compress + encrypt entire folders into `.ecrypt` files
- 🎨 **Beautiful CLI**: Interactive menu or traditional command-line interface
- 🛡️ **Secure by Default**: Argon2id KDF with 256MB memory, 3 iterations
- ⚡ **Fast & Lightweight**: Single binary, no dependencies

---

## 📥 Installation

### Download Pre-built Binary (Recommended)

1. Go to [Releases](https://github.com/pandarudra/ecrypto/releases)
2. Download `ecrypto.exe` for Windows
3. Run it!

### Build from Source

**Requirements**: Go 1.21+

```powershell
git clone https://github.com/pandarudra/ecrypto.git
cd ecrypto
go build -o ecrypto.exe
```

---

## 🚀 Quick Start

### Interactive Mode (Beginner-Friendly)

Just run the executable with no arguments:

```powershell
.\ecrypto.exe
```

You'll see a beautiful interactive menu:

```
██████ ▄█████ █████▄  ██  ██ █████▄ ██████ ▄████▄
██▄▄   ██     ██▄▄██▄  ▀██▀  ██▄▄█▀   ██   ██  ██
██▄▄▄▄ ▀█████ ██   ██   ██   ██       ██   ▀████▀

[1] [ENCRYPT]  Encrypt a Folder
[2] [DECRYPT]  Decrypt a File
[3] [KEYGEN]   Generate Encryption Key
[4] [INFO]     View Container Info
[5] [EXIT]     Quit Application
```

### Command-Line Mode (Advanced)

#### Encrypt a folder with passphrase:

```powershell
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --pass "strong-passphrase"
```

#### Decrypt a container:

```powershell
.\ecrypto.exe decrypt --in "backup.ecrypt" --out "restored" --pass "strong-passphrase"
```

#### Generate a random 32-byte key:

```powershell
.\ecrypto.exe keygen --out mykey.txt
```

#### Encrypt with raw key file:

```powershell
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --key-file mykey.txt
```

#### View container info (no decryption):

```powershell
.\ecrypto.exe info --file backup.ecrypt
```

---

## 📖 Usage Examples

### Example 1: Backup Personal Documents

```powershell
# Encrypt your documents folder
.\ecrypto.exe encrypt --in "C:\Users\YourName\Documents" --out "documents_backup.ecrypt" --pass "MySecurePassword123"

# Later, restore to a new location
.\ecrypto.exe decrypt --in "documents_backup.ecrypt" --out "C:\Restored\Documents" --pass "MySecurePassword123"
```

### Example 2: Secure File Transfer

```powershell
# Generate a random key
.\ecrypto.exe keygen --out transfer_key.txt

# Encrypt with the key
.\ecrypto.exe encrypt --in "C:\SensitiveData" --out "transfer.ecrypt" --key-file transfer_key.txt

# Send transfer.ecrypt + transfer_key.txt separately
# Recipient decrypts:
.\ecrypto.exe decrypt --in "transfer.ecrypt" --out "received_data" --key-file transfer_key.txt
```

### Example 3: Adjust Argon2 Settings (for slower machines)

```powershell
# Reduce memory usage to 128MB, 2 iterations
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --pass "pass" --argon-m 131072 --argon-t 2
```

---

## 🔒 Security

### Cryptography

- **Cipher**: XChaCha20-Poly1305 (AEAD, 256-bit key)
- **KDF**: Argon2id (winner of Password Hashing Competition)
  - Default: 256 MB memory, 3 iterations, 1 thread
- **Nonce**: 24-byte random nonce per container (never reused)
- **AAD**: Header authenticated to prevent tampering

### Threat Model

✅ **Protects against:**

- Unauthorized access to encrypted files
- Filename/metadata leakage
- Tampering detection (authentication tag)

❌ **Does NOT protect against:**

- Physical key extraction from memory (use disk encryption)
- Weak passphrases (use strong, unique passphrases)
- Malware on the system during encryption/decryption

### Best Practices

1. **Use strong passphrases**: Minimum 16 characters, mix of letters/numbers/symbols
2. **Keep key files safe**: Store in password managers (1Password, Bitwarden)
3. **Backup your keys**: Losing the key = permanent data loss
4. **Don't reuse passphrases**: Use unique passphrases per container
5. **Test decryption**: Always verify you can decrypt before deleting originals

---

## 🛠️ Command Reference

### `encrypt`

| Flag         | Description             | Default        |
| ------------ | ----------------------- | -------------- |
| `--in`       | Input folder path       | (required)     |
| `--out`      | Output .ecrypt file     | (required)     |
| `--pass`     | Passphrase (Argon2id)   | -              |
| `--key-file` | 32-byte Base64 key file | -              |
| `--argon-m`  | Argon2 memory (KiB)     | 262144 (256MB) |
| `--argon-t`  | Argon2 iterations       | 3              |
| `--argon-p`  | Argon2 parallelism      | 1              |

### `decrypt`

| Flag         | Description        | Default    |
| ------------ | ------------------ | ---------- |
| `--in`       | Input .ecrypt file | (required) |
| `--out`      | Output folder path | (required) |
| `--pass`     | Passphrase         | -          |
| `--key-file` | Key file           | -          |

### `keygen`

| Flag    | Description     | Default            |
| ------- | --------------- | ------------------ |
| `--out` | Output key file | (prints to stdout) |

### `info`

| Flag     | Description       | Default    |
| -------- | ----------------- | ---------- |
| `--file` | .ecrypt file path | (required) |

---

## 🏗️ Architecture

```
.ecrypt Container Format (v1):
┌────────────────────────────────────────┐
│ Header (59 bytes)                      │
│  - Magic: "ECRYPT01"                   │
│  - Version: 1                          │
│  - KDF: 0=raw, 1=Argon2id              │
│  - Argon2 params (m, t, p)             │
│  - Salt (16 bytes)                     │
│  - Nonce (24 bytes)                    │
├────────────────────────────────────────┤
│ Encrypted Data (XChaCha20-Poly1305)    │
│  - Compressed folder (ZIP)             │
│  - Authentication Tag (16 bytes)       │
└────────────────────────────────────────┘
```

**File Flow:**

```
Input Folder → ZIP Archive → Encrypt (XChaCha20) → .ecrypt Container
```

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE) file for details.

---

## 🐛 Troubleshooting

### "Access is denied" error when encrypting

**Problem**: Output path is a directory, not a file.

**Solution**: Specify a file path ending with `.ecrypt`:

```powershell
# ❌ Wrong
.\ecrypto.exe encrypt --in folder --out D:\backup

# ✅ Correct
.\ecrypto.exe encrypt --in folder --out D:\backup.ecrypt
```

### "Decryption failed: authentication tag mismatch"

**Problem**: Wrong passphrase/key or corrupted file.

**Solution**:

- Double-check your passphrase
- Verify file integrity (use `info` command)
- Ensure file wasn't modified

### "File not found" in interactive mode

**Problem**: Path with spaces not recognized.

**Solution**: Use quotes around paths:

```
Enter folder path: "C:\My Documents\Folder"
```

---

## 🔗 Links

- [Report Bug](https://github.com/pandarudra/ecrypto/issues)
- [Request Feature](https://github.com/pandarudra/ecrypto/issues)
- [Documentation](https://github.com/pandarudra/ecrypto/wiki)

---

## ⭐ Star History

If you find this useful, please star the repo!

---

**Made with ❤️ for privacy and security**
