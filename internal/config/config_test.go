package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.env")
	cfg := Config{
		ConfigPath: path,
		SubscriptionURL: "https://example.com/s/a?x=1&y=2",
		Country: "🇨🇭",
		OutputPath: "/tmp/out.json",
		StatePath: "/tmp/state.json",
		XKeenCommand: "xkeen",
		Interval: time.Hour,
		HTTPTimeout: 20 * time.Second,
		UserAgent: "v2rayNG/1.8.19",
		InsecureTLS: true,
	}
	if err := SaveFile(cfg); err != nil {
		t.Fatal(err)
	}
	values, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if values["XKEEN_SUBSCRIPTION_URL"] != cfg.SubscriptionURL ||
		values["XKEEN_COUNTRY"] != cfg.Country ||
		values["XKEEN_USER_AGENT"] != cfg.UserAgent ||
		values["XKEEN_INSECURE_TLS"] != "true" {
		t.Fatalf("unexpected values: %#v", values)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}
}
