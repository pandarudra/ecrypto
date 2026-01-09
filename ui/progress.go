package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ProgressTracker tracks operation progress
type ProgressTracker struct {
	Total       int
	Current     int
	CurrentFile string
	StartTime   time.Time
	Operation   string
	done        chan bool
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(operation string, total int) *ProgressTracker {
	return &ProgressTracker{
		Operation: operation,
		Total:     total,
		Current:   0,
		StartTime: time.Now(),
		done:      make(chan bool),
	}
}

// Start begins showing progress
func (p *ProgressTracker) Start() {
	go p.animate()
}

// Update increments progress
func (p *ProgressTracker) Update(filename string) {
	p.Current++
	p.CurrentFile = filename
}

// Stop stops the progress display
func (p *ProgressTracker) Stop() {
	p.done <- true
	time.Sleep(50 * time.Millisecond) // Let final update show
	fmt.Print("\r" + strings.Repeat(" ", 100) + "\r") // Clear line
}

// animate shows animated progress
func (p *ProgressTracker) animate() {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			p.render(spinner[idx%len(spinner)])
			idx++
		}
	}
}

// render draws the progress bar
func (p *ProgressTracker) render(spinner string) {
	percent := 0
	if p.Total > 0 {
		percent = (p.Current * 100) / p.Total
	}

	// Progress bar
	barWidth := 30
	filled := (percent * barWidth) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Elapsed time
	elapsed := time.Since(p.StartTime)
	elapsedStr := fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)

	// File name truncation
	displayFile := p.CurrentFile
	if len(displayFile) > 40 {
		displayFile = "..." + displayFile[len(displayFile)-37:]
	}

	// Build status line
	status := fmt.Sprintf("%s %s [%s] %d%% (%d/%d files) | %s | %s",
		spinner,
		p.Operation,
		bar,
		percent,
		p.Current,
		p.Total,
		elapsedStr,
		displayFile,
	)

	// Style and print
	styledStatus := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render(status)

	fmt.Print("\r" + styledStatus)
}

// CalculateFolderSize returns size and file count
func CalculateFolderSize(path string) (int64, int, error) {
	var size int64
	var count int

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return nil
	})

	return size, count, err
}

// FormatBytes returns human-readable byte size
func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIdx := 0

	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}

	if unitIdx == 0 {
		return fmt.Sprintf("%d %s", int64(size), units[unitIdx])
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIdx])
}

// ShowSimpleProgress shows a simple spinner for unknown duration tasks
func ShowSimpleProgress(operation string) func() {
	done := make(chan bool)
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
				return
			case <-ticker.C:
				status := lipgloss.NewStyle().
					Foreground(ColorPrimary).
					Render(fmt.Sprintf("%s %s...", spinner[idx%len(spinner)], operation))
				fmt.Print("\r" + status)
				idx++
			}
		}
	}()

	// Return stop function
	return func() {
		done <- true
		time.Sleep(50 * time.Millisecond)
	}
}
