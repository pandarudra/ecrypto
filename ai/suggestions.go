package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Suggestion represents an AI-powered suggestion
type Suggestion struct {
	Text        string
	Confidence  float64 // 0.0 to 1.0
	Type        string  // "path", "output", "password", "option"
	Description string
}

// SuggestOutputPath generates smart output path suggestions based on input
func SuggestOutputPath(inputPath string) []Suggestion {
	suggestions := []Suggestion{}

	// Get base name without extension
	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	if ext != "" {
		baseName = strings.TrimSuffix(baseName, ext)
	}

	// Get directory of input
	inputDir := filepath.Dir(inputPath)

	// Suggestion 1: Same location with timestamp
	timestamp := time.Now().Format("20060102_150405")
	sameDirPath := filepath.Join(inputDir, fmt.Sprintf("%s_%s.ecrypt", baseName, timestamp))
	suggestions = append(suggestions, Suggestion{
		Text:        sameDirPath,
		Confidence:  0.9,
		Type:        "output",
		Description: "Same location with timestamp",
	})

	// Suggestion 2: Simple .ecrypt extension
	simplePath := filepath.Join(inputDir, baseName+".ecrypt")
	suggestions = append(suggestions, Suggestion{
		Text:        simplePath,
		Confidence:  0.85,
		Type:        "output",
		Description: "Simple naming",
	})

	// Suggestion 3: Desktop backup
	desktopPath := filepath.Join(getUserHome(), "Desktop", fmt.Sprintf("%s_backup.ecrypt", baseName))
	suggestions = append(suggestions, Suggestion{
		Text:        desktopPath,
		Confidence:  0.7,
		Type:        "output",
		Description: "Quick desktop backup",
	})

	return suggestions
}

// SuggestRecentPaths returns recently used paths from history
func SuggestRecentPaths(historyType string, limit int) []Suggestion {
	history := GetHistory()
	suggestions := []Suggestion{}

	count := 0
	for i := len(history.Operations) - 1; i >= 0 && count < limit; i-- {
		op := history.Operations[i]
		
		// Filter by type if specified
		if historyType != "" && op.Type != historyType {
			continue
		}

		// Add input path suggestion
		if op.InputPath != "" {
			suggestions = append(suggestions, Suggestion{
				Text:        op.InputPath,
				Confidence:  0.8 - (float64(len(history.Operations)-i-1) * 0.1),
				Type:        "path",
				Description: fmt.Sprintf("Recent %s (%s)", op.Type, op.Timestamp.Format("Jan 02, 15:04")),
			})
			count++
		}
	}

	return suggestions
}

// AnalyzePasswordStrength evaluates password strength and provides suggestions
func AnalyzePasswordStrength(password string) (string, []string, float64) {
	if len(password) == 0 {
		return "Empty", []string{"Password cannot be empty"}, 0.0
	}

	score := 0.0
	suggestions := []string{}

	// Length check
	if len(password) >= 12 {
		score += 0.3
	} else if len(password) >= 8 {
		score += 0.15
		suggestions = append(suggestions, "💡 Use at least 12 characters for better security")
	} else {
		suggestions = append(suggestions, "⚠️ Password too short - use at least 12 characters")
	}

	// Character variety checks
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	varietyCount := 0
	if hasUpper {
		varietyCount++
		score += 0.15
	} else {
		suggestions = append(suggestions, "💡 Add uppercase letters (A-Z)")
	}

	if hasLower {
		varietyCount++
		score += 0.15
	} else {
		suggestions = append(suggestions, "💡 Add lowercase letters (a-z)")
	}

	if hasDigit {
		varietyCount++
		score += 0.15
	} else {
		suggestions = append(suggestions, "💡 Add numbers (0-9)")
	}

	if hasSpecial {
		varietyCount++
		score += 0.25
	} else {
		suggestions = append(suggestions, "💡 Add special characters (!@#$%^&*)")
	}

	// Common patterns check
	commonPatterns := []string{"password", "123456", "qwerty", "abc123", "letmein", "welcome"}
	lowerPass := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lowerPass, pattern) {
			score -= 0.3
			suggestions = append(suggestions, "⚠️ Avoid common patterns like '"+pattern+"'")
			break
		}
	}

	// Determine strength level
	var strength string
	if score >= 0.75 {
		strength = "Strong"
	} else if score >= 0.5 {
		strength = "Medium"
		suggestions = append(suggestions, "💡 Consider using a key file for maximum security")
	} else if score >= 0.25 {
		strength = "Weak"
		suggestions = append(suggestions, "⚠️ This password is weak - strongly consider using a key file instead")
	} else {
		strength = "Very Weak"
		suggestions = append(suggestions, "🚨 CRITICAL: Use a much stronger password or generate a key file!")
	}

	// Ensure score is between 0 and 1
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return strength, suggestions, score
}

// SuggestCommonPaths returns commonly accessed system folders
func SuggestCommonPaths() []Suggestion {
	home := getUserHome()
	
	commonPaths := []struct {
		Path        string
		Description string
		Confidence  float64
	}{
		{filepath.Join(home, "Documents"), "Documents folder", 0.9},
		{filepath.Join(home, "Downloads"), "Downloads folder", 0.85},
		{filepath.Join(home, "Desktop"), "Desktop folder", 0.8},
		{filepath.Join(home, "Pictures"), "Pictures folder", 0.75},
		{filepath.Join(home, "Videos"), "Videos folder", 0.7},
		{filepath.Join(home, "Music"), "Music folder", 0.65},
	}

	suggestions := []Suggestion{}
	for _, cp := range commonPaths {
		suggestions = append(suggestions, Suggestion{
			Text:        cp.Path,
			Confidence:  cp.Confidence,
			Type:        "path",
			Description: cp.Description,
		})
	}

	return suggestions
}

// SuggestNextAction recommends what the user should do next
func SuggestNextAction(context string) []Suggestion {
	suggestions := []Suggestion{}

	switch context {
	case "after_encrypt":
		suggestions = append(suggestions, Suggestion{
			Text:        "Verify encryption with 'info' command",
			Confidence:  0.9,
			Type:        "option",
			Description: "Check encrypted file metadata",
		})
		suggestions = append(suggestions, Suggestion{
			Text:        "Backup the encryption key/passphrase securely",
			Confidence:  0.95,
			Type:        "option",
			Description: "Store in password manager",
		})

	case "after_decrypt":
		suggestions = append(suggestions, Suggestion{
			Text:        "Verify decrypted files integrity",
			Confidence:  0.85,
			Type:        "option",
			Description: "Check if all files restored correctly",
		})

	case "weak_password":
		suggestions = append(suggestions, Suggestion{
			Text:        "Generate a key file instead",
			Confidence:  0.9,
			Type:        "option",
			Description: "Use 'keygen' command for maximum security",
		})
	}

	return suggestions
}

// getUserHome returns the user's home directory
func getUserHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "C:\\Users\\Default"
	}
	return home
}
