package ui

import (
	"fmt"
	"strings"
)

// PrintBox prints a formatted box with title and content lines
func PrintBox(title string, lines ...string) {
	// Calculate max width needed
	maxWidth := len(title) + 6 // title + " " + padding
	for _, line := range lines {
		// Account for 2 space indentation
		if len(line)+2 > maxWidth {
			maxWidth = len(line) + 2
		}
	}

	// Top border
	dashes := strings.Repeat("─", maxWidth-len(title)-2)
	fmt.Printf("╭─ %s %s╮\n", title, dashes)

	// Content
	for _, line := range lines {
		fmt.Printf("  %s\n", line)
	}

	// Bottom border
	fmt.Printf("╰%s╯\n", strings.Repeat("─", maxWidth+1))
}

// PrintSuccess prints a success box
func PrintSuccess(lines ...string) {
	PrintBox("Success", lines...)
}

// PrintWarning prints a warning box
func PrintWarning(lines ...string) {
	PrintBox("Warning", lines...)
}

// PrintErrors prints an errors box
func PrintErrors(lines ...string) {
	PrintBox("Errors Encountered", lines...)
}
