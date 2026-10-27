# Установка на Keenetic / Entware

На роутере не нужны Go, make, git или gcc. Используется готовый статический MIPSLE-бинарник из GitHub Releases.

## 1. Скачать бинарник

### wget

```sh
cd /opt/bin

wget -O xkeen-autoreload-vless \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless
```

### curl

```sh
curl -fL \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle \
  -o /opt/bin/xkeen-autoreload-vless

chmod +x /opt/bin/xkeen-autoreload-vless
```

Проверка версии:

```sh
/opt/bin/xkeen-autoreload-vless version
```

## 2. Первый запуск

```sh
/opt/bin/xkeen-autoreload-vless setup
```

Приложение попросит URL подписки и автоматически определит формат:

- Base64 со строками `vless://...`;
- JSON-массив Xray-конфигураций.

Для нового JSON-формата запрос выполняется с User-Agent:

```text
v2rayNG/1.8.19
```

Далее приложение покажет список стран:

```text
Доступные страны:
  1) 🇦🇹 Австрия
  2) 🇩🇪 Германия
  3) 🇳🇱 Нидерланды
  4) 🇵🇱 Польша

Выберите страну [1-4]:
```

Выбор сохраняется в:

```text
/opt/etc/xkeen-autoreload-vless.env
```

с правами `0600`.

Если сертификат подписки не проходит проверку, мастер предложит повторить запрос без TLS-проверки. Это соответствует поведению `curl -k`.

## 3. Проверка

```sh
/opt/bin/xkeen-autoreload-vless check
```

Команда генерирует конфиг и выводит его в консоль, но не записывает файл и не перезапускает XKeen.

## 4. Разовое применение

```sh
/opt/bin/xkeen-autoreload-vless update
```

Принудительное обновление:

```sh
/opt/bin/xkeen-autoreload-vless update --force
```

## 5. Установка как сервис Entware

### wget

```sh
wget -O /opt/etc/init.d/S99xkeen-autoreload \
  https://raw.githubusercontent.com/vladsaraykin/xkeen-autoreload-vless/main/scripts/S99xkeen-autoreload

chmod +x /opt/etc/init.d/S99xkeen-autoreload
```

### curl

```sh
curl -fL \
  https://raw.githubusercontent.com/vladsaraykin/xkeen-autoreload-vless/main/scripts/S99xkeen-autoreload \
  -o /opt/etc/init.d/S99xkeen-autoreload

chmod +x /opt/etc/init.d/S99xkeen-autoreload
```

Запуск:

```sh
/opt/etc/init.d/S99xkeen-autoreload start
```

Остановка:

```sh
/opt/etc/init.d/S99xkeen-autoreload stop
```

Перезапуск:

```sh
/opt/etc/init.d/S99xkeen-autoreload restart
```

Проверка:

```sh
ps | grep xkeen-autoreload-vless
```

По умолчанию подписка проверяется раз в час.

## 6. Обновление приложения

Сначала остановите сервис:

```sh
/opt/etc/init.d/S99xkeen-autoreload stop
```

### wget

```sh
wget -O /opt/bin/xkeen-autoreload-vless.new \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle
```

### curl

```sh
curl -fL \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle \
  -o /opt/bin/xkeen-autoreload-vless.new
```

Затем:

```sh
chmod +x /opt/bin/xkeen-autoreload-vless.new

mv /opt/bin/xkeen-autoreload-vless.new \
   /opt/bin/xkeen-autoreload-vless

/opt/bin/xkeen-autoreload-vless version

/opt/etc/init.d/S99xkeen-autoreload start
```

Использование временного файла предотвращает ошибку:

```text
Text file busy
```
