package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FileBrowserItem represents a file or folder
type FileBrowserItem struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
}

// BrowseForFolder shows an interactive folder browser
func BrowseForFolder(title string) string {
	return browseFileSystem(title, true)
}

// BrowseForFile shows an interactive file browser
func BrowseForFile(title string) string {
	return browseFileSystem(title, false)
}

func browseFileSystem(title string, foldersOnly bool) string {
	helpText := lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Render("↑/↓: Navigate | Enter: Select | P: Paste Path | B: Back | Q: Cancel")
	
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("\n" + title))
	fmt.Println(helpText)
	fmt.Println()

	choice := PromptUser("Choose [1] Browse or [2] Paste Path", "1")
	
	if choice == "2" {
		path := PromptUser("Enter path", "")
		path = strings.Trim(path, "\"")
		if path != "" {
			info, err := os.Stat(path)
			if err == nil {
				if foldersOnly && info.IsDir() {
					return path
				} else if !foldersOnly && !info.IsDir() {
					return path
				} else if !foldersOnly && info.IsDir() {
					fmt.Println(ErrorStyle.Render("✗ Please select a file, not a folder"))
					return browseFileSystem(title, foldersOnly)
				} else if foldersOnly && !info.IsDir() {
					fmt.Println(ErrorStyle.Render("✗ Please select a folder, not a file"))
					return browseFileSystem(title, foldersOnly)
				}
			} else {
				fmt.Println(ErrorStyle.Render(fmt.Sprintf("✗ Path not found: %s", path)))
				return browseFileSystem(title, foldersOnly)
			}
		}
		return ""
	}

	// Start by showing drives on Windows
	return interactiveBrowse("", foldersOnly)
}

func interactiveBrowse(startPath string, foldersOnly bool) string {
	currentPath := startPath
	maxAttempts := 50 // Prevent infinite loops

	for attempt := 0; attempt < maxAttempts; attempt++ {
		ClearScreen()
		
		// Show current path or drives
		if currentPath == "" {
			fmt.Println(lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Render("💾 Select Drive"))
		} else {
			fmt.Println(lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Render(fmt.Sprintf("📂 Current: %s", currentPath)))
		}
		fmt.Println()

		// List items in current directory or show drives
		var items []FileBrowserItem
		var err error
		
		if currentPath == "" {
			items = listDrives()
		} else {
			items, err = listDirectory(currentPath, foldersOnly)
			if err != nil {
				fmt.Println(ErrorStyle.Render(fmt.Sprintf("Error reading directory: %v", err)))
				Pause()
				return ""
			}
		}

		if len(items) == 0 {
			fmt.Println(HelpStyle.Render("(Empty directory)"))
		}

		// Show items with numbers
		for i, item := range items {
			icon := "📄"
			if item.IsDir {
				icon = "📁"
			}
			if item.Name == ".." {
				icon = "⬆️"
			}
			
			itemStyle := MenuItemStyle
			prefix := fmt.Sprintf("  [%d] %s %s", i+1, icon, item.Name)
			
			if !item.IsDir && item.Size > 0 {
				prefix += fmt.Sprintf(" (%s)", FormatBytes(item.Size))
			}
			
			fmt.Println(itemStyle.Render(prefix))
		}

		fmt.Println()
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorDark).
			Italic(true).
			Render("Enter number to select | [S]elect current folder | [P]aste path | [Q]uit"))
		
		input := PromptUser("Choice", "")
		input = strings.ToUpper(strings.TrimSpace(input))

		// Handle special commands
		if input == "Q" {
			return ""
		}
		
		if input == "S" {
			if currentPath == "" {
				fmt.Println(ErrorStyle.Render("✗ Please select a drive first"))
				Pause()
				continue
			}
			if foldersOnly {
				return currentPath
			} else {
				fmt.Println(ErrorStyle.Render("✗ Please select a file"))
				Pause()
				continue
			}
		}

		if input == "P" {
			path := PromptUser("Enter path", "")
			path = strings.Trim(path, "\"")
			if path != "" {
				info, err := os.Stat(path)
				if err == nil {
					if foldersOnly && info.IsDir() {
						return path
					} else if !foldersOnly && !info.IsDir() {
						return path
					}
				}
				fmt.Println(ErrorStyle.Render("✗ Invalid path"))
				Pause()
			}
			continue
		}

		// Try to parse as number
		var selected int
		_, err = fmt.Sscanf(input, "%d", &selected)
		if err != nil || selected < 1 || selected > len(items) {
			fmt.Println(ErrorStyle.Render("✗ Invalid selection"))
			Pause()
			continue
		}

		selectedItem := items[selected-1]

		if selectedItem.IsDir {
			// Navigate into directory
			currentPath = selectedItem.Path
		} else {
			// File selected
			if foldersOnly {
				fmt.Println(ErrorStyle.Render("✗ Please select a folder, not a file"))
				Pause()
			} else {
				return selectedItem.Path
			}
		}
	}

	return ""
}

func listDirectory(path string, foldersOnly bool) ([]FileBrowserItem, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []FileBrowserItem

	// Add parent directory option
	parentPath := filepath.Dir(path)
	if parentPath != path {
		// Check if we're at a drive root (e.g., C:\)
		if len(path) == 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
			// At drive root, go back to drives view
			items = append(items, FileBrowserItem{
				Name:  "..",
				Path:  "", // Empty path means drives view
				IsDir: true,
			})
		} else {
			items = append(items, FileBrowserItem{
				Name:  "..",
				Path:  parentPath,
				IsDir: true,
			})
		}
	}

	// Add current directory items
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip hidden files on Windows (starting with .)
		if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".." {
			continue
		}

		// If folders only, skip files
		if foldersOnly && !entry.IsDir() {
			continue
		}

		items = append(items, FileBrowserItem{
			Name:  entry.Name(),
			Path:  filepath.Join(path, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}

	// Sort: directories first, then files, alphabetically
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == ".." {
			return true
		}
		if items[j].Name == ".." {
			return false
		}
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

// listDrives returns available drives on Windows
func listDrives() []FileBrowserItem {
	var drives []FileBrowserItem
	
	// Check drives A-Z
	for drive := 'A'; drive <= 'Z'; drive++ {
		drivePath := fmt.Sprintf("%c:\\", drive)
		if _, err := os.Stat(drivePath); err == nil {
			driveName := fmt.Sprintf("%c:", drive)
			drives = append(drives, FileBrowserItem{
				Name:  driveName,
				Path:  drivePath,
				IsDir: true,
			})
		}
	}
	
	return drives
}

// QuickPathSuggestions shows common folder suggestions
func QuickPathSuggestions() []string {
	var suggestions []string
	
	home, _ := os.UserHomeDir()
	if home != "" {
		suggestions = append(suggestions,
			filepath.Join(home, "Documents"),
			filepath.Join(home, "Downloads"),
			filepath.Join(home, "Pictures"),
			filepath.Join(home, "Desktop"),
			home,
		)
	}

	// Add current directory
	cwd, _ := os.Getwd()
	if cwd != "" {
		suggestions = append(suggestions, cwd)
	}

	// Filter out non-existent paths
	var validSuggestions []string
	for _, path := range suggestions {
		if _, err := os.Stat(path); err == nil {
			validSuggestions = append(validSuggestions, path)
		}
	}

	return validSuggestions
}
