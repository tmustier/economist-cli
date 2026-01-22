package ui

import (
	"strings"
	"testing"

	"github.com/tmustier/economist-tui/internal/article"
)

func TestRenderArticleIndent(t *testing.T) {
	art := &article.Article{
		Overtitle: "Section",
		Title:     "Headline",
		Subtitle:  "Subhead",
		DateLine:  "Jan 1st 2024",
		Content:   "This is a paragraph that should wrap across multiple lines to verify indentation is consistent across wraps.",
		URL:       "https://example.com/test",
	}

	out, err := RenderArticle(art, ArticleRenderOptions{
		NoColor:   true,
		WrapWidth: 40,
	})
	if err != nil {
		t.Fatalf("render article: %v", err)
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	indent := strings.Repeat(" ", bodyIndent)
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, indent) {
			t.Fatalf("expected indent %q, got %q", indent, line)
		}
	}
}

func TestReflowArticleBodyColumns(t *testing.T) {
	styles := NewArticleStyles(true)
	base := strings.Repeat("word ", 80)
	opts := ArticleRenderOptions{
		NoColor:   true,
		WrapWidth: 80,
		TwoColumn: true,
	}
	layout := ResolveArticleLayoutWithContent(base, opts)
	if !layout.UseColumns {
		t.Fatalf("expected columns to be enabled")
	}

	body := ReflowArticleBodyWithLayout(base, styles, opts, layout)
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")
	indent := strings.Repeat(" ", layout.Indent)
	gap := strings.Repeat(" ", columnGap)

	for _, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, indent) {
			t.Fatalf("expected indent %q, got %q", indent, line)
		}
		trimmed := strings.TrimPrefix(line, indent)
		for col := 1; col < layout.ColumnCount; col++ {
			gapStart := col*layout.ColumnWidth + (col-1)*columnGap
			if len(trimmed) > gapStart {
				if len(trimmed) < gapStart+columnGap {
					t.Fatalf("expected column gap, got %q", trimmed)
				}
				if trimmed[gapStart:gapStart+columnGap] != gap {
					t.Fatalf("expected column gap %q, got %q", gap, trimmed)
				}
			}
		}
	}
}

func TestResolveArticleLayoutAddsColumnsWhenWide(t *testing.T) {
	opts := ArticleRenderOptions{
		NoColor:   true,
		TermWidth: 220,
		TwoColumn: true,
	}
	layout := ResolveArticleLayout(opts)
	if !layout.UseColumns {
		t.Fatalf("expected columns to be enabled")
	}
	if layout.ColumnCount < 3 {
		t.Fatalf("expected at least 3 columns, got %d", layout.ColumnCount)
	}
	if layout.ColumnWidth > maxColumnWidth {
		t.Fatalf("expected column width <= %d, got %d", maxColumnWidth, layout.ColumnWidth)
	}
}

func TestResolveArticleLayoutCapsColumnsByHeight(t *testing.T) {
	base := strings.Repeat("line\n", 10)
	opts := ArticleRenderOptions{
		NoColor:   true,
		TermWidth: 240,
		TwoColumn: true,
		Center:    true,
	}
	layout := ResolveArticleLayoutWithContent(base, opts)
	if !layout.UseColumns {
		t.Fatalf("expected columns to be enabled")
	}
	if layout.ColumnCount != 2 {
		t.Fatalf("expected columns to cap at 2, got %d", layout.ColumnCount)
	}
	if layout.ContentWidth >= layout.TermWidth {
		t.Fatalf("expected content width to be capped for centering")
	}
}

func TestHighlightTrailingMarker(t *testing.T) {
	styles := NewArticleStyles(false)
	input := "Ends with marker ■"
	output := HighlightTrailingMarker(input, styles)

	styled := styles.Title.Render("■")
	if !strings.Contains(output, styled) {
		t.Fatalf("expected styled marker, got %q", output)
	}

	replaced := strings.Replace(input, "■", styled, 1)
	if output != replaced {
		t.Fatalf("expected marker replacement, got %q", output)
	}
}
