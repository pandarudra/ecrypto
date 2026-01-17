package gui

import (
	"bytes"
	"ecrypto/ai"
	"ecrypto/cmd"
	"ecrypto/crypto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Server struct {
	port            int
	progressClients sync.Map // map[string]chan ProgressUpdate
}

type ProgressUpdate struct {
	OperationID string `json:"operationId"`
	Current     int    `json:"current"`
	Total       int    `json:"total"`
	Filename    string `json:"filename"`
	Percentage  int    `json:"percentage"`
}

type EncryptRequest struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	Password   string `json:"password,omitempty"`
	KeyFile    string `json:"keyFile,omitempty"`
	UseKey     bool   `json:"useKey"`
}

type DecryptRequest struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	Password   string `json:"password,omitempty"`
	KeyFile    string `json:"keyFile,omitempty"`
	UseKey     bool   `json:"useKey"`
}

type KeygenRequest struct {
	OutputPath string `json:"outputPath"`
}

type InfoRequest struct {
	FilePath string `json:"filePath"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	// Enable CORS for development
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/encrypt", s.handleEncrypt)
	mux.HandleFunc("/decrypt", s.handleDecrypt)
	mux.HandleFunc("/keygen", s.handleKeygen)
	mux.HandleFunc("/info", s.handleInfo)
	mux.HandleFunc("/history", s.handleHistory)
	mux.HandleFunc("/undo", s.handleUndo)
	mux.HandleFunc("/suggest-path", s.handleSuggestPath)
	mux.HandleFunc("/check-password", s.handleCheckPassword)
	mux.HandleFunc("/progress", s.handleProgressSSE)
	mux.HandleFunc("/health", s.handleHealth)

	handler := corsMiddleware(mux)

	addr := fmt.Sprintf("localhost:%d", s.port)
	log.Printf("Server started on http://%s\n", addr)
	return http.ListenAndServe(addr, handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":    "Ecrypto API Server",
		"version": "1.0",
		"status":  "running",
		"endpoints": []string{
			"POST /encrypt",
			"POST /decrypt",
			"POST /keygen",
			"POST /info",
			"GET  /history",
			"POST /undo",
			"POST /suggest-path",
			"POST /check-password",
			"GET  /progress",
			"GET  /health",
		},
	})
}

func (s *Server) handleEncrypt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.InputPath == "" || req.OutputPath == "" {
		sendError(w, "inputPath and outputPath are required", http.StatusBadRequest)
		return
	}

	// Check if input is a file or folder
	info, err := os.Stat(req.InputPath)
	if err != nil {
		sendError(w, fmt.Sprintf("Cannot access input path: %v", err), http.StatusBadRequest)
		return
	}

	// Progress callback
	progressCb := func(filename string) {
		// Could broadcast to SSE clients here
		log.Printf("Progress: %s\n", filename)
	}

	var encryptErr error
	if info.IsDir() {
		// Folder encryption
		if req.UseKey {
			if req.KeyFile == "" {
				sendError(w, "keyFile is required when useKey is true", http.StatusBadRequest)
				return
			}
			encryptErr = cmd.EncryptWithKeyFile(req.InputPath, req.OutputPath, req.KeyFile, progressCb)
		} else {
			if req.Password == "" {
				sendError(w, "password is required when useKey is false", http.StatusBadRequest)
				return
			}
			encryptErr = cmd.EncryptWithPassphrase(req.InputPath, req.OutputPath, req.Password, progressCb)
		}
	} else {
		// Single file encryption
		if req.UseKey {
			if req.KeyFile == "" {
				sendError(w, "keyFile is required when useKey is true", http.StatusBadRequest)
				return
			}
			encryptErr = cmd.EncryptFileWithKeyFile(req.InputPath, req.OutputPath, req.KeyFile, progressCb)
		} else {
			if req.Password == "" {
				sendError(w, "password is required when useKey is false", http.StatusBadRequest)
				return
			}
			encryptErr = cmd.EncryptFileWithPassphrase(req.InputPath, req.OutputPath, req.Password, progressCb)
		}
	}

	if encryptErr != nil {
		ai.AddOperation("encrypt", req.InputPath, req.OutputPath, getMethodName(req.UseKey), false)
		sendError(w, fmt.Sprintf("Encryption failed: %v", encryptErr), http.StatusInternalServerError)
		return
	}

	// Save to history
	ai.AddOperation("encrypt", req.InputPath, req.OutputPath, getMethodName(req.UseKey), true)

	sendSuccess(w, "Encryption completed successfully", map[string]interface{}{
		"outputPath": req.OutputPath,
	})
}

func (s *Server) handleDecrypt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.InputPath == "" || req.OutputPath == "" {
		sendError(w, "inputPath and outputPath are required", http.StatusBadRequest)
		return
	}

	progressCb := func(filename string) {
		log.Printf("Progress: %s\n", filename)
	}

	var decryptErr error
	if req.UseKey {
		if req.KeyFile == "" {
			sendError(w, "keyFile is required when useKey is true", http.StatusBadRequest)
			return
		}
		decryptErr = cmd.DecryptWithKeyFile(req.InputPath, req.OutputPath, req.KeyFile, progressCb)
	} else {
		if req.Password == "" {
			sendError(w, "password is required when useKey is false", http.StatusBadRequest)
			return
		}
		decryptErr = cmd.DecryptWithPassphrase(req.InputPath, req.OutputPath, req.Password, progressCb)
	}

	if decryptErr != nil {
		ai.AddOperation("decrypt", req.InputPath, req.OutputPath, getMethodName(req.UseKey), false)
		sendError(w, fmt.Sprintf("Decryption failed: %v", decryptErr), http.StatusInternalServerError)
		return
	}

	// Save to history
	ai.AddOperation("decrypt", req.InputPath, req.OutputPath, getMethodName(req.UseKey), true)

	sendSuccess(w, "Decryption completed successfully", map[string]interface{}{
		"outputPath": req.OutputPath,
	})
}

func (s *Server) handleKeygen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req KeygenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	key, err := cmd.GenerateKey()
	if err != nil {
		sendError(w, fmt.Sprintf("Key generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Write key to file if output path is provided
	if req.OutputPath != "" {
		if err := os.WriteFile(req.OutputPath, []byte(key), 0o600); err != nil {
			sendError(w, fmt.Sprintf("Failed to save key: %v", err), http.StatusInternalServerError)
			return
		}
	}

	sendSuccess(w, "Key generated successfully", map[string]interface{}{
		"key":        key,
		"outputPath": req.OutputPath,
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FilePath == "" {
		sendError(w, "filePath is required", http.StatusBadRequest)
		return
	}

	// Read the info to send as JSON
	data, err := os.ReadFile(req.FilePath)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to read container: %v", err), http.StatusInternalServerError)
		return
	}

	h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to decode header: %v", err), http.StatusInternalServerError)
		return
	}

	kdfType := "Raw Key"
	if h.KDF == 1 {
		kdfType = "Argon2id"
	}

	info := map[string]interface{}{
		"magic":            string(h.Magic[:]),
		"version":          h.Version,
		"kdfType":          kdfType,
		"argonMemory":      h.ArgonM,
		"argonTime":        h.ArgonT,
		"argonParallelism": h.ArgonP,
		"size":             len(data),
		"headerSize":       crypto.HeaderSize(),
		"encryptedSize":    len(data) - crypto.HeaderSize(),
	}

	sendSuccess(w, "Container info retrieved successfully", info)
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history, err := ai.LoadHistory()
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to load history: %v", err), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, "History retrieved successfully", history)
}

func (s *Server) handleUndo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OperationID string `json:"operationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find the operation
	op, err := ai.FindOperationByID(req.OperationID)
	if err != nil {
		sendError(w, "Operation not found", http.StatusNotFound)
		return
	}

	// Verify it's an encryption operation
	if op.Type != "encrypt" {
		sendError(w, "Only encryption operations can be undone", http.StatusBadRequest)
		return
	}

	// Check if the encrypted file still exists
	if _, err := os.Stat(op.OutputPath); os.IsNotExist(err) {
		sendError(w, "Encrypted file no longer exists", http.StatusNotFound)
		return
	}

	// Delete the encrypted file
	if err := os.Remove(op.OutputPath); err != nil {
		sendError(w, fmt.Sprintf("Failed to delete encrypted file: %v", err), http.StatusInternalServerError)
		return
	}

	// Remove the operation from history
	if err := ai.RemoveOperation(req.OperationID); err != nil {
		// File was deleted but history update failed - still consider it success
		log.Printf("Warning: Failed to remove operation from history: %v", err)
	}

	sendSuccess(w, "Operation undone successfully - encrypted file deleted", map[string]interface{}{
		"deletedFile": op.OutputPath,
	})
}

func (s *Server) handleSuggestPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	suggestions := ai.SuggestOutputPath(req.Path)
	sendSuccess(w, "Suggestions generated", suggestions)
}

func (s *Server) handleCheckPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
// Simple password strength check
	result := analyzePasswordStrength(req.Password)
	sendSuccess(w, "Password strength checked", result)
}

func analyzePasswordStrength(password string) map[string]interface{} {
	length := len(password)
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	score := 0
	if length >= 8 { score++ }
	if length >= 12 { score++ }
	if length >= 16 { score++ }
	if hasUpper { score++ }
	if hasLower { score++ }
	if hasDigit { score++ }
	if hasSpecial { score++ }

	var strength string
	var suggestions []string

	if score <= 3 {
		strength = "Weak"
		suggestions = append(suggestions, "Use at least 12 characters")
		if !hasUpper { suggestions = append(suggestions, "Add uppercase letters") }
		if !hasLower { suggestions = append(suggestions, "Add lowercase letters") }
		if !hasDigit { suggestions = append(suggestions, "Add numbers") }
		if !hasSpecial { suggestions = append(suggestions, "Add special characters") }
	} else if score <= 5 {
		strength = "Medium"
		if length < 16 { suggestions = append(suggestions, "Consider using 16+ characters") }
	} else if score <= 6 {
		strength = "Strong"
	} else {
		strength = "Very Strong"
	}

	return map[string]interface{}{
		"strength":    strength,
		"score":       score,
		"suggestions": suggestions,
	}
}

func (s *Server) handleProgressSSE(w http.ResponseWriter, r *http.Request) {
	// Server-Sent Events for real-time progress
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	
	// Keep connection open
	flusher, ok := w.(http.Flusher)
	if !ok {
		sendError(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send keep-alive
	fmt.Fprintf(w, "data: {\"status\": \"connected\"}\n\n")
	flusher.Flush()

	// Keep alive for now
	select {
	case <-r.Context().Done():
		return
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, "Server is healthy", map[string]interface{}{
		"status":  "ok",
		"version": "1.0",
	})
}

func getMethodName(useKey bool) string {
	if useKey {
		return "keyfile"
	}
	return "passphrase"
}

func sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}
