package ui

import (
	"ecrypto/cmd"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RunInteractiveMenu starts the interactive menu-driven interface
func RunInteractiveMenu() error {
	for {
		PrintMenu()

		options := []string{
			"[ENCRYPT]  Encrypt a Folder/File",
			"[DECRYPT]  Decrypt a Folder/File",
			"[KEYGEN]   Generate Encryption Key",
			"[INFO]     View Container Info",
			"[UNDO]     Undo Recent Operation",
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
			if err := undoInteractive(); err != nil {
				PrintError(err.Error())
				Pause()
			}
		case 5:
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
	fmt.Println(headerBox.Render("🔒 ENCRYPT"))
	fmt.Println()
	PrintInfo("Let's encrypt your data securely")

	// Step 0: Choose between file or folder
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("\nStep 1: What do you want to encrypt?"))
	
	encryptOpts := []string{
		"📁 Folder (recommended for multiple files)",
		"📄 Single File",
	}
	encryptType, detectedPath, isPathProvided := SelectOptionOrPath("Encryption Type", encryptOpts)

	var inPath string
	var isFolder bool
	
	if isPathProvided {
		// User pasted a path directly
		inPath = detectedPath
		isFolder = (encryptType == 0)
		if isFolder {
			fmt.Println(HelpStyle.Render("✓ Detected folder: " + filepath.Base(inPath)))
		} else {
			fmt.Println(HelpStyle.Render("✓ Detected file: " + filepath.Base(inPath)))
		}
	} else if encryptType == 0 {
		// Encrypt folder
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorDark).
			Italic(true).
			Render("Tip: Copy-paste or type the full path to your folder"))
		inPath = SelectFolderEnhanced("Select folder to encrypt")
		isFolder = true
	} else {
		// Encrypt file
		fmt.Println(lipgloss.NewStyle().
			Foreground(ColorDark).
			Italic(true).
			Render("Tip: Copy-paste or type the full path to your file"))
		inPath = SelectFileEnhanced("Select file to encrypt")
		isFolder = false
	}

	if inPath == "" {
		return nil
	}

	// Show source info
	var size int64
	var fileCount int
	var displayName string
	
	if isFolder {
		displayName = filepath.Base(inPath)
		size, fileCount, _ = CalculateFolderSize(inPath)
		fmt.Println(HelpStyle.Render(fmt.Sprintf("📁 Folder: %s | 📄 Files: %d | 💾 Size: %s", 
			displayName, fileCount, FormatBytes(size))))
	} else {
		displayName = filepath.Base(inPath)
		info, err := os.Stat(inPath)
		if err == nil {
			size = info.Size()
			fileCount = 1
		}
		fmt.Println(HelpStyle.Render(fmt.Sprintf("📄 File: %s | 💾 Size: %s", 
			displayName, FormatBytes(size))))
	}
	
	// Step 2: Output location
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 2: Choose Output Location"))
	
	defaultOut := filepath.Join(filepath.Dir(inPath), displayName+".ecrypt")
	fmt.Println(HelpStyle.Render(fmt.Sprintf("Default: %s", defaultOut)))
	
	outFile := strings.Trim(PromptUser("Output file path (press Enter to use default)", defaultOut), "\"")
	
	// Make absolute path
	if !filepath.IsAbs(outFile) {
		if filepath.Dir(outFile) == "." {
			saveDir := strings.Trim(PromptUser("Save location (directory)", filepath.Dir(inPath)), "\"")
			outFile = filepath.Join(saveDir, outFile)
		} else {
			absPath, err := filepath.Abs(outFile)
			if err == nil {
				outFile = absPath
			}
		}
	}
	
	// Validate output
	if info, err := os.Stat(outFile); err == nil && info.IsDir() {
		PrintError("Output path is a directory! Please specify a file path ending with .ecrypt")
		Pause()
		return nil
	}

	// Step 3: Key method
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 3: Choose Security Method"))
	
	keyModeOpts := []string{
		"💡 Passphrase (easier to remember)",
		"🔑 Random Key (maximum security)",
	}
	keyMode := SelectOption("Security Method", keyModeOpts)

	var pass, keyFile string
	if keyMode == 0 {
		// Passphrase mode
		fmt.Println()
		fmt.Println(HelpStyle.Render("💡 Use a strong passphrase (16+ characters with symbols)"))
		pass = PromptPassphrase("Enter passphrase")
		confirm := PromptPassphrase("Confirm passphrase")
		if pass != confirm {
			PrintError("Passphrases do not match!")
			Pause()
			return nil
		}
	} else {
		// Key file mode
		fmt.Println()
		keyActionOpts := []string{"Generate new key (Recommended)", "Use existing key file"}
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
				Render("⚠ IMPORTANT: Save this key file in a secure location!"))
			fmt.Println()
			fmt.Println(lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSecondary).
				Render(key))
			fmt.Println()
			
			defaultKeyFile := filepath.Join(filepath.Dir(outFile), "encryption_key.txt")
			keyFile = strings.Trim(PromptUser("Save key to file", defaultKeyFile), "\"")
			
			// Make absolute
			if !filepath.IsAbs(keyFile) {
				absPath, err := filepath.Abs(keyFile)
				if err == nil {
					keyFile = absPath
				}
			}
			
			if err := os.WriteFile(keyFile, []byte(key), 0o600); err != nil {
				PrintError(fmt.Sprintf("Failed to save key: %v", err))
				Pause()
				return nil
			}
			PrintSuccess(fmt.Sprintf("Key saved to: %s", keyFile))
		} else {
			// Use existing key file
			keyFile = SelectFileEnhanced("Select existing key file")
			if keyFile == "" {
				PrintError("No key file selected.")
				Pause()
				return nil
			}
		}
	}

	// Confirmation before encryption
	fmt.Println()
	confirmMsg := ""
	if isFolder {
		confirmMsg = fmt.Sprintf("📁 Encrypting folder with %d file(s) (%s)\n▶ Output: %s", fileCount, FormatBytes(size), filepath.Base(outFile))
	} else {
		confirmMsg = fmt.Sprintf("📄 Encrypting file (%s)\n▶ Output: %s", FormatBytes(size), filepath.Base(outFile))
	}
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDark).
		Render(confirmMsg))
	
	if !ConfirmAction("Ready to encrypt?") {
		Pause()
		return nil
	}

	PrintInfo("Encrypting your data...")
	
	// Create and start progress tracker
	var progress *ProgressTracker
	if fileCount > 0 {
		progress = NewProgressTracker("Encrypting", fileCount)
		progress.Start()
	}

	var encErr error
	if keyMode == 0 {
		if isFolder {
			encErr = cmd.EncryptWithPassphrase(inPath, outFile, pass, func(filename string) {
				if progress != nil {
					progress.Update(filename)
				}
			})
		} else {
			encErr = cmd.EncryptFileWithPassphrase(inPath, outFile, pass, func(filename string) {
				if progress != nil {
					progress.Update(filename)
				}
			})
		}
	} else {
		if isFolder {
			encErr = cmd.EncryptWithKeyFile(inPath, outFile, keyFile, func(filename string) {
				if progress != nil {
					progress.Update(filename)
				}
			})
		} else {
			encErr = cmd.EncryptFileWithKeyFile(inPath, outFile, keyFile, func(filename string) {
				if progress != nil {
					progress.Update(filename)
				}
			})
		}
	}
	
	if progress != nil {
		progress.Stop()
		fmt.Println()
	}
	
	if encErr != nil {
		PrintError(fmt.Sprintf("Encryption failed: %v", encErr))
		// Log failed operation
		history := NewOperationHistory()
		history.AddOperation(Operation{
			Type:       "encrypt",
			SourcePath: inPath,
			OutputPath: outFile,
			Size:       size,
			FileCount:  fileCount,
			KeyMethod:  map[int]string{0: "passphrase", 1: "keyfile"}[keyMode],
			KeyPath:    keyFile,
			Status:     "failed",
			Error:      encErr.Error(),
		})
		Pause()
		return nil
	}

	// Log successful operation
	history := NewOperationHistory()
	history.AddOperation(Operation{
		Type:       "encrypt",
		SourcePath: inPath,
		OutputPath: outFile,
		Size:       size,
		FileCount:  fileCount,
		KeyMethod:  map[int]string{0: "passphrase", 1: "keyfile"}[keyMode],
		KeyPath:    keyFile,
		Status:     "success",
	})

	successMsg := ""
	if isFolder {
		successMsg = fmt.Sprintf("Folder encrypted successfully!\n▶ Output: %s", outFile)
	} else {
		successMsg = fmt.Sprintf("File encrypted successfully!\n▶ Output: %s", outFile)
	}
	PrintSuccess(successMsg)
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
	fmt.Println(headerBox.Render("🔓 DECRYPT FILE"))
	fmt.Println()
	PrintInfo("Let's decrypt your .ecrypt file in a few steps")

	// Step 1: Select file
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 1: Select .ecrypt File"))
	fmt.Println(HelpStyle.Render("Tip: Copy-paste the full path to your .ecrypt file"))
	
	inFile := SelectFileEnhanced("Select .ecrypt file to decrypt")
	if inFile == "" {
		return nil
	}

	// Show file info
	fileInfo, _ := os.Stat(inFile)
	fmt.Println(HelpStyle.Render(fmt.Sprintf("📦 File: %s | 💾 Size: %s", 
		filepath.Base(inFile), GetFileSize(fileInfo.Size()))))
	
	// Step 2: Output location
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 2: Choose Output Location"))
	
	defaultOutDir := filepath.Join(filepath.Dir(inFile), "restored")
	fmt.Println(HelpStyle.Render(fmt.Sprintf("Default: %s", defaultOutDir)))
	
	outDir := strings.Trim(PromptUser("Output folder path (press Enter to use default)", defaultOutDir), "\"")
	
	// Make absolute if relative
	if !filepath.IsAbs(outDir) {
		absPath, err := filepath.Abs(outDir)
		if err == nil {
			outDir = absPath
		}
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 3: Read Container"))
	
	PrintInfo("Reading container info...")
	if err := cmd.InfoPrint(inFile); err != nil {
		PrintError(fmt.Sprintf("Could not read container: %v", err))
		Pause()
		return nil
	}

	// Step 4: Choose decryption method
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 4: Choose Decryption Method"))
	
	keyModeOpts := []string{"🔐 Use Passphrase", "🔑 Use Key File"}
	keyMode := SelectOption("Decryption Method", keyModeOpts)

	var pass, keyFile string
	if keyMode == 0 {
		pass = PromptPassphrase("Enter passphrase")
	} else {
		keyFile = SelectFileEnhanced("Select key file")
		if keyFile == "" {
			PrintError("No key file selected.")
			Pause()
			return nil
		}
	}

	// Final confirmation
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDark).
		Render(fmt.Sprintf("📦 Decrypting: %s\n▶ Output: %s", filepath.Base(inFile), filepath.Base(outDir))))
	
	if !ConfirmAction("Ready to decrypt?") {
		Pause()
		return nil
	}

	PrintInfo("Decrypting your file...")
	
	stopSpinner := ShowSimpleProgress("Decrypting")
	
	var decErr error
	if keyMode == 0 {
		decErr = cmd.DecryptWithPassphrase(inFile, outDir, pass, nil)
	} else {
		decErr = cmd.DecryptWithKeyFile(inFile, outDir, keyFile, nil)
	}
	
	stopSpinner()
	fmt.Println()
	
	if decErr != nil {
		// Better error handling for common decryption errors
		errMsg := decErr.Error()
		if strings.Contains(errMsg, "authentication tag") {
			PrintError("Authentication failed! This usually means:\n  • Wrong passphrase or key file\n  • File is corrupted\n\nDouble-check your passphrase/key and try again.")
		} else {
			PrintError(fmt.Sprintf("Decryption failed: %v", decErr))
		}
		Pause()
		return nil
	}

	PrintSuccess(fmt.Sprintf("File decrypted successfully!\n▶ Output: %s", outDir))
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
	fmt.Println(headerBox.Render("🔑 GENERATE ENCRYPTION KEY"))
	fmt.Println()

	PrintInfo("Generating a random 32-byte encryption key for maximum security...")
	fmt.Println()

	key, err := cmd.GenerateKey()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to generate key: %v", err))
		Pause()
		return nil
	}

	fmt.Println(lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		Render("⚠ YOUR ENCRYPTION KEY (SAVE THIS SAFELY):"))
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Render(key))
	fmt.Println()
	fmt.Println(HelpStyle.Render("💡 This key is 256 bits of random data. Without it, you cannot decrypt your files."))
	fmt.Println()

	saveOpt := []string{
		"💾 Save to file (Recommended)",
		"📋 Just copy (I'll paste elsewhere)",
	}
	choice := SelectOption("How to save this key?", saveOpt)

	if choice == 0 {
		outFile := strings.Trim(PromptUser("Enter key file path", "encryption_key.txt"), "\"")
		if err := os.WriteFile(outFile, []byte(key), 0o600); err != nil {
			PrintError(fmt.Sprintf("Failed to save key: %v", err))
			Pause()
			return nil
		}
		PrintSuccess(fmt.Sprintf("Key saved to: %s", outFile))
		fmt.Println(HelpStyle.Render("📌 Remember: Keep this file in a secure location (password manager, encrypted USB, etc.)"))
	} else {
		PrintWarning("Key is in clipboard. Please save it somewhere secure!")
		fmt.Println(HelpStyle.Render("💡 Store in: Password manager, encrypted USB, or printed copy in safe"))
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
	fmt.Println(headerBox.Render("📊 CONTAINER INFORMATION"))
	fmt.Println()
	PrintInfo("View details about an encrypted container (no decryption needed)")

	fmt.Println()
	inFile := SelectFileEnhanced("Select .ecrypt file to inspect")
	if inFile == "" {
		return nil
	}

	PrintInfo("Reading container metadata...")
	fmt.Println()

	if err := cmd.InfoPrint(inFile); err != nil {
		PrintError(fmt.Sprintf("Failed to read container: %v", err))
		Pause()
		return nil
	}

	Pause()
	return nil
}

// undoInteractive handles undoing recent operations
func undoInteractive() error {
	ClearScreen()
	headerBox := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorWarning).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(60)
	fmt.Println(headerBox.Render("↶ UNDO RECENT OPERATION"))
	fmt.Println()
	PrintInfo("Review and undo your recent encryption operations")

	// Load operation history
	history := NewOperationHistory()
	recentOps := history.GetRecentOperations(10)

	if len(recentOps) == 0 {
		PrintWarning("No operations to undo yet. Encrypt a folder first!")
		Pause()
		return nil
	}

	// Filter to only undoable operations (successful encryptions)
	var undoableOps []Operation
	for _, op := range recentOps {
		if op.IsUndoable() {
			undoableOps = append(undoableOps, op)
		}
	}

	if len(undoableOps) == 0 {
		PrintWarning("No undoable operations found. (Encrypted files may have been deleted)")
		Pause()
		return nil
	}

	// Display recent operations
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Recent Operations:"))
	fmt.Println()

	options := []string{}
	for i, op := range undoableOps {
		size := FormatBytes(op.Size)
		timestamp := op.FormatTime()
		displayOption := fmt.Sprintf("[%d] %s | %s | %d files | %s",
			i+1,
			op.SourcePath,
			size,
			op.FileCount,
			timestamp,
		)
		if len(displayOption) > 80 {
			displayOption = displayOption[:77] + "..."
		}
		options = append(options, displayOption)
	}
	options = append(options, "[CANCEL]  Go Back")

	choice := SelectOption("Select operation to undo", options)

	if choice >= len(undoableOps) {
		return nil // User selected cancel
	}

	selectedOp := undoableOps[choice]

	// Show confirmation with details
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorDark).
		Italic(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDark).
		Render(fmt.Sprintf("Original folder: %s\n🔒 Encrypted file: %s\n📊 Files: %d | 💾 Size: %s\n📅 Date: %s",
			selectedOp.SourcePath,
			filepath.Base(selectedOp.OutputPath),
			selectedOp.FileCount,
			selectedOp.FormatSize(),
			selectedOp.FormatTime())))

	fmt.Println()
	fmt.Println(HelpStyle.Render("⚠ Undoing will decrypt the file. You'll need the original passphrase/key."))

	if !ConfirmAction("Decrypt to restore original folder?") {
		Pause()
		return nil
	}

	// Ask for output location
	fmt.Println()
	restoredName := filepath.Base(selectedOp.SourcePath) + "_restored"
	defaultRestoreDir := filepath.Join(filepath.Dir(selectedOp.OutputPath), restoredName)

	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 1: Choose Output Location"))
	fmt.Println(HelpStyle.Render(fmt.Sprintf("Default: %s", defaultRestoreDir)))

	restoreDir := strings.Trim(PromptUser("Restore to folder (press Enter for default)", defaultRestoreDir), "\"")

	if !filepath.IsAbs(restoreDir) {
		absPath, err := filepath.Abs(restoreDir)
		if err == nil {
			restoreDir = absPath
		}
	}

	// Ask for decryption method
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Step 2: Choose Decryption Method"))

	keyModeOpts := []string{"🔐 Use Passphrase", "🔑 Use Key File"}
	keyMode := SelectOption("Decryption Method", keyModeOpts)

	var pass, keyFile string
	if keyMode == 0 {
		pass = PromptPassphrase("Enter original passphrase")
	} else {
		keyFile = SelectFileEnhanced("Select key file")
		if keyFile == "" {
			PrintError("No key file selected.")
			Pause()
			return nil
		}
	}

	// Final confirmation
	fmt.Println()
	if !ConfirmAction("Ready to decrypt and restore?") {
		Pause()
		return nil
	}

	PrintInfo("Decrypting and restoring your folder...")
	stopSpinner := ShowSimpleProgress("Restoring")

	var decErr error
	if keyMode == 0 {
		decErr = cmd.DecryptWithPassphrase(selectedOp.OutputPath, restoreDir, pass, nil)
	} else {
		decErr = cmd.DecryptWithKeyFile(selectedOp.OutputPath, restoreDir, keyFile, nil)
	}

	stopSpinner()
	fmt.Println()

	if decErr != nil {
		errMsg := decErr.Error()
		if strings.Contains(errMsg, "authentication tag") {
			PrintError("Authentication failed! Wrong passphrase or key file.")
		} else {
			PrintError(fmt.Sprintf("Restoration failed: %v", decErr))
		}
		Pause()
		return nil
	}

	PrintSuccess(fmt.Sprintf("Folder restored successfully!\n▶ Output: %s", restoreDir))
	Pause()
	return nil
}
