# wss-chat

Simple websocket chat backend (MVP)

- [wss-chat](#wss-chat)
  - [Задача](#задача)
    - [ТЗ](#тз)
    - [Результат](#результат)
  - [Documentation](#documentation)
  - [Start server](#start-server)
  - [Usage](#usage)
    - [API](#api)
    - [Chat usage](#chat-usage)

## Задача

### ТЗ

Необходимо реализовать чат сервис на WSS:

- Чат должен быть разделен на комнаты
- Пользователь, при входе в комнату, видит сообщения за последние N минут
- Участники комнаты видят, что зашел новый пользователь
- Каждый участник чата может посмотреть общий список комнат
- Каждый участник чата может заходить сразу в несколько комнат

Суть тестового задания - написать упрощенный Телеграмм. В чат можно пересылать любые текстовые сообщения. Загрузка файлов не требуется.

### Результат

- [x] Чат разделен на комнаты. Участники чата, не получают уведомления о действия в другом чате
- [x] Пользователь, при подключении к комнате, видит сообщения за запрошенное количество времени
- [x] О подлючении нового пользователя сообщается всем участникам комнаты
- [x] По API доступен запрос для получения всех комнат
- [x] Любой пользователь имеет возможность подключения к нескольким комнатам
- [x] База задокументирована
- [x] API задокументировано
- [ ] Приложение представляет собой минимально жизнеспособный продукт
  - база и приложение никак не сохраняют состояние
  - API минимально и реализует полный CRUDL

## Documentation

- [DB docs](https://dbdocs.io/miromax42/wss-chat)
- API swagger [endpoint](http://localhost:8080/swagger/index.html): `/swagger/index.html`

## Start server

> depencies: docker

```bash
docker compose up
```

## Usage

### API

>Examples for localhost:8080

- Get rooms

  ```bash
  curl --location --request GET 'http://localhost:8080/rooms'
  ```

### Chat usage

1. Connect to chat
   - Connect `ws://localhost:8080/ws?room=public`, where `public` is room you want to connect
   - Default time for message history is `1 minute`. But you can set it via form `ws://localhost:8080/ws?room=public&time=1h` (message history = 1 Hour)
2. Send message: json with required fields(`sender`, `payload`)

   ```json
    {
       "sender":"your_username",
       "payload":"test messasage!"
    }
   ```
