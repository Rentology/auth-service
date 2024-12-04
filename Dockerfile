# Начальная стадия: загрузка модулей
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Загружаем зависимости
RUN go mod download

# Устанавливаем необходимые пакеты для работы с proto
RUN apk add --no-cache make protobuf-dev curl bash unzip

RUN curl --location --silent -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v29.0/protoc-29.0-linux-aarch_64.zip && \
    unzip protoc.zip -d /usr/local/ && \
    rm -fr protoc.zip

# Установка Go плагинов для protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Добавьте Go бинарники в PATH
ENV PATH=$PATH:/go/bin



RUN make generate

# Компиляция приложения
RUN go build -o /app/main ./cmd/auth-service

# Финальная стадия: запуск приложения
FROM golang:1.23-alpine AS runner

# Копируем скомпилированное приложение и другие файлы
COPY --from=builder /app /app

# Копируем скрипт wait-for-it в контейнер
COPY wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

# Устанавливаем make в контейнере для выполнения миграций
RUN apk add --no-cache make

WORKDIR /app

# Устанавливаем переменные окружения
ENV CONFIG_PATH=./config/production.yaml
ENV DATABASE_URL=postgres://postgres:123@auth_db:5432/auth?sslmode=disable

# Ожидаем пока БД будет готова и выполняем миграции, затем запускаем приложение
CMD ["sh", "-c", "wait-for-it auth_db:5432 -- make migrate && ./main"]

EXPOSE 8080
