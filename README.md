# REST service to send XMPP messages

### .env example

```
SERVICE_PORT=8080
XMPP_SERVER=localhost
XMPP_USERNAME=username
XMPP_PASSWORD=password
```

### build

```bash
go build -o out/xmpp-sender-amd64
```

### run

```bash
set -a && source .env && set +a && ./xmpp-sender
```


### extra

repo contains Makefile to build and run the service on amd64 and arm64 platforms