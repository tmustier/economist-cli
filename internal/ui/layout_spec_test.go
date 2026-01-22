package ui

import "testing"

func TestLayoutSpecVisibleLinesMin(t *testing.T) {
	spec := LayoutSpec{
		HeaderLines:     2,
		FooterLines:     3,
		FooterPadding:   1,
		FooterGapLines:  0,
		MinVisibleLines: 5,
	}

	visible := spec.VisibleLines(6)
	if visible != 5 {
		t.Fatalf("expected min visible lines 5, got %d", visible)
	}
}

func TestPageSizeUsesVisibleLines(t *testing.T) {
	spec := LayoutSpec{MinVisibleLines: 1}
	page := PageSize(20, 5, spec)
	if page != 4 {
		t.Fatalf("expected page size 4, got %d", page)
	}
}
