package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

func TestColumnizeTrimsBlankLines(t *testing.T) {
	input := "\n\nFirst line\nSecond line\nThird line\n\n"
	out := columnize(input, 20)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected output lines")
	}
	if isLineBlank(lines[0]) {
		t.Fatalf("expected leading blank lines trimmed")
	}
	if isLineBlank(lines[len(lines)-1]) {
		t.Fatalf("expected trailing blank lines trimmed")
	}
}

func TestColumnizeANSIWidth(t *testing.T) {
	bold := lipgloss.NewStyle().Bold(true).Render("LEFT")
	input := strings.Join([]string{
		bold,
		"RIGHT",
		"LEFT2",
		"RIGHT2",
	}, "\n")

	out := columnize(input, 6)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected output lines")
	}

	first := lines[0]
	if ansi.PrintableRuneWidth(first) < 6+columnGap {
		t.Fatalf("expected padded columns, got %q", first)
	}
	gap := strings.Repeat(" ", columnGap)
	plain := stripANSI(first)
	if len(plain) < 6+columnGap {
		t.Fatalf("expected padded columns, got %q", plain)
	}
	if plain[6:6+columnGap] != gap {
		t.Fatalf("expected column gap %q, got %q", gap, plain)
	}
}

func stripANSI(text string) string {
	var b strings.Builder
	inANSI := false
	for _, r := range text {
		switch {
		case inANSI:
			if ansi.IsTerminator(r) {
				inANSI = false
			}
			continue
		case r == ansi.Marker:
			inANSI = true
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
