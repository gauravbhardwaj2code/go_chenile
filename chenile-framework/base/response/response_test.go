package response

import (
	"net/http"
	"testing"
)

func TestSuccessBuildsGenericResponse(t *testing.T) {
	result := Success(map[string]string{"id": "1"})

	if !result.Success {
		t.Fatal("expected success response")
	}
	if result.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", result.Code)
	}
	if result.Payload == nil {
		t.Fatal("expected payload")
	}
	if len(result.Errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(result.Errors))
	}
}

func TestFailureBuildsErrorResponse(t *testing.T) {
	result := Failure(http.StatusBadRequest, "invalid request", 1200)

	if result.Success {
		t.Fatal("expected failure response")
	}
	if result.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", result.Code)
	}
	if result.Severity != Error {
		t.Fatalf("expected severity ERROR, got %q", result.Severity)
	}
	if result.SubErrorCode != 1200 {
		t.Fatalf("expected sub error 1200, got %d", result.SubErrorCode)
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected one error, got %d", len(result.Errors))
	}
	if result.Errors[0].Description != "invalid request" {
		t.Fatalf("unexpected error description %q", result.Errors[0].Description)
	}
}

func TestFailureDefaultsToInternalServerError(t *testing.T) {
	result := Failure(0, "boom", 0)

	if result.Code != http.StatusInternalServerError {
		t.Fatalf("expected default status 500, got %d", result.Code)
	}
}
