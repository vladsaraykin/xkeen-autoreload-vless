#!/bin/sh
set -eu

BIN_SOURCE=${1:-./xkeen-autoreload-vless-linux-mipsle}
BIN_TARGET=/opt/bin/xkeen-autoreload-vless
INIT_TARGET=/opt/etc/init.d/S99xkeen-autoreload
ENV_TARGET=/opt/etc/xkeen-autoreload-vless.env

[ -f "$BIN_SOURCE" ] || { echo "binary not found: $BIN_SOURCE" >&2; exit 1; }

install -m 0755 "$BIN_SOURCE" "$BIN_TARGET"
install -m 0755 ./scripts/S99xkeen-autoreload "$INIT_TARGET"

if [ ! -f "$ENV_TARGET" ]; then
    cat > "$ENV_TARGET" <<'ENV'
export XKEEN_SUBSCRIPTION_URL='https://replace-me'
export XKEEN_COUNTRY='Швейцария'
export XKEEN_CITY='Цюрих'
export XKEEN_OUTBOUND_PATH='/opt/etc/xray/configs/04_outbounds.json'
export XKEEN_STATE_PATH='/opt/var/lib/xkeen-autoreload-vless/state.json'
export XKEEN_COMMAND='xkeen'
export XKEEN_INTERVAL='1h'
export XKEEN_HTTP_TIMEOUT='20s'
ENV
    chmod 0600 "$ENV_TARGET"
    echo "created $ENV_TARGET; edit subscription URL before starting"
fi

echo "installed $BIN_TARGET and $INIT_TARGET"
echo "start with: $INIT_TARGET start"
