FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
