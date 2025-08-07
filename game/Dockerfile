FROM golang:1.24.1 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЕ файлы (включая внутренние пакеты)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]