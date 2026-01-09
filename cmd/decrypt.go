package cmd

import (
	"bytes"
	"ecrypto/archive"
	"ecrypto/crypto"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
    decInFile   string
    decOutDir   string
    decPass     string
    decKeyFile  string
)

var decryptCmd = &cobra.Command{
    Use:   "decrypt",
    Short: "Decrypt a .ecrypt container to a folder",
    Long: `Decrypt a .ecrypt container and extract to a folder.
Use the same passphrase or key file used during encryption.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if decInFile == "" || decOutDir == "" {
            return errors.New("--in and --out are required")
        }

        // Read container
        fmt.Fprintf(os.Stderr, "Reading container...\n")
        data, err := os.ReadFile(decInFile)
        if err != nil {
            return err
        }

        // Parse header
        h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
        if err != nil {
            return err
        }

        headerLen := crypto.HeaderSize()
        if len(data) < headerLen {
            return errors.New("file too small to contain header")
        }

        fmt.Fprintf(os.Stderr, "Container version: %d, KDF: %d\n", h.Version, h.KDF)

        var key []byte

        // Derive or load key
        if h.KDF == 1 {
            if decPass == "" {
                return errors.New("passphrase required (Argon2id)")
            }
            key = crypto.DeriveKeyArgon2id(decPass, h.Salt[:], h.ArgonM, h.ArgonT, h.ArgonP)
            fmt.Fprintf(os.Stderr, "Derived key from passphrase\n")
        } else {
            if decKeyFile == "" {
                return errors.New("key file required (raw key mode)")
            }
            var err error
            key, err = crypto.ReadKeyFromFile(decKeyFile)
            if err != nil {
                return err
            }
            fmt.Fprintf(os.Stderr, "Loaded 32-byte key from file\n")
        }

        // Decrypt
        fmt.Fprintf(os.Stderr, "Decrypting...\n")
        aad := data[:headerLen]
        ct := data[headerLen:]
        pt, err := crypto.DecryptAEAD(key, ct, aad, h.Nonce[:])
        if err != nil {
            return err
        }

        // Extract
        fmt.Fprintf(os.Stderr, "Extracting...\n")
        if err := os.MkdirAll(decOutDir, 0o755); err != nil {
            return err
        }
        if err := archive.UnzipTo(decOutDir, pt); err != nil {
            return err
        }

        fmt.Fprintf(os.Stderr, "✓ Decrypted to: %s\n", decOutDir)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(decryptCmd)
    decryptCmd.Flags().StringVar(&decInFile, "in", "", "Input .ecrypt file")
    decryptCmd.Flags().StringVar(&decOutDir, "out", "", "Output folder")
    decryptCmd.Flags().StringVar(&decPass, "pass", "", "Passphrase (Argon2id)")
    decryptCmd.Flags().StringVar(&decKeyFile, "key-file", "", "32-byte Base64(URL) key file")
}