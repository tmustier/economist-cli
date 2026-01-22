package ui

import (
	"strings"
	"testing"

	"github.com/muesli/reflow/ansi"
)

func TestVisibleRange(t *testing.T) {
	start, end := VisibleRange(0, 3, 10)
	if start != 0 || end != 3 {
		t.Fatalf("expected range 0-3, got %d-%d", start, end)
	}

	start, end = VisibleRange(5, 3, 10)
	if start != 3 || end != 6 {
		t.Fatalf("expected range 3-6, got %d-%d", start, end)
	}
}

func TestRenderListAlignsRightColumn(t *testing.T) {
	items := []ListItem{{Title: "Hello", Right: "R"}}
	out := RenderList(items, ListOptions{
		Width:            20,
		PrefixWidth:      3,
		RightColumnWidth: 4,
		TitleLines:       1,
		SubtitleLines:    0,
		ItemGapLines:     0,
		Start:            0,
		End:              1,
		Prefix: func(int) string {
			return "1. "
		},
	}, ListStyles{})

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatalf("expected output line, got none")
	}
	if width := ansi.PrintableRuneWidth(lines[0]); width != 20 {
		t.Fatalf("expected line width 20, got %d", width)
	}
}
