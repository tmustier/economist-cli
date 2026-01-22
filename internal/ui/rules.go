package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

const (
	RuleLight  = "─"
	RuleHeavy  = "━"
	RuleDouble = "═"
	RuleDotted = "┄"
)

func DrawRule(width int, char string, style lipgloss.Style) string {
	if width <= 0 {
		return ""
	}
	return style.Render(strings.Repeat(char, width))
}

func AccentRule(width int, styles Styles) string {
	return DrawRule(width, RuleHeavy, styles.RuleAccent)
}

func SectionRule(width int, styles Styles) string {
	return DrawRule(width, RuleLight, styles.Rule)
}

func SectionBadge(section string, styles Styles) string {
	return styles.Overline.Render(strings.ToUpper(section))
}

func IsRuleLine(line string) bool {
	stripped := strings.TrimSpace(StripANSI(line))
	if stripped == "" {
		return false
	}
	for _, r := range stripped {
		switch string(r) {
		case RuleLight, RuleHeavy, RuleDouble, RuleDotted:
			continue
		default:
			return false
		}
	}
	return true
}

func StripANSI(text string) string {
	var b strings.Builder
	inSeq := false
	for _, r := range text {
		if inSeq {
			if ansi.IsTerminator(r) {
				inSeq = false
			}
			continue
		}
		if r == ansi.Marker {
			inSeq = true
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
