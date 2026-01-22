package ui

import "strings"

// PadBlockRight pads each line in a block to the target width, preserving ANSI codes.
func PadBlockRight(text string, width int) string {
	if width <= 0 || text == "" {
		return text
	}

	lines := strings.SplitAfter(text, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		hasNewline := strings.HasSuffix(line, "\n")
		if hasNewline {
			line = strings.TrimSuffix(line, "\n")
		}
		line = padRightANSI(line, width)
		if hasNewline {
			line += "\n"
		}
		lines[i] = line
	}

	return strings.Join(lines, "")
}
