# 🔐 ecrypto-cli

**Military-Grade Folder & File Encryption Tool**

Encrypt entire folders or individual files into secure containers using **XChaCha20-Poly1305** (AEAD) and **Argon2id** (KDF).

[![npm version](https://img.shields.io/npm/v/ecrypto-cli.svg)](https://www.npmjs.com/package/ecrypto-cli)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/pandarudra/ecrypto/blob/main/LICENSE)
[![Downloads](https://img.shields.io/npm/dm/ecrypto-cli.svg)](https://www.npmjs.com/package/ecrypto-cli)

---

## ✨ Features

- 🔒 **Military-Grade Encryption**: XChaCha20-Poly1305 AEAD cipher (256-bit keys)
- 📁 **Folder & File Support**: Encrypt entire folders or individual files
- 🎨 **Interactive File Browser**: Navigate your file system with arrow keys
- 🔑 **Flexible Authentication**: Passphrase or raw 32-byte key files
- 📦 **Single Secure Container**: Compressed + encrypted `.ecrypt` files
- 🛡️ **Secure by Default**: Argon2id KDF (256MB memory, 3 iterations)
- ↶ **Undo Feature**: Easily restore recently encrypted data
- ⚡ **Fast & Lightweight**: Single binary, zero runtime dependencies
- 🌐 **Cross-Platform**: Windows, macOS, and Linux

---

## 📥 Installation

### Global Installation (Recommended)

Install globally to use from anywhere:

```bash
npm install -g ecrypto-cli
```

Then use it:

```bash
ecrypto
```

### Local Installation

Install in your project:

```bash
npm install ecrypto-cli
```

Run with npx:

```bash
npx ecrypto
```

---

## 🚀 Quick Start

### Interactive Mode

Launch the interactive menu:

```bash
ecrypto
```

You'll see a beautiful interface:

```
██████ ▄█████ █████▄  ██  ██ █████▄ ██████ ▄████▄
██▄▄   ██     ██▄▄██▄  ▀██▀  ██▄▄█▀   ██   ██  ██
██▄▄▄▄ ▀█████ ██   ██   ██   ██       ██   ▀████▀

  XChaCha20-Poly1305 | Argon2id | Military-Grade Security

 Main Menu
───────────

   > [1] [ENCRYPT]  Encrypt a Folder/File
     [2] [DECRYPT]  Decrypt a Folder/File
     [3] [KEYGEN]   Generate Encryption Key
     [4] [INFO]     View Container Info
     [5] [UNDO]     Undo Recent Operation
     [6] [EXIT]     Quit Application
```

### Command-Line Mode

#### Encrypt a Folder

```bash
# Using a passphrase
ecrypto encrypt --in ./my-documents --out encrypted.ecrypt --pass

# Using a key file
ecrypto encrypt --in ./my-documents --out encrypted.ecrypt --key-file key.txt
```

#### Encrypt a Single File

```bash
# Using a passphrase
ecrypto encrypt --in photo.jpg --out photo.jpg.ecrypt --pass

# Using a key file
ecrypto encrypt --in document.pdf --out document.pdf.ecrypt --key-file key.txt
```

#### Decrypt

```bash
# Using a passphrase
ecrypto decrypt --in encrypted.ecrypt --out ./restored --pass

# Using a key file
ecrypto decrypt --in encrypted.ecrypt --out ./restored --key-file key.txt
```

#### Generate Encryption Key

```bash
ecrypto keygen --out mykey.txt
```

#### View Container Info

```bash
ecrypto info --in encrypted.ecrypt
```

---

## 🎯 Interactive Features

### 🗂️ Interactive File Browser

No need to type paths manually! Navigate with a visual file browser:

```
💾 Select Drive

  [1] 💾 C:
  [2] 💾 D:

📂 Current: C:\Users\YourName

  [1] ⬆️  ..
  [2] 📁 Documents
  [3] 📁 Downloads
  [4] 📁 Pictures
  [5] 📄 report.pdf (2.5 MB)
```

**Features:**

- ✅ Browse drives and folders interactively
- ✅ See file sizes before selecting
- ✅ Quick access to common folders
- ✅ Still supports pasting paths directly

### 📋 Smart Path Detection

Paste paths directly at any prompt:

```bash
› Select option: C:\Users\YourName\Documents\MyFolder
✓ Detected folder: MyFolder
```

Works with:

- Paths with or without quotes
- Spaces in names
- Forward or backward slashes

---

## 🔒 Security Features

### Encryption

- **Cipher**: XChaCha20-Poly1305 (AEAD)
  - 256-bit keys
  - 192-bit nonces (extended nonce space)
  - Authenticated encryption (tamper-proof)

### Key Derivation

- **KDF**: Argon2id (winner of Password Hashing Competition)
  - Memory: 256 MB
  - Iterations: 3
  - Parallelism: 1
  - Random 128-bit salt per encryption

### Container Format

```
[Header 59 bytes] + [Encrypted Data]
```

Header includes:

- Magic bytes: `ECRYPT01`
- Version: 1
- KDF type (0=raw key, 1=Argon2id)
- Salt (128 bits)
- Nonce (192 bits)
- Argon2id parameters

---

## 📚 Usage Examples

### Example 1: Encrypt Personal Documents

```bash
# Interactive mode
ecrypto

# Select [1] Encrypt
# Choose [1] Browse or paste: ~/Documents/Personal
# Enter passphrase (hidden input)
# Done! ✓
```

### Example 2: Batch Encrypt Multiple Files

```bash
# Encrypt each file
ecrypto encrypt --in file1.txt --out file1.ecrypt --pass
ecrypto encrypt --in file2.pdf --out file2.ecrypt --pass
ecrypto encrypt --in file3.jpg --out file3.ecrypt --pass
```

### Example 3: Generate and Use Key File

```bash
# Generate key
ecrypto keygen --out secret.key

# Encrypt with key
ecrypto encrypt --in ./sensitive-data --out backup.ecrypt --key-file secret.key

# Later, decrypt with same key
ecrypto decrypt --in backup.ecrypt --out ./restored --key-file secret.key
```

### Example 4: Quick Undo

```bash
# In interactive mode
# Select [5] Undo
# Shows recent operations
# Select operation to reverse
```

---

## 🔧 Advanced Options

### Encryption Options

```bash
ecrypto encrypt \
  --in <folder-or-file> \
  --out <output.ecrypt> \
  --pass                    # Use passphrase (prompts securely)
  --key-file <key.txt>      # Use raw 32-byte key file
  --argon-m <memory-KB>     # Argon2id memory (default: 262144 = 256MB)
  --argon-t <iterations>    # Argon2id iterations (default: 3)
  --argon-p <parallelism>   # Argon2id parallelism (default: 1)
```

### Decryption Options

```bash
ecrypto decrypt \
  --in <encrypted.ecrypt> \
  --out <output-folder> \
  --pass                    # Use passphrase
  --key-file <key.txt>      # Use raw key file
```

---

## 🆚 CLI vs Interactive Mode

| Feature      | CLI Mode       | Interactive Mode     |
| ------------ | -------------- | -------------------- |
| Speed        | ⚡ Fastest     | 🎨 User-friendly     |
| Automation   | ✅ Scriptable  | ❌ Manual            |
| File Browser | ❌ Type paths  | ✅ Visual navigation |
| Progress     | 📊 Text        | 📊 Visual + Icons    |
| Undo         | ❌ Manual      | ✅ Built-in          |
| Best For     | Scripts, CI/CD | Daily use, beginners |

---

## 🛠️ System Requirements

- **Node.js**: 12.x or higher (for npm installation only)
- **OS**: Windows, macOS, Linux
- **Disk Space**: ~10 MB
- **Memory**: Minimum 512 MB RAM (256 MB for Argon2id)

---

## 🔄 Updates

Keep ecrypto-cli up to date:

```bash
# Check current version
ecrypto --version

# Update to latest
npm update -g ecrypto-cli
```

---

## 📖 Documentation

Full documentation: [https://github.com/pandarudra/ecrypto](https://github.com/pandarudra/ecrypto)

### Help Commands

```bash
ecrypto --help              # Show all commands
ecrypto encrypt --help      # Encryption help
ecrypto decrypt --help      # Decryption help
ecrypto keygen --help       # Key generation help
```

---

## 🐛 Troubleshooting

### Permission Errors

**macOS/Linux:**

```bash
sudo npm install -g ecrypto-cli
```

### Binary Not Found

**After global install:**

```bash
# Add npm global bin to PATH
export PATH="$PATH:$(npm config get prefix)/bin"
```

### Old Version Cached

```bash
npm cache clean --force
npm install -g ecrypto-cli@latest
```

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](https://github.com/pandarudra/ecrypto/blob/main/CONTRIBUTING.md)

---

## 📄 License

MIT © [pandarudra](https://github.com/pandarudra)

---

## 🔗 Links

- **GitHub**: [https://github.com/pandarudra/ecrypto](https://github.com/pandarudra/ecrypto)
- **Issues**: [https://github.com/pandarudra/ecrypto/issues](https://github.com/pandarudra/ecrypto/issues)
- **npm**: [https://www.npmjs.com/package/ecrypto-cli](https://www.npmjs.com/package/ecrypto-cli)
- **Releases**: [https://github.com/pandarudra/ecrypto/releases](https://github.com/pandarudra/ecrypto/releases)

---

<div align="center">

**Made with ❤️ for Privacy & Security**

⭐ Star us on GitHub if you find this useful!

</div>
