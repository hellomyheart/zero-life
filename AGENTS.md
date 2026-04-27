# AGENTS.md

## 项目

Zero-Life — 个人财务管理系统，Firefly III 的 Go + Vue 重写版。README 和 UI 为中文，代码注释为中文。

## 命令

```bash
# 后端（在 server/ 目录下执行）
go run ./cmd/server/          # 开发服务器 :8080
go vet ./...                  # 代码检查
go build -o server ./cmd/server/  # 编译

# 前端（在 web/ 目录下执行）
npm run dev                   # 开发服务器 :5173（Vite 代理 /api → :8080）
npm run build                 # vue-tsc 类型检查 + vite 构建（输出到 dist/）
# 无单独的 lint 或 test 命令

# Swagger 文档（在 server/ 目录下执行）
swag init -g cmd/server/main.go -o docs --outputTypes go,json

# Docker（在项目根目录执行）
docker compose up -d          # api(:8080) + web(:80)
docker compose up -d --build  # 代码更新后重新构建
```

## 架构

**后端** (`server/`) — Go 模块 `github.com/hellomyheart/zero-life/server`
- 入口：`cmd/server/main.go` — 手动依赖注入（无 DI 框架）
- 分层：`controller` → `service` → `repository` → GORM 模型
- DTO 拆分为 `dto/request/` 和 `dto/response/`
- 共享包在 `internal/pkg/`：`errcode`、`jwt`、`pagination`、`totp`、`validator`、`webhook`、`email`、`hash`
- 配置：Viper 读取 `config.yaml`，支持环境变量覆盖（`viper.AutomaticEnv()`）
- 数据库：SQLite + WAL 模式，启动时自动迁移（`migrations/` 目录为空，无手动迁移文件）
- 认证：JWT Bearer Token；`middleware.Auth` 将 `user_id`/`email` 注入 Gin 上下文；`middleware.Admin` 检查 `user.Role == "admin"`
- 限流：使用 `kv_store` 表（非 Redis）
- 定时任务端点：`GET /api/v1/cron/:token` — 基于 token 验证，非 JWT

**前端** (`web/`) — Vue 3 + TypeScript + Vite
- 路径别名：`@` → `src/`
- API 层：`src/api/` — 每个领域一个文件，使用 `src/utils/request.ts`（Axios，自动刷新 Token）
- 状态管理：Pinia stores 在 `src/stores/`
- 国际化：中/英在 `src/i18n/locales/`
- 所有页面懒加载路由；认证守卫在 `src/router/index.ts`
- 无测试框架

## 关键约定

- **错误码**：按模块分段定义在 `internal/pkg/errcode/errcode.go`（1xxxx=认证, 2xxxx=账户, 3xxxx=交易, 4xxxx=分类, 5xxxx=标签, 6xxxx=预算, 7xxxx=账单, 8xxxx=规则, 9xxxx=导入, 10xxxx=定期交易, 11xxxx=Webhook, 12xxxx=对象组, 13xxxx=交易链接, 14xxxx=偏好, 15xxxx=对账, 16xxxx=MFA）
- **API 响应格式**：统一为 `{ "code": 0, "message": "success", "data": ... }`，通过 `controller.Success()` / `controller.Error()` 返回
- **金额处理**：Go 用 `shopspring/decimal`，TS 用 `decimal.js` — 禁止用浮点数表示金额
- **无测试套件**：前后端均未配置测试
- **Swagger**：生成到 `server/docs/`；修改 API 后必须重新生成

## 注意事项

- SQLite `max_idle_conns` 和 `max_open_conns` 默认为 1 — 这是 SQLite 单写模式的刻意设计，不要修改
- `server/cmd/server/config.yaml` 和 `server/config.yaml` 同时存在 — 运行时从 CWD 读取配置，所以必须在 `server/` 目录下启动
- 数据库文件（`*.db`、`*.db-shm`、`*.db-wal`）和 `*.exe` 已加入 gitignore
- Docker Compose 从 `.env` 读取 `JWT_SECRET`（默认值为不安全的 `change-me-in-production`）
- 前端 `nginx.conf` 硬编码 `proxy_pass http://api:8080` — 在 Docker Compose 外使用需修改
- `Recurrence` 和 `RecurringTransaction` 是两个独立的模型/控制器，路由不同（`/recurrences` vs `/recurring-transactions`），不要混淆
