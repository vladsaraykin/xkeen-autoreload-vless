package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/config"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/state"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/subscription"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/xray"
)

type Updater struct {
	cfg config.Config
	sub *subscription.Client
}

func New(cfg config.Config) *Updater {
	return &Updater{cfg: cfg, sub: subscription.New(cfg.HTTPTimeout)}
}

func (u *Updater) Execute(ctx context.Context) (bool, error) {
	lines, err := u.sub.Fetch(ctx, u.cfg.SubscriptionURL)
	if err != nil {
		return false, err
	}
	node, err := vless.Select(lines, u.cfg.Country, u.cfg.City)
	if err != nil {
		return false, err
	}
	rendered, err := xray.Render(node)
	if err != nil {
		return false, err
	}
	if u.cfg.DryRun {
		fmt.Println(string(rendered))
		return false, nil
	}
	st, err := state.Load(u.cfg.StatePath)
	if err != nil {
		return false, fmt.Errorf("load state: %w", err)
	}
	key := node.StableKey()
	if !u.cfg.Force && st.StableKey == key {
		fmt.Printf("unchanged: %s:%d (%s); SNI/SID rotation ignored\n", node.Address, node.Port, node.Name)
		return false, nil
	}
	if err := atomicReplaceWithBackup(u.cfg.OutputPath, rendered); err != nil {
		return false, err
	}
	if err := restart(ctx, u.cfg.XKeenCommand); err != nil {
		_ = restoreBackup(u.cfg.OutputPath)
		_ = restart(ctx, u.cfg.XKeenCommand)
		return false, fmt.Errorf("restart xkeen: %w", err)
	}
	if err := state.Save(u.cfg.StatePath, key); err != nil {
		return true, fmt.Errorf("save state: %w", err)
	}
	fmt.Printf("updated: %s:%d (%s)\n", node.Address, node.Port, node.Name)
	return true, nil
}

func atomicReplaceWithBackup(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		current, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path+".bak", current, 0o600); err != nil {
			return err
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
func restoreBackup(path string) error {
	b, err := os.ReadFile(path + ".bak")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}
func restart(ctx context.Context, command string) error {
	cmd := exec.CommandContext(ctx, command, "-restart")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
