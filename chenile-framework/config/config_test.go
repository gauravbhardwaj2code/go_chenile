package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileSupportsPropertiesAndYamlStyle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "application.yaml")
	if err := os.WriteFile(path, []byte("server.port: 8081\nservice.name=inventory\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.String("server.port", "") != "8081" {
		t.Fatalf("expected server.port from file")
	}
	if cfg.String("service.name", "") != "inventory" {
		t.Fatalf("expected service.name from file")
	}
}

func TestLoadEnvNormalizesKeys(t *testing.T) {
	t.Setenv("CHENILE_PROFILE", "test")

	cfg := Empty()
	cfg.LoadEnv()

	if cfg.String("chenile.profile", "") != "test" {
		t.Fatalf("expected env key override")
	}
}
