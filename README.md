# xkeen-autoreload-vless

Небольшой Go-сервис для Keenetic/Entware, который:

- скачивает VLESS-подписку в Base64;
- при первом запуске предлагает выбрать страну;
- генерирует `04_outbounds.json` для XKeen;
- перезапускает XKeen только при реальном изменении ноды.

Провайдер может менять `sni` и `sid` при каждом запросе подписки. Чтобы не делать лишние перезапуски, стабильная идентичность ноды определяется по:

`UUID | address | port | publicKey`

Используйте `--force`, если нужно принудительно получить свежие `sni`/`sid`, даже если сама нода не изменилась.

## Первый запуск и настройка

Ручное редактирование env-файла не требуется.

Запустите:

```sh
/opt/bin/xkeen-autoreload-vless setup
```

Мастер настройки:

1. попросит ссылку на подписку;
2. скачает Base64-ответ;
3. декодирует и распарсит VLESS-ноды;
4. найдёт уникальные страны по названиям нод;
5. покажет пронумерованный список стран;
6. сохранит выбранную страну и ссылку в `/opt/etc/xkeen-autoreload-vless.env` с правами `0600`.

Пример:

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

Если запустить:

```sh
xkeen-autoreload-vless
```

или:

```sh
xkeen-autoreload-vless run
```

без сохранённой конфигурации, мастер настройки также запустится автоматически, если приложение работает в интерактивном терминале.

Чтобы позже поменять ссылку или страну, снова выполните:

```sh
xkeen-autoreload-vless setup
```

## Команды

```sh
xkeen-autoreload-vless setup
xkeen-autoreload-vless run
xkeen-autoreload-vless update
xkeen-autoreload-vless update --force
xkeen-autoreload-vless check
xkeen-autoreload-vless version
```

Подробнее об установке на Keenetic/Entware: [INSTALL.md](INSTALL.md).

## Сборка для Keenetic / MIPSLE

```sh
make test
make build-mipsle
```

Результат:

```text
dist/xkeen-autoreload-vless-linux-mipsle
```

## Безопасность

Конфигурационный файл читается самим Go-приложением как данные и не выполняется через shell.

Перед заменой конфигурации XKeen текущий файл сохраняется в:

```text
04_outbounds.json.bak
```

Если `xkeen -restart` завершится ошибкой, приложение восстановит предыдущий конфиг и повторно перезапустит XKeen.
