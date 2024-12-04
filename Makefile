# Папки для .proto файлов и сгенерированных файлов
PROTO_DIR := proto
GEN_DIR := gen/go

# Путь к дополнительным .proto файлам (включая google/api)
INCLUDE_DIR := proto

# URL базы данных и путь к миграциям
DATABASE_URL := postgres://postgres:123@localhost:5432/postgres?sslmode=disable
MIGRATIONS_PATH := ./migrations

ifeq ($(ENV),production)
    DATABASE_URL := postgres://postgres:123@auth_db:5432/auth?sslmode=disable
endif


# Файлы .proto, которые будут сгенерированы
PROTO_FILES := $(wildcard $(PROTO_DIR)/**/*.proto)

# Преобразуем каждый .proto файл в соответствующие pb.go, pb.grpc.go и pb.gw.go файлы в папке GEN_DIR
GEN_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(GEN_DIR)/%.pb.go,$(PROTO_FILES))
GEN_GRPC_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(GEN_DIR)/%.pb.grpc.go,$(PROTO_FILES))
GEN_GATEWAY_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(GEN_DIR)/%.pb.gw.go,$(PROTO_FILES))

# Команда генерации
PROTOC := protoc
PROTOC_FLAGS := --go_out=$(GEN_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative
PROTOC_GATEWAY_FLAGS := --grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=paths=source_relative,logtostderr=true

# Основная цель по умолчанию - сгенерировать все файлы
all: generate

# Генерация файлов
generate: $(GEN_FILES) $(GEN_GRPC_FILES) $(GEN_GATEWAY_FILES)

# Правило генерации .pb.go, .pb.grpc.go и .pb.gw.go файлов
$(GEN_DIR)/%.pb.go $(GEN_DIR)/%.pb.grpc.go $(GEN_DIR)/%.pb.gw.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(GEN_DIR)
	$(PROTOC) -I $(PROTO_DIR) -I $(INCLUDE_DIR) $(PROTOC_FLAGS) $(PROTOC_GATEWAY_FLAGS) $<

# Очистка сгенерированных файлов
clean:
	rm -rf $(GEN_DIR)

# Регенерация всех файлов
rebuild: clean generate

# Запуск миграций
migrate:
	go run ./cmd/migrator -database-url "$(DATABASE_URL)" -migrations-path "$(MIGRATIONS_PATH)"
