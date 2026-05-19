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

func TestBuilderCapturesStructuredErrorDetails(t *testing.T) {
	err := Builder().
		Status(400).
		Code(2001).
		SubErrorCode(7).
		MessageKey("inventory.name.required").
		Description("name is required").
		Param("field", "name").
		Field("name", "required", "missing name").
		Build()

	if err.Status != 400 || err.Code != 2001 || err.SubErrorCode != 7 {
		t.Fatalf("unexpected codes: %#v", err)
	}
	if err.MessageKey != "inventory.name.required" {
		t.Fatalf("unexpected message key %q", err.MessageKey)
	}
	if err.Params["field"] != "name" {
		t.Fatalf("unexpected params %#v", err.Params)
	}
	if len(err.Fields) != 1 || err.Fields[0].Field != "name" {
		t.Fatalf("unexpected fields %#v", err.Fields)
	}
}
