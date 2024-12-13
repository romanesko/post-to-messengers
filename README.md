# REST service to send XMPP/Telegram/Matrix messages


## Prepare and run

create `docker-compose.yml` file with next content:

```yaml
services:
  http-xmpp:
    container_name: http-xmpp
    image: savarez/http-xmpp:latest
    environment:
      - XMPP_SERVER=xmpp.server
      - XMPP_USERNAME=xmpp_username
      - XMPP_PASSWORD=xmpp_password
      - TELEGRAM_BOT_TOKEN=123456789:BotToken
      - TELEGRAM_WELCOME_MESSAGE=Добро пожаловать!\n\nВаш id <code>{{tg_chat_id}}</code>\n\nПерешлите его нам, чтобы начать получать сообщения
      - MATRIX_SERVER=matrix.server
      - MATRIX_USERNAME=matrix_username
      - MATRIX_PASSWORD=matrix_password
      - USE_XMPP=1
      - USE_TELEGRAM=1
      - USE_MATRIX=1
    ports:
      - "8080:8080"
    restart: always
```

run 
```bash
docker-compose pull && docker-compose up -d
```

to see logs:
```bash
docker-compose logs -f
```


## Usage

#### XMPP:

```bash
curl -d '{"recipient": "test1@server","message": "hello there"}' 'http://localhost:8080/xmpp'
```

#### Telegram:

```bash
curl -d '{"recipient": 123456789,"message": "hello there"}' 'http://localhost:8080/telegram'
```

#### Matrix:

```bash
curl -d '{"recipient": "@test1:server","message": "hello there"}' 'http://localhost:8080/matrix'
```

### Service responses

`200 - Success`

```json
{
  "status": "sent",
  "message":"hello there"
}
```

`400 - Bad request`

```json
{
  "code": 1,
  "message": "Invalid request body"

}
```


| Code | Description                       |
|-----|-----------------------------------|
| 1   | Invalid request body              |
| 2   | Missing required fields           |
| 3   | Wrong Recipient format            |
| 4   | Unexpected messenger error        |
| 5   | The recipient has blocked the bot |
