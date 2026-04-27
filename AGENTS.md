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
- 限流：使用 `kv_store` 表（非 Redis），详见下方「SQLite 替代 Redis 方案」
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

## SQLite 替代 Redis 方案

用一张 `kv_store` 表模拟 Redis KV 存储，通过 `expires_at` 字段实现 TTL 过期。

**表结构** (`model/kv_store.go`):

| 字段 | 类型 | 说明 |
|---|---|---|
| `key` | string (PK, size:200) | 键名，如 `login_fail:user@example.com` |
| `value` | text | 值内容，统一存字符串 |
| `expires_at` | *time.Time (indexed) | 过期时间，NULL = 永不过期 |

**API 映射** (`repository/kv_repository.go`):

| Redis 命令 | KVRepository 方法 | 实现方式 |
|---|---|---|
| `SET key value [EX ttl]` | `Set(key, value, ttl)` | GORM `Save` 实现 UPSERT |
| `GET key` | `Get(key)` | 查询，不存在返回 `gorm.ErrRecordNotFound` |
| `DEL key [key ...]` | `Del(keys ...)` | 批量删除 |
| `INCR key` | `Incr(key)` | 事务内读-改-写，键不存在从 1 开始 |
| `EXISTS key` | `Exists(key)` | Count 查询 |
| `EXPIRE key ttl` | `Expire(key, ttl)` | 更新 `expires_at` |

**TTL 过期机制 — 惰性清理（lazy eviction）**：
- 每次操作前调用 `cleanExpired()`：`DELETE FROM kv_store WHERE expires_at IS NOT NULL AND expires_at < now()`
- 没有后台定期清理线程或定时任务
- 从未被再次访问的过期 key 永远不会被清理，表会缓慢膨胀

**原子性**：`Incr` 使用 GORM 事务保证一致性。在 SQLite 单写模式（`max_open_conns=1`）下事务串行执行，等价于原子操作。换 PostgreSQL/MySQL 时需改用 `SELECT ... FOR UPDATE`。

**三大业务场景**：

1. **登录失败计数 + 账户锁定** (`auth_service.go`)
   - `login_fail:{email}` → 失败次数，TTL 30min，`Incr` 递增
   - `login_lock:{email}` → 锁定标记 "1"，TTL 30min，`Exists` 检查
   - 失败 ≥ 5 次 → 写锁定标记 + 删除计数器；密码正确 → 删除计数器

2. **密码重置令牌** (`auth_service.go`)
   - `reset_token:{token}` → 用户ID，TTL 24h，一次性令牌
   - `ForgotPassword` 写入 → `ResetPassword` 读取验证后删除

3. **请求限流** (`middleware/ratelimit.go`)
   - `ratelimit:{IP}` → 请求计数，首次 `Incr` 后设 `Expire(window)`
   - count > limit → 返回 429

**依赖注入链路**：`kvRepo` 仅注入到 `AuthService` 和 `middleware.RateLimit`，职责边界清晰。

**局限**：
- 惰性清理死角：从未再被访问的过期 key 不会被删除
- `cleanExpired()` 每次操作前执行 DELETE 全表扫描，高频场景有性能开销（`expires_at` 有索引缓解）
- `Incr` 对非数字值静默重置为 1，无类型校验
- SQLite 单进程部署，无法水平扩展

## 用户认证体系

### 数据模型 (`model/user.go`)

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| Email | string (unique, size:255) | 登录标识，唯一索引 |
| Password | string (size:255, json:"-") | bcrypt(cost=10) 哈希，不返回前端 |
| Nickname | string (size:100) | 显示昵称 |
| Role | string (default:"user") | 角色：`user` 或 `admin` |
| IsLocked | bool (default:false) | 管理员手动锁定标记 |
| MFASecret | string (size:255, json:"-") | TOTP 密钥（Base32） |
| MFAEnabled | bool (default:false) | 是否启用 MFA |
| Language | string (default:"zh-CN") | 语言偏好 |
| Timezone | string (default:"Asia/Shanghai") | 时区偏好 |
| DeletedAt | gorm.DeletedAt | 软删除 |

### API 路由

**公开路由**（无需 Token）：

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/auth/register` | POST | 注册（成功即返回 TokenPair，自动登录） |
| `/api/v1/auth/login` | POST | 登录 |
| `/api/v1/auth/refresh` | POST | 刷新 Token |
| `/api/v1/auth/forgot-password` | POST | 发送重置邮件 |
| `/api/v1/auth/reset-password` | POST | 用令牌重置密码 |

**认证后路由**（需 Bearer Token）：

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/auth/profile` | GET/PUT | 获取/更新个人信息 |
| `/api/v1/auth/password` | PUT | 修改密码（需旧密码） |
| `/api/v1/mfa/setup` | POST | 初始化 MFA（生成密钥+二维码URL） |
| `/api/v1/mfa/enable` | POST | 启用 MFA（验证 TOTP 码） |
| `/api/v1/mfa/disable` | POST | 禁用 MFA（验证 TOTP 码+清除密钥） |
| `/api/v1/mfa/verify` | POST | 验证 MFA 码 |
| `/api/v1/mfa/status` | GET | MFA 状态 |

**管理员路由**（需 admin 角色，`/api/v1/users/*`）：

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/users` | GET | 用户列表 |
| `/api/v1/users/:id` | GET/PUT/DELETE | 用户详情/更新/软删除 |
| `/api/v1/users/:id/role` | PUT | 变更角色（user ↔ admin） |
| `/api/v1/users/:id/lock` | POST | 锁定用户（`is_locked=true`） |
| `/api/v1/users/:id/unlock` | POST | 解锁用户 |

### 核心流程

**注册**：检查邮箱唯一 → bcrypt 加密 → INSERT → 返回 TokenPair（注册即登录）

**登录**（含暴力破解防护）：
1. 检查 `login_lock:{email}` 是否存在（kv_store 临时锁定）→ 存在则拒绝
2. 查找用户 → 不存在返回 `ErrInvalidCredential`（不暴露"用户不存在"）
3. 密码错误 → `Incr(login_fail:{email})`，首次设 Expire 30min；失败 ≥ 5 次 → 写锁定标记 + 删除计数器
4. 密码正确 → `Del(login_fail:{email})` → 生成 TokenPair

**Token 机制**：双 Token，结构相同（`Claims{UserID, Email}`），签名算法 HMAC-SHA256，密钥相同
- AccessToken：TTL 15min，用于 API 认证
- RefreshToken：TTL 168h（7天），用于续期
- 前端自动续期：401 → 用 refreshToken 调 `/auth/refresh` → 其他请求排队等待 → 新 Token 到达后统一重发

**密码重置**：
- `ForgotPassword` → 生成 32 字节随机 token → `kv_store: Set(reset_token:{token}, userID, 24h)` → 发邮件（发送失败不阻断）
- `ResetPassword` → `kv_store: Get` 取出 userID → 验证 → 更新密码 → `Del` 令牌（一次性）

**MFA**：基于 TOTP（`pquerna/otp`），兼容 Google Authenticator
- Setup → 生成 20 字节随机密钥 → Base32 编码 → 写入 `user.mfa_secret`（未启用） → 返回 `otpauth://totp` URL
- Enable → 验证 TOTP 6 位码 → `mfa_enabled=true`
- Disable → 验证码 → `mfa_enabled=false` + 清除 `mfa_secret`

### 已知问题

1. **`IsLocked` 未在登录流程中检查**：管理员通过 `/users/:id/lock` 设置的 `is_locked=true` 不会阻止登录，只有 kv_store 的临时锁定（连续失败 5 次）才生效。两套锁定机制未打通。

2. **MFA 未接入登录流程**：`Login` 方法没有检查 `user.MFAEnabled`，也没有要求二次验证。启用 MFA 的用户仍只需密码即可登录。

3. **备用码功能未实现**：`model/backup_code.go` 已定义且已注册 AutoMigrate，但无生成/验证的业务代码。

4. **AccessToken 和 RefreshToken 结构完全相同**：仅靠过期时间区分，`ParseAccessToken` 和 `ParseRefreshToken` 内部调用同一个 `parseToken`。未过期的 RefreshToken 也能当作 AccessToken 使用。

5. **管理员权限校验不一致**：`UserController` 在每个方法内检查 `ctx.GetString("role")`，但 `middleware.Auth` 只注入 `user_id` 和 `email`，没有注入 `role`。`/users/*` 路由组也未挂载 `middleware.Admin`。因此管理员用户管理的权限检查可能失效。
