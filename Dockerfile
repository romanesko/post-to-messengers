# The build stage
FROM golang:1.23.2-alpine AS builder
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY internal internal
COPY main.go .
COPY go.mod .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app /app/main.go

# The run stage
FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app .
CMD ["./app"]