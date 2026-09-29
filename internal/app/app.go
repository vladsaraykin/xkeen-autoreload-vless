package app

import (
	"context"
	"fmt"
	"time"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/config"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/service"
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
