package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SelectFolderEnhanced displays folder browser with option to browse or paste
func SelectFolderEnhanced(title string) string {
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("\n" + title))
	
	// Show quick suggestions
	suggestions := QuickPathSuggestions()
	if len(suggestions) > 0 {
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorDark).
			Italic(true).
			Render("\nQuick access:"))
		for i, path := range suggestions {
			if i < 5 { // Show max 5 suggestions
				fmt.Println(HelpStyle.Render(fmt.Sprintf("  • %s", path)))
			}
		}
		fmt.Println()
	}
	
	choice := PromptUser("Choose [1] Browse interactively or [2] Paste/Type path", "1")
	
	if choice == "1" {
		return BrowseForFolder(title)
	}
	
	// Paste path mode with attempts
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser("Enter folder path", "")
		// Remove surrounding quotes if present
		path = strings.Trim(path, "\"")
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return path
		}
		remaining := maxAttempts - i - 1
		if remaining > 0 {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("✗ Folder not found. %d attempt(s) remaining.", remaining)))
		} else {
			fmt.Println(ErrorStyle.Render("✗ Too many invalid attempts. Returning to menu."))
			Pause()
			return ""
		}
	}
	return ""
}

// SelectFileEnhanced displays file browser with option to browse or paste
func SelectFileEnhanced(title string) string {
	choice := PromptUser("Choose [1] Browse interactively or [2] Paste/Type path", "1")
	
	if choice == "1" {
		return BrowseForFile(title)
	}
	
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser("Enter file path", "")
		// Remove surrounding quotes if present
		path = strings.Trim(path, "\"")
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path
		}
		remaining := maxAttempts - i - 1
		if remaining > 0 {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("✗ File not found. %d attempt(s) remaining.", remaining)))
		} else {
			fmt.Println(ErrorStyle.Render("✗ Too many invalid attempts. Returning to menu."))
			Pause()
			return ""
		}
	}
	return ""
}
