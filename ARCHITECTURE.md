# ECRYPTO - Architecture & Project Flow Documentation

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Core Libraries & Dependencies](#core-libraries--dependencies)
4. [Project Structure](#project-structure)
5. [Data Flow](#data-flow)
6. [Module Deep Dive](#module-deep-dive)
7. [Cryptographic Components](#cryptographic-components)
8. [User Interaction Flow](#user-interaction-flow)

---

## 📌 Project Overview

**ECRYPTO** is a military-grade folder encryption CLI tool written in Go that encrypts entire folders into secure `.ecrypt` containers. It combines:

- **XChaCha20-Poly1305 AEAD** - Modern authenticated encryption (256-bit keys)
- **Argon2id Key Derivation** - Password-based key derivation (winner of Password Hashing Competition)
- **ZIP Compression** - Before encryption for storage efficiency
- **Beautiful UI** - Interactive menu system with styled output using Charmbracelet's Lipgloss

The tool supports two operational modes:

1. **Interactive Mode** - User-friendly menu-driven interface (default when no arguments)
2. **CLI Mode** - Traditional command-line interface with flags

---

## 🏗️ Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         ECRYPTO CLI TOOL                         │
└─────────────────────────────────────────────────────────────────┘

                              ┌─────────────┐
                              │   main.go   │
                              └──────┬──────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                                 │
                    ▼                                 ▼
            ┌────────────────┐           ┌─────────────────────┐
            │  Interactive   │           │   CLI Mode (Cobra)  │
            │   UI (Lipgloss)│           └──────────┬──────────┘
            └────────┬───────┘                      │
                     │                              │
                     └──────────────┬────────────────┘
                                    │
                        ┌───────────┴───────────┐
                        │    cmd/root.go        │
                        │   (Cobra Command)     │
                        └───────────┬───────────┘
                                    │
        ┌───────────────┬───────────┼──────────────┬──────────────┐
        │               │           │              │              │
        ▼               ▼           ▼              ▼              ▼
   ┌─────────┐   ┌────────────┐ ┌────────┐  ┌────────┐  ┌─────────┐
   │Encrypt  │   │ Decrypt    │ │Keygen  │  │ Info   │  │  Undo   │
   │Command  │   │ Command    │ │Command │  │Command │  │Command  │
   └────┬────┘   └─────┬──────┘ └───┬────┘  └──┬─────┘  └────┬────┘
        │              │            │         │           │
        └──────────────┼────────────┼─────────┴───────────┘
                       │            │
        ┌──────────────┴────────────┴──────────────┐
        │                                          │
        ▼                                          ▼
   ┌──────────────┐                        ┌────────────────┐
   │ crypto/      │                        │  archive/      │
   │ ├─ cipher.go │◄──────┐               │  ├─ zip.go      │
   │ ├─ derive.go │       │               │  └─ manifest.go │
   │ └─ types.go  │       │               └────────┬────────┘
   └──────────────┘       │                        │
                          │                        │
                    ┌─────┴──────────┐             │
                    │                │             │
                    ▼                ▼             ▼
            ┌────────────────┐ ┌──────────────────────┐
            │  XChaCha20-    │ │  File Compression    │
            │  Poly1305 AEAD │ │  & Archiving         │
            └────────────────┘ └──────────────────────┘
                    │                    │
                    └────────┬───────────┘
                             │
                    ┌────────▼────────┐
                    │  Output: .ecrypt │
                    │  Container File  │
                    └──────────────────┘

┌──────────────────────────────────────────────────────────────┐
│                  External Libraries (Charmbracelet)          │
│                                                              │
│  lipgloss - Styled terminal output & UI components           │
│  (colors, borders, padding, layouts)                         │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔧 Core Libraries & Dependencies

### Direct Dependencies

| Library                             | Version | Purpose                                       |
| ----------------------------------- | ------- | --------------------------------------------- |
| `golang.org/x/crypto`               | v0.46.0 | XChaCha20-Poly1305 AEAD cipher & Argon2id KDF |
| `github.com/spf13/cobra`            | v1.10.2 | CLI framework for command structure           |
| `github.com/charmbracelet/lipgloss` | v1.1.0  | Styled terminal UI rendering                  |

### Standard Library (Go Built-in)

| Module                 | Functions Used                                 |
| ---------------------- | ---------------------------------------------- |
| `crypto/rand`          | `rand.Read()` - Random nonce & salt generation |
| `encoding/binary`      | Binary serialization for header encoding       |
| `encoding/base64`      | Base64URL encoding for key files               |
| `archive/zip`          | ZIP compression for folder archiving           |
| `io`, `os`, `filepath` | File system operations                         |
| `bufio`                | User input reading                             |
| `encoding/json`        | Operation history serialization                |

---

## 📂 Project Structure

```
ecrypto/
├── main.go                 # Entry point - route to interactive UI or CLI
├── go.mod                  # Module dependencies
├── go.sum                  # Dependency checksums
│
├── cmd/                    # CLI Commands (Cobra Framework)
│   ├── root.go            # Root command configuration
│   ├── encrypt.go         # Encrypt folder command
│   ├── decrypt.go         # Decrypt container command
│   ├── keygen.go          # Key generation command
│   ├── info.go            # Container info command
│   ├── interactive_wrappers.go  # Wrapper functions for interactive mode
│   └── undo.go            # Undo recent operation command
│
├── crypto/                # Cryptographic Core
│   ├── types.go          # HeaderV1 struct & serialization
│   ├── cipher.go         # XChaCha20-Poly1305 AEAD encrypt/decrypt
│   └── derive.go         # Argon2id key derivation & key file reading
│
├── archive/               # File Compression & Archiving
│   ├── zip.go            # Folder to ZIP conversion with progress
│   └── manifest.go       # Archive metadata handling
│
├── ui/                    # User Interface & Styling
│   ├── interactive.go    # User input prompts & file selection
│   ├── menu.go           # Main menu & operation workflows
│   ├── history.go        # Operation history tracking
│   ├── styles.go         # Lipgloss styling definitions
│   └── progress.go       # Progress visualization
│
├── assets/                # Static assets
│   └── image.png         # Project logo/image
│
└── docs/                  # Documentation
    ├── index.html        # Web documentation
    └── style.css         # Documentation styling
```

---

## 🔄 Data Flow

### Encryption Flow

```
┌────────────────────┐
│  User Input        │
│  • Folder Path     │
│  • Output File     │
│  • Passphrase/Key  │
└─────────┬──────────┘
          │
          ▼
┌────────────────────────────────┐
│ 1. Key Generation              │
│    • If passphrase:            │
│      - Generate random salt    │
│      - Argon2id(pass, salt)    │
│      - Derive 32-byte key      │
│    • If key file:              │
│      - Read & decode Base64    │
└─────────┬──────────────────────┘
          │
          ▼
┌────────────────────────────────┐
│ 2. Header Preparation          │
│    • Magic: "ECRYPT01"         │
│    • Version: 1                │
│    • Salt (16 bytes)           │
│    • Generate random nonce     │  (24 bytes for XChaCha20)
│    • Store KDF parameters      │
└─────────┬──────────────────────┘
          │
          ▼
┌────────────────────────────────┐
│ 3. Archive Creation            │
│    • Walk folder tree          │
│    • Compress files to ZIP     │
│    • Track file metadata       │
└─────────┬──────────────────────┘
          │
          ▼
┌────────────────────────────────┐
│ 4. Encryption                  │
│    • Serialize header to bytes │
│    • Use header as AAD         │  (Additional Authenticated Data)
│    • XChaCha20-Poly1305        │
│    • Encrypt ZIP with key      │
└─────────┬──────────────────────┘
          │
          ▼
┌────────────────────────────────┐
│ 5. Output                      │
│    • Write header to file      │
│    • Write ciphertext + tag    │
│    • Atomic rename (.tmp)      │
│    • Create .ecrypt container  │
└────────────────────────────────┘
```

### Decryption Flow

```
┌──────────────────────────┐
│  User Input              │
│  • Container File        │
│  • Output Folder         │
│  • Passphrase/Key File   │
└──────────┬───────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ 1. Container Reading             │
│    • Read entire .ecrypt file    │
│    • Parse header (59 bytes)     │
│    • Extract ciphertext + tag    │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ 2. Key Derivation                │
│    • Read KDF method from header │
│    • If Argon2id:                │
│      - Extract salt from header  │
│      - Derive key from passphrase│
│    • If raw key:                 │
│      - Read key file             │
│      - Validate 32-byte length   │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ 3. Decryption                    │
│    • Get nonce from header       │
│    • Use header as AAD           │
│    • XChaCha20-Poly1305 decrypt  │
│    • Verify authentication tag   │
│    • Extract ZIP plaintext       │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ 4. Archive Extraction            │
│    • Create output directory     │
│    • Extract ZIP to folder       │
│    • Restore file metadata       │
│    • Verify file integrity       │
└──────────┬───────────────────────┘
           │
           ▼
┌──────────────────────────────────┐
│ 5. Completion                    │
│    • Report success              │
│    • Show output location        │
│    • Record in history           │
└──────────────────────────────────┘
```

---

## 🔬 Module Deep Dive

### 1. **crypto/ Package** - Cryptographic Core

#### `types.go` - Container Format Definition

```go
type HeaderV1 struct {
    Magic   [8]byte   // "ECRYPT01" - file identifier
    Version uint8     // 1 - format version for future compatibility
    KDF     uint8     // 0=raw key, 1=Argon2id
    ArgonM  uint32    // Argon2 memory parameter (KiB) - 256MB = 262144
    ArgonT  uint32    // Argon2 time cost - iterations (3)
    ArgonP  uint8     // Argon2 parallelism factor (1)
    Salt    [16]byte  // Random salt for KDF
    Nonce   [24]byte  // Random nonce for XChaCha20
}
// Total: 59 bytes
```

**Functions:**

- `Encode()` - Serialize header to binary format
- `DecodeHeaderV1()` - Deserialize header from binary
- `HeaderSize()` - Returns 59 bytes constant

#### `cipher.go` - AEAD Encryption/Decryption

**Functions:**

- `EncryptAEAD(key, plaintext, aad, nonce)` - XChaCha20-Poly1305 encryption

  - Input: 32-byte key, plaintext bytes, AAD (header), 24-byte nonce
  - Output: Ciphertext with 16-byte authentication tag appended
  - Error on invalid key/nonce length

- `DecryptAEAD(key, ciphertext, aad, nonce)` - XChaCha20-Poly1305 decryption
  - Input: 32-byte key, ciphertext+tag, AAD (header), 24-byte nonce
  - Output: Plaintext bytes
  - Error if authentication fails (tampering detected)

#### `derive.go` - Key Management

**Functions:**

- `DeriveKeyArgon2id(pass, salt, m, t, p)` - Argon2id KDF

  - Parameters:
    - `pass` - User passphrase (string)
    - `salt` - 16-byte random salt
    - `m` - Memory cost (256 × 1024 = 256 MB)
    - `t` - Time cost (3 iterations)
    - `p` - Parallelism (1 thread)
  - Output: 32-byte key

- `ReadKeyFromFile(path)` - Load key from Base64URL-encoded file

  - Tries Base64URL first, then standard Base64
  - Validates 32-byte length
  - Returns error if format invalid

- `KeySize()` - Returns 32 (XChaCha20-Poly1305 key size)

---

### 2. **archive/ Package** - Compression & Archiving

#### `zip.go` - Folder Compression

**Functions:**

- `ZipFolder(root)` - Compress folder to ZIP bytes

  - Walks entire directory tree
  - Compresses with zip.Deflate method
  - Preserves modification times
  - Returns []byte containing ZIP data

- `ZipFolderWithProgress(root, callback)` - Compression with progress reporting
  - `callback ProgressCallback` - Called for each file
  - Enables UI progress updates
  - Same compression as ZipFolder

**Process:**

1. Walk directory tree recursively
2. Skip directories (zip only files)
3. Create relative paths (maintain structure)
4. Use Deflate compression method
5. Preserve file modification times
6. Report progress for each file

#### `manifest.go` - Archive Metadata

Handles ZIP manifest information for displaying archive contents.

---

### 3. **cmd/ Package** - CLI Commands (Cobra Framework)

#### `root.go` - Command Root

```go
rootCmd = &cobra.Command{
    Use: "ecrypto",
    Short: "Encrypt and decrypt folders into secure .ecrypt containers"
}
```

All subcommands registered via `init()` functions.

#### `encrypt.go` - Encryption Command

**Flags:**

- `--in` - Input folder path (required)
- `--out` - Output .ecrypt file (required)
- `--pass` - Passphrase (Argon2id mode)
- `--key-file` - 32-byte key file (raw key mode)

**Argon2id Parameters (hardcoded):**

- Memory: 256 MB (262144 KiB)
- Time: 3 iterations
- Parallelism: 1

**Workflow:**

1. Validate inputs
2. Initialize HeaderV1
3. Generate random salt & nonce
4. Derive key from passphrase OR load from file
5. Compress folder to ZIP
6. Encrypt ZIP with XChaCha20-Poly1305
7. Write header + ciphertext to .ecrypt file
8. Atomic rename (.tmp → final)

#### `decrypt.go` - Decryption Command

**Flags:**

- `--in` - Input .ecrypt file (required)
- `--out` - Output folder (required)
- `--pass` - Passphrase (for Argon2id-derived keys)
- `--key-file` - 32-byte key file (for raw keys)

**Workflow:**

1. Validate inputs
2. Read .ecrypt container
3. Parse HeaderV1
4. Detect KDF method from header
5. Derive key (passphrase) OR load key file
6. Extract ciphertext after header
7. Decrypt with XChaCha20-Poly1305
8. Extract ZIP to output folder
9. Display success message

#### `keygen.go` - Key Generation Command

**Flags:**

- `--out` - Output key file (optional)

**Workflow:**

1. Generate 32 random bytes
2. Encode to Base64URL format
3. Print to stdout
4. Optionally save to file with 0o600 permissions

#### `info.go` - Container Inspection Command

**Flags:**

- `--in` - Input .ecrypt file

**Displays:**

- Magic & version
- KDF method (Argon2id or raw key)
- Argon2 parameters
- Container size
- Encrypted data size

#### `interactive_wrappers.go` - Interactive Mode Wrappers

Wraps CLI commands for interactive menu:

- `encryptInteractive()`
- `decryptInteractive()`
- `keygenInteractive()`
- `infoInteractive()`
- `undoInteractive()`

---

### 4. **ui/ Package** - User Interface

#### `styles.go` - Lipgloss Styling

Defines color scheme and styled components using Charmbracelet Lipgloss:

- Color constants (Primary, Secondary, Success, Error, etc.)
- Predefined styles (HeaderStyle, MenuItemStyle, SelectedItemStyle, etc.)
- Border styles and padding

#### `interactive.go` - Input/Selection

**Functions:**

- `PromptUser(label, defaultVal)` - Display prompt & read input

  - Shows label with optional default value
  - Returns trimmed user input or default

- `SelectOption(title, options)` - Menu selection

  - Displays title with border
  - Lists numbered options
  - Returns selected index (0-based)

- `SelectFile(title)` - File browser

  - Prompts for path
  - Validates file existence
  - Retries up to 3 times
  - Returns valid path

- `SelectFolder(title)` - Folder browser

  - Similar to SelectFile
  - Validates directory existence

- `Pause()` - Wait for user keypress

#### `menu.go` - Main Menu & Workflows

**Main Menu Options:**

1. ENCRYPT - Encrypt a folder
2. DECRYPT - Decrypt a file
3. KEYGEN - Generate encryption key
4. INFO - View container info
5. UNDO - Undo recent operation
6. EXIT - Quit

**Interactive Workflows:**

**encryptInteractive():**

1. Select folder to encrypt
2. Choose key method (passphrase or key file)
3. Select output location
4. Show confirmation & file stats
5. Execute encryption with progress
6. Display success & save to history

**decryptInteractive():**

1. Select .ecrypt file
2. Choose key method (must match encryption)
3. Select output folder
4. Execute decryption
5. Display success

**keygenInteractive():**

1. Generate random key
2. Display key in Base64URL format
3. Option to save to file

**infoInteractive():**

1. Select .ecrypt container
2. Parse & display header info
3. Show KDF parameters
4. Display file size details

**undoInteractive():**

1. Show recent operations (last 5)
2. Allow user to select one
3. Prompt to delete encrypted/decrypted output
4. Confirm deletion

#### `history.go` - Operation Tracking

**Operation Struct:**

```go
type Operation struct {
    ID         string    // Timestamp-based ID
    Type       string    // "encrypt" or "decrypt"
    SourcePath string    // Original folder/container path
    OutputPath string    // Encrypted/decrypted output path
    Timestamp  time.Time // When operation occurred
    Size       int64     // Output file size
    FileCount  int       // Number of files encrypted
    KeyMethod  string    // "passphrase" or "keyfile"
    KeyPath    string    // Path to key file (if applicable)
    Status     string    // "success" or "failed"
    Error      string    // Error message if failed
}
```

**Functions:**

- `NewOperationHistory()` - Create history tracker

  - Stores in `~/.ecrypto/operations.json`
  - Auto-creates directory with 0o700 permissions

- `Load()` - Read history from disk

  - JSON unmarshaling
  - Creates empty if file doesn't exist

- `Save()` - Write history to disk

  - JSON marshaling with indentation
  - Secure file permissions (0o600)

- `AddOperation(op)` - Record new operation

  - Auto-generates ID if empty
  - Auto-sets timestamp
  - Keeps only last 50 operations

- `GetRecentOperations(count)` - Retrieve N most recent

#### `progress.go` - Progress Visualization

Progress bar display during compression/encryption/decryption operations.

---

## 🔐 Cryptographic Components

### XChaCha20-Poly1305 (AEAD)

**Purpose:** Authenticated Encryption with Associated Data

**Key Characteristics:**

- **Key Size:** 32 bytes (256-bit)
- **Nonce Size:** 24 bytes (192-bit) - larger than ChaCha20
- **Authentication Tag:** 16 bytes
- **AAD Support:** Encrypts + authenticates header (prevents tampering)
- **Stream Cipher:** Enables encryption without size limitations

**Why XChaCha20?**

- Extended nonce (24 bytes vs 12 for Poly1305) = more random nonces safely
- Modern, NIST-approved successor to AES-GCM
- Resistant to side-channel attacks
- No key expansion needed (direct use)

**Security Properties:**

- IND-CPA: Indistinguishable from random
- INT-CTXT: Integrity (authentication tags prevent forgery)
- Any tampering with header or ciphertext detected

### Argon2id (KDF)

**Purpose:** Passphrase-Based Key Derivation

**Parameters (Hardcoded):**

```
Memory:       256 MB (262144 KiB)
Time Cost:    3 iterations
Parallelism:  1 thread
Output:       32 bytes (for XChaCha20-Poly1305)
Salt:         16 bytes (random per encryption)
```

**Why Argon2id?**

- **Winner** of Password Hashing Competition (2015)
- **Memory-hard:** Resists GPU/ASIC brute-force attacks
- **Time-hard:** Multiple passes prevent rapid guessing
- **Side-channel resistant:** Designed for security
- **Versioning:** Compatible with cryptanalysis updates

**Attack Resistance:**

- GPU Attack: 256 MB memory requirement slows brute-force 1000x
- Dictionary Attack: Memory requirement = cost multiplier
- Rainbow Tables: Unique salt per encryption = no precomputation
- Timing Attack: Constant-time operations

---

## 🖥️ User Interaction Flow

### Entry Point Decision

```
main.go
│
├─ --version or -v flag?
│  └─ Print version → Exit
│
├─ No arguments?
│  └─ Launch ui.RunInteractiveMenu()
│
└─ Has arguments?
   └─ cmd.Execute() (Cobra CLI)
```

### Interactive Mode Flow

```
RunInteractiveMenu()
    ↓
PrintMenu() (Display banner)
    ↓
┌─────────────────────────┐
│  Main Menu Displayed    │
│  (6 options)            │
└────────┬────────────────┘
         │
    User Selection
         │
    ┌────┴────┬────────┬──────┬──────┬──────┐
    ▼         ▼        ▼      ▼      ▼      ▼
 Encrypt  Decrypt  Keygen  Info   Undo   Exit
    │         │        │      │      │      │
    └─────────┼────────┴──────┴──────┘      │
              │                             │
         Execute Flow                   Return/Exit
              │
          ┌───▼──────────┐
          │ Success?     │
          │   Yes → Save │
          │   No  → Show │
          │   Error      │
          └────┬─────────┘
               │
            Pause()
               │
            Back to Menu

```

### CLI Mode Flow (Cobra)

```
Command Line Input
    ↓
Cobra Parser
    ↓
Identify Command (encrypt/decrypt/keygen/info/undo)
    ↓
Parse Flags
    ↓
Validate Arguments
    ↓
Execute RunE Function
    ↓
Handle Errors & Exit Code
```

---

## 📊 Key Design Decisions

### 1. **Header as Additional Authenticated Data (AAD)**

- Header is authenticated but **not** encrypted
- Allows reading metadata (version, KDF type) without decryption key
- Prevents tampering with header values
- Format version enables future compatibility

### 2. **Atomic File Operations**

- Encrypt to `.tmp` file first
- Only rename when complete
- Prevents partial/corrupted output

### 3. **Passphrase vs Raw Key**

- Passphrase: User-friendly, Argon2id protects against brute-force
- Raw Key: Maximum security, no memory cost, 32 random bytes
- Selection stored in header (KDF field)

### 4. **Argon2id Parameters**

- 256 MB memory: Expensive for attackers, reasonable for users (~1-2 seconds)
- 3 iterations: Multiple passes increase time cost
- 1 parallelism: Conservative (single-threaded) for compatibility

### 5. **ZIP + Encrypt (Compression Before Encryption)**

- ZIP compression reduces data size
- Encryption obscures patterns compression would reveal
- Standard ZIP format (no custom compression)

### 6. **Operation History**

- Tracks all operations locally
- Enables "Undo" feature (delete recent outputs)
- JSON format for readability & portability
- Only last 50 operations kept

### 7. **UI Library Choice (Lipgloss)**

- Cross-platform (Windows, macOS, Linux)
- No external binary dependencies
- Styled terminal output (colors, borders, spacing)
- Part of Charmbracelet ecosystem

---

## 🔄 Function Reference by Use Case

### Encrypt a Folder

**CLI:**

```
ecrypto encrypt --in /path/to/folder --out output.ecrypt --pass "my passphrase"
```

**Functions Involved:**

1. `encryptCmd.RunE()` - Parse & validate
2. `HeaderV1{}` - Initialize container format
3. `crypto.DeriveKeyArgon2id()` - Derive 32-byte key
4. `crypto.EncryptAEAD()` - Encrypt ZIP with XChaCha20-Poly1305
5. `archive.ZipFolder()` - Compress folder
6. File I/O - Write .ecrypt file

### Decrypt a Container

**CLI:**

```
ecrypto decrypt --in output.ecrypt --out /path/to/extract --pass "my passphrase"
```

**Functions Involved:**

1. `decryptCmd.RunE()` - Parse & validate
2. `crypto.DecodeHeaderV1()` - Read container format
3. `crypto.DeriveKeyArgon2id()` - Derive key from passphrase
4. `crypto.DecryptAEAD()` - Decrypt with XChaCha20-Poly1305
5. `archive.UnzipTo()` - Extract ZIP to folder

### Generate Encryption Key

**CLI:**

```
ecrypto keygen --out mykey.txt
```

**Functions Involved:**

1. `keygenCmd.RunE()` - Parse flags
2. `crypto.KeySize()` - Get 32-byte requirement
3. `crypto/rand.Read()` - Generate random bytes
4. `base64.RawURLEncoding` - Encode for storage
5. File I/O - Save to file

---

## 🛡️ Security Model Summary

```
┌─────────────────────────────────────┐
│  Threat Model & Mitigations         │
├─────────────────────────────────────┤
│ Weak Passphrase        → Argon2id   │
│                          (memory-hard)
│                                     │
│ Brute-Force Attack     → 256 MB     │
│                          memory cost
│                                     │
│ Tampering Detection    → Auth Tag   │
│                          (Poly1305)
│                                     │
│ Ciphertext Forgery     → 16-byte    │
│                          auth tag
│                                     │
│ Random Nonce Reuse     → 24-byte    │
│                          random per op
│                                     │
│ Key Leakage            → Separate   │
│                          key files
│                                     │
│ File Metadata Leakage  → Encrypted  │
│                          in ZIP
└─────────────────────────────────────┘
```

---

## 📦 Deployment & Binaries

**Build Command:**

```bash
go build -o ecrypto
```

**Cross-Platform Builds:**

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o ecrypto.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o ecrypto

# Linux
GOOS=linux GOARCH=amd64 go build -o ecrypto
```

**Single Binary:**

- No external dependencies at runtime
- All cryptography included (golang.org/x/crypto)
- Charmbracelet libraries statically linked
- Zero configuration needed

---

## 📝 Summary

ECRYPTO is a production-grade encryption tool with:

✅ **Military-grade cryptography** (XChaCha20-Poly1305)  
✅ **Strong key derivation** (Argon2id winner)  
✅ **User-friendly interface** (Interactive + CLI)  
✅ **Tamper detection** (Authentication tags)  
✅ **Cross-platform** (Windows/macOS/Linux)  
✅ **Single binary** (Zero dependencies)  
✅ **Operation tracking** (Undo feature)

All components work together to provide a secure, efficient, and user-centric encryption experience.
