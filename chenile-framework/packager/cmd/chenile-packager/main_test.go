package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckManifestReturnsNilForExistingManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chenile-packager.yaml")
	if err := os.WriteFile(path, []byte("name: customer\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := checkManifest(path); err != nil {
		t.Fatalf("expected manifest to exist: %v", err)
	}
}

func TestCheckManifestReturnsErrorForMissingManifest(t *testing.T) {
	err := checkManifest(filepath.Join(t.TempDir(), "missing.yaml"))

	if err == nil {
		t.Fatal("expected missing manifest error")
	}
}

func TestRunReturnsSuccessForExistingManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chenile-packager.yaml")
	if err := os.WriteFile(path, []byte("name: customer\n"), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"--manifest", path}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "exists") {
		t.Fatalf("expected success output, got %q", stdout.String())
	}
}

func TestRunReturnsFailureForMissingManifest(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"--manifest", filepath.Join(t.TempDir(), "missing.yaml")}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Fatalf("expected not found error, got %q", stderr.String())
	}
}
