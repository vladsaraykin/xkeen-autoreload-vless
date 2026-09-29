package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	SubscriptionURL string
	Country         string
	City            string
	OutputPath      string
	StatePath       string
	XKeenCommand    string
	Interval        time.Duration
	HTTPTimeout     time.Duration
	Force           bool
	Once            bool
	DryRun          bool
}

func Parse(args []string) (Config, string, error) {
	if len(args) == 0 {
		args = []string{"run"}
	}
	cmd := args[0]
	switch cmd {
	case "run", "check", "update", "version":
	default:
		return Config{}, "", fmt.Errorf("unknown command %q", cmd)
	}
	if cmd == "version" {
		return Config{}, cmd, nil
	}

	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	cfg := Config{}
	fs.StringVar(&cfg.SubscriptionURL, "subscription-url", getenv("XKEEN_SUBSCRIPTION_URL", ""), "subscription URL")
	fs.StringVar(&cfg.Country, "country", getenv("XKEEN_COUNTRY", "Швейцария"), "country name to select")
	fs.StringVar(&cfg.City, "city", getenv("XKEEN_CITY", "Цюрих"), "optional city name to prefer")
	fs.StringVar(&cfg.OutputPath, "output", getenv("XKEEN_OUTBOUND_PATH", "/opt/etc/xray/configs/04_outbounds.json"), "04_outbounds.json path")
	fs.StringVar(&cfg.StatePath, "state", getenv("XKEEN_STATE_PATH", "/opt/var/lib/xkeen-autoreload-vless/state.json"), "state file path")
	fs.StringVar(&cfg.XKeenCommand, "xkeen-command", getenv("XKEEN_COMMAND", "xkeen"), "xkeen executable")
	fs.DurationVar(&cfg.Interval, "interval", getenvDuration("XKEEN_INTERVAL", time.Hour), "service polling interval")
	fs.DurationVar(&cfg.HTTPTimeout, "http-timeout", getenvDuration("XKEEN_HTTP_TIMEOUT", 20*time.Second), "HTTP timeout")
	fs.BoolVar(&cfg.Force, "force", false, "force update even when stable node identity is unchanged")
	fs.BoolVar(&cfg.Once, "once", false, "run one iteration and exit")
	fs.BoolVar(&cfg.DryRun, "dry-run", false, "print generated config without writing or restarting")
	if err := fs.Parse(args[1:]); err != nil {
		return Config{}, "", err
	}
	if cfg.SubscriptionURL == "" {
		return Config{}, "", errors.New("subscription URL is required: use --subscription-url or XKEEN_SUBSCRIPTION_URL")
	}
	if cmd == "check" {
		cfg.DryRun = true
		cfg.Once = true
	}
	if cmd == "update" {
		cfg.Once = true
	}
	return cfg, cmd, nil
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
func getenvDuration(k string, fallback time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
