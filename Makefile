PROTO_DIR := proto
GEN_DIR := gen/go

DATABASE_URL := postgres://postgres:123@localhost:5432/postgres?sslmode=disable
MIGRATIONS_PATH := ./migrations

# Файлы .proto, которые будут сгенерированы
PROTO_FILES := $(wildcard $(PROTO_DIR)/**/*.proto)

# Преобразуем каждый .proto файл в соответствующие pb.go и pb.grpc.go файлы в папке GEN_DIR
GEN_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(GEN_DIR)/%.pb.go,$(PROTO_FILES))
GEN_GRPC_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(GEN_DIR)/%.pb.grpc.go,$(PROTO_FILES))

# Команда генерации
PROTOC := protoc
PROTOC_FLAGS := --go_out=$(GEN_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative

# Основная цель по умолчанию - сгенерировать все файлы
all: $(GEN_FILES) $(GEN_GRPC_FILES)

# Правило генерации .pb.go и .pb.grpc.go файлов
$(GEN_DIR)/%.pb.go $(GEN_DIR)/%.pb.grpc.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(GEN_DIR)
	$(PROTOC) -I $(PROTO_DIR) $(PROTOC_FLAGS) $<

# Очистка сгенерированных файлов
clean:
	rm -rf $(GEN_DIR)

# Регенерация всех файлов
rebuild: clean all

migrate:
	go run ./cmd/migrator -database-url "$(DATABASE_URL)" -migrations-path "$(MIGRATIONS_PATH)"