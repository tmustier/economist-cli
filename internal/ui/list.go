package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ListItem represents a row with an optional right-aligned column.
type ListItem struct {
	Title    string
	Subtitle string
	Right    string
}

// ListStyles controls how list rows are styled.
type ListStyles struct {
	Title         lipgloss.Style
	Subtitle      lipgloss.Style
	Selected      lipgloss.Style
	Right         lipgloss.Style
	RightSelected lipgloss.Style
}

// ListOptions configures list rendering.
type ListOptions struct {
	Width            int
	PrefixWidth      int
	RightColumnWidth int
	TitleLines       int
	SubtitleLines    int
	ItemGapLines     int
	SelectedIndex    int
	Start            int
	End              int
	Prefix           func(index int) string
}

// VisibleRange returns the start/end indices for a cursor within a max-visible window.
func VisibleRange(cursor, maxVisible, total int) (int, int) {
	if total <= 0 || maxVisible <= 0 {
		return 0, 0
	}
	if maxVisible > total {
		maxVisible = total
	}
	start := 0
	if cursor >= maxVisible {
		start = cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > total {
		end = total
	}
	return start, end
}

type listLayout struct {
	Width       int
	PrefixWidth int
	TitleWidth  int
	RightWidth  int
}

func newListLayout(width, prefixWidth, rightWidth int) listLayout {
	if rightWidth < 0 {
		rightWidth = 0
	}
	titleWidth := width - prefixWidth - rightWidth
	minWidth := prefixWidth + rightWidth + MinTitleWidth
	if titleWidth < MinTitleWidth && width >= minWidth {
		titleWidth = MinTitleWidth
	}
	if titleWidth < 1 {
		titleWidth = 1
	}
	return listLayout{
		Width:       width,
		PrefixWidth: prefixWidth,
		TitleWidth:  titleWidth,
		RightWidth:  rightWidth,
	}
}

// RenderList renders list items with wrapping and optional right column.
func RenderList(items []ListItem, opts ListOptions, styles ListStyles) string {
	if opts.Start < 0 {
		opts.Start = 0
	}
	if opts.End > len(items) {
		opts.End = len(items)
	}
	if opts.End < opts.Start {
		opts.End = opts.Start
	}
	layout := newListLayout(opts.Width, opts.PrefixWidth, opts.RightColumnWidth)
	prefixPad := strings.Repeat(" ", layout.PrefixWidth)
	gapLines := opts.ItemGapLines
	if gapLines < 0 {
		gapLines = 0
	}

	var b strings.Builder
	for i := opts.Start; i < opts.End; i++ {
		item := items[i]
		lineStyle := styles.Title
		rightStyle := styles.Right
		if i == opts.SelectedIndex {
			lineStyle = styles.Selected
			rightStyle = styles.RightSelected
		}

		prefix := ""
		if opts.Prefix != nil {
			prefix = opts.Prefix(i)
		}

		titleLines := LimitLines(WrapLines(item.Title, layout.TitleWidth), opts.TitleLines, layout.TitleWidth)
		if len(titleLines) == 0 {
			titleLines = []string{""}
		}
		for lineIdx, line := range titleLines {
			if lineIdx == 0 {
				paddedTitle := fmt.Sprintf("%-*s", layout.TitleWidth, line)
				b.WriteString(prefix)
				b.WriteString(lineStyle.Render(paddedTitle))
				if layout.RightWidth > 0 {
					rightColumn := fmt.Sprintf("%*s", layout.RightWidth, item.Right)
					b.WriteString(rightStyle.Render(rightColumn))
				}
				b.WriteString("\n")
				continue
			}
			if line == "" {
				b.WriteString(prefixPad + "\n")
				continue
			}
			b.WriteString(fmt.Sprintf("%s%s\n", prefixPad, lineStyle.Render(line)))
		}

		subtitleLines := LimitLines(WrapLines(item.Subtitle, layout.TitleWidth), opts.SubtitleLines, layout.TitleWidth)
		for _, line := range subtitleLines {
			if line == "" {
				b.WriteString(prefixPad + "\n")
				continue
			}
			b.WriteString(fmt.Sprintf("%s%s\n", prefixPad, styles.Subtitle.Render(line)))
		}

		if gapLines > 0 {
			b.WriteString(strings.Repeat("\n", gapLines))
		}
	}

	return b.String()
}
