# Используем официальный образ Golang для сборки
FROM golang:1.21 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod tidy

# Собираем бинарный файл
RUN go build -o orchestrator ./cmd/orchestrator

# Создаем минимальный финальный образ
FROM debian:latest
WORKDIR /root/
COPY --from=builder /app/orchestrator .

# Открываем порт для API
EXPOSE 8080

# Запускаем оркестратор
CMD ["./orchestrator"]
