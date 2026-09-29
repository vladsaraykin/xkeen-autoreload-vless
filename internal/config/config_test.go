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
		Country: "Швейцария",
		OutputPath: "/tmp/out.json",
		StatePath: "/tmp/state.json",
		XKeenCommand: "xkeen",
		Interval: time.Hour,
		HTTPTimeout: 20 * time.Second,
	}
	if err := SaveFile(cfg); err != nil {
		t.Fatal(err)
	}
	values, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if values["XKEEN_SUBSCRIPTION_URL"] != cfg.SubscriptionURL || values["XKEEN_COUNTRY"] != cfg.Country {
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
