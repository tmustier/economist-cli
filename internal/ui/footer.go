package ui

import "strings"

// BuildFooter builds a footer block with a leading blank line and divider.
func BuildFooter(divider string, lines ...string) string {
	footerLines := []string{"", divider}
	for _, line := range lines {
		if line == "" {
			continue
		}
		footerLines = append(footerLines, line)
	}
	return strings.Join(footerLines, "\n")
}
