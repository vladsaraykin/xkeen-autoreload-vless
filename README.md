# xkeen-autoreload-vless

Small Go service for Keenetic/Entware that downloads a Base64 VLESS subscription, selects a country/city node, renders XKeen `04_outbounds.json`, and restarts XKeen only when the stable node identity changes.

The provider may rotate `sni` and `sid` on every subscription request. To avoid needless restarts, the stable identity is:

`UUID | address | port | publicKey`

Use `--force` when you intentionally want fresh `sni`/`sid` even if that stable identity is unchanged.

## Commands

```sh
xkeen-autoreload-vless run
xkeen-autoreload-vless update --force
xkeen-autoreload-vless check
xkeen-autoreload-vless version
```

Configuration can be supplied with flags or environment variables. See `.env.example`.

## Build for Keenetic / MIPSLE

```sh
make test
make build-mipsle
```

The result is `dist/xkeen-autoreload-vless-linux-mipsle`.

## Install on Entware

```sh
cp dist/xkeen-autoreload-vless-linux-mipsle /opt/bin/xkeen-autoreload-vless
chmod +x /opt/bin/xkeen-autoreload-vless
```

Create `/opt/etc/xkeen-autoreload-vless.env`:

```sh
export XKEEN_SUBSCRIPTION_URL='https://...'
export XKEEN_COUNTRY='Швейцария'
export XKEEN_CITY='Цюрих'
export XKEEN_INTERVAL='1h'
```

For Entware init scripts, either export these variables in the wrapper or create a small launcher script that sources the env file before starting the binary.

## One-shot update

```sh
XKEEN_SUBSCRIPTION_URL='https://...' xkeen-autoreload-vless update
```

Force refresh of rotating SNI/SID:

```sh
XKEEN_SUBSCRIPTION_URL='https://...' xkeen-autoreload-vless update --force
```

## Dry-run / validation

```sh
XKEEN_SUBSCRIPTION_URL='https://...' xkeen-autoreload-vless check
```

This prints the generated JSON without writing files or restarting XKeen.

## Safety

Before replacing the config, the current file is copied to `04_outbounds.json.bak`. If `xkeen -restart` fails, the previous config is restored and XKeen is restarted again.
