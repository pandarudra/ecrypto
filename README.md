# 🔐 ECRYPTO

![ecrypto](assets/image.png)

<div align="center">

**Military-Grade Folder Encryption Tool**

Encrypt entire folders into a single secure container using **XChaCha20-Poly1305** (AEAD) and **Argon2id** (KDF).  
Protects filenames, metadata, and contents with cutting-edge cryptography.

[![Release](https://img.shields.io/github/v/release/pandarudra/ecrypto)](https://github.com/pandarudra/ecrypto/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24+-00ADD8.svg)](https://golang.org)

[Download](#-installation) • [Quick Start](#-quick-start) • [Documentation](#-usage-examples) • [Security](#-security)

</div>

---

## ✨ Features

- 🔒 **Military-Grade Encryption**: XChaCha20-Poly1305 AEAD cipher (256-bit keys)
- 🔑 **Flexible Key Management**: Use passphrases or raw 32-byte key files
- 📦 **Single Secure Container**: Compress + encrypt entire folders into `.ecrypt` files
- 🎨 **Beautiful Interactive UI**: User-friendly menu system or powerful command-line interface
- 🛡️ **Secure by Default**: Argon2id KDF (256MB memory, 3 iterations) - winner of Password Hashing Competition
- ↶ **Undo Feature**: Easily decrypt and restore recently encrypted folders
- ⚡ **Fast & Lightweight**: Single binary, zero dependencies, cross-platform ready
- 🔍 **Tamper Detection**: Authentication tags prevent file modifications
- 📊 **Progress Tracking**: Real-time visual feedback during operations
- 🌐 **Cross-Platform**: Works on Windows, macOS, and Linux

---

## 📥 Installation

### Option 1: Download Pre-built Binary (Recommended)

#### Windows

```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/pandarudra/ecrypto/releases/latest/download/ecrypto-windows-amd64.exe" -OutFile "ecrypto.exe"

# Run it
.\ecrypto.exe
```

**Or manually:**

1. Go to **[Releases](https://github.com/pandarudra/ecrypto/releases/latest)** 📦
2. Download `ecrypto-windows-amd64.exe`
3. Rename to `ecrypto.exe` and run!

#### macOS

```bash
# Download the latest release
curl -L -o ecrypto "https://github.com/pandarudra/ecrypto/releases/latest/download/ecrypto-darwin-amd64"
chmod +x ecrypto

# Move to PATH (optional)
sudo mv ecrypto /usr/local/bin/
```

#### Linux

```bash
# Download the latest release
wget -O ecrypto "https://github.com/pandarudra/ecrypto/releases/latest/download/ecrypto-linux-amd64"
chmod +x ecrypto

# Move to PATH (optional)
sudo mv ecrypto /usr/local/bin/
```

### Option 2: Install via npm (Coming Soon)

```bash
npm install -g ecrypto
```

### Option 3: Build from Source

**Requirements**: Go 1.24+

```bash
git clone https://github.com/pandarudra/ecrypto.git
cd ecrypto
go build -o ecrypto
```

---

## 🚀 Quick Start

### 🎯 Interactive Mode (Perfect for Beginners!)

The easiest way to use ECRYPTO - just run the executable with no arguments:

```powershell
.\ecrypto.exe
```

You'll see an intuitive interactive menu:

```
██████ ▄█████ █████▄  ██  ██ █████▄ ██████ ▄████▄
██▄▄   ██     ██▄▄██▄  ▀██▀  ██▄▄█▀   ██   ██  ██
██▄▄▄▄ ▀█████ ██   ██   ██   ██       ██   ▀████▀

[1] [ENCRYPT]  Encrypt a Folder
[2] [DECRYPT]  Decrypt a File
[3] [KEYGEN]   Generate Encryption Key
[4] [INFO]     View Container Info
[5] [UNDO]     Undo Recent Operation
[6] [EXIT]     Quit Application
```

**Step-by-step walkthrough:**

1. Select `[1] ENCRYPT`
2. Enter your folder path (e.g., `C:\MyDocuments`)
3. Choose output location (e.g., `D:\backup.ecrypt`)
4. Enter a strong passphrase
5. Done! Your folder is now encrypted

### ⚡ Command-Line Mode (For Power Users)

#### Encrypt a folder with passphrase:

```powershell
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --pass "YourStrongPassphrase123!"
```

#### Decrypt a container:

```powershell
.\ecrypto.exe decrypt --in "backup.ecrypt" --out "restored" --pass "YourStrongPassphrase123!"
```

#### Generate a random 32-byte encryption key:

```powershell
.\ecrypto.exe keygen --out mykey.txt
```

#### Encrypt with a raw key file (maximum security):

```powershell
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --key-file mykey.txt
```

#### View container information (no decryption required):

```powershell
.\ecrypto.exe info --file backup.ecrypt
```

**Pro Tip:** Key files provide stronger security than passphrases. Store them in a password manager!

---

## 📖 Usage Examples

### 💼 Example 1: Backup Personal Documents

Perfect for creating encrypted backups of your important files:

```powershell
# Encrypt your documents folder
.\ecrypto.exe encrypt --in "C:\Users\YourName\Documents" --out "documents_backup.ecrypt" --pass "MySecurePassword123!"

# Later, restore to a new location
.\ecrypto.exe decrypt --in "documents_backup.ecrypt" --out "C:\Restored\Documents" --pass "MySecurePassword123!"
```

**Use Case:** Regular backups, cloud storage, disaster recovery

### 🔐 Example 2: Secure File Transfer

Share sensitive data securely by encrypting with a key file:

```powershell
# Step 1: Generate a random encryption key
.\ecrypto.exe keygen --out transfer_key.txt

# Step 2: Encrypt your sensitive data
.\ecrypto.exe encrypt --in "C:\SensitiveData" --out "transfer.ecrypt" --key-file transfer_key.txt

# Step 3: Send transfer.ecrypt via one channel (email/cloud)
#         Send transfer_key.txt via a DIFFERENT secure channel (Signal/WhatsApp)

# Recipient decrypts with the key:
.\ecrypto.exe decrypt --in "transfer.ecrypt" --out "received_data" --key-file transfer_key.txt
```

**Use Case:** Confidential file sharing, client deliverables, HIPAA/GDPR compliance

### 🖥️ Example 3: Optimize for Different Hardware

Adjust Argon2 parameters for slower machines or faster encryption:

```powershell
# Lower settings for older machines (128MB memory, 2 iterations)
.\ecrypto.exe encrypt --in "C:\MyFolder" --out "backup.ecrypt" --pass "password" --argon-m 131072 --argon-t 2

# Higher security for critical data (512MB memory, 5 iterations)
.\ecrypto.exe encrypt --in "C:\TopSecret" --out "critical.ecrypt" --pass "password" --argon-m 524288 --argon-t 5
```

**Use Case:** Performance tuning, high-security requirements, legacy systems

### 📂 Example 4: Batch Operations

Encrypt multiple folders programmatically:

```powershell
# Encrypt multiple project folders
$folders = @("C:\Project1", "C:\Project2", "C:\Project3")
foreach ($folder in $folders) {
    $name = Split-Path $folder -Leaf
    .\ecrypto.exe encrypt --in $folder --out "D:\Backups\$name.ecrypt" --pass "YourPassword"
}
```

**Use Case:** Automated backups, CI/CD pipelines, scheduled tasks

### ↶ Example 5: Undo & Restore

Accidentally encrypted something? Restore it with one click:

```powershell
# Encrypted a folder
.\ecrypto.exe encrypt --in "C:\MyFiles" --out "backup.ecrypt" --pass "password"

# Later: Need to undo the encryption
.\ecrypto.exe
→ [5] [UNDO] Undo Recent Operation
→ Select: C:\MyFiles | 450 files | 1.24 GB
→ Passphrase: password
→ ✓ Restored to: C:\MyFiles_restored
```

**Use Case:** Testing encryption settings, accidental encryption, backup verification

---

## 🔒 Security

### 🛡️ Cryptography Details

ECRYPTO uses industry-leading cryptographic standards:

- **Encryption Cipher**: XChaCha20-Poly1305

  - AEAD (Authenticated Encryption with Associated Data)
  - 256-bit keys for maximum security
  - ChaCha20 stream cipher + Poly1305 MAC
  - [RFC 8439](https://tools.ietf.org/html/rfc8439) compliant

- **Key Derivation**: Argon2id

  - Winner of the Password Hashing Competition (2015)
  - Resistant to GPU/ASIC attacks
  - Default: 256 MB memory, 3 iterations, 1 thread
  - [RFC 9106](https://datatracker.ietf.org/doc/html/rfc9106) compliant

- **Random Generation**: Cryptographically secure (Go's `crypto/rand`)

  - 24-byte XChaCha20 nonce (never reused)
  - 16-byte Argon2 salt
  - True randomness from OS entropy sources

- **Data Integrity**: Poly1305 authentication tag
  - 16-byte MAC prevents tampering
  - Header authenticated as AAD (Additional Authenticated Data)

### ✅ What ECRYPTO Protects Against

| Threat                    | Protection                                  |
| ------------------------- | ------------------------------------------- |
| Unauthorized file access  | ✅ Strong 256-bit encryption                |
| Filename/metadata leakage | ✅ Everything encrypted in container        |
| Brute-force attacks       | ✅ Argon2id makes cracking impractical      |
| File tampering            | ✅ Authentication tag detects modifications |
| Rainbow table attacks     | ✅ Unique salt per container                |
| Nonce reuse attacks       | ✅ Random 24-byte nonce per encryption      |

### ⚠️ Limitations & Threat Model

ECRYPTO **does NOT** protect against:

| Threat                       | Mitigation                                         |
| ---------------------------- | -------------------------------------------------- |
| Physical memory extraction   | Use full-disk encryption (BitLocker, FileVault)    |
| Weak passphrases             | Use 16+ character passphrases or key files         |
| Malware/keyloggers on system | Keep OS updated, use antivirus software            |
| Loss of encryption key       | **Always backup your keys/passphrases!**           |
| Side-channel attacks         | Not designed for hostile multi-tenant environments |

### 🔐 Security Best Practices

#### 1. **Strong Passphrases**

```
❌ Weak:   password123, qwerty, admin
✅ Strong: Correct-Horse-Battery-Staple-2026!
✅ Better: Use a key file generated with `keygen`
```

Recommendations:

- Minimum 16 characters
- Mix uppercase, lowercase, numbers, symbols
- Use a password manager (1Password, Bitwarden, KeePass)
- Never reuse passphrases across containers

#### 2. **Key File Security**

```powershell
# Generate key files for maximum security
.\ecrypto.exe keygen --out project_key.txt

# Store in password manager or encrypted USB drive
# Never store keys next to encrypted files
```

#### 3. **Backup Strategy**

```
Original Files → Encrypt → .ecrypt container
     ↓              ↓            ↓
  (Delete)      (Backup)    (Store securely)
                   ↓
         Cloud Storage / External Drive
```

**Critical:** Test decryption BEFORE deleting original files!

#### 4. **Secure Deletion of Originals**

Windows:

```powershell
# Use SDelete (Sysinternals) for secure deletion
sdelete -p 3 "C:\OriginalFolder"
```

Linux/macOS:

```bash
# Use shred for secure deletion
shred -vfz -n 3 /path/to/file
```

#### 5. **Container Storage**

- ✅ Cloud storage (Google Drive, Dropbox) - encrypted container is safe
- ✅ External drives with additional disk encryption
- ✅ Network shares with proper access controls
- ❌ Public file-sharing sites (risk of corruption)

### 🔍 Verification & Testing

```powershell
# 1. Check container integrity
.\ecrypto.exe info --file backup.ecrypt

# 2. Test decryption to temporary location
.\ecrypto.exe decrypt --in backup.ecrypt --out test_restore --pass "YourPassword"

# 3. Verify files are intact
# Compare checksums or spot-check files

# 4. Delete test restore
Remove-Item -Recurse test_restore
```

### 📊 Security Audit

ECRYPTO has been designed with security in mind:

- ✅ No hardcoded secrets or backdoors
- ✅ Open-source code available for review
- ✅ Standard cryptographic libraries (Go's `crypto/*`)
- ✅ Minimal dependencies reduce attack surface
- ✅ No network connections or telemetry

**Want to contribute to security?** Report vulnerabilities via [GitHub Issues](https://github.com/pandarudra/ecrypto/issues) (use "Security" label)

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

## 🐛 Troubleshooting

### Common Issues & Solutions

<details>
<summary><strong>❌ "Access is denied" error when encrypting</strong></summary>

**Problem:** Output path is a directory instead of a file.

**Solution:** Specify a file path ending with `.ecrypt`:

```powershell
# ❌ Wrong - this is a directory
.\ecrypto.exe encrypt --in folder --out D:\backup

# ✅ Correct - this is a file
.\ecrypto.exe encrypt --in folder --out D:\backup.ecrypt
```

</details>

<details>
<summary><strong>🔑 "Decryption failed: authentication tag mismatch"</strong></summary>

**Problem:** Wrong passphrase/key or corrupted file.

**Solutions:**

1. Double-check your passphrase (case-sensitive!)
2. Verify you're using the correct key file
3. Check file integrity with `info` command:
   ```powershell
   .\ecrypto.exe info --file backup.ecrypt
   ```
4. Ensure the file wasn't modified or corrupted during transfer
5. Try re-downloading the file if transferred over network

**Prevention:** Always test decryption immediately after encryption!

</details>

<details>
<summary><strong>📁 "File not found" in interactive mode</strong></summary>

**Problem:** Paths with spaces not recognized.

**Solution:** Use quotes around paths:

```
Enter folder path: "C:\My Documents\Folder"
```

Or use paths without spaces:

```
Enter folder path: C:\Users\John\Documents
```

</details>

<details>
<summary><strong>⚠️ "Out of memory" during encryption</strong></summary>

**Problem:** Argon2 memory settings too high for your system.

**Solution:** Reduce Argon2 memory parameter:

```powershell
# Default is 256MB - reduce to 128MB
.\ecrypto.exe encrypt --in folder --out backup.ecrypt --pass "password" --argon-m 131072

# For very low-memory systems (64MB)
.\ecrypto.exe encrypt --in folder --out backup.ecrypt --pass "password" --argon-m 65536
```

**Note:** Lower memory = faster cracking, so use strongest your system can handle.

</details>

<details>
<summary><strong>🐌 Encryption/decryption is very slow</strong></summary>

**Causes & Solutions:**

1. **Large Argon2 parameters:** Reduce `--argon-m` and `--argon-t`
2. **Large folders:** This is expected - compression + encryption takes time
3. **Slow storage:** Move to SSD instead of HDD
4. **Antivirus scanning:** Add exception for ecrypto or `.ecrypt` files

**Performance tips:**

- Use key files instead of passphrases (skips Argon2)
- Split large folders into smaller containers
- Disable real-time antivirus scanning during operations
</details>

<details>
<summary><strong>❓ "Invalid command" errors</strong></summary>

**Solution:** Check command syntax:

```powershell
# Correct syntax
.\ecrypto.exe <command> --flag value

# Common mistakes
.\ecrypto.exe --in folder encrypt  # ❌ Command must come first
.\ecrypto.exe encrypt -in folder   # ❌ Use -- for flags
```

Use `--help` for syntax help:

```powershell
.\ecrypto.exe --help
.\ecrypto.exe encrypt --help
```

</details>

<details>
<summary><strong>🔄 "Cannot restore folder structure"</strong></summary>

**Problem:** Permissions issues or invalid output path.

**Solutions:**

1. Run as Administrator (right-click → Run as administrator)
2. Ensure output directory exists and is writable
3. Check available disk space
4. Avoid network drives if experiencing issues
</details>

### Still Having Issues?

1. **Check the logs:** Look for error messages in terminal output
2. **Verify system requirements:** Go 1.24+ if building from source
3. **Test with small files first:** Isolate whether issue is size-related
4. **Report bugs:** [Open an issue](https://github.com/pandarudra/ecrypto/issues) with:
   - Operating system & version
   - ECRYPTO version (`.\ecrypto.exe --version`)
   - Command you ran (redact sensitive info)
   - Full error message

---

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

### Ways to Contribute

- 🐛 **Report bugs** - Found an issue? [Create a bug report](https://github.com/pandarudra/ecrypto/issues/new?labels=bug)
- 💡 **Suggest features** - Have an idea? [Open a feature request](https://github.com/pandarudra/ecrypto/issues/new?labels=enhancement)
- 📖 **Improve docs** - Fix typos, add examples, clarify instructions
- 🔐 **Security audits** - Review cryptographic implementation
- 💻 **Code contributions** - Fix bugs, implement features

### Development Setup

```bash
# Clone the repository
git clone https://github.com/pandarudra/ecrypto.git
cd ecrypto

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o ecrypto
```

### Contribution Workflow

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Make** your changes with clear commit messages
4. **Test** thoroughly (`go test ./...`)
5. **Commit** changes (`git commit -m 'Add amazing feature'`)
6. **Push** to your fork (`git push origin feature/amazing-feature`)
7. **Open** a Pull Request with detailed description

### Code Guidelines

- Follow Go best practices and `gofmt` formatting
- Add tests for new features
- Update documentation for user-facing changes
- Keep commits atomic and well-described
- Ensure backwards compatibility with `.ecrypt` format

### Security Contributions

Found a security vulnerability? Please:

1. **DO NOT** open a public issue
2. Email details to [security contact] or use GitHub Security Advisories
3. Include: description, steps to reproduce, potential impact
4. We'll respond within 48 hours

---

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

**TL;DR:** You can use, modify, and distribute this software freely. Just include the license notice.

---

## 🔗 Links & Resources

### Project Links

- 📦 **[Download Latest Release](https://github.com/pandarudra/ecrypto/releases/latest)**
- 🐛 **[Report a Bug](https://github.com/pandarudra/ecrypto/issues/new?labels=bug)**
- 💡 **[Request a Feature](https://github.com/pandarudra/ecrypto/issues/new?labels=enhancement)**
- 📖 **[Documentation & Wiki](https://github.com/pandarudra/ecrypto/wiki)**
- 💬 **[Discussions](https://github.com/pandarudra/ecrypto/discussions)**

### Documentation

- 📘 **[How It Works](docs/index.html)** - Visual explanation of encryption pipeline
- ↶ **[Undo Feature](docs/UNDO_FEATURE.md)** - Restore encrypted folders easily
- 🎨 **[UI Enhancements](docs/UI_ENHANCEMENTS.md)** - User-friendly terminal experience

### Related Resources

- 📚 [XChaCha20-Poly1305 Specification](https://tools.ietf.org/html/rfc8439)
- 🔐 [Argon2 Password Hashing](https://datatracker.ietf.org/doc/html/rfc9106)
- 🛡️ [OWASP Cryptographic Storage](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html)
- 🔑 [Password Manager Recommendations](https://www.privacyguides.org/passwords/)

---

## ⭐ Support This Project

If you find ECRYPTO useful, please consider:

- ⭐ **Starring the repository** - Helps others discover the project
- 🐛 **Reporting bugs** - Makes the tool better for everyone
- 📢 **Sharing** - Tell colleagues and friends
- 💻 **Contributing** - Submit PRs or improve documentation
- ☕ **Sponsoring** - Support ongoing development

### Star History

[![Star History Chart](https://api.star-history.com/svg?repos=pandarudra/ecrypto&type=Date)](https://star-history.com/#pandarudra/ecrypto&Date)

---

## 🙏 Acknowledgments

ECRYPTO is built on top of excellent open-source libraries:

- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)** - Cryptographic implementations
- **[github.com/spf13/cobra](https://github.com/spf13/cobra)** - CLI framework
- **[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling

Special thanks to the cryptography community for their research and implementations.

---

## 📊 Project Stats

![GitHub release (latest by date)](https://img.shields.io/github/v/release/pandarudra/ecrypto)
![GitHub all releases](https://img.shields.io/github/downloads/pandarudra/ecrypto/total)
![GitHub stars](https://img.shields.io/github/stars/pandarudra/ecrypto)
![GitHub issues](https://img.shields.io/github/issues/pandarudra/ecrypto)
![GitHub license](https://img.shields.io/github/license/pandarudra/ecrypto)

---

## 🎯 Roadmap

Planned features for future releases:

- [ ] npm package for easy installation
- [ ] GUI application (Windows/macOS/Linux)
- [ ] Compression algorithm selection (ZSTD, GZIP)
- [ ] Multiple key support (multi-party encryption)
- [ ] Hardware security module (HSM) integration
- [ ] Cloud storage integration (S3, Azure Blob)
- [ ] Automated backup scheduling
- [ ] Mobile app (iOS/Android)

Vote for features in [Discussions](https://github.com/pandarudra/ecrypto/discussions/categories/ideas)!

---

<div align="center">

**Made with ❤️ for Privacy and Security**

_"Your data, your control, your peace of mind."_

[⬆ Back to Top](#-ecrypto)

</div>
