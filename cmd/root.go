package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ecrypto",
	Short: "Encrypt and decrypt folders into secure .ecrypt containers",
    Long: `ecrypto encrypts entire folders into a single encrypted container using XChaCha20-Poly1305.
Decrypt with a passphrase or raw key. Filenames and metadata are protected.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {}