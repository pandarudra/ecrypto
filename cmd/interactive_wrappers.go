package cmd

import (
	"bytes"
	"crypto/rand"
	"ecrypto/archive"
	"ecrypto/crypto"
	"encoding/base64"
	"fmt"
	"os"
)

// EncryptWithPassphrase encrypts folder with passphrase
func EncryptWithPassphrase(inDir, outFile, pass string) error {
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

	zipBytes, err := archive.ZipFolder(inDir)
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
func EncryptWithKeyFile(inDir, outFile, keyFile string) error {
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

	zipBytes, err := archive.ZipFolder(inDir)
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
func DecryptWithPassphrase(inFile, outDir, pass string) error {
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

	return archive.UnzipTo(outDir, pt)
}

// DecryptWithKeyFile decrypts file with key file
func DecryptWithKeyFile(inFile, outDir, keyFile string) error {
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

	return archive.UnzipTo(outDir, pt)
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
