FROM golang:1.21-alpine

WORKDIR /app

# Устанавливаем зависимости для компиляции и netcat
RUN apk add --no-cache gcc musl-dev postgresql-client netcat-openbsd

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]