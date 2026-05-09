package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeriveBuildsNamesFromServiceName(t *testing.T) {
	d := derive("order-item", "../..")

	if d.Package != "orderitem" {
		t.Fatalf("expected package orderitem, got %q", d.Package)
	}
	if d.TypeName != "OrderItem" {
		t.Fatalf("expected type OrderItem, got %q", d.TypeName)
	}
	if d.ServiceID != "orderitemService" {
		t.Fatalf("expected service id orderitemService, got %q", d.ServiceID)
	}
	if d.RouteBase != "/order-items" {
		t.Fatalf("expected route /order-items, got %q", d.RouteBase)
	}
}

func TestGenerateWritesRunnableServiceSkeleton(t *testing.T) {
	dir := t.TempDir()
	d := derive("customer", "../..")
	target := filepath.Join(dir, d.BinaryName)

	if err := generate(target, d); err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		"go.mod",
		"cmd/customer-service/main.go",
		"customer/controller.go",
		"customer/controller_test.go",
		"customer/service.go",
		"customer/service_test.go",
		"customer/model.go",
		"customer/module.go",
		"test/customer_service_test.go",
		"test/features/customer.feature",
		"test/fixtures/create_customer.json",
	}
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(filepath.Join(target, expectedFile)); err != nil {
			t.Fatalf("expected generated file %s: %v", expectedFile, err)
		}
	}
}

func TestUpdateWorkspaceAddsGeneratedService(t *testing.T) {
	root := t.TempDir()
	workFile := filepath.Join(root, "go.work")
	if err := os.WriteFile(workFile, []byte("go 1.22\n\nuse (\n\t./base\n)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "examples", "customer-service")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}

	if err := updateWorkspace(target, data{FrameworkRoot: "../.."}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(workFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "./examples/customer-service") {
		t.Fatalf("expected generated service in go.work, got:\n%s", string(content))
	}
}

func TestRunGeneratesServiceAndUpdatesWorkspace(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.work"), []byte("go 1.22\n\nuse (\n\t./base\n)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "examples")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"new", "--name", "invoice", "--out", out, "--framework-root", "../.."}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(out, "invoice-service", "go.mod")); err != nil {
		t.Fatalf("expected generated service: %v", err)
	}
	if !strings.Contains(stdout.String(), "go test ./...") {
		t.Fatalf("expected next steps, got %q", stdout.String())
	}
}

func TestRunRequiresName(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"new"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "--name is required") {
		t.Fatalf("expected missing name error, got %q", stderr.String())
	}
}
