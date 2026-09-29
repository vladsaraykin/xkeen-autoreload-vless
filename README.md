# xkeen-autoreload-vless

Small Go service for Keenetic/Entware that downloads a Base64 VLESS subscription, lets the user choose a country, renders XKeen `04_outbounds.json`, and restarts XKeen only when the stable node identity changes.

The provider may rotate `sni` and `sid` on every subscription request. To avoid needless restarts, the stable identity is:

`UUID | address | port | publicKey`

Use `--force` when you intentionally want fresh `sni`/`sid` even if that stable identity is unchanged.

## First-run setup

No manual env-file editing is required.

Run:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

The wizard will:

1. Ask for the subscription URL.
2. Download the Base64 subscription.
3. Decode and parse VLESS nodes.
4. Discover unique countries from node labels.
5. Show a numbered country list.
6. Save the selected country and subscription URL to `/opt/etc/xkeen-autoreload-vless.env` with mode `0600`.

Example:

```text
Первичная настройка xkeen-autoreload-vless

Вставьте URL подписки: https://example.com/s/...
Загружаю и декодирую подписку...

Доступные страны:
  1) Австрия
  2) Германия
  3) Нидерланды
  4) Швейцария

Выберите страну [1-4]: 4

Готово. Выбрана страна: Швейцария
Конфигурация сохранена: /opt/etc/xkeen-autoreload-vless.env
```

Running `xkeen-autoreload-vless` or `xkeen-autoreload-vless run` without a saved configuration also starts the wizard automatically when stdin is interactive.

To change the subscription or country later, run `setup` again.

## Commands

```sh
xkeen-autoreload-vless setup
xkeen-autoreload-vless run
xkeen-autoreload-vless update
xkeen-autoreload-vless update --force
xkeen-autoreload-vless check
xkeen-autoreload-vless version
```

See [INSTALL.md](INSTALL.md) for Keenetic/Entware installation.

## Build for Keenetic / MIPSLE

```sh
make test
make build-mipsle
```

The result is `dist/xkeen-autoreload-vless-linux-mipsle`.

## Safety

The config file is parsed by the Go application as data and is not sourced as shell code.

Before replacing the XKeen config, the current file is copied to `04_outbounds.json.bak`. If `xkeen -restart` fails, the previous config is restored and XKeen is restarted again.
