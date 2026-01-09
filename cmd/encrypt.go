package cmd

import (
	"crypto/rand"
	"ecrypto/archive"
	"ecrypto/crypto"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
    encInDir   string
    encOutFile string
    encPass    string
    encKeyFile string
    encArgonM  uint32 = 256 * 1024 // 256 MB in KiB
    encArgonT  uint32 = 3
    encArgonP  uint8  = 1
)

var encryptCmd = &cobra.Command{
    Use:   "encrypt",
    Short: "Encrypt a folder into a .ecrypt container",
    Long: `Encrypt a folder into a secure .ecrypt container.
Use --pass for passphrase (Argon2id KDF) or --key-file for a raw 32-byte key.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if encInDir == "" || encOutFile == "" {
            return errors.New("--in and --out are required")
        }

        // Initialize header
        h := &crypto.HeaderV1{
            Magic:   [8]byte{'E', 'C', 'R', 'Y', 'P', 'T', '0', '1'},
            Version: 1,
        }

        var key []byte

        // Derive or load key
        if encPass != "" {
            h.KDF = 1
            if _, err := rand.Read(h.Salt[:]); err != nil {
                return err
            }
            key = crypto.DeriveKeyArgon2id(encPass, h.Salt[:], encArgonM, encArgonT, encArgonP)
            h.ArgonM, h.ArgonT, h.ArgonP = encArgonM, encArgonT, encArgonP
            fmt.Fprintf(os.Stderr, "Using Argon2id: m=%d, t=%d, p=%d\n", encArgonM, encArgonT, encArgonP)
        } else if encKeyFile != "" {
            h.KDF = 0
            var err error
            key, err = crypto.ReadKeyFromFile(encKeyFile)
            if err != nil {
                return err
            }
            fmt.Fprintf(os.Stderr, "Loaded 32-byte key from file\n")
        } else {
            return errors.New("provide --pass or --key-file")
        }

        // Generate random nonce
        if _, err := rand.Read(h.Nonce[:]); err != nil {
            return err
        }

        // Zip folder
        fmt.Fprintf(os.Stderr, "Compressing folder...\n")
        zipBytes, err := archive.ZipFolder(encInDir)
        if err != nil {
            return err
        }
        fmt.Fprintf(os.Stderr, "Compressed size: %d bytes\n", len(zipBytes))

        // Encrypt
        fmt.Fprintf(os.Stderr, "Encrypting...\n")
        aad := h.Encode()
        ct, err := crypto.EncryptAEAD(key, zipBytes, aad, h.Nonce[:])
        if err != nil {
            return err
        }

        // Write to file (atomic: write to .tmp, then rename)
        tmp := encOutFile + ".tmp"
        f, err := os.Create(tmp)
        if err != nil {
            return err
        }
        defer f.Close()

        if _, err := f.Write(aad); err != nil {
            return err
        }
        if _, err := f.Write(ct); err != nil {
            return err
        }
        if err := f.Close(); err != nil {
            return err
        }

        if err := os.Rename(tmp, encOutFile); err != nil {
            return err
        }

        fmt.Fprintf(os.Stderr, "✓ Encrypted to: %s\n", encOutFile)
        fmt.Fprintf(os.Stderr, "  Encrypted size: %d bytes\n", len(aad)+len(ct))
        return nil
    },
}

func init() {
    rootCmd.AddCommand(encryptCmd)
    encryptCmd.Flags().StringVar(&encInDir, "in", "", "Input folder to encrypt")
    encryptCmd.Flags().StringVar(&encOutFile, "out", "", "Output .ecrypt file")
    encryptCmd.Flags().StringVar(&encPass, "pass", "", "Passphrase (Argon2id KDF)")
    encryptCmd.Flags().StringVar(&encKeyFile, "key-file", "", "32-byte Base64(URL) key file")
    encryptCmd.Flags().Uint32Var(&encArgonM, "argon-m", encArgonM, "Argon2 memory (KiB)")
    encryptCmd.Flags().Uint32Var(&encArgonT, "argon-t", encArgonT, "Argon2 iterations")
    encryptCmd.Flags().Uint8Var(&encArgonP, "argon-p", encArgonP, "Argon2 parallelism")
}