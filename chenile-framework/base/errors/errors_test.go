package errors

import "testing"

func TestChenileErrorUsesDescriptionAsErrorMessage(t *testing.T) {
	err := New(400, 10, "bad input")

	if err.Error() != "bad input" {
		t.Fatalf("unexpected error message %q", err.Error())
	}
}

func TestChenileErrorHasFallbackMessage(t *testing.T) {
	err := New(404, 0, "")

	if err.Error() == "" {
		t.Fatal("expected fallback error message")
	}
}
