package ui

// LayoutSpec describes the fixed layout lines reserved outside the main content area.
type LayoutSpec struct {
	HeaderLines     int
	FooterLines     int
	FooterPadding   int
	FooterGapLines  int
	MinVisibleLines int
}

// ReservedLines returns the total number of lines reserved by the layout.
func (spec LayoutSpec) ReservedLines() int {
	return spec.HeaderLines + spec.FooterLines + spec.FooterPadding + spec.FooterGapLines
}

// VisibleLines returns the visible content lines for the given height.
func (spec LayoutSpec) VisibleLines(height int) int {
	visible := height - spec.ReservedLines()
	if visible < spec.MinVisibleLines {
		visible = spec.MinVisibleLines
	}
	return visible
}

// PageSize returns how many items fit in the view given an item height.
func PageSize(height, itemHeight int, spec LayoutSpec) int {
	visible := spec.VisibleLines(height)
	page := visible / itemHeight
	if page < 1 {
		page = 1
	}
	return page
}
