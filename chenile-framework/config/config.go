package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	values map[string]string
}

func Empty() Config {
	return Config{values: map[string]string{}}
}

func Load(files ...string) (Config, error) {
	cfg := Empty()
	for _, file := range files {
		if file == "" {
			continue
		}
		if err := cfg.LoadFile(file); err != nil {
			return Config{}, err
		}
	}
	cfg.LoadEnv()
	return cfg, nil
}

func (c *Config) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			key, value, ok = strings.Cut(line, ":")
		}
		if !ok {
			return fmt.Errorf("parse config %q line %d: expected key=value or key: value", path, lineNumber)
		}
		c.Set(strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"`))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}
	return nil
}

func (c *Config) LoadEnv() {
	for _, env := range os.Environ() {
		key, value, ok := strings.Cut(env, "=")
		if ok {
			c.Set(envKey(key), value)
		}
	}
}

func (c Config) Get(key string) (string, bool) {
	value, ok := c.values[key]
	return value, ok
}

func (c Config) String(key string, fallback string) string {
	if value, ok := c.Get(key); ok {
		return value
	}
	return fallback
}

func (c Config) MustString(key string) (string, error) {
	value, ok := c.Get(key)
	if !ok || value == "" {
		return "", fmt.Errorf("config key %q is required", key)
	}
	return value, nil
}

func (c *Config) Set(key string, value string) {
	if c.values == nil {
		c.values = map[string]string{}
	}
	c.values[key] = value
}

func envKey(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "_", "."))
}
