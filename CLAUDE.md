# RoadBarber Backend（后端 API）

Go Fiber + GORM + PostgreSQL 后端服务，按业务模块化组织（customer / merchant / admin / common）。

## 常用命令

```bash
go mod tidy           # 同步依赖
go run cmd/main.go    # 启动服务（默认 :8080）
go build ./...        # 编译检查
go test ./...         # 运行测试
```

## 分支规范

| 分支 | 用途 |
|------|------|
| `main` | 生产环境分支，始终保持可发布状态，仅接受来自 `release-*` 或 `hotfix-*` 的合并 |
| `develop` | 主开发分支，包含所有已完成的功能 |
| `feature/*` | 功能开发分支，从 `develop` 切出，完成后合并回 `develop` |
| `release/*` | 发布准备分支，用于版本号 bump / 小修，从 `develop` 切出 |
| `hotfix/*` | 紧急修复分支，从 `main` 切出，修复后同时合并回 `main` 和 `develop` |

分支命名示例：
- `feature/customer-booking-flow`
- `feature/admin-merchant-verify`
- `release/v1.0.0`
- `hotfix/fix-jwt-expire`

## Commit 提交规范（Angular 规范）

所有 commit 消息必须使用以下格式，**标题使用中文简体**：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type（必填）

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复 bug |
| `docs` | 仅文档变更 |
| `style` | 代码格式调整（不影响逻辑） |
| `refactor` | 重构（既非新功能也非 bug 修复） |
| `perf` | 性能优化 |
| `test` | 增加或修改测试 |
| `build` | 构建系统或外部依赖变更 |
| `ci` | CI 配置文件变更 |
| `chore` | 其他不修改 src 或测试文件的变更 |

### Scope（可选，描述影响范围）

常用 scope：`auth`、`customer`、`merchant`、`admin`、`common`、`model`、`middleware`、`config`、`db`

### Subject（必填，简短描述）

- 使用中文简体
- 不超过 50 个字符
- 不使用句号结尾
- 使用动词开头

### 示例

```bash
git commit -m "feat(customer): 新增预约创建接口"
git commit -m "fix(auth): 修复 JWT token 过期时间错误"
git commit -m "refactor(merchant): 重构排班服务的事务逻辑"
git commit -m "docs: 更新 API 接口文档"
git commit -m "chore: 升级 fiber 到 v2.52.0"
```

带 body 和 footer 的完整示例：

```
feat(customer): 新增商家收藏接口

支持顾客收藏 / 取消收藏商家，列表查询收藏商家。
使用唯一索引防止重复收藏。

Closes #123
```

## 目录约定

```
backend/
├── cmd/                 # 入口
├── internal/
│   ├── config/          # 配置
│   ├── models/          # GORM 模型
│   ├── middleware/      # 中间件
│   ├── modules/         # 业务模块
│   │   ├── common/      # 公共：认证、地区
│   │   ├── customer/    # 顾客端：商家/预约/评价/收藏
│   │   ├── merchant/    # 商家端：入驻/排班/预约
│   │   └── admin/       # 管理员端：仪表盘/审核
│   └── routes/          # 路由注册
├── pkg/                 # 公共包
└── migrations/          # SQL 迁移
```