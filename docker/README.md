# ============================================================
# GoWork Docker 环境使用说明
# ============================================================

## 目录结构

```
GoWork/
├── docker/                    # Docker 配置（统一环境）
│   ├── Dockerfile            # 开发镜像
│   ├── docker-compose.yml    # 统一配置
│   ├── Makefile             # 快捷命令
│   └── README.md
│
├── gotour/                   # 项目1
├── todo-app/                 # 项目2（可扩展）
└── pkg/                     # 公共包
```

## 架构说明

```
┌─────────────────────────────────────────────────┐
│              Docker Network (gowork_net)         │
│                                                   │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────┐  │
│  │  MySQL  │  │  Redis  │  │    Postgres      │  │
│  │ :3306   │  │ :6379   │  │    :5432         │  │
│  └────┬────┘  └────┬────┘  └────────┬────────┘  │
│       │            │                  │           │
│       └────────────┴──────────────────┘           │
│                    │                              │
│       ┌────────────┼──────────────────┐           │
│       │            │                  │           │
│  ┌────▼────┐  ┌───▼────┐  ┌────────▼────────┐    │
│  │ gotour  │  │todo-app│  │   other app     │    │
│  │  :8080  │  │  :8081 │  │                 │    │
│  └─────────┘  └─────────┘  └─────────────────┘    │
│                                                   │
└───────────────────────────────────────────────────┘
```

## 快速开始

### 1. 启动统一环境（基础设施 + gotour）
```bash
cd GoWork/docker
docker compose up -d
```

### 2. 进入项目容器开发
```bash
# gotour 项目
docker exec -it gotour_dev /bin/sh

# 在容器内运行
go run ./gotour/ch18
```

### 3. 停止环境
```bash
docker compose down
```

## 添加新项目

编辑 `docker-compose.yml`，添加新服务：

```yaml
services:
  # ... 现有服务 ...

  myapp:
    build:
      context: ../myapp
      dockerfile: ../docker/Dockerfile.dev
    container_name: myapp_dev
    ports:
      - "8081:8080"      # 换一个端口
    volumes:
      - ../myapp:/app
    depends_on:
      - mysql
      - redis
    networks:
      - gowork_net
```

## 数据库连接

在项目代码中使用：

```go
// MySQL
dsn := "root:root123@tcp(gowork_mysql:3306)/gowork?charset=utf8mb4"

// Redis
redisAddr := "gowork_redis:6379"

// 使用容器名作为 host
```

## 常用命令

```bash
cd GoWork/docker

# 启动所有服务
docker compose up -d

# 只启动基础设施
docker compose up -d mysql redis

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f

# 进入容器
docker exec -it gotour_dev /bin/sh

# 重新构建
docker compose build --no-cache gotour

# 停止
docker compose down

# 清理数据
docker compose down -v
```

## 热重载

gotour 容器使用 `air` 工具，代码修改后自动重编译。

连接 VSCode 调试：
```json
{
  "name": "Docker Delve",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remoteHost": "localhost",
  "remotePort": "2345"
}
```
