package ai

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// PathPattern represents a detected pattern in a file path
type PathPattern struct {
	Type        string  // "project", "backup", "documents", "media", "system"
	Confidence  float64
	Description string
}

// DetectPathPattern analyzes a path and identifies its likely purpose
func DetectPathPattern(path string) PathPattern {
	lowerPath := strings.ToLower(path)
	
	// Project/Development patterns
	if containsAny(lowerPath, []string{"project", "code", "src", "github", "git", "repo", "dev"}) {
		return PathPattern{
			Type:        "project",
			Confidence:  0.85,
			Description: "Development project",
		}
	}

	// Backup patterns
	if containsAny(lowerPath, []string{"backup", "archive", "old", "bak"}) {
		return PathPattern{
			Type:        "backup",
			Confidence:  0.9,
			Description: "Backup or archive",
		}
	}

	// Document patterns
	if containsAny(lowerPath, []string{"document", "docs", "papers", "reports", "work"}) {
		return PathPattern{
			Type:        "documents",
			Confidence:  0.8,
			Description: "Document folder",
		}
	}

	// Media patterns
	if containsAny(lowerPath, []string{"photo", "picture", "image", "video", "music", "media"}) {
		return PathPattern{
			Type:        "media",
			Confidence:  0.85,
			Description: "Media files",
		}
	}

	// System/Important patterns
	if containsAny(lowerPath, []string{"system", "program", "windows", "linux"}) {
		return PathPattern{
			Type:        "system",
			Confidence:  0.95,
			Description: "System folder (encrypt with caution)",
		}
	}

	return PathPattern{
		Type:        "general",
		Confidence:  0.5,
		Description: "General folder",
	}
}

// SuggestOutputName generates a smart output filename based on input pattern
func SuggestOutputName(inputPath string, pattern PathPattern) string {
	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	if ext != "" {
		baseName = strings.TrimSuffix(baseName, ext)
	}

	// Add descriptive suffix based on pattern
	switch pattern.Type {
	case "project":
		return baseName + "_project_encrypted.ecrypt"
	case "backup":
		return baseName + "_backup.ecrypt"
	case "documents":
		return baseName + "_docs_secure.ecrypt"
	case "media":
		return baseName + "_media_encrypted.ecrypt"
	default:
		return baseName + ".ecrypt"
	}
}

// GetSmartBackupLocation suggests best backup location based on OS and available drives
func GetSmartBackupLocation() string {
	// Try to find a secondary drive or external storage
	if runtime.GOOS == "windows" {
		// Check for D:, E:, F: drives
		for _, drive := range []string{"D:\\", "E:\\", "F:\\"} {
			if _, err := os.Stat(drive); err == nil {
				return filepath.Join(drive, "Backups")
			}
		}
	}

	// Fallback to user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(home, "Backups")
}

// IsCommonFolder checks if a path is a commonly accessed system folder
func IsCommonFolder(path string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	commonFolders := []string{
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Pictures"),
		filepath.Join(home, "Videos"),
		filepath.Join(home, "Music"),
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for _, common := range commonFolders {
		if strings.EqualFold(absPath, common) {
			return true
		}
	}

	return false
}

// PredictUserIntent analyzes recent history to predict what user wants to do
func PredictUserIntent() string {
	recent := GetRecentOperations(5)
	
	if len(recent) == 0 {
		return "new_user"
	}

	// Count recent operation types
	encryptCount := 0
	decryptCount := 0

	for _, op := range recent {
		if op.Type == "encrypt" {
			encryptCount++
		} else if op.Type == "decrypt" {
			decryptCount++
		}
	}

	// Predict based on patterns
	if encryptCount > decryptCount*2 {
		return "frequent_encryptor"
	} else if decryptCount > encryptCount*2 {
		return "frequent_decryptor"
	}

	return "balanced_user"
}

// GetPathSuggestionScore calculates relevance score for a path suggestion
func GetPathSuggestionScore(suggestedPath string, currentContext string) float64 {
	score := 0.5 // Base score

	// Boost score if path exists
	if _, err := os.Stat(suggestedPath); err == nil {
		score += 0.2
	}

	// Boost for common folders
	if IsCommonFolder(suggestedPath) {
		score += 0.15
	}

	// Boost if matches current context
	if strings.Contains(strings.ToLower(suggestedPath), strings.ToLower(currentContext)) {
		score += 0.15
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// FilterRelevantSuggestions removes low-confidence suggestions
func FilterRelevantSuggestions(suggestions []Suggestion, minConfidence float64) []Suggestion {
	filtered := []Suggestion{}
	
	for _, s := range suggestions {
		if s.Confidence >= minConfidence {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// containsAny checks if string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// GetContextualHint provides a helpful hint based on current operation context
func GetContextualHint(context string) string {
	hints := map[string]string{
		"encrypt_input":       "💡 Tip: You can drag & drop folders or browse with the file picker",
		"encrypt_output":      "💡 Tip: Save to a different drive for better backup security",
		"decrypt_input":       "💡 Tip: Recent encrypted files appear at the top",
		"decrypt_output":      "💡 Tip: Decrypt to a temporary location first to verify integrity",
		"password_entry":      "💡 Tip: Use a passphrase with 12+ characters or generate a key file",
		"keyfile_select":      "💡 Tip: Key files provide maximum security - store them safely",
		"operation_complete":  "💡 Tip: Verify your files and backup the encryption key",
		"weak_password":       "🚨 Warning: Consider using a key file for better security",
	}

	if hint, ok := hints[context]; ok {
		return hint
	}

	return ""
}
