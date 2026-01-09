package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PromptUser displays a styled prompt and reads input
func PromptUser(label string, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	prompt := label
	if defaultVal != "" {
		prompt += fmt.Sprintf(" [%s]", defaultVal)
	}
	prompt += ": "

	// Use inline style without border for single-line input
	inlineStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)
	
	fmt.Print(inlineStyle.Render(prompt))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultVal != "" {
		return defaultVal
	}
	return input
}

// SelectOption displays a menu and returns selected index
func SelectOption(title string, options []string) int {
	headerStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorPrimary)
	
	fmt.Println(headerStyle.Render(title))
	fmt.Println()

	for i, opt := range options {
		marker := "   "
		style := MenuItemStyle
		if i == 0 {
			marker = " > "
			style = SelectedItemStyle
		}
		fmt.Println(style.Render(fmt.Sprintf("%s[%d] %s", marker, i+1, opt)))
	}
	fmt.Println()

	for {
		input := PromptUser("Select option", "1")
		if choice, err := parseChoice(input, len(options)); err == nil {
			return choice
		}
		fmt.Println(ErrorStyle.Render("Invalid choice. Try again."))
	}
}

// SelectFile displays file browser
func SelectFile(title string) string {
	for {
		path := PromptUser(title, "")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		fmt.Println(ErrorStyle.Render("File not found. Try again."))
	}
}

// SelectFileOrSkip displays file browser with ability to cancel
func SelectFileOrSkip(title string) string {
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser(title, "")
		if path == "" {
			return ""
		}
		if _, err := os.Stat(path); err == nil {
			return path
		}
		fmt.Println(ErrorStyle.Render(fmt.Sprintf("File not found. %d attempts remaining (or press Enter to cancel)", maxAttempts-i-1)))
	}
	return ""
}

// SelectFolder displays folder browser
func SelectFolder(title string) string {
	for {
		path := PromptUser(title, "")
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return path
		}
		fmt.Println(ErrorStyle.Render("Folder not found. Try again."))
	}
}

// PrintSuccess displays success message
func PrintSuccess(msg string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Padding(0, 2).
		Bold(true)
	fmt.Println(box.Render("[SUCCESS] " + msg))
	fmt.Println()
}

// PrintError displays error message
func PrintError(msg string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Foreground(ColorError).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorError).
		Padding(0, 2).
		Bold(true)
	fmt.Println(box.Render("[ERROR] " + msg))
	fmt.Println()
}

// PrintWarning displays warning message
func PrintWarning(msg string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Foreground(ColorWarning).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorWarning).
		Padding(0, 2).
		Bold(true)
	fmt.Println(box.Render("[WARNING] " + msg))
	fmt.Println()
}

// PrintInfo displays info message
func PrintInfo(msg string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 2)
	fmt.Println(box.Render("[INFO] " + msg))
	fmt.Println()
}

// ClearScreen clears terminal
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Pause waits for user to press Enter
func Pause() {
	fmt.Print(HelpStyle.Render("Press Enter to continue..."))
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// parseChoice converts user input to valid option index
func parseChoice(input string, maxOptions int) (int, error) {
	input = strings.TrimSpace(input)
	var choice int
	_, err := fmt.Sscanf(input, "%d", &choice)
	if err != nil || choice < 1 || choice > maxOptions {
		return 0, fmt.Errorf("invalid")
	}
	return choice - 1, nil
}

// PromptPassphrase securely prompts for passphrase (no echo)
func PromptPassphrase(label string) string {
	reader := bufio.NewReader(os.Stdin)
	inlineStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)
	fmt.Print(inlineStyle.Render(label + ": "))
	pass, _ := reader.ReadString('\n')
	fmt.Println()
	return strings.TrimSpace(pass)
}

// PrintBanner displays colorful banner
func PrintBanner() {
	
	banner := `                                                  
                                                  
██████ ▄█████ █████▄  ██  ██ █████▄ ██████ ▄████▄ 
██▄▄   ██     ██▄▄██▄  ▀██▀  ██▄▄█▀   ██   ██  ██ 
██▄▄▄▄ ▀█████ ██   ██   ██   ██       ██   ▀████▀ 
                                                  `
	fmt.Println(TitleStyle.Copy().Width(0).Padding(0).Foreground(ColorPrimary).Bold(true).Render(banner))
}

// PrintMenu displays main menu
func PrintMenu() {
	ClearScreen()
	PrintBanner()
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Align(lipgloss.Center).
		Width(60).
		Render("XChaCha20-Poly1305 | Argon2id | Military-Grade Security"))
	fmt.Println()
	fmt.Println()
}

// ConfirmAction asks for confirmation
func ConfirmAction(msg string) bool {
	inlineStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)
	fmt.Print(inlineStyle.Render(msg + " (yes/no): "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(input)) == "yes"
}

// GetFileSize returns human-readable file size
func GetFileSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIdx := 0
	for size > 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIdx])
}

// PrintFileInfo displays file information
func PrintFileInfo(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	fmt.Printf("File: %s\n", filepath.Base(path))
	fmt.Printf("Size: %s\n", GetFileSize(info.Size()))
	fmt.Printf("Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
}
