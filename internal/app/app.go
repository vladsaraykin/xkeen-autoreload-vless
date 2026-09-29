package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/config"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/service"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/setup"
)

func Run(ctx context.Context, args []string, version string) error {
	cfg, cmd, err := config.Parse(args)
	if err != nil {
		return err
	}
	if cmd == "version" {
		fmt.Println(version)
		return nil
	}

	if cmd == "setup" {
		_, err := setup.Run(ctx, cfg, os.Stdin, os.Stdout)
		return err
	}

	if cfg.SubscriptionURL == "" || cfg.Country == "" {
		if !config.IsInteractive() {
			return fmt.Errorf("configuration is missing; run 'xkeen-autoreload-vless setup' interactively")
		}
		cfg, err = setup.Run(ctx, cfg, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
	}

	updater := service.New(cfg)
	runOnce := func() error {
		_, err := updater.Execute(ctx)
		return err
	}
	if cfg.Once || cmd == "check" || cmd == "update" {
		return runOnce()
	}

	if err := runOnce(); err != nil {
		fmt.Println("update failed:", err)
	}
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := runOnce(); err != nil {
				fmt.Println("update failed:", err)
			}
		}
	}
}
