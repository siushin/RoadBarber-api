# RoadBarber API

> Go Fiber + GORM + PostgreSQL 模块化单体后端。
> 业务模块：`common`（认证/地区）/ `customer` / `merchant` / `admin`

## 目录结构

```
api/
├── cmd/                    # 入口（main.go 启动时自动迁移）
├── internal/
│   ├── config/             # 全局配置 + DB/Redis 初始化
│   ├── migrate/            # golang-migrate 库封装
│   ├── models/             # 12 张表的 GORM 模型
│   ├── middleware/         # Auth / CORS / Logger
│   ├── modules/            # 业务模块（每个含 handler/service/router）
│   │   ├── common/
│   │   ├── customer/
│   │   ├── merchant/
│   │   └── admin/
│   └── routes/             # 路由聚合注册
├── pkg/                    # 公共包（response / utils）
├── migrations/             # SQL 迁移（golang-migrate 管理）
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_seed_demo.up.sql
│   └── 000002_seed_demo.down.sql
├── Makefile                # 迁移 / 运行快捷命令
├── .env                    # 本地配置（不入 git）
├── .env.example            # 配置模板
└── go.mod
```

## 常用命令

```bash
# 同步依赖
go mod tidy

# 启动服务（默认 :8080；启动时自动执行 migrate up）
go run cmd/main.go

# 编译检查
go build ./...
```

## 数据库迁移

使用 [golang-migrate](https://github.com/golang-migrate/migrate) 管理 schema。

### 工具安装（首次）

```bash
# 方式 A：brew
brew install golang-migrate

# 方式 B：go install（备选）
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# 产物在 $(go env GOPATH)/bin/migrate，建议把 ~/go/bin 加进 PATH
```

### Makefile 命令

```bash
cd api

# 升级到最新版本（启动后端也会自动执行）
make migrate/up

# 回滚所有迁移（删除所有表 + 清空种子）
make migrate/down

# 一键彻底重置：down -all + up（开发期最常用）
make migrate/reset

# 仅灌种子（增量，等同 make migrate/up）
make migrate/seed

# 查看当前版本
make migrate/version

# 强制把库标记成版本 N（出错修复用，慎用）
make migrate/force V=1

# 物理删库 + 重建空库（不会自动迁移，仅清理）
make migrate/drop
```

### 手动使用 migrate CLI

```bash
# 连接串从 .env 读取
export DATABASE_URL="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

migrate -path ./migrations -database "$DATABASE_URL" up
migrate -path ./migrations -database "$DATABASE_URL" down -all
migrate -path ./migrations -database "$DATABASE_URL" version
migrate -path ./migrations -database "$DATABASE_URL" force 1
```

### 新增迁移

文件命名格式：`<版本号6位>_<名称>.{up|down}.sql`

```bash
# 示例：新增第三个迁移
touch migrations/000003_add_xxx.up.sql
touch migrations/000003_add_xxx.down.sql
```

down 迁移要可逆（删除表/字段/数据），否则留空文件即可。

### 演示账号（执行 `make migrate/up` 后自动写入）

| 角色 | 手机号 | 密码 | 备注 |
|------|--------|------|------|
| 管理员 | `13800000000` | `admin123` | 可登录管理后台 |
| 顾客 | `13900000001` / `13900000002` | `customer123` | 演示下单 |
| 商家 | `13900000011` (Tony) | `merchant123` | 南山店，已发布排班 |
| 商家 | `13900000012` (Kevin) | `merchant123` | 天河店，已发布排班 |

演示数据还包括：
- 广东省 / 深圳市 / 广州市 / 南山区 / 天河区 地区节点
- 3 个服务项目（精剪男 / 精剪女 / 冷烫）
- 5 个未来排班时段

## 注意事项

- **后端进程占用数据库时 `DROP DATABASE` 会失败**：迁移前先停服务，或在 psql 里执行 `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'barber';`。
- **migrate 启动失败留下 dirty 状态**：`make migrate/down` 中途失败时，库可能处于 dirty 状态。可用 `make migrate/force V=1`（填上一个干净版本）修复后重试。
- **.env 不入 git**：使用前先 `cp .env.example .env` 并填入本机数据库账号密码。示例：

```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=barber
DB_PASSWORD=<your-password>
DB_NAME=barber
```

## 快速开始（新机器）

```bash
# 1. 装 PostgreSQL 并创建数据库
brew install postgresql@15
brew services start postgresql@15
createuser -s barber             # 或在 psql 里 CREATE ROLE barber WITH LOGIN SUPERUSER PASSWORD '...';
createdb -O barber barber         # 或 CREATE DATABASE barber OWNER barber;

# 2. 装 migrate CLI
brew install golang-migrate

# 3. 配置 + 启动
cd api
cp .env.example .env
# 修改 .env 里的 DB_USER / DB_PASSWORD / DB_NAME

go mod tidy
make migrate/reset        # 建表 + 灌种子
make run                  # 启动后端
```

服务跑起来后访问 `http://localhost:8080/health` 验证。