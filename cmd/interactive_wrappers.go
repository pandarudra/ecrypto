package cmd

import (
	"bytes"
	"crypto/rand"
	"ecrypto/archive"
	"ecrypto/crypto"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EncryptWithPassphrase encrypts folder with passphrase
func EncryptWithPassphrase(inDir, outFile, pass string, progressCallback archive.ProgressCallback) error {
	h := &crypto.HeaderV1{
		Magic:   [8]byte{'E', 'C', 'R', 'Y', 'P', 'T', '0', '1'},
		Version: 1,
		KDF:     1,
	}

	if _, err := rand.Read(h.Salt[:]); err != nil {
		return err
	}

	key := crypto.DeriveKeyArgon2id(pass, h.Salt[:], encArgonM, encArgonT, encArgonP)
	h.ArgonM, h.ArgonT, h.ArgonP = encArgonM, encArgonT, encArgonP

	if _, err := rand.Read(h.Nonce[:]); err != nil {
		return err
	}

	zipBytes, err := archive.ZipFolderWithProgress(inDir, progressCallback)
	if err != nil {
		return err
	}

	aad := h.Encode()
	ct, err := crypto.EncryptAEAD(key, zipBytes, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	tmp := outFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := f.Write(aad); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(ct); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, outFile)
}

// EncryptWithKeyFile encrypts folder with key file
func EncryptWithKeyFile(inDir, outFile, keyFile string, progressCallback archive.ProgressCallback) error {
	h := &crypto.HeaderV1{
		Magic:   [8]byte{'E', 'C', 'R', 'Y', 'P', 'T', '0', '1'},
		Version: 1,
		KDF:     0,
	}

	key, err := crypto.ReadKeyFromFile(keyFile)
	if err != nil {
		return err
	}

	if _, err := rand.Read(h.Nonce[:]); err != nil {
		return err
	}

	zipBytes, err := archive.ZipFolderWithProgress(inDir, progressCallback)
	if err != nil {
		return err
	}

	aad := h.Encode()
	ct, err := crypto.EncryptAEAD(key, zipBytes, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	tmp := outFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := f.Write(aad); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(ct); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, outFile)
}

// DecryptWithPassphrase decrypts file with passphrase
func DecryptWithPassphrase(inFile, outDir, pass string, progressCallback archive.ProgressCallback) error {
	data, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
	if err != nil {
		return err
	}

	key := crypto.DeriveKeyArgon2id(pass, h.Salt[:], h.ArgonM, h.ArgonT, h.ArgonP)

	headerLen := crypto.HeaderSize()
	aad := data[:headerLen]
	ct := data[headerLen:]

	pt, err := crypto.DecryptAEAD(key, ct, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	// Try to unzip first (for folder encryption)
	err = archive.UnzipToWithProgress(outDir, pt, progressCallback)
	if err != nil {
		// If unzip fails, it's likely a single file encryption
		// Extract original filename from input file (remove .ecrypt extension)
		originalName := filepath.Base(inFile)
		if strings.HasSuffix(originalName, ".ecrypt") {
			originalName = strings.TrimSuffix(originalName, ".ecrypt")
		}
		
		// Write decrypted data as a single file
		outputPath := filepath.Join(outDir, originalName)
		if writeErr := os.WriteFile(outputPath, pt, 0o644); writeErr != nil {
			return fmt.Errorf("failed to unzip and failed to write as file: %v (original unzip error: %v)", writeErr, err)
		}
		
		if progressCallback != nil {
			progressCallback(originalName)
		}
	}
	
	return nil
}

// DecryptWithKeyFile decrypts file with key file
func DecryptWithKeyFile(inFile, outDir, keyFile string, progressCallback archive.ProgressCallback) error {
	data, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
	if err != nil {
		return err
	}

	key, err := crypto.ReadKeyFromFile(keyFile)
	if err != nil {
		return err
	}

	headerLen := crypto.HeaderSize()
	aad := data[:headerLen]
	ct := data[headerLen:]

	pt, err := crypto.DecryptAEAD(key, ct, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	// Try to unzip first (for folder encryption)
	err = archive.UnzipToWithProgress(outDir, pt, progressCallback)
	if err != nil {
		// If unzip fails, it's likely a single file encryption
		// Extract original filename from input file (remove .ecrypt extension)
		originalName := filepath.Base(inFile)
		if strings.HasSuffix(originalName, ".ecrypt") {
			originalName = strings.TrimSuffix(originalName, ".ecrypt")
		}
		
		// Write decrypted data as a single file
		outputPath := filepath.Join(outDir, originalName)
		if writeErr := os.WriteFile(outputPath, pt, 0o644); writeErr != nil {
			return fmt.Errorf("failed to unzip and failed to write as file: %v (original unzip error: %v)", writeErr, err)
		}
		
		if progressCallback != nil {
			progressCallback(originalName)
		}
	}
	
	return nil
}

// GenerateKey creates a random 32-byte key
func GenerateKey() (string, error) {
	key := make([]byte, crypto.KeySize())
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(key), nil
}

// InfoPrint prints container info
func InfoPrint(inFile string) error {
	data, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
	if err != nil {
		return err
	}

	fmt.Printf("Container: %s\n", inFile)
	fmt.Printf("Magic: %s\n", string(h.Magic[:]))
	fmt.Printf("Version: %d\n", h.Version)

	kdfName := "Raw Key"
	if h.KDF == 1 {
		kdfName = "Argon2id"
	}
	fmt.Printf("KDF: %s\n", kdfName)

	if h.KDF == 1 {
		fmt.Printf("  Memory: %d KiB\n", h.ArgonM)
		fmt.Printf("  Time: %d iterations\n", h.ArgonT)
		fmt.Printf("  Parallelism: %d\n", h.ArgonP)
	}

	fmt.Printf("Total file size: %d bytes\n", len(data))
	fmt.Printf("Header size: %d bytes\n", crypto.HeaderSize())
	fmt.Printf("Encrypted data size: %d bytes\n", len(data)-crypto.HeaderSize())

	return nil
}

// EncryptFileWithPassphrase encrypts a single file with passphrase
func EncryptFileWithPassphrase(filePath, outFile, pass string, progressCallback archive.ProgressCallback) error {
	h := &crypto.HeaderV1{
		Magic:   [8]byte{'E', 'C', 'R', 'Y', 'P', 'T', '0', '1'},
		Version: 1,
		KDF:     1,
	}

	if _, err := rand.Read(h.Salt[:]); err != nil {
		return err
	}

	key := crypto.DeriveKeyArgon2id(pass, h.Salt[:], encArgonM, encArgonT, encArgonP)
	h.ArgonM, h.ArgonT, h.ArgonP = encArgonM, encArgonT, encArgonP

	if _, err := rand.Read(h.Nonce[:]); err != nil {
		return err
	}

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if progressCallback != nil {
		progressCallback(filepath.Base(filePath))
	}

	aad := h.Encode()
	ct, err := crypto.EncryptAEAD(key, fileData, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	tmp := outFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := f.Write(aad); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(ct); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, outFile)
}

// EncryptFileWithKeyFile encrypts a single file with key file
func EncryptFileWithKeyFile(filePath, outFile, keyFile string, progressCallback archive.ProgressCallback) error {
	h := &crypto.HeaderV1{
		Magic:   [8]byte{'E', 'C', 'R', 'Y', 'P', 'T', '0', '1'},
		Version: 1,
		KDF:     0,
	}

	key, err := crypto.ReadKeyFromFile(keyFile)
	if err != nil {
		return err
	}

	if _, err := rand.Read(h.Nonce[:]); err != nil {
		return err
	}

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if progressCallback != nil {
		progressCallback(filepath.Base(filePath))
	}

	aad := h.Encode()
	ct, err := crypto.EncryptAEAD(key, fileData, aad, h.Nonce[:])
	if err != nil {
		return err
	}

	tmp := outFile + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := f.Write(aad); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(ct); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, outFile)
}
