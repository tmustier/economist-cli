package demo

import (
	"strings"
	"testing"
)

func TestDemoSource(t *testing.T) {
	source := NewSource()
	title, items, err := source.Section("")
	if err != nil {
		t.Fatalf("section: %v", err)
	}
	if title == "" {
		t.Fatalf("expected title")
	}
	if len(items) == 0 {
		t.Fatalf("expected items")
	}

	art, err := source.Article(items[0].Link)
	if err != nil {
		t.Fatalf("article: %v", err)
	}
	if art.Title == "" || art.Content == "" {
		t.Fatalf("expected article content")
	}
	if !strings.Contains(strings.ToLower(art.Content), "destroyers") {
		snippet := art.Content
		if len(snippet) > 80 {
			snippet = snippet[:80]
		}
		t.Fatalf("expected fixture content, got %q", snippet)
	}
	if !strings.HasSuffix(strings.TrimSpace(art.Content), "■") {
		t.Fatalf("expected trailing marker")
	}
}
