package ui

import (
	"ecrypto/cmd"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// RunInteractiveMenu starts the interactive menu-driven interface
func RunInteractiveMenu() error {
	for {
		PrintMenu()

		options := []string{
			"[ENCRYPT]  Encrypt a Folder",
			"[DECRYPT]  Decrypt a File",
			"[KEYGEN]   Generate Encryption Key",
			"[INFO]     View Container Info",
			"[EXIT]     Quit Application",
		}

		choice := SelectOption("Main Menu", options)

		switch choice {
		case 0:
			if err := encryptInteractive(); err != nil {
				PrintError(err.Error())
				Pause()
			}
		case 1:
			if err := decryptInteractive(); err != nil {
				PrintError(err.Error())
				Pause()
			}
		case 2:
			if err := keygenInteractive(); err != nil {
				PrintError(err.Error())
				Pause()
			}
		case 3:
			if err := infoInteractive(); err != nil {
				PrintError(err.Error())
				Pause()
			}
		case 4:
			ClearScreen()
			fmt.Println()
			exitBox := lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 4).
				Align(lipgloss.Center).
				Width(60)
			fmt.Println(exitBox.Render("Thank you for using ECRYPTO!\\nYour data is secure."))
			fmt.Println()
			return nil
		}
	}
}

// encryptInteractive handles interactive encryption
func encryptInteractive() error {
	ClearScreen()
	headerBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSuccess).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorSuccess).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(60)
	fmt.Println(headerBox.Render("ENCRYPT FOLDER"))
	fmt.Println()

	inDir := SelectFolder("Enter folder path to encrypt")
	
	// Generate smart default with absolute path
	folderName := filepath.Base(inDir)
	defaultOut := filepath.Join(filepath.Dir(inDir), folderName+".ecrypt")
	
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Render("💡 Tip: Specify full path like D:\\backup\\myfiles.ecrypt or just filename for current dir"))
	fmt.Println()
	
	outFile := PromptUser("Output file path (full path or filename)", defaultOut)
	
	// If user provided relative path or just filename, make it absolute
	if !filepath.IsAbs(outFile) {
		// Check if it's just a filename or relative path
		if filepath.Dir(outFile) == "." {
			// Just a filename - ask where to save
			fmt.Println()
			saveDir := PromptUser("Save location (directory)", filepath.Dir(inDir))
			outFile = filepath.Join(saveDir, outFile)
		} else {
			// Relative path - make absolute from current dir
			absPath, err := filepath.Abs(outFile)
			if err == nil {
				outFile = absPath
			}
		}
	}
	
	// Validate output is not a directory
	if info, err := os.Stat(outFile); err == nil && info.IsDir() {
		PrintError("Output path is a directory! Please specify a file path ending with .ecrypt")
		Pause()
		return nil
	}

	fmt.Println()
	keyModeOpts := []string{"Passphrase (easier to remember)", "Random Key (more secure)"}
	keyMode := SelectOption("Choose key method", keyModeOpts)

	var pass, keyFile string
	if keyMode == 0 {
		pass = PromptPassphrase("Enter passphrase")
		confirm := PromptPassphrase("Confirm passphrase")
		if pass != confirm {
			return errors.New("passphrases do not match")
		}
	} else {
		// Random key mode - generate or use existing
		fmt.Println()
		keyActionOpts := []string{"Generate new key", "Use existing key file"}
		keyAction := SelectOption("Key file action", keyActionOpts)
		
		if keyAction == 0 {
			// Generate new key
			key, err := cmd.GenerateKey()
			if err != nil {
				return err
			}
			
			fmt.Println()
			fmt.Println(lipgloss.NewStyle().
				Foreground(lipgloss.Color("11")).
				Bold(true).
				Render("Generated Key (SAVE THIS!):"))
			fmt.Println()
			fmt.Println(lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSecondary).
				Render(key))
			fmt.Println()
			
			defaultKeyFile := filepath.Join(filepath.Dir(outFile), "encryption_key.txt")
			keyFile = PromptUser("Save key to file", defaultKeyFile)
			
			// Make absolute
			if !filepath.IsAbs(keyFile) {
				absPath, err := filepath.Abs(keyFile)
				if err == nil {
					keyFile = absPath
				}
			}
			
			if err := os.WriteFile(keyFile, []byte(key), 0o600); err != nil {
				return err
			}
			PrintSuccess(fmt.Sprintf("Key saved to: %s", keyFile))
		} else {
			// Use existing key file
			keyFile = SelectFileOrSkip("Select existing key file (or type path)")
			if keyFile == "" {
				PrintError("No key file selected.")
				Pause()
				return nil
			}
		}
	}

	PrintInfo("Encrypting your folder...")

	if keyMode == 0 {
		if err := cmd.EncryptWithPassphrase(inDir, outFile, pass); err != nil {
			return err
		}
	} else {
		if err := cmd.EncryptWithKeyFile(inDir, outFile, keyFile); err != nil {
			return err
		}
	}

	PrintSuccess(fmt.Sprintf("Folder encrypted successfully!\nOutput: %s", outFile))
	Pause()
	return nil
}

// decryptInteractive handles interactive decryption
func decryptInteractive() error {
	ClearScreen()
	headerBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(60)
	fmt.Println(headerBox.Render("DECRYPT FILE"))
	fmt.Println()

	inFile := SelectFile("Select .ecrypt file to decrypt")
	
	// Generate smart default based on file location
	defaultOutDir := filepath.Join(filepath.Dir(inFile), "restored")
	
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Render("💡 Tip: Specify full path like D:\\restored\\myfiles or relative path"))
	fmt.Println()
	
	outDir := PromptUser("Output folder path (where to extract)", defaultOutDir)
	
	// Make absolute if relative
	if !filepath.IsAbs(outDir) {
		absPath, err := filepath.Abs(outDir)
		if err == nil {
			outDir = absPath
		}
	}

	PrintInfo("Reading container...")
	// Show container info briefly
	if err := cmd.InfoPrint(inFile); err != nil {
		return err
	}

	fmt.Println()
	keyModeOpts := []string{"Use Passphrase", "Use Key File"}
	keyMode := SelectOption("Choose decryption method", keyModeOpts)

	var pass, keyFile string
	if keyMode == 0 {
		pass = PromptPassphrase("Enter passphrase")
	} else {
		keyFile = SelectFile("Select key file")
	}

	if !ConfirmAction("Proceed with decryption?") {
		PrintWarning("Decryption cancelled.")
		Pause()
		return nil
	}

	PrintInfo("Decrypting your file...")

	if keyMode == 0 {
		if err := cmd.DecryptWithPassphrase(inFile, outDir, pass); err != nil {
			return err
		}
	} else {
		if err := cmd.DecryptWithKeyFile(inFile, outDir, keyFile); err != nil {
			return err
		}
	}

	PrintSuccess(fmt.Sprintf("File decrypted successfully!\nOutput: %s", outDir))
	Pause()
	return nil
}

// keygenInteractive handles interactive key generation
func keygenInteractive() error {
	ClearScreen()
	headerBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorWarning).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(60)
	fmt.Println(headerBox.Render("GENERATE ENCRYPTION KEY"))
	fmt.Println()

	PrintInfo("Generating random 32-byte encryption key...")

	key, err := cmd.GenerateKey()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		Render("Your Key (save this safely):"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Render(key))
	fmt.Println()

	saveOpt := []string{"Save to file", "Just copy (skip saving)"}
	choice := SelectOption("Save key?", saveOpt)

	if choice == 0 {
		outFile := PromptUser("Enter key file path", "mykey.txt")
		if err := os.WriteFile(outFile, []byte(key), 0o600); err != nil {
			return err
		}
		PrintSuccess(fmt.Sprintf("Key saved to: %s", outFile))
	} else {
		PrintWarning("Remember to save this key somewhere safe!")
	}

	Pause()
	return nil
}

// infoInteractive handles interactive container info
func infoInteractive() error {
	ClearScreen()
	headerBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorSecondary).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(60)
	fmt.Println(headerBox.Render("CONTAINER INFORMATION"))
	fmt.Println()

	inFile := SelectFile("Select .ecrypt file")

	PrintInfo("Reading container info...")
	fmt.Println()

	if err := cmd.InfoPrint(inFile); err != nil {
		return err
	}

	Pause()
	return nil
}
