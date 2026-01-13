# Interactive File Browser - User Guide

## Features

The new interactive file browser provides a user-friendly way to select files and folders without typing paths manually.

### ✨ Key Features

1. **Interactive Navigation**
   - Browse through your file system with numbered selections
   - Navigate into folders by selecting them
   - Go back to parent directory with ".." option
2. **Quick Path Suggestions**
   - Shows common folders (Documents, Downloads, Pictures, Desktop)
   - Quick access to frequently used locations
3. **Dual Input Mode**
   - **Option 1**: Browse interactively with numbered menu
   - **Option 2**: Paste/Type path directly (like before)
4. **Smart Features**
   - Folders shown first, then files
   - File sizes displayed
   - Hidden files filtered out
   - Visual icons (📁 for folders, 📄 for files)

## Usage

### For Folder Selection

When prompted "Choose [1] Browse interactively or [2] Paste/Type path":

**Option 1 - Browse:**

```
📂 Current: C:\Users\YourName

  [1] ⬆️  ..
  [2] 📁 Documents
  [3] 📁 Downloads
  [4] 📁 Pictures
  [5] 📁 Desktop

Enter number to select | [S]elect current folder | [P]aste path | [Q]uit
Choice: 2
```

**Option 2 - Paste:**

```
Choice: 2
Enter folder path: C:\Users\YourName\Documents
```

### For File Selection

Same interface but shows both files and folders:

```
📂 Current: C:\Users\YourName\Documents

  [1] ⬆️  ..
  [2] 📁 Projects
  [3] 📄 report.pdf (2.5 MB)
  [4] 📄 notes.txt (15 KB)

Enter number to select | [P]aste path | [Q]uit
Choice: 3
```

## Commands

While browsing:

- **Number (1-9)**: Select item
- **S**: Select current folder (when selecting folders)
- **P**: Switch to paste/type mode
- **Q**: Quit/Cancel

## Integration

To use in your code:

```go
// For folder selection with enhanced UX
folder := ui.SelectFolderEnhanced("Select folder to encrypt")

// For file selection with enhanced UX
file := ui.SelectFileEnhanced("Select file to encrypt")
```

## Updating Existing Code

To enable the new file browser in menu.go, replace:

- `SelectFolder(...)` → `SelectFolderEnhanced(...)`
- `SelectFile(...)` → `SelectFileEnhanced(...)`

Example:

```go
// Old way
inPath = SelectFolder("Enter folder path")

// New way (with interactive browser)
inPath = SelectFolderEnhanced("Enter folder path")
```

## Benefits

1. ✅ No need to remember exact paths
2. ✅ Visual exploration of file system
3. ✅ Prevents typos in paths
4. ✅ Still supports pasting paths for power users
5. ✅ Quick access to common folders
6. ✅ See file sizes before selecting

## Future Enhancements

Possible improvements:

- Search/filter functionality
- Bookmarks/favorites
- Recent files list
- Multi-select for batch operations
- Arrow key navigation (↑/↓)
