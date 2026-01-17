package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Operation represents a single encryption/decryption operation
type Operation struct {
	ID         string    `json:"id"`          // Unique identifier
	Type       string    `json:"type"`        // "encrypt" or "decrypt"
	InputPath  string    `json:"input_path"`  // Source file/folder path
	OutputPath string    `json:"output_path"` // Destination path
	Method     string    `json:"method"`      // "passphrase" or "keyfile"
	Timestamp  time.Time `json:"timestamp"`
	Success    bool      `json:"success"`
}

// History stores all operations
type History struct {
	Operations []Operation `json:"operations"`
	MaxSize    int         `json:"max_size"` // Maximum operations to keep
}

var (
	historyFile = ".ecrypto_history.json"
	maxHistory  = 100 // Keep last 100 operations
)

// GetHistoryPath returns the path to the history file in user's home directory
func GetHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return historyFile
	}
	return filepath.Join(home, historyFile)
}

// LoadHistory loads the operation history from disk
func LoadHistory() (*History, error) {
	histPath := GetHistoryPath()

	// If file doesn't exist, return empty history
	if _, err := os.Stat(histPath); os.IsNotExist(err) {
		return &History{
			Operations: []Operation{},
			MaxSize:    maxHistory,
		}, nil
	}

	// Read and parse file
	data, err := os.ReadFile(histPath)
	if err != nil {
		return nil, err
	}

	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	// Set default max size if not set
	if history.MaxSize == 0 {
		history.MaxSize = maxHistory
	}

	return &history, nil
}

// SaveHistory saves the operation history to disk
func SaveHistory(history *History) error {
	histPath := GetHistoryPath()

	// Trim history if it exceeds max size
	if len(history.Operations) > history.MaxSize {
		history.Operations = history.Operations[len(history.Operations)-history.MaxSize:]
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(histPath, data, 0600)
}

// AddOperation adds a new operation to history
func AddOperation(opType, inputPath, outputPath, method string, success bool) error {
	history, err := LoadHistory()
	if err != nil {
		// If we can't load history, create new one
		history = &History{
			Operations: []Operation{},
			MaxSize:    maxHistory,
		}
	}

	// Add new operation
	timestamp := time.Now()
	op := Operation{
		ID:         generateOperationID(timestamp),
		Type:       opType,
		InputPath:  inputPath,
		OutputPath: outputPath,
		Method:     method,
		Timestamp:  timestamp,
		Success:    success,
	}

	history.Operations = append(history.Operations, op)

	// Save updated history
	return SaveHistory(history)
}

// GetHistory returns the current history (convenience function)
func GetHistory() *History {
	history, err := LoadHistory()
	if err != nil {
		// Return empty history on error
		return &History{
			Operations: []Operation{},
			MaxSize:    maxHistory,
		}
	}
	return history
}

// GetRecentOperations returns the N most recent operations
func GetRecentOperations(n int) []Operation {
	history := GetHistory()
	
	if len(history.Operations) == 0 {
		return []Operation{}
	}

	// Return last N operations
	start := len(history.Operations) - n
	if start < 0 {
		start = 0
	}

	return history.Operations[start:]
}

// GetRecentPaths returns unique paths from recent operations
func GetRecentPaths(opType string, n int) []string {
	history := GetHistory()
	pathMap := make(map[string]bool)
	paths := []string{}

	// Iterate backwards (most recent first)
	for i := len(history.Operations) - 1; i >= 0 && len(paths) < n; i-- {
		op := history.Operations[i]
		
		// Filter by operation type if specified
		if opType != "" && op.Type != opType {
			continue
		}

		// Add input path if not already added
		if op.InputPath != "" && !pathMap[op.InputPath] {
			paths = append(paths, op.InputPath)
			pathMap[op.InputPath] = true
		}

		// Add output path if not already added
		if op.OutputPath != "" && !pathMap[op.OutputPath] {
			paths = append(paths, op.OutputPath)
			pathMap[op.OutputPath] = true
		}
	}

	return paths
}

// generateOperationID creates a unique ID based on timestamp
func generateOperationID(t time.Time) string {
	return t.Format("20060102-150405.000000")
}

// FindOperationByID finds an operation by its ID
func FindOperationByID(id string) (*Operation, error) {
	history, err := LoadHistory()
	if err != nil {
		return nil, err
	}

	for _, op := range history.Operations {
		if op.ID == id {
			return &op, nil
		}
	}

	return nil, os.ErrNotExist
}

// RemoveOperation removes an operation by ID
func RemoveOperation(id string) error {
	history, err := LoadHistory()
	if err != nil {
		return err
	}

	// Find and remove the operation
	newOps := []Operation{}
	found := false
	for _, op := range history.Operations {
		if op.ID != id {
			newOps = append(newOps, op)
		} else {
			found = true
		}
	}

	if !found {
		return os.ErrNotExist
	}

	history.Operations = newOps
	return SaveHistory(history)
}

// ClearHistory removes all history
func ClearHistory() error {
	histPath := GetHistoryPath()
	
	// Remove file if it exists
	if _, err := os.Stat(histPath); err == nil {
		return os.Remove(histPath)
	}
	
	return nil
}

// GetStats returns statistics about operations
func GetStats() map[string]interface{} {
	history := GetHistory()
	
	stats := map[string]interface{}{
		"total_operations": len(history.Operations),
		"encryptions":      0,
		"decryptions":      0,
		"successes":        0,
		"failures":         0,
		"passphrase_ops":   0,
		"keyfile_ops":      0,
	}

	for _, op := range history.Operations {
		if op.Type == "encrypt" {
			stats["encryptions"] = stats["encryptions"].(int) + 1
		} else if op.Type == "decrypt" {
			stats["decryptions"] = stats["decryptions"].(int) + 1
		}

		if op.Success {
			stats["successes"] = stats["successes"].(int) + 1
		} else {
			stats["failures"] = stats["failures"].(int) + 1
		}

		if op.Method == "passphrase" {
			stats["passphrase_ops"] = stats["passphrase_ops"].(int) + 1
		} else if op.Method == "keyfile" {
			stats["keyfile_ops"] = stats["keyfile_ops"].(int) + 1
		}
	}

	return stats
}
