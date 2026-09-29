# Install on Keenetic / Entware

The router does **not** need Go, make, git, or gcc. Use the prebuilt static MIPSLE binary from the latest GitHub Release.

## 1. Download

```sh
cd /opt/bin

wget -O xkeen-autoreload-vless \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless
```

If `wget` is unavailable:

```sh
curl -fL \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle \
  -o /opt/bin/xkeen-autoreload-vless

chmod +x /opt/bin/xkeen-autoreload-vless
```

Check the version:

```sh
/opt/bin/xkeen-autoreload-vless version
```

## 2. First-run interactive setup

You no longer need to create `/opt/etc/xkeen-autoreload-vless.env` manually.

Run:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

The application asks for the subscription URL:

```text
Первичная настройка xkeen-autoreload-vless

Вставьте URL подписки:
```

Paste your subscription URL and press Enter. The application downloads the response, decodes Base64, parses VLESS nodes, and extracts the countries from labels such as:

```text
🇳🇱 Амстердам, Нидерланды, Extra
🇩🇪 Берлин, Германия, Extra
🇨🇭 Цюрих, Швейцария, Extra
```

Then it displays a numbered list:

```text
Доступные страны:
  1) Германия
  2) Нидерланды
  3) Швейцария

Выберите страну [1-3]: 3
```

After the choice:

```text
Готово. Выбрана страна: Швейцария
Конфигурация сохранена: /opt/etc/xkeen-autoreload-vless.env
```

The file is stored with permission `0600`. The application reads it directly; it is not executed or sourced as a shell script.

If you start:

```sh
/opt/bin/xkeen-autoreload-vless
```

without a configuration, the same wizard starts automatically when running in an interactive terminal.

To change the URL or selected country later:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

## 3. Validate without changing XKeen

```sh
/opt/bin/xkeen-autoreload-vless check
```

This downloads the subscription and prints the generated `04_outbounds.json`, but does not write the file and does not restart XKeen.

## 4. Apply once

```sh
/opt/bin/xkeen-autoreload-vless update
```

Force refresh of rotating SNI/SID:

```sh
/opt/bin/xkeen-autoreload-vless update --force
```

## 5. Install as Entware service

Download the init script:

```sh
wget -O /opt/etc/init.d/S99xkeen-autoreload \
  https://raw.githubusercontent.com/vladsaraykin/xkeen-autoreload-vless/main/scripts/S99xkeen-autoreload

chmod +x /opt/etc/init.d/S99xkeen-autoreload
```

Start:

```sh
/opt/etc/init.d/S99xkeen-autoreload start
```

Stop:

```sh
/opt/etc/init.d/S99xkeen-autoreload stop
```

Restart:

```sh
/opt/etc/init.d/S99xkeen-autoreload restart
```

Verify:

```sh
ps | grep xkeen-autoreload-vless
```

The service checks the subscription every hour by default. If only the provider's rotating SNI/SID changes, XKeen is not restarted.

## 6. Update the application later

Stop the service first and download to a temporary file to avoid `Text file busy`:

```sh
/opt/etc/init.d/S99xkeen-autoreload stop

wget -O /opt/bin/xkeen-autoreload-vless.new \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless.new
mv /opt/bin/xkeen-autoreload-vless.new /opt/bin/xkeen-autoreload-vless

/opt/bin/xkeen-autoreload-vless version
/opt/etc/init.d/S99xkeen-autoreload start
```
