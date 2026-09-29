# Install on Keenetic / Entware

The router does **not** need Go or make. Download the prebuilt static MIPSLE binary from the latest GitHub Release.

## 1. Download

```sh
cd /opt/bin

wget -O xkeen-autoreload-vless \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless
```

If `wget` is unavailable, use curl:

```sh
curl -fL \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle \
  -o /opt/bin/xkeen-autoreload-vless

chmod +x /opt/bin/xkeen-autoreload-vless
```

## 2. Check binary

```sh
/opt/bin/xkeen-autoreload-vless version
```

Optional architecture check:

```sh
uname -m
file /opt/bin/xkeen-autoreload-vless
```

## 3. Create config

```sh
cat > /opt/etc/xkeen-autoreload-vless.env <<'EOF'
export XKEEN_SUBSCRIPTION_URL='PASTE_SUBSCRIPTION_URL_HERE'
export XKEEN_COUNTRY='Швейцария'
export XKEEN_CITY='Цюрих'

export XKEEN_OUTBOUND_PATH='/opt/etc/xray/configs/04_outbounds.json'
export XKEEN_STATE_PATH='/opt/var/lib/xkeen-autoreload-vless/state.json'

export XKEEN_COMMAND='xkeen'
export XKEEN_INTERVAL='1h'
export XKEEN_HTTP_TIMEOUT='20s'
EOF

chmod 600 /opt/etc/xkeen-autoreload-vless.env
```

## 4. Test without changing XKeen

```sh
. /opt/etc/xkeen-autoreload-vless.env
/opt/bin/xkeen-autoreload-vless check
```

The command should print the generated `04_outbounds.json` only.

## 5. One-time update

```sh
. /opt/etc/xkeen-autoreload-vless.env
/opt/bin/xkeen-autoreload-vless update
```

Force refresh of rotating SNI/SID:

```sh
/opt/bin/xkeen-autoreload-vless update --force
```

## 6. Install as Entware service

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

## 7. Verify process

```sh
ps | grep xkeen-autoreload-vless
```

## 8. Update to latest binary later

```sh
/opt/etc/init.d/S99xkeen-autoreload stop

wget -O /opt/bin/xkeen-autoreload-vless.new \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless.new
mv /opt/bin/xkeen-autoreload-vless.new /opt/bin/xkeen-autoreload-vless

/opt/etc/init.d/S99xkeen-autoreload start
```

Using a temporary file avoids replacing an executable while it is running and prevents `Text file busy`.
