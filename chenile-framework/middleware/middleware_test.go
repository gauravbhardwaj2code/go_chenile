package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core"
)

func TestRequestIDGeneratesMissingHeader(t *testing.T) {
	m := &RequestID{Generate: func() string { return "req-1" }}
	exchange := &core.Exchange{}

	if err := m.Before(context.Background(), exchange); err != nil {
		t.Fatal(err)
	}

	if exchange.Headers["X-Request-Id"] != "req-1" {
		t.Fatalf("expected generated request id")
	}
}

type invalidInput struct{}

func (invalidInput) Validate() error { return errors.New("invalid") }

func TestValidateInputDelegatesToPayload(t *testing.T) {
	err := ValidateInput{}.Before(context.Background(), &core.Exchange{Input: invalidInput{}})

	if err == nil || err.Error() != "invalid" {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestRecoverConvertsPanicToError(t *testing.T) {
	err := Recover(func() error {
		panic("boom")
	})

	if err == nil {
		t.Fatal("expected panic error")
	}
}
