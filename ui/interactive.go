package ui

import (
	"bufio"
	"ecrypto/ai"
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
	
	fmt.Print(inlineStyle.Render("› " + prompt))
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
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser(title, "")
		// Remove surrounding quotes if present
		path = strings.Trim(path, "\"")
		if _, err := os.Stat(path); err == nil {
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

// SelectFileOrSkip displays file browser with ability to cancel
func SelectFileOrSkip(title string) string {
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser(title, "")
		if path == "" {
			return ""
		}
		// Remove surrounding quotes if present
		path = strings.Trim(path, "\"")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		fmt.Println(ErrorStyle.Render(fmt.Sprintf("File not found. %d attempts remaining (or press Enter to cancel)", maxAttempts-i-1)))
	}
	return ""
}

// SelectFolder displays folder browser
func SelectFolder(title string) string {
	// Show contextual hint
	hint := ai.GetContextualHint("encrypt_input")
	if hint != "" {
		fmt.Println(lipgloss.NewStyle().Foreground(ColorPrimary).Render(hint))
	}
	
	// Show AI suggestions for recent and common paths
	suggestions := ai.SuggestRecentPaths("encrypt", 3)
	if len(suggestions) > 0 {
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Render("💡 Recent paths:"))
		for i, sug := range suggestions {
			if i < 3 {
				fmt.Println(lipgloss.NewStyle().Foreground(ColorDark).Render(fmt.Sprintf("   %s", sug.Text)))
			}
		}
		fmt.Println()
	}
	
	maxAttempts := 3
	for i := 0; i < maxAttempts; i++ {
		path := PromptUser(title, "")
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

// PrintSuccess displays success message
func PrintSuccess(msg string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Padding(0, 2).
		Bold(true)
	fmt.Println(box.Render("✓ SUCCESS: " + msg))
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
	fmt.Println(box.Render("✗ ERROR: " + msg))
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
	fmt.Println(box.Render("⚠ WARNING: " + msg))
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
	fmt.Println(box.Render("ℹ INFO: " + msg))
	fmt.Println()
}

// ClearScreen clears terminal
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Pause waits for user to press Enter
func Pause() {
	fmt.Println()
	fmt.Print(HelpStyle.Render("Press Enter to continue..."))
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
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
	
	// Show contextual hint
	hint := ai.GetContextualHint("password_entry")
	if hint != "" {
		fmt.Println(lipgloss.NewStyle().Foreground(ColorPrimary).Render(hint))
		fmt.Println()
	}
	
	fmt.Print(inlineStyle.Render(label + ": "))
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)
	
	// Analyze password strength and show suggestions
	if pass != "" {
		strength, suggestions, score := ai.AnalyzePasswordStrength(pass)
		
		// Color code based on strength
		var strengthColor lipgloss.Color
		switch strength {
		case "Strong":
			strengthColor = ColorSuccess
		case "Medium":
			strengthColor = ColorWarning
		case "Weak", "Very Weak":
			strengthColor = ColorError
		default:
			strengthColor = ColorDark
		}
		
		// Show strength indicator
		strengthStyle := lipgloss.NewStyle().Foreground(strengthColor).Bold(true)
		fmt.Println(strengthStyle.Render(fmt.Sprintf("\n  Strength: %s (%.0f%%)", strength, score*100)))
		
		// Show suggestions if any
		if len(suggestions) > 0 && strength != "Strong" {
			fmt.Println()
			for _, suggestion := range suggestions {
				fmt.Println(lipgloss.NewStyle().Foreground(ColorDark).Render("  " + suggestion))
			}
		}
		fmt.Println()
	}
	
	return pass
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
		Foreground(ColorWarning).
		Bold(true)
	fmt.Println()
	fmt.Print(inlineStyle.Render("⚠ " + msg + " (type 'yes' to confirm): "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	result := strings.ToLower(strings.TrimSpace(input)) == "yes"
	if !result {
		fmt.Println(HelpStyle.Render("✓ Operation cancelled."))
	}
	fmt.Println()
	return result
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

// SelectOptionOrPath displays a menu but also accepts a file/folder path
// Returns (choice index, detected path, is path detected)
func SelectOptionOrPath(title string, options []string) (int, string, bool) {
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
		
		// Try to parse as a menu choice first
		if choice, err := parseChoice(input, len(options)); err == nil {
			return choice, "", false
		}
		
		// Check if it looks like a path
		path := strings.Trim(input, "\"")
		if strings.Contains(path, "\\") || strings.Contains(path, "/") || strings.Contains(path, ":") {
			// Validate the path exists
			info, err := os.Stat(path)
			if err == nil {
				// Auto-detect if it's a folder or file
				if info.IsDir() {
					return 0, path, true // 0 = folder option
				} else {
					return 1, path, true // 1 = file option
				}
			}
			fmt.Println(ErrorStyle.Render("✗ Path not found. Enter a number (1-2) or paste a valid path."))
		} else {
			fmt.Println(ErrorStyle.Render("Invalid choice. Try again."))
		}
	}
}

// SelectOutputPath displays output path selector with AI suggestions
func SelectOutputPath(inputPath string) string {
	// Generate smart output suggestions
	suggestions := ai.SuggestOutputPath(inputPath)
	
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true).Render("💡 Suggested output paths:"))
	for i, sug := range suggestions {
		if i < 3 {
			confidence := fmt.Sprintf("%.0f%%", sug.Confidence*100)
			fmt.Println(lipgloss.NewStyle().Foreground(ColorDark).Render(
				fmt.Sprintf("   [%d] %s (%s) - %s", i+1, filepath.Base(sug.Text), confidence, sug.Description)))
		}
	}
	fmt.Println()
	
	// Show contextual hint
	hint := ai.GetContextualHint("encrypt_output")
	if hint != "" {
		fmt.Println(lipgloss.NewStyle().Foreground(ColorPrimary).Render(hint))
		fmt.Println()
	}
	
	// Allow user to pick suggestion or enter custom path
	input := PromptUser("Select suggestion (1-3) or enter custom path", "1")
	
	// Check if it's a number (selecting a suggestion)
	if choice, err := parseChoice(input, len(suggestions)); err == nil && choice < 3 {
		return suggestions[choice].Text
	}
	
	// Otherwise treat as custom path
	path := strings.Trim(input, "\"")
	return path
}
