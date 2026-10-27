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

	fmt.Fprintln(out, "Загружаю подписку и определяю формат...")
	client := subscription.NewWithOptions(cfg.HTTPTimeout, cfg.UserAgent, cfg.InsecureTLS)
	nodes, err := client.Fetch(ctx, subscriptionURL)
	if err != nil && !cfg.InsecureTLS && isCertificateError(err) {
		answer, promptErr := prompt(reader, out, "Проверка TLS-сертификата не прошла. Повторить без проверки сертификата? [y/N]: ")
		if promptErr != nil {
			return cfg, promptErr
		}
		if strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes") || strings.EqualFold(answer, "д") || strings.EqualFold(answer, "да") {
			cfg.InsecureTLS = true
			client = subscription.NewWithOptions(cfg.HTTPTimeout, cfg.UserAgent, true)
			nodes, err = client.Fetch(ctx, subscriptionURL)
		}
	}
	if err != nil {
		return cfg, fmt.Errorf("не удалось загрузить подписку: %w", err)
	}

	choices := vless.Choices(nodes)
	if len(choices) == 0 {
		return cfg, fmt.Errorf("в подписке не найдено ни одной VLESS-ноды с распознаваемой страной")
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Доступные страны:")
	for i, choice := range choices {
		fmt.Fprintf(out, "  %d) %s\n", i+1, choice.Label)
	}

	var selected vless.Choice
	for selected.Key == "" {
		raw, err := prompt(reader, out, fmt.Sprintf("Выберите страну [1-%d]: ", len(choices)))
		if err != nil {
			return cfg, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil || n < 1 || n > len(choices) {
			fmt.Fprintln(out, "Некорректный номер. Попробуйте ещё раз.")
			continue
		}
		selected = choices[n-1]
	}

	cfg.SubscriptionURL = subscriptionURL
	cfg.Country = selected.Key
	cfg.City = ""
	if err := config.SaveFile(cfg); err != nil {
		return cfg, fmt.Errorf("сохранить конфигурацию: %w", err)
	}

	fmt.Fprintln(out)
	fmt.Fprintf(out, "Готово. Выбрана страна: %s\n", selected.Label)
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

func isCertificateError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "x509") ||
		strings.Contains(msg, "certificate") ||
		strings.Contains(msg, "tls")
}
