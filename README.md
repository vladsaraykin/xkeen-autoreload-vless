# xkeen-autoreload-vless

Небольшой Go-сервис для Keenetic/Entware, который автоматически обновляет VLESS-конфигурацию XKeen.

Поддерживаются два формата подписок:

1. Base64 со списком ссылок `vless://...`;
2. JSON-массив готовых Xray-конфигураций, например тот, который возвращается при запросе с User-Agent `v2rayNG/1.8.19`.

Приложение:

- скачивает подписку;
- автоматически определяет её формат;
- при первом запуске предлагает выбрать страну;
- генерирует `04_outbounds.json`;
- сохраняет исходные `grpcSettings`, `tcpSettings`, `mux` и Reality-параметры JSON-подписки;
- перезапускает XKeen только при реальном изменении ноды.

Стабильная идентичность ноды определяется по:

`UUID | address | port | publicKey`

Изменение только `sni` или `sid` не вызывает лишний restart.

## Первый запуск

```sh
/opt/bin/xkeen-autoreload-vless setup
```

Пример:

```text
Первичная настройка xkeen-autoreload-vless

Вставьте URL подписки: https://connect.example.org/...
Загружаю подписку и определяю формат...

Доступные страны:
  1) 🇦🇹 Австрия
  2) 🇩🇪 Германия
  3) 🇳🇱 Нидерланды
  4) 🇨🇭 Швейцария

Выберите страну [1-4]: 3

Готово. Выбрана страна: 🇳🇱 Нидерланды
Конфигурация сохранена: /opt/etc/xkeen-autoreload-vless.env
```

Для JSON-подписок страна группируется прежде всего по emoji-флагу в поле `remarks`. Поэтому варианты вроде `🇳🇱 Нидерланды · ультра`, `🇳🇱 Премиум 1` и `🇳🇱 Нидерланды • LTE` считаются одной страной.

По умолчанию HTTP-запрос отправляется с:

```text
User-Agent: v2rayNG/1.8.19
```

Если TLS-сертификат подписки не проходит проверку, интерактивный мастер предложит повторить запрос без проверки сертификата. Это аналог `curl -k`, но он не включается автоматически.

## Команды

```sh
xkeen-autoreload-vless setup
xkeen-autoreload-vless run
xkeen-autoreload-vless update
xkeen-autoreload-vless update --force
xkeen-autoreload-vless check
xkeen-autoreload-vless version
```

Подробнее об установке: [INSTALL.md](INSTALL.md).

## Сборка

```sh
make test
make vet
make build-mipsle
```

Результат:

```text
dist/xkeen-autoreload-vless-linux-mipsle
```

## Безопасность

Конфигурация читается Go-приложением как данные и не выполняется shell-ом.

Перед заменой `04_outbounds.json` создаётся backup:

```text
04_outbounds.json.bak
```

Если `xkeen -restart` завершается ошибкой, предыдущий конфиг восстанавливается.
