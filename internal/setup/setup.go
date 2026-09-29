package setup

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/config"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/subscription"
	"github.com/vladsaraykin/xkeen-autoreload-vless/internal/vless"
)

func Run(ctx context.Context, cfg config.Config, in io.Reader, out io.Writer) (config.Config, error) {
	reader := bufio.NewReader(in)

	fmt.Fprintln(out, "Первичная настройка xkeen-autoreload-vless")
	fmt.Fprintln(out)

	subscriptionURL, err := prompt(reader, out, "Вставьте URL подписки: ")
	if err != nil {
		return cfg, err
	}
	if subscriptionURL == "" {
		return cfg, fmt.Errorf("URL подписки не может быть пустым")
	}

	fmt.Fprintln(out, "Загружаю и декодирую подписку...")
	client := subscription.New(cfg.HTTPTimeout)
	lines, err := client.Fetch(ctx, subscriptionURL)
	if err != nil {
		return cfg, fmt.Errorf("не удалось загрузить подписку: %w", err)
	}

	countries := vless.Countries(lines)
	if len(countries) == 0 {
		return cfg, fmt.Errorf("в подписке не найдено ни одной VLESS-ноды с распознаваемой страной")
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Доступные страны:")
	for i, country := range countries {
		fmt.Fprintf(out, "  %d) %s\n", i+1, country)
	}

	var selected string
	for selected == "" {
		raw, err := prompt(reader, out, fmt.Sprintf("Выберите страну [1-%d]: ", len(countries)))
		if err != nil {
			return cfg, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil || n < 1 || n > len(countries) {
			fmt.Fprintln(out, "Некорректный номер. Попробуйте ещё раз.")
			continue
		}
		selected = countries[n-1]
	}

	cfg.SubscriptionURL = subscriptionURL
	cfg.Country = selected
	cfg.City = ""
	if err := config.SaveFile(cfg); err != nil {
		return cfg, fmt.Errorf("сохранить конфигурацию: %w", err)
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Готово. Выбрана страна: %s\n", selected)
	fmt.Fprintf(out, "Конфигурация сохранена: %s\n", cfg.ConfigPath)
	return cfg, nil
}

func prompt(reader *bufio.Reader, out io.Writer, text string) (string, error) {
	fmt.Fprint(out, text)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
