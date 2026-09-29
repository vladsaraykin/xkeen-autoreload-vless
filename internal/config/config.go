package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DefaultPath = "/opt/etc/xkeen-autoreload-vless.env"

type Config struct {
	ConfigPath      string
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
	case "run", "check", "update", "setup", "version":
	default:
		return Config{}, "", fmt.Errorf("unknown command %q", cmd)
	}
	if cmd == "version" {
		return Config{}, cmd, nil
	}

	configPath := os.Getenv("XKEEN_CONFIG_PATH")
	if configPath == "" {
		configPath = DefaultPath
	}
	fileValues, err := LoadFile(configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, "", fmt.Errorf("load config file: %w", err)
	}

	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	cfg := Config{ConfigPath: configPath}
	fs.StringVar(&cfg.ConfigPath, "config", configPath, "configuration file path")
	fs.StringVar(&cfg.SubscriptionURL, "subscription-url", value("XKEEN_SUBSCRIPTION_URL", fileValues, ""), "subscription URL")
	fs.StringVar(&cfg.Country, "country", value("XKEEN_COUNTRY", fileValues, ""), "country name to select")
	fs.StringVar(&cfg.City, "city", value("XKEEN_CITY", fileValues, ""), "optional city name to prefer")
	fs.StringVar(&cfg.OutputPath, "output", value("XKEEN_OUTBOUND_PATH", fileValues, "/opt/etc/xray/configs/04_outbounds.json"), "04_outbounds.json path")
	fs.StringVar(&cfg.StatePath, "state", value("XKEEN_STATE_PATH", fileValues, "/opt/var/lib/xkeen-autoreload-vless/state.json"), "state file path")
	fs.StringVar(&cfg.XKeenCommand, "xkeen-command", value("XKEEN_COMMAND", fileValues, "xkeen"), "xkeen executable")
	fs.DurationVar(&cfg.Interval, "interval", durationValue("XKEEN_INTERVAL", fileValues, time.Hour), "service polling interval")
	fs.DurationVar(&cfg.HTTPTimeout, "http-timeout", durationValue("XKEEN_HTTP_TIMEOUT", fileValues, 20*time.Second), "HTTP timeout")
	fs.BoolVar(&cfg.Force, "force", false, "force update even when stable node identity is unchanged")
	fs.BoolVar(&cfg.Once, "once", false, "run one iteration and exit")
	fs.BoolVar(&cfg.DryRun, "dry-run", false, "print generated config without writing or restarting")
	if err := fs.Parse(args[1:]); err != nil {
		return Config{}, "", err
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

func LoadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make(map[string]string)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if len(v) >= 2 {
			if (v[0] == '\'' && v[len(v)-1] == '\'') || (v[0] == '"' && v[len(v)-1] == '"') {
				v = v[1 : len(v)-1]
			}
		}
		values[k] = v
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func SaveFile(cfg Config) error {
	if cfg.SubscriptionURL == "" {
		return errors.New("subscription URL is empty")
	}
	if cfg.Country == "" {
		return errors.New("country is empty")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.ConfigPath), 0o755); err != nil {
		return err
	}
	content := strings.Join([]string{
		"XKEEN_SUBSCRIPTION_URL=" + strconv.Quote(cfg.SubscriptionURL),
		"XKEEN_COUNTRY=" + strconv.Quote(cfg.Country),
		"XKEEN_CITY=" + strconv.Quote(cfg.City),
		"XKEEN_OUTBOUND_PATH=" + strconv.Quote(cfg.OutputPath),
		"XKEEN_STATE_PATH=" + strconv.Quote(cfg.StatePath),
		"XKEEN_COMMAND=" + strconv.Quote(cfg.XKeenCommand),
		"XKEEN_INTERVAL=" + strconv.Quote(cfg.Interval.String()),
		"XKEEN_HTTP_TIMEOUT=" + strconv.Quote(cfg.HTTPTimeout.String()),
		"",
	}, "\n")
	tmp := cfg.ConfigPath + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		return err
	}
	if err := os.Chmod(tmp, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, cfg.ConfigPath)
}

func IsInteractive() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func value(key string, fileValues map[string]string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if v := fileValues[key]; v != "" {
		return v
	}
	return fallback
}

func durationValue(key string, fileValues map[string]string, fallback time.Duration) time.Duration {
	if raw := value(key, fileValues, ""); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil {
			return d
		}
	}
	return fallback
}
