package cmd

import (
	"bytes"
	"ecrypto/crypto"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var infoFile string

var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Print .ecrypt container header (no decryption)",
    Long: `Read and display the header of a .ecrypt container without decrypting.
Shows magic, version, KDF type, and Argon2 parameters.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if infoFile == "" {
            return errors.New("--file is required")
        }

        data, err := os.ReadFile(infoFile)
        if err != nil {
            return err
        }

        h, err := crypto.DecodeHeaderV1(bytes.NewReader(data))
        if err != nil {
            return err
        }

        fmt.Printf("Container: %s\n", infoFile)
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
    },
}

func init() {
    rootCmd.AddCommand(infoCmd)
    infoCmd.Flags().StringVar(&infoFile, "file", "", "Path to .ecrypt container")
}