package config

import "testing"

func TestGetEnv_DefaultAndValue(t *testing.T) {
	t.Setenv("CFG_TEST_KEY", "present")
	if got := getEnv("CFG_TEST_KEY", "default"); got != "present" {
		t.Fatalf("getEnv existing got %q", got)
	}

	if got := getEnv("CFG_MISSING_KEY", "default"); got != "default" {
		t.Fatalf("getEnv default got %q", got)
	}
}

func TestLoad_ReadsEnvironment(t *testing.T) {
	t.Setenv("EXTERNAL_API_URL", "https://api.example.com")
	t.Setenv("EXTERNAL_API_TOKEN", "token")
	t.Setenv("DB_HOST", "db-host")
	t.Setenv("DB_PORT", "26258")
	t.Setenv("DB_USER", "db-user")
	t.Setenv("DB_PASSWORD", "db-pass")
	t.Setenv("DB_NAME", "db-name")
	t.Setenv("DB_SCHEMA", "db-schema")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("API_PORT", "9999")
	t.Setenv("API_HOST", "0.0.0.0")
	t.Setenv("FETCH_INTERVAL", "120")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := Load()

	if cfg.ExternalAPIURL != "https://api.example.com" || cfg.DBPort != 26258 || cfg.APIPort != "9999" || cfg.FetchInterval != 120 || cfg.LogLevel != "debug" {
		t.Fatalf("unexpected config loaded: %+v", cfg)
	}
}
