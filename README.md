# REST service to send XMPP messages

### .env file example

```
SERVICE_PORT=8080
XMPP_SERVER=localhost
XMPP_USERNAME=username
XMPP_PASSWORD=password
```

### build

```bash
go build -o out/xmpp-sender
```

### run

```bash
set -a && source .env && ./xmpp-sender
```


### extra

repo contains `Makefile` to build and run the service on amd64 and arm64 platforms