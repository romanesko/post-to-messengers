build-amd:
	GOOS=linux GOARCH=amd64 go build -o out/xmpp-sender-amd64

build-arm:
	GOOS=darwin GOARCH=arm64 go build -o out/xmpp-sender-arm64

run-amd:
	set -a && source .env && set +a && ./out/xmpp-sender-amd64

run-arm:
	set -a && source .env && set +a && ./out/xmpp-sender-arm64