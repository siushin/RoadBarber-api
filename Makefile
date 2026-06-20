# ============================================================
# RoadBarber Backend Makefile
# ============================================================

# 加载 .env 中的变量
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATE       ?= migrate
MIGRATIONS    ?= ./migrations

# ============================================================
# 默认目标：打印帮助
# ============================================================
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run               - 启动后端"
	@echo "  make build             - 编译"
	@echo "  make migrate/up        - 升级到最新版本"
	@echo "  make migrate/down      - 回滚所有迁移"
	@echo "  make migrate/reset     - 彻底重置库（down all + up）"
	@echo "  make migrate/version   - 查看当前版本"
	@echo "  make migrate/force N   - 强制把库标记成版本 N"
	@echo "  make migrate/seed      - 同 migrate/up，单独别名"

# ============================================================
# 运行
# ============================================================
.PHONY: run build tidy
run:
	go run cmd/main.go

build:
	go build ./...

tidy:
	go mod tidy

# ============================================================
# 数据库迁移
# ============================================================
.PHONY: migrate/up migrate/down migrate/reset migrate/version migrate/force migrate/seed migrate/drop

migrate/up:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DATABASE_URL)" up

migrate/down:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DATABASE_URL)" down -all

migrate/reset: migrate/down migrate/up
	@echo "✅ Database reset complete"

migrate/version:
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DATABASE_URL)" version

migrate/force:
	@if [ -z "$(V)" ]; then echo "Usage: make migrate/force V=1"; exit 1; fi
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DATABASE_URL)" force $(V)

migrate/seed: migrate/up

# 物理删除整个库（危险操作）
migrate/drop:
	@echo "⚠️  Dropping database $(DB_NAME)..."
	psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"
	@echo "✅ Database dropped and recreated"