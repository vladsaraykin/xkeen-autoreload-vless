# Установка на Keenetic / Entware

На роутере **не нужны** Go, make, git или gcc. Используется готовый статический бинарник MIPSLE из последнего GitHub Release.

## 1. Скачать бинарник

```sh
cd /opt/bin

wget -O xkeen-autoreload-vless \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle

chmod +x /opt/bin/xkeen-autoreload-vless
```

Если `wget` недоступен, используйте `curl`:

```sh
curl -fL \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle \
  -o /opt/bin/xkeen-autoreload-vless

chmod +x /opt/bin/xkeen-autoreload-vless
```

Проверить версию:

```sh
/opt/bin/xkeen-autoreload-vless version
```

## 2. Первый запуск и интерактивная настройка

Создавать `/opt/etc/xkeen-autoreload-vless.env` вручную больше не нужно.

Запустите:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

Приложение попросит ссылку на подписку:

```text
Первичная настройка xkeen-autoreload-vless

Вставьте URL подписки:
```

Вставьте ссылку и нажмите Enter.

Далее приложение:

- скачает подписку;
- декодирует Base64;
- распарсит VLESS-ноды;
- извлечёт страны из названий вида:

```text
🇳🇱 Амстердам, Нидерланды, Extra
🇩🇪 Берлин, Германия, Extra
🇨🇭 Цюрих, Швейцария, Extra
```

После этого появится список:

```text
Доступные страны:
  1) Германия
  2) Нидерланды
  3) Швейцария

Выберите страну [1-3]: 3
```

После выбора:

```text
Готово. Выбрана страна: Швейцария
Конфигурация сохранена: /opt/etc/xkeen-autoreload-vless.env
```

Файл конфигурации сохраняется с правами `0600`.

Приложение читает этот файл напрямую и не выполняет его как shell-скрипт.

Если запустить:

```sh
/opt/bin/xkeen-autoreload-vless
```

без конфигурации, тот же мастер настройки запустится автоматически при интерактивном запуске.

Чтобы позже поменять ссылку или страну:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

## 3. Проверка без изменения XKeen

```sh
/opt/bin/xkeen-autoreload-vless check
```

Команда:

- скачает подписку;
- найдёт выбранную страну;
- сгенерирует `04_outbounds.json`;
- выведет его в консоль;
- не будет записывать файл;
- не будет перезапускать XKeen.

## 4. Разовое применение конфигурации

```sh
/opt/bin/xkeen-autoreload-vless update
```

Принудительно обновить ротируемые `SNI/SID`:

```sh
/opt/bin/xkeen-autoreload-vless update --force
```

## 5. Установка как сервис Entware

Скачайте init-скрипт:

```sh
wget -O /opt/etc/init.d/S99xkeen-autoreload \
  https://raw.githubusercontent.com/vladsaraykin/xkeen-autoreload-vless/main/scripts/S99xkeen-autoreload

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

Проверка процесса:

```sh
ps | grep xkeen-autoreload-vless
```

По умолчанию сервис проверяет подписку раз в час.

Если у провайдера изменились только ротируемые `SNI/SID`, XKeen не будет перезапущен.

## 6. Обновление приложения

Сначала остановите сервис:

```sh
/opt/etc/init.d/S99xkeen-autoreload stop
```

Скачайте новую версию во временный файл:

```sh
wget -O /opt/bin/xkeen-autoreload-vless.new \
  https://github.com/vladsaraykin/xkeen-autoreload-vless/releases/latest/download/xkeen-autoreload-vless-linux-mipsle
```

Дайте права:

```sh
chmod +x /opt/bin/xkeen-autoreload-vless.new
```

Замените бинарник:

```sh
mv /opt/bin/xkeen-autoreload-vless.new \
   /opt/bin/xkeen-autoreload-vless
```

Проверьте версию:

```sh
/opt/bin/xkeen-autoreload-vless version
```

Запустите сервис:

```sh
/opt/etc/init.d/S99xkeen-autoreload start
```

Использование временного файла позволяет избежать ошибки:

```text
Text file busy
```
