# ============================================================
# GoWork Docker Makefile
# ============================================================
.PHONY: help up down restart logs clean status

help:
	@echo "GoWork Docker 环境"
	@echo ""
	@echo "  make up       - 启动所有服务"
	@echo "  make infra    - 只启动基础设施(MySQL/Redis/Postgres)"
	@echo "  make down     - 停止所有服务"
	@echo "  make restart  - 重启所有服务"
	@echo "  make logs     - 查看日志"
	@echo "  make status   - 查看状态"
	@echo "  make clean    - 清理数据"

up:
	docker compose up -d
	@echo "✓ 服务已启动"
	@echo "  gotour:   http://localhost:8080"
	@echo "  MySQL:    localhost:3306"
	@echo "  Redis:    localhost:6379"
	@echo "  Postgres: localhost:5432"

infra:
	docker compose up -d mysql redis postgres
	@echo "✓ 基础设施已启动"

down:
	docker compose down
	@echo "✓ 服务已停止"

restart:
	docker compose restart

logs:
	docker compose logs -f

logs-gotour:
	docker compose logs -f gotour

status:
	docker compose ps

clean:
	docker compose down -v
	@echo "✓ 已清理所有数据和容器"
