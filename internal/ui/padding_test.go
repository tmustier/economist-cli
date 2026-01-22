package ui

import "testing"

func TestPadBlockRight(t *testing.T) {
	input := "hi\nok"
	out := PadBlockRight(input, 4)
	expected := "hi  \nok  "
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestPadBlockRightPreservesTrailingNewline(t *testing.T) {
	input := "hi\n"
	out := PadBlockRight(input, 4)
	expected := "hi  \n"
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}
