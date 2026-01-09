package cmd

import (
	"crypto/rand"
	"ecrypto/crypto"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var keygenOutFile string

var keygenCmd = &cobra.Command{
    Use:   "keygen",
    Short: "Generate a random 32-byte encryption key",
    Long: `Generate a random 32-byte key and print it in Base64URL format.
Optionally save to a file with --out.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        key := make([]byte, crypto.KeySize())
        if _, err := rand.Read(key); err != nil {
            return err
        }

        encoded := base64.RawURLEncoding.EncodeToString(key)
        fmt.Println(encoded)

        if keygenOutFile != "" {
            if err := os.WriteFile(keygenOutFile, []byte(encoded), 0o600); err != nil {
                return err
            }
            fmt.Fprintf(os.Stderr, "✓ Key saved to: %s\n", keygenOutFile)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(keygenCmd)
    keygenCmd.Flags().StringVar(&keygenOutFile, "out", "", "Output key file (optional)")
}