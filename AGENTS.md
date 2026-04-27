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
- 认证：JWT Bearer Token；`middleware.Auth` 将 `user_id`/`email` 注入 Gin 上下文；`middleware.Admin` 查数据库校验 `user.Role == "admin"` 并注入 `role` 到上下文，挂载在 `/users/*` 路由组
- 限流：使用 `kv_store` 表（非 Redis），详见下方「SQLite 替代 Redis 方案」
- 定时任务：内置 `robfig/cron` 调度器，每天 00:00 执行循环交易和到期账单，每天 09:00 生成预算历史快照

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
| `/api/v1/auth/login` | POST | 登录（MFA 用户返回临时令牌） |
| `/api/v1/auth/refresh` | POST | 刷新 Token |
| `/api/v1/auth/forgot-password` | POST | 发送重置邮件 |
| `/api/v1/auth/reset-password` | POST | 用令牌重置密码 |
| `/api/v1/auth/mfa-verify` | POST | MFA 登录二次验证（用临时令牌+TOTP码换 TokenPair） |

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

**注册**：检查邮箱唯一 → 判断是否首个用户（`authRepo.Count() == 0` 则 `Role = "admin"`）→ bcrypt 加密 → INSERT → 返回 TokenPair（注册即登录）。首个注册用户自动成为管理员，后续注册用户默认为普通用户。

**登录**（含暴力破解防护 + MFA）：
1. 根据邮箱查找用户 → 不存在返回 `ErrInvalidCredential`（不暴露"用户不存在"）
2. 检查 `user.IsLocked`（管理员手动锁定）→ 锁定则拒绝
3. 检查 `login_lock:{email}` 是否存在（kv_store 临时锁定）→ 存在则拒绝
4. 密码错误 → `Incr(login_fail:{email})`，首次设 Expire 30min；失败 ≥ 5 次 → 写锁定标记 + 删除计数器
5. 密码正确 → `Del(login_fail:{email})`
6. 若 `user.MFAEnabled` → 生成 32 字节随机 mfa_token → `kv_store: Set(mfa_token:{token}, userID, 5min)` → 返回 `{mfa_required: true, access_token: mfa_token}`
7. 若未启用 MFA → 直接生成 TokenPair 返回

**MFA 登录二次验证**：
- 前端收到 `mfa_required: true` 后展示 TOTP 码输入框
- 用户输入 6 位码后调 `POST /auth/mfa-verify {mfa_token, code}`
- 后端从 kv_store 取出 userID → 验证 TOTP 码 → 删除临时令牌 → 返回真正的 TokenPair

**Token 机制**：双 Token，`Claims{UserID, Email, TokenType}`，签名算法 HMAC-SHA256，密钥相同
- AccessToken：TTL 15min，TokenType="access"，用于 API 认证
- RefreshToken：TTL 168h（7天），TokenType="refresh"，用于续期
- `ParseAccessToken` 和 `ParseRefreshToken` 会校验 `token_type` 字段，RefreshToken 不能当作 AccessToken 使用
- 前端自动续期：401 → 用 refreshToken 调 `/auth/refresh` → 其他请求排队等待 → 新 Token 到达后统一重发

**密码重置**：
- `ForgotPassword` → 生成 32 字节随机 token → `kv_store: Set(reset_token:{token}, userID, 24h)` → 发邮件（发送失败不阻断）
- `ResetPassword` → `kv_store: Get` 取出 userID → 验证 → 更新密码 → `Del` 令牌（一次性）

**MFA**：基于 TOTP（`pquerna/otp`），兼容 Google Authenticator
- Setup → 生成 20 字节随机密钥 → Base32 编码 → 写入 `user.mfa_secret`（未启用） → 返回 `otpauth://totp` URL
- Enable → 验证 TOTP 6 位码 → `mfa_enabled=true` → 生成 10 个备用码（bcrypt 哈希存储，明文仅此一次返回前端）
- Disable → 验证码 → `mfa_enabled=false` + 清除 `mfa_secret` + 清除备用码
- 备用码格式：`XXXX-XXXX`（8位大写字母数字），登录 MFA 验证时 TOTP 失败会自动尝试备用码
- 重新生成备用码：`POST /mfa/backup-codes`，旧码全部失效

### 预算管理

**数据模型**：
- `Budget`：`name`、`amount`（限额）、`period`（daily/weekly/monthly/quarterly/yearly）、`is_enabled`，与 Category 多对多
- `BudgetCategory`：多对多关联表，复合主键 `(budget_id, category_id)`
- `BudgetHistory`：历史快照，`period_start`/`period_end`、`amount`（限额）、`spent`（实际支出）

**支出计算**：`calculateSpent` 根据当前周期范围查询关联分类（含后代展开）的 withdrawal 交易金额总和。`is_enabled=false` 的预算跳过计算，返回零值。

**状态判定**：使用率 ≥100% → `overspent`，≥80% → `warning`，<80% → `normal`

**历史快照**：`CronService` 每天 09:00 执行 `SnapshotHistory()`，遍历所有启用预算，为近2年内所有已结束周期生成/更新快照（UPSERT），当前周期不生成（可看实时数据）。前端 `BudgetDetailPage` 折线图展示历史趋势（X轴=`period_start`，Y轴=限额+支出双线）。

### 定时任务管理

**调度器**：`CronService` 使用 `robfig/cron/v3`，在 `main.go` 中启动，无需外部触发。

**内置任务**：

| 任务ID | 名称 | 周期 | 说明 |
|---|---|---|---|
| `recurrences_bills` | 循环交易与到期账单 | 每天 00:00 | 执行循环交易生成交易记录，处理到期账单 |
| `budget_snapshot` | 预算历史快照 | 每天 09:00 | 遍历启用预算，为近2年已结束周期生成/更新快照 |

**管理 API**（需 Admin 权限）：

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/cron` | GET | 列出所有任务（ID、名称、描述、周期） |
| `/api/v1/cron/:id/run` | POST | 手动执行指定任务 |

**任务注册模式**：`CronService` 维护 `[]CronTask` 注册表，新增任务只需在 `newCronService()` 中追加 `CronTask` 并注册调度函数。前端 `CronPage` 展示任务列表并支持手动触发。

### 已知问题

1. **分类可重复关联多个预算**：没有校验，可能导致报表重复计算

## 账户管理

### 核心概念

基于**复式记账法**，账户分为四种类型：

| 类型 | 说明 | 示例 |
|---|---|---|
| `asset`（资产） | 你拥有的钱 | 银行存款、现金、投资 |
| `expense`（支出） | 花钱的分类 | 餐饮、交通、购物 |
| `revenue`（收入） | 赚钱的分类 | 工资、奖金 |
| `liability`（负债） | 你欠的钱 | 信用卡、贷款 |

### 数据模型

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:255) | 账户名称，同类型下唯一 |
| AccountNumber | string (size:100) | 账户号，同一用户内唯一（含软删除），创建后不可修改 |
| Type | string (size:20) | 账户类型：asset/expense/revenue/liability |
| CurrencyID | uint64 | 关联货币，外键 |
| InitialBalance | decimal(19,4) | 初始余额，创建后不可修改 |
| CurrentBalance | decimal(19,4) | 当前余额，系统自动计算 = 初始余额 + 收入 - 支出 |
| IsVirtual | bool | 虚拟账户标记，不代表真实资金 |
| Notes | text | 备注 |
| Currency | Currency | 关联货币对象（Belongs To） |
| DeletedAt | gorm.DeletedAt | 软删除 |

**唯一索引**：`idx_user_account_num_del` — `(user_id, account_number, deleted_at)` 三字段联合唯一，确保同一用户内账户号不重复（含软删除记录）

### 业务规则

1. **同类型下账户名唯一**：同一用户的同类型账户不能重名
2. **账户号用户内唯一**：同一用户的账户号不能重复（含软删除记录），创建后不可修改
3. **不可变字段**：账户号、账户类型、货币、初始余额创建后不可修改（修改会破坏交易余额计算）
4. **删除保护**：有关联交易的账户不能删除
5. **金额精度**：Go 用 `shopspring/decimal`，TS 用 `decimal.js`，禁止浮点数

### API 端点

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/accounts` | POST | 创建账户 |
| `/api/v1/accounts` | GET | 列表（支持按类型/搜索/排序/分页） |
| `/api/v1/accounts/:id` | GET | 详情 |
| `/api/v1/accounts/:id` | PUT | 更新（仅 name/notes/is_virtual） |
| `/api/v1/accounts/:id` | DELETE | 删除（软删除） |

### 前端页面

- **账户列表页**（`/accounts`）：四个 Tab 按类型筛选，展示名称、余额、货币、虚拟标记，支持编辑/删除
- **账户表单页**（`/accounts/create` 和 `/accounts/:id/edit`）：创建时可设置所有字段，编辑时类型/货币/初始余额置灰不可改
