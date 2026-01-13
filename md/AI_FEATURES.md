# 🤖 AI-Powered Suggestions - MVP Implementation

## Overview

ECRYPTO now includes **local AI-powered suggestions** to enhance user experience and make encryption more intuitive. This MVP implementation uses pattern matching, history analysis, and rule-based intelligence—no external AI APIs or internet required!

## ✨ Features Implemented

### 1. 💡 Smart Output Path Suggestions

When encrypting files or folders, the AI suggests intelligent output paths:

```
💡 Suggested output paths:
   [1] project_20260113_150405.ecrypt (90%) - Same location with timestamp
   [2] project.ecrypt (85%) - Simple naming
   [3] C:\Users\shree\Desktop\project_backup.ecrypt (70%) - Quick desktop backup
```

**How it works:**

- Analyzes input path and generates contextually relevant suggestions
- Adds timestamps for version control
- Suggests alternative backup locations (desktop, secondary drives)
- Calculates confidence scores for each suggestion

### 2. 📊 Real-Time Password Strength Analysis

As you type your passphrase, get instant feedback:

```
Enter passphrase: MyPass123!

  Strength: Medium (60%)

  💡 Use at least 12 characters for better security
  💡 Add special characters (!@#$%^&*)
  💡 Consider using a key file for maximum security
```

**Strength indicators:**

- ✅ **Strong** (75%+): Green, ready to go
- ⚠️ **Medium** (50-74%): Yellow, consider improvements
- 🚨 **Weak/Very Weak** (<50%): Red, strongly recommend key file

**Analysis factors:**

- Password length (8/12/16+ characters)
- Character variety (uppercase, lowercase, numbers, symbols)
- Common pattern detection (password, 123456, qwerty, etc.)

### 3. 🕐 Recent Path Suggestions

The AI remembers your recent operations:

```
💡 Recent paths:
   C:\Users\shree\Documents\project (Recent encrypt, Jan 13, 15:04)
   D:\Backups\data (Recent encrypt, Jan 13, 14:30)
   C:\Downloads\files (Recent encrypt, Jan 13, 13:15)
```

**History tracking:**

- Stores last 100 operations in `.ecrypto_history.json` (in home directory)
- Suggests frequently used paths
- Shows operation type and timestamp
- Privacy-first: All data stored locally

### 4. 💬 Contextual Hints

Smart tips appear at the right moment:

- **During encryption:** "💡 Tip: Save to a different drive for better backup security"
- **During decryption:** "💡 Tip: Recent encrypted files appear at the top"
- **Password entry:** "💡 Tip: Use a passphrase with 12+ characters or generate a key file"
- **After encryption:** "💡 Tip: Verify your files and backup the encryption key"

### 5. 🎯 Next Action Recommendations

After completing an operation, AI suggests what to do next:

```
✓ SUCCESS: File encrypted successfully!
▶ Output: C:\project\data.ecrypt

💡 What's next?
   • Verify encryption with 'info' command
   • Backup the encryption key/passphrase securely
```

## 📂 Architecture

### File Structure

```
ecrypto/
├── ai/
│   ├── suggestions.go    # Core suggestion engine
│   ├── history.go         # Operation history tracking
│   └── patterns.go        # Pattern matching & analysis
├── ui/
│   ├── interactive.go     # Enhanced with AI suggestions
│   └── menu.go            # History tracking integration
└── .ecrypto_history.json  # Local history storage (user's home dir)
```

### Key Components

**1. Suggestion Engine (`ai/suggestions.go`)**

- `SuggestOutputPath()` - Smart output naming
- `AnalyzePasswordStrength()` - Real-time password analysis
- `SuggestRecentPaths()` - History-based suggestions
- `SuggestCommonPaths()` - System folder suggestions
- `SuggestNextAction()` - Post-operation recommendations

**2. History Manager (`ai/history.go`)**

- `AddOperation()` - Log operations (encrypt/decrypt)
- `LoadHistory()` - Retrieve past operations
- `GetRecentPaths()` - Extract frequently used paths
- `GetStats()` - Usage statistics

**3. Pattern Analyzer (`ai/patterns.go`)**

- `DetectPathPattern()` - Identify folder types (project/backup/media)
- `GetSmartBackupLocation()` - Find optimal backup drives
- `PredictUserIntent()` - Understand user behavior patterns
- `GetContextualHint()` - Context-aware tips

## 🔐 Privacy & Security

✅ **100% Local Processing** - No cloud AI, no external APIs
✅ **Offline-First** - Works without internet connection
✅ **User Control** - History stored in user's home directory
✅ **Transparent** - All suggestions are rule-based and explainable
✅ **Secure Storage** - History file has 0600 permissions (owner-only read/write)

## 🚀 Usage Examples

### Example 1: First-Time User

```powershell
# Launch interactive mode
.\ecrypto.exe

# AI shows contextual help
💡 Tip: You can drag & drop folders or browse with the file picker

# During password entry
Enter passphrase: pass

  Strength: Very Weak (15%)

  ⚠️ Password too short - use at least 12 characters
  💡 Add uppercase letters (A-Z)
  💡 Add numbers (0-9)
  💡 Add special characters (!@#$%^&*)
  🚨 CRITICAL: Use a much stronger password or generate a key file!
```

### Example 2: Returning User

```powershell
# AI remembers your previous operations
💡 Recent paths:
   C:\Projects\webapp (Recent encrypt, Jan 13, 15:04)
   D:\Backups\data (Recent encrypt, Jan 13, 14:30)

# Smart output suggestions based on input
Enter: C:\Projects\webapp

💡 Suggested output paths:
   [1] webapp_20260113_153020.ecrypt (90%) - Same location with timestamp
   [2] webapp.ecrypt (85%) - Simple naming
   [3] D:\Backups\webapp_backup.ecrypt (70%) - Secondary drive backup

Select suggestion (1-3) or enter custom path [1]:
```

### Example 3: Power User Workflow

```powershell
# After encryption completes
✓ SUCCESS: Folder encrypted successfully!
▶ Output: D:\Backups\webapp_20260113.ecrypt

💡 What's next?
   • Verify encryption with 'info' command
   • Backup the encryption key/passphrase securely
```

## 📊 History File Format

The AI stores operation history in `~/.ecrypto_history.json`:

```json
{
  "operations": [
    {
      "type": "encrypt",
      "input_path": "C:\\Projects\\webapp",
      "output_path": "D:\\Backups\\webapp_20260113.ecrypt",
      "method": "passphrase",
      "timestamp": "2026-01-13T15:30:20Z",
      "success": true
    }
  ],
  "max_size": 100
}
```

**Management:**

- Auto-rotates at 100 operations
- Can be manually cleared (future feature)
- Respects privacy—no password data stored

## 🎨 Visual Enhancements

All AI suggestions use color-coded displays:

- 💚 **Green** - Strong/Success
- 💛 **Yellow** - Medium/Warning
- 🔴 **Red** - Weak/Critical
- 🔵 **Blue** - Info/Tips
- ⚪ **Gray** - Secondary info

## 🔮 Future Enhancements (Beyond MVP)

Potential additions for future versions:

1. **Machine Learning Integration**

   - Learn from user behavior patterns
   - Personalized suggestions based on usage

2. **Advanced Pattern Recognition**

   - File type detection (code/media/documents)
   - Project structure analysis
   - Dependency detection

3. **Smart Key Management**

   - Suggest key rotation schedules
   - Detect weak key storage practices
   - Integration with password managers

4. **Collaborative Features**

   - Team workspace detection
   - Shared encryption patterns
   - Multi-user security policies

5. **Cloud Intelligence (Optional)**
   - Opt-in cloud AI for advanced suggestions
   - Privacy-preserving federated learning
   - Encrypted telemetry

## 🧪 Testing

To test the AI features:

```powershell
# 1. Fresh start - no history
rm ~\.ecrypto_history.json
.\ecrypto.exe

# 2. Encrypt a file with weak password
# Watch password strength analysis in action

# 3. Complete encryption
# See "what's next" suggestions

# 4. Run again - see recent path suggestions
.\ecrypto.exe

# 5. Check history file
cat ~\.ecrypto_history.json
```

## 📝 Configuration

Currently, AI features are always enabled. Future versions may include:

```json
{
  "ai_suggestions": true,
  "password_strength_check": true,
  "history_tracking": true,
  "max_history_size": 100,
  "suggestion_confidence_threshold": 0.5
}
```

## 🤝 Contributing

Want to improve the AI? Focus areas:

- Better pattern detection algorithms
- More contextual hints
- Improved password strength heuristics
- Cross-platform path suggestions
- Performance optimizations

## 📄 License

Same as ECRYPTO main license (see LICENSE file)

---

**Status:** ✅ MVP Complete - Ready for testing
**Version:** 1.0.8-dev (AI features)
**Last Updated:** January 13, 2026
