# Этап 1: сборка бинарника
FROM golang:1.25 AS builder

WORKDIR /app

# Копируем go-файлы и зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Сборка для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o app 

# Этап 2: финальный образ
FROM ubuntu:latest

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник и frontend
COPY --from=builder /app/app .
COPY --from=builder /app/web ./web

# Указываем порт
EXPOSE 7540

# Установка переменных окружения по умолчанию (можно переопределять при запуске)
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/todo.db
ENV TODO_PASSWORD=

# Точка входа
CMD ["./app"]