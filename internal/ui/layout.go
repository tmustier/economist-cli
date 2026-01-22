package ui

import "strings"

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func LayoutWithFooter(content, footer string, height, bottomPadding int) string {
	content = strings.TrimRight(content, "\n")
	footer = strings.TrimRight(footer, "\n")
	if footer == "" || height <= 0 {
		if content == "" {
			return footer
		}
		if footer == "" {
			return content
		}
		return content + "\n" + footer
	}

	contentLines := lineCount(content)
	footerLines := lineCount(footer)
	gap := height - contentLines - footerLines - bottomPadding
	if content != "" && gap < 1 {
		gap = 1
	}
	if content == "" && gap < 0 {
		gap = 0
	}

	var b strings.Builder
	if content != "" {
		b.WriteString(content)
	}
	if gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	b.WriteString(footer)
	if bottomPadding > 0 {
		b.WriteString(strings.Repeat("\n", bottomPadding))
	}
	return b.String()
}

func LayoutWithFooterDivider(content, divider, footer string, height, bottomPadding int) string {
	divider = strings.TrimRight(divider, "\n")
	if divider == "" {
		return LayoutWithFooter(content, footer, height, bottomPadding)
	}

	content = strings.TrimRight(content, "\n")
	footer = strings.TrimRight(footer, "\n")
	contentLines := lineCount(content)
	dividerLines := lineCount(divider)
	footerLines := lineCount(footer)

	gap := height - contentLines - dividerLines - footerLines - bottomPadding
	if gap < 0 {
		gap = 0
	}

	var b strings.Builder
	if content != "" {
		b.WriteString(content)
		b.WriteString("\n")
	}
	b.WriteString(divider)
	if footer != "" {
		b.WriteString("\n")
	}
	if gap > 0 {
		b.WriteString(strings.Repeat("\n", gap))
	}
	if footer != "" {
		b.WriteString(footer)
	}
	if bottomPadding > 0 {
		b.WriteString(strings.Repeat("\n", bottomPadding))
	}
	return b.String()
}

func lineCount(text string) int {
	if text == "" {
		return 0
	}
	return strings.Count(text, "\n") + 1
}
