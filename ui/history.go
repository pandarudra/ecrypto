package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Operation represents a recorded encryption/decryption operation
type Operation struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "encrypt" or "decrypt"
	SourcePath string   `json:"source_path"`
	OutputPath string   `json:"output_path"`
	Timestamp time.Time `json:"timestamp"`
	Size      int64     `json:"size"`
	FileCount int       `json:"file_count"`
	KeyMethod string    `json:"key_method"` // "passphrase" or "keyfile"
	KeyPath   string    `json:"key_path,omitempty"` // Only for keyfile method
	Status    string    `json:"status"` // "success" or "failed"
	Error     string    `json:"error,omitempty"`
}

// OperationHistory manages a local history of operations
type OperationHistory struct {
	HistoryFile string
	Operations  []Operation
}

// NewOperationHistory creates a new history tracker
func NewOperationHistory() *OperationHistory {
	// Store history in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	historyDir := filepath.Join(homeDir, ".ecrypto")
	os.MkdirAll(historyDir, 0o700)
	historyFile := filepath.Join(historyDir, "operations.json")

	oh := &OperationHistory{
		HistoryFile: historyFile,
		Operations:  []Operation{},
	}
	oh.Load()
	return oh
}

// Load reads operation history from disk
func (oh *OperationHistory) Load() error {
	data, err := os.ReadFile(oh.HistoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's fine
		}
		return err
	}

	return json.Unmarshal(data, &oh.Operations)
}

// Save writes operation history to disk
func (oh *OperationHistory) Save() error {
	data, err := json.MarshalIndent(oh.Operations, "", "  ")
	if err != nil {
		return err
	}

	// Make sure directory exists
	os.MkdirAll(filepath.Dir(oh.HistoryFile), 0o700)
	
	return os.WriteFile(oh.HistoryFile, data, 0o600)
}

// AddOperation records a new operation
func (oh *OperationHistory) AddOperation(op Operation) error {
	if op.ID == "" {
		op.ID = fmt.Sprintf("%d", time.Now().Unix())
	}
	if op.Timestamp.IsZero() {
		op.Timestamp = time.Now()
	}

	oh.Operations = append(oh.Operations, op)

	// Keep only last 50 operations to avoid large history file
	if len(oh.Operations) > 50 {
		oh.Operations = oh.Operations[len(oh.Operations)-50:]
	}

	return oh.Save()
}

// GetRecentOperations returns the N most recent operations
func (oh *OperationHistory) GetRecentOperations(count int) []Operation {
	if len(oh.Operations) == 0 {
		return []Operation{}
	}

	start := len(oh.Operations) - count
	if start < 0 {
		start = 0
	}

	// Return in reverse order (most recent first)
	recent := oh.Operations[start:]
	for i, j := 0, len(recent)-1; i < j; i, j = i+1, j-1 {
		recent[i], recent[j] = recent[j], recent[i]
	}
	return recent
}

// GetOperationByID retrieves an operation by ID
func (oh *OperationHistory) GetOperationByID(id string) *Operation {
	for i := range oh.Operations {
		if oh.Operations[i].ID == id {
			return &oh.Operations[i]
		}
	}
	return nil
}

// Clear removes all operations from history
func (oh *OperationHistory) Clear() error {
	oh.Operations = []Operation{}
	return oh.Save()
}

// FormatTime formats timestamp for display
func (op *Operation) FormatTime() string {
	return op.Timestamp.Format("2006-01-02 15:04:05")
}

// FormatSize returns human-readable size
func (op *Operation) FormatSize() string {
	return FormatBytes(op.Size)
}

// IsUndoable returns true if operation can be undone
func (op *Operation) IsUndoable() bool {
	// Can undo if:
	// - It's an encryption operation
	// - The encrypted file still exists
	// - The status was success
	if op.Type != "encrypt" || op.Status != "success" {
		return false
	}

	_, err := os.Stat(op.OutputPath)
	return err == nil
}
