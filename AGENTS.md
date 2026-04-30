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
- 数据库：SQLite + WAL 模式，启动时自动迁移（`migrations/` 目录为空，无手动迁移文件）。读写分离：写库 1 连接（串行保证安全），读库 100 连接（`?mode=ro` 只读，利用 WAL 并发读）
- 认证：JWT Bearer Token；`middleware.Auth` 将 `user_id`/`email` 注入 Gin 上下文；`middleware.Admin` 查数据库校验 `user.Role == "admin"` 并注入 `role` 到上下文，挂载在 `/users/*` 路由组
- 限流：使用 `kv_store` 表（非 Redis），详见下方「SQLite 替代 Redis 方案」
- 定时任务：内置 `robfig/cron` 调度器，每天 00:00 执行到期循环交易，每天 09:00 生成预算历史快照，每 6 小时执行 WAL 检查点

**前端** (`web/`) — Vue 3 + TypeScript + Vite
- 路径别名：`@` → `src/`
- API 层：`src/api/` — 每个领域一个文件，使用 `src/utils/request.ts`（Axios，自动刷新 Token）
- 状态管理：Pinia stores 在 `src/stores/`
- 国际化：中/英在 `src/i18n/locales/`
- 所有页面懒加载路由；认证守卫在 `src/router/index.ts`
- 无测试框架

## 关键约定

- **错误码**：按模块分段定义在 `internal/pkg/errcode/errcode.go`（1xxxx=认证, 2xxxx=账户, 3xxxx=交易, 4xxxx=分类, 5xxxx=标签, 6xxxx=预算, 7xxxx=循环交易, 8xxxx=规则, 9xxxx=导入, 10xxxx=定期交易, 11xxxx=Webhook, 14xxxx=偏好, 15xxxx=对账, 16xxxx=MFA, 17xxxx=储蓄罐）
- **API 响应格式**：统一为 `{ "code": 0, "message": "success", "data": ... }`，通过 `controller.Success()` / `controller.Error()` 返回
- **金额处理**：Go 用 `shopspring/decimal`，TS 用 `decimal.js` — 禁止用浮点数表示金额
- **无测试套件**：前后端均未配置测试
- **Swagger**：生成到 `server/docs/`；修改 API 后必须重新生成

## 注意事项

- SQLite 写库 `max_idle_conns` 和 `max_open_conns` 默认为 1 — 这是 SQLite 单写模式的刻意设计，不要修改。读库使用 `?mode=ro` 只读连接，`max_open_conns=100`，利用 WAL 并发读
- `server/cmd/server/config.yaml` 和 `server/config.yaml` 同时存在 — 运行时从 CWD 读取配置，所以必须在 `server/` 目录下启动
- 数据库文件（`*.db`、`*.db-shm`、`*.db-wal`）和 `*.exe` 已加入 gitignore
- Docker Compose 从 `.env` 读取 `JWT_SECRET`（默认值为不安全的 `change-me-in-production`）
- 前端 `nginx.conf` 硬编码 `proxy_pass http://api:8080` — 在 Docker Compose 外使用需修改
- `Recurrence` 和 `RecurringTransaction` 是两个独立的模型/控制器，路由不同（`/recurrences` vs `/recurring-transactions`），不要混淆 — **已合并**：删除了 `Recurrence` 和 `Bill`，统一使用 `RecurringTransaction`

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

## 预算管理

### 核心概念

预算用于设定分类支出上限，跟踪每个周期的支出使用情况。支持日度、周度、月度、季度、年度五种周期。一个预算可关联多个分类（多对多），选择父分类时自动包含子分类的支出。

### 数据模型

**Budget** (`model/budget.go`)

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:100) | 预算名称 |
| Amount | decimal(19,4) | 预算金额上限 |
| Period | string (size:20) | 预算周期：daily/weekly/monthly/quarterly/yearly |
| IsEnabled | bool (default:true) | 是否启用，禁用时不跟踪支出 |
| Categories | []Category | 多对多关联，通过 budget_categories 中间表 |
| DeletedAt | gorm.DeletedAt | 软删除 |

**BudgetCategory** (`model/budget.go`)：多对多关联表，复合主键 `(BudgetID, CategoryID)`，纯关联表无额外字段。

**BudgetHistory** (`model/budget.go`)

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| BudgetID | uint64 | 关联预算ID，索引 |
| PeriodStart | time.Time | 周期开始日期 |
| PeriodEnd | time.Time | 周期结束日期 |
| Amount | decimal(19,4) | 该周期的预算金额 |
| Spent | decimal(19,4) | 该周期的实际支出 |
| CreatedAt | time.Time | 创建时间 |
| UpdatedAt | time.Time | 更新时间 |

**唯一索引**：`idx_budget_history_budget_period` — `(budget_id, period_start)` 联合唯一索引，确保同一预算同一周期只有一条快照记录

### API 端点

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/budgets` | POST | 创建预算 |
| `/api/v1/budgets` | GET | 列表（含使用率计算） |
| `/api/v1/budgets/:id` | GET | 详情（含使用率计算） |
| `/api/v1/budgets/:id` | PUT | 更新（名称/金额/周期/启用状态/分类） |
| `/api/v1/budgets/:id` | DELETE | 删除（级联删除关联和历史） |
| `/api/v1/budgets/:id/history` | GET | 历史快照列表 |

### 请求/响应 DTO

**CreateBudgetReq**：`name`(必填)、`amount`(必填)、`period`(必填，oneof=daily/weekly/monthly/quarterly/yearly)、`category_ids`(必填，min=1)

**UpdateBudgetReq**：`name`(可选)、`amount`(可选)、`period`(可选，omitempty oneof)、`category_ids`(可选，nil=保留现有)、`is_enabled`(*bool，可选)

**BudgetResp**：`id`、`name`、`amount`(string, StringFixed(4))、`period`、`is_enabled`、`categories`([]CategoryResp)、`spent`(string)、`remaining`(string)、`usage_rate`(float64, 比率0~1+)、`status`(normal/warning/overspent)、`created_at`、`updated_at`

**BudgetHistoryResp**：`id`、`period_start`、`period_end`、`amount`(string)、`spent`(string)、`usage_rate`(float64, 比率0~1+)、`created_at`、`updated_at`

### 业务规则

1. **金额校验**：创建和更新时金额必须大于零（`ErrBudgetAmountInvalid`，60001）
2. **分类关联**：创建时至少关联1个分类；更新时若未提供 `category_ids` 则保留现有分类（全量替换语义）
3. **分类校验**：`validateCategoryIDs` 会去重后校验所有分类ID存在且属于当前用户（`ErrBudgetCategoryInvalid`，60002；`ErrBudgetCategoryRequired`，60003）
4. **启用/禁用**：`is_enabled=false` 的预算跳过支出计算，返回零值
5. **级联删除**：删除预算时同时删除 `budget_categories` 关联 + `budget_history` 历史记录，事务保证原子性
6. **空字符串语义**：`UpdateBudgetReq` 中 `name`/`amount`/`period` 为空字符串表示"不修改"；`is_enabled` 使用 `*bool` 指针区分"未提供"和"设为false"

### 支出计算 (`calculateSpent`)

1. 根据 `budgetPeriodRange(period, now)` 计算当前周期的起止时间
2. 收集预算关联的所有分类ID → 调用 `GetDescendantIDs` 展开子分类
3. 若展开后分类ID为空，返回零值（不查询交易）
4. 查询时间范围内、展开后分类下的 withdrawal 交易
5. 累加交易金额

`calculateSpent` 返回 `(decimal.Decimal, error)`，错误会向上传播到 `toRespWithUsage` 并返回 `errcode.ErrInternal`。`calculateSpentInRangeWithError` 同理，用于历史快照支出计算。

**`budgetPeriodRange` 函数**：根据周期类型计算起止时间，支持5种周期：
- daily：当天 00:00:00 ~ 23:59:59
- weekly：本周一 ~ 周日
- monthly：本月1日 ~ 月末
- quarterly：本季度首月1日 ~ 季末
- yearly：本年1月1日 ~ 12月31日

**使用场景**（5处支出计算，逻辑必须一致）：
- `BudgetService.calculateSpent`：预算列表/详情的使用率计算
- `BudgetService.calculateSpentInRangeWithError`：历史快照的支出计算
- `DashboardService.Get`：仪表盘预算预警
- `ChartService.BudgetSpending`：预算支出图表
- `ReportService.Budget`：预算报表

**错误处理**：所有5处支出计算均向上传播错误（不再静默吞错），`BudgetService.calculateSpent` 返回 `(decimal.Decimal, error)`，`DashboardService`/`ChartService`/`ReportService` 在 `GetDescendantIDs` 或 `ListAll` 失败时返回 `errcode.ErrInternal`。

### 状态判定

使用率 ≥100% → `overspent`，≥80% → `warning`，<80% → `normal`。除零保护：金额为零时使用率为0。

**所有5处支出计算的状态判定均使用 decimal 比较**（非 float64），避免浮点精度问题：`usageRateDecimal.GreaterThanOrEqual(decimal.NewFromInt(1))` 和 `decimal.NewFromFloat(0.8)`。

### 历史快照 (`SnapshotHistory`)

`CronService` 每天 09:00 执行，遍历所有启用预算，为近2年内所有已结束周期生成/更新快照（UPSERT），当前周期不生成（可看实时数据）。`UpsertHistory` 按 `budget_id + period_start` 判重，存在则只更新 `period_end`/`spent`，不更新 `amount`。

**Amount 保留策略**：已有快照的 `Amount` 不会被覆盖，保留该周期首次生成时的预算金额。只有新建快照时才写入当前预算金额。这样用户修改预算金额后，历史快照中的限额仍反映该周期实际设定的值。

**UpsertHistory 实现**：在事务中先查询 `WHERE budget_id = ? AND period_start = ?`，若 `errors.Is(err, gorm.ErrRecordNotFound)` 则创建新记录，否则更新已有记录。显式判断 `ErrRecordNotFound` 而非 `err != nil`，避免其他数据库错误被误判为"记录不存在"而创建重复记录。

### 交易查询

支出计算、报表统计等场景使用 `TransactionRepository.ListAll` 方法，内部循环分页（每批 5000 条）查询直到取完所有数据，无条数限制。

### 前端实现

- **类型** (`types/budget.ts`)：`Budget` 接口含 `spent`/`remaining`/`usage_rate`/`status` 字段，`BudgetPeriod` 枚举定义5种周期
- **列表页** (`pages/budgets/BudgetListPage.vue`)：表格展示预算列表，含使用率进度条和状态标签；创建/编辑对话框含分类树形选择器；金额输入添加前端校验（`isNaN(numAmount) || numAmount <= 0`）；表单初始金额为空字符串
- **详情页** (`pages/budgets/BudgetDetailPage.vue`)：三卡片展示金额/已用/剩余；使用率进度条；分类标签；历史趋势折线图（LineChart）；当前周期关联交易列表（分页）
- **报表页** (`pages/reports/BudgetReport.vue`)：柱状图对比预算金额vs已用金额；明细表格；无日期范围选择器（预算报表按各预算自身周期计算，不接受日期范围参数）
- **无独立 Store**：预算页面直接调用 API，数据存在组件本地 ref 中

### 定时任务管理

**调度器**：`CronService` 使用 `robfig/cron/v3`，在 `main.go` 中启动，无需外部触发。

**内置任务**：

| 任务ID | 名称 | 周期 | 说明 |
|---|---|---|---|
| `recurring_transactions` | 到期循环交易 | 每天 00:00 | 执行循环交易生成交易记录 |
| `budget_snapshot` | 预算历史快照 | 每天 09:00 | 遍历启用预算，为近2年已结束周期生成/更新快照 |
| `wal_checkpoint` | WAL检查点 | 每6小时 | 执行 `PRAGMA wal_checkpoint(TRUNCATE)`，将 WAL 文件合并回主数据库 |

**管理 API**（需 Admin 权限）：

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/cron` | GET | 列出所有任务（ID、名称、描述、周期） |
| `/api/v1/cron/:id/run` | POST | 手动执行指定任务 |

**任务注册模式**：`CronService` 维护 `[]CronTask` 注册表，新增任务只需在 `NewCronService()` 中追加 `CronTask` 并注册调度函数。前端 `CronPage` 展示任务列表并支持手动触发。


## 分类管理

### 核心概念

分类采用**最多5级树形结构**：通过 `parent_id` 形成层级关系，最多5层深度。分类用于标记交易的支出/收入类别，也用于预算关联和报表统计。

### 数据模型 (`model/category.go`)

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:100) | 分类名称，用户内全局唯一（非同父下唯一） |
| ParentID | *uint64 (indexed) | 父分类ID，NULL = 一级分类 |
| Icon | string (size:50) | 图标标识 |
| Notes | text | 备注 |
| SortOrder | int (default:0) | 排序权重，越小越靠前 |
| Children | []Category | 虚拟字段，GORM `foreignKey:ParentID`，不持久化 |
| DeletedAt | gorm.DeletedAt | 软删除 |

### API 端点

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/categories` | POST | 创建分类 |
| `/api/v1/categories` | GET | 列表（返回树形结构） |
| `/api/v1/categories/:id` | GET | 详情 |
| `/api/v1/categories/:id` | PUT | 更新（仅 name/icon/notes/sort_order） |
| `/api/v1/categories/:id` | DELETE | 删除（级联删除子分类） |

### 请求/响应 DTO

**CreateCategoryReq**：`name`(必填)、`parent_id`(可选，null=一级)、`icon`、`notes`

**UpdateCategoryReq**：`name`、`icon`、`notes`、`sort_order`(*int，指针区分0和未设置)。**不含 `parent_id`**，不支持移动分类到其他父级。

**CategoryResp**：`id`、`name`、`parent_id`、`icon`、`notes`、`sort_order`、`children`(递归，omitempty)、`created_at`、`updated_at`。不含 `user_id`。

### 业务规则

1. **最多5级深度限制**：仅创建时校验。沿 parent 链向上计算深度，若新分类将超过5级则拒绝（`ErrCategoryTooDeep`，40002）。更新时不校验，因为 `UpdateCategoryReq` 不含 `parent_id`，无法改变层级。

2. **名称唯一性**：创建和更新时均校验。用户内全局唯一（非同父下唯一），通过加载全部分类后线性扫描检查。更新时仅在新名称与原名称不同时才校验。

3. **级联删除**：在数据库事务中执行。使用 `GetDescendantIDs` 收集所有后代分类ID（含自身），然后：置空所有被删分类关联交易的 `category_id` → 清理 `budget_categories` 关联表 → 删除所有后代分类及自身。

4. **空字符串语义**：`UpdateCategoryReq` 中 `name`/`icon`/`notes` 为空字符串表示"不修改"，无法主动清空为空字符串。`sort_order` 使用 `*int` 指针，可区分"未提供"和"设为0"。

### 树构建算法 (`buildTree`)

O(n) 两遍扫描，`Children` 使用 `[]*CategoryResp` 指针切片避免值副本导致深层子节点丢失：
1. 遍历所有分类，创建 `map[uint64]*CategoryResp` 用于 O(1) 查找
2. 遍历所有分类，将有 `parent_id` 的节点**指针**挂到父节点的 `Children` 上
3. 收集 `parent_id == nil` 的根节点

**值语义陷阱**：如果 `Children` 使用 `[]CategoryResp`（值切片），`append` 时子节点是值副本。当三级节点被挂到二级节点时，一级节点中已持有的二级节点副本的 `Children` 不会更新，导致三级及更深层节点丢失。改用 `[]*CategoryResp` 指针切片后，所有引用指向同一份数据，深层子节点正确显示。

**孤儿节点处理**：如果 `parent_id` 指向不存在的分类（被删除或数据不一致），该节点既不会出现在根节点中，也不会作为子节点，会被静默丢弃。

**排序**：`List` 查询使用 `ORDER BY sort_order ASC, id ASC`，树构建保持此顺序。

### 后代展开 (`GetDescendantIDs`)

迭代 BFS 算法：输入一组分类ID → 包含自身 → 查询 `WHERE parent_id IN (当前层)` → 收集子ID → 重复直到无更多子节点 → 返回所有ID（含输入）。

**使用场景**（5处）：
- `TransactionService.List`：按分类筛选交易时，选中父分类自动包含子分类的交易
- `BudgetService.calculateSpent`：计算预算支出时，展开关联分类的子分类
- `DashboardService`：仪表盘预算告警支出计算
- `ChartService`：报表分类图表数据聚合

### 前端实现

- **类型** (`types/category.ts`)：`Category` 接口含递归 `children: Category[]`，与后端 `CategoryResp` 字段完全对齐
- **Store** (`stores/category.ts`)：仅缓存分类树，`fetchCategories()` 调用 `list()` API，失败时保持空数组
- **列表页** (`pages/categories/CategoryListPage.vue`)：使用 `el-tree` 展示树形结构，支持创建子分类（预填 `parent_id`）、编辑、删除
- **编辑限制**：编辑对话框不包含 `parent_id` 选择器和 `sort_order` 字段，与后端 `UpdateCategoryReq` 一致

## 标签管理

### 核心概念

标签采用**最多5级树形结构**：通过 `parent_id` 形成层级关系，最多5层深度。标签用于对交易进行灵活标记和分组，与分类互补。

**与分类的区别**：
- 分类是树形结构（最多5级），标签也是树形结构（最多5级）
- 一个交易只能有一个分类，但可以有多个标签
- 分类用于预算控制，标签用于灵活标记

### 数据模型 (`model/tag.go`)

| 字段 | 类型 | 说明 |
|---|---|---|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:100) | 标签名称，用户内全局唯一 |
| Color | string (size:7, default:#409EFF) | 标签颜色，十六进制格式 |
| ParentID | *uint64 (indexed) | 父标签ID，NULL = 顶级标签 |
| DeletedAt | gorm.DeletedAt | 软删除 |

**TransactionTag** (`model/tag.go`)：多对多关联表，复合主键 `(TransactionID, TagID)`，纯关联表无额外字段。

### API 端点

| 端点 | 方法 | 说明 |
|---|---|---|
| `/api/v1/tags` | POST | 创建标签 |
| `/api/v1/tags` | GET | 列表（返回树形结构） |
| `/api/v1/tags/:id` | GET | 详情 |
| `/api/v1/tags/:id` | PUT | 更新（名称、颜色、父标签ID） |
| `/api/v1/tags/:id` | DELETE | 删除（级联删除子标签） |

### 请求/响应 DTO

**CreateTagReq**：`name`(必填)、`color`(可选，默认 `#409EFF`)、`parent_id`(可选，null=顶级)

**UpdateTagReq**：`name`(可选)、`color`(可选)、`parent_id`(*uint64，可选)。**支持修改 `parent_id`**，可移动标签到其他父级。

**TagResp**：`id`、`name`、`color`、`parent_id`、`transaction_count`（关联交易数）、`children`(递归，omitempty)、`created_at`、`updated_at`。不含 `user_id`。

### 业务规则

1. **最多5级深度限制**：创建和更新时均校验。创建时沿 parent 链向上计算深度，若新标签将超过5级则拒绝（`ErrTagTooDeep`，50003）。更新时若指定了 `parent_id`，需满足三个条件：(a) 不能将自己设为自己的子标签；(b) 新父标签不能本身也是子标签（否则会形成三级）；(c) 若当前标签已有子标签，则不能将其变为子标签（否则会形成三级）。

2. **名称唯一性**：创建时校验。用户内全局唯一，通过加载全部标签后线性扫描检查。

3. **级联删除**：删除父标签时，先删除所有子标签（及子标签的 `transaction_tags` 关联），再删除父标签自身及其关联。不使用事务包裹整个级联删除（每个子标签单独调用 `tagRepo.Delete`，该方法内部有事务）。

4. **空字符串语义**：`UpdateTagReq` 中 `name`/`color` 为空字符串表示"不修改"，无法主动清空为空字符串。`parent_id` 使用 `*uint64` 指针，可区分"未提供"和"设为null（顶级）"。

### 树构建算法 (`buildTree`)

O(n) 两遍扫描，`Children` 使用 `[]*TagResp` 指针切片避免值副本导致深层子节点丢失：
1. 预计算每个标签的关联交易数（`CountTransactions`）
2. 遍历所有标签，创建 `map[uint64]*TagResp` 用于 O(1) 查找
3. 遍历所有标签，将有 `parent_id` 的节点**指针**挂到父节点的 `Children` 上
4. 收集 `parent_id == nil` 的根节点

**值语义陷阱**：如果 `Children` 使用 `[]TagResp`（值切片），`append` 时子节点是值副本。当三级节点被挂到二级节点时，一级节点中已持有的二级节点副本的 `Children` 不会更新，导致三级及更深层节点丢失。改用 `[]*TagResp` 指针切片后，所有引用指向同一份数据，深层子节点正确显示。

**孤儿节点处理**：如果 `parent_id` 指向不存在的标签（被删除或数据不一致），该节点既不会出现在根节点中，也不会作为子节点，会被静默丢弃。

### 后代展开 (`GetDescendantIDs`)

迭代 BFS 算法：输入一组标签ID → 包含自身 → 查询 `WHERE parent_id IN (当前层)` → 收集子ID → 重复直到无更多子节点 → 返回所有ID（含输入）。

### 前端实现

- **类型** (`types/tag.ts`)：`Tag` 接口含递归 `children?: Tag[]`，与后端 `TagResp` 字段完全对齐
- **Store** (`stores/tag.ts`)：缓存标签树，`fetchTags()` 调用 `list()` API，失败时保持空数组
- **列表页** (`pages/tags/TagListPage.vue`)：使用 `el-table` + `tree-props` 展示树形结构，支持创建子标签（预填 `parent_id`）、编辑、删除
- **颜色选择**：使用自定义 `ColorPicker` 组件，支持预设调色板 + 自定义颜色
- **父标签选择器**：编辑时可选择父标签（所有顶级标签可选），编辑自身时禁选自身

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
| CurrentBalance | decimal(19,4) | 当前余额，系统自动计算。资产=初始+收入-支出；负债=初始+欠款增加-还款；收入/支出=不追踪 |
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

## 交易管理

### 核心概念

基于**复式记账法**，每笔交易是**单条记录**（不是拆成两条），通过 `source_id` 和 `destination_id` 同时引用两个账户：

| 类型 | source_id → destination_id | 余额影响 |
|------|---------------------------|----------|
| withdrawal（支出） | 资产/负债账户 → 支出/负债账户 | 见下方余额规则 |
| deposit（收入） | 收入/负债账户 → 资产/负债账户 | 见下方余额规则 |
| transfer（转账） | 资产/负债账户 → 资产/负债账户 | 见下方余额规则 |

### 账户类型与交易使用规则

四种账户类型在交易中的角色：

| 账户类型 | 作为 source | 作为 destination | 余额含义 |
|----------|------------|-----------------|---------|
| asset（资产） | 取款/转账的付款方 | 存款/转账的收款方 | 正值=拥有的钱 |
| expense（支出） | 不适用 | 取款的支出分类 | 不追踪余额 |
| revenue（收入） | 存款的收入来源 | 不适用 | 不追踪余额 |
| liability（负债） | 取款（信用卡消费）/转账/存款（取现） | 存款（借款）/取款（还款）/转账 | 正值=欠款金额 |

**前端账户选择器规则**（`TransactionFormPage.vue`）：
- withdrawal：source = 资产+负债账户，destination = 支出+负债账户
- deposit：source = 收入+负债账户，destination = 资产+负债账户
- transfer：source = 资产+负债账户，destination = 资产+负债账户

**典型负债场景**：
- 信用卡消费：withdrawal（负债→支出），负债余额增加
- 还信用卡：withdrawal（资产→负债），资产减少、负债减少
- 借款/贷款：deposit（收入→负债），负债余额增加
- 信用卡取现：deposit（负债→资产），负债减少、资产增加

### 数据模型 (`model/transaction.go`)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint64 PK | 主键自增 |
| UserID | uint64 | 所属用户，与 Date 复合索引 |
| Type | string(20) | 交易类型：deposit/withdrawal/transfer |
| Date | time.Time | 交易日期，与 UserID 复合索引 |
| Description | string(500) | 描述 |
| Amount | decimal(19,4) | 金额，15位整数+4位小数 |
| SourceID | uint64 | 源账户ID，与 UserID 复合索引 |
| DestinationID | *uint64 | 目标账户ID（可空） |
| CategoryID | *uint64 | 分类ID（可空） |
| Notes | text | 备注 |
| BillID | *uint64 | 关联账单ID |
| ParentID | *uint64 | 父交易ID（拆分用，nil=顶级交易） |
| IsReconciled | bool | 是否已对账 |
| DeletedAt | gorm.DeletedAt | 软删除（仅模型定义，实际删除使用硬删除 Unscoped） |

**关联关系**：Source（Belongs To Account）、Destination（Belongs To *Account）、Category（Belongs To *Category）、Tags（many2many:transaction_tags）、Splits（Has Many，foreignKey:ParentID）

**`transaction_tags` 中间表**：复合主键 `(TransactionID, TagID)`

### API 端点

**基础交易**（`/api/v1/transactions`）：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/transactions` | POST | 创建交易 |
| `/transactions` | GET | 交易列表（分页+过滤） |
| `/transactions/search` | GET | 高级搜索 |
| `/transactions/:id` | GET/PUT/DELETE | 详情/更新/删除 |
| `/transactions/:id/split` | POST | 拆分交易 |
| `/transactions/:id/splits` | GET | 获取拆分列表 |
| `/transactions/:id/merge` | POST | 合并拆分 |

**批量操作**（`/api/v1/transactions/bulk`）：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/transactions/bulk/edit` | POST | 批量编辑（分类/备注/追加标签） |
| `/transactions/bulk/delete` | POST | 批量删除 |
| `/transactions/bulk/:id/convert` | POST | 类型转换 |
| `/transactions/bulk/:id/clone` | POST | 克隆交易 |

### 余额计算 (`calculateBalanceChanges`)

返回 `map[uint64]decimal.Decimal`，key 为账户ID，value 为变更量（正数=增加，负数=减少）。

**函数签名**：`calculateBalanceChanges(txnType, amount, sourceID, destID, sourceType, destType)` — 需要传入源/目标账户类型，因为负债账户的余额方向与资产账户相反。

**余额变更规则**（按账户类型）：

| 账户类型 | 作为 source 时 | 作为 destination 时 |
|----------|---------------|-------------------|
| asset（资产） | `-amount`（钱流出，余额减少） | `+amount`（钱流入，余额增加） |
| liability（负债） | `+amount`（欠款增加，如信用卡消费） | `-amount`（欠款减少，如还款） |
| revenue（收入） | 不更新（仅分类用途） | 不适用 |
| expense（支出） | 不适用 | 不更新（仅分类用途） |

**各交易类型的余额变更**：

| 交易类型 | source 变更 | destination 变更 |
|----------|-----------|-----------------|
| withdrawal | 资产：`-amount`，负债：`+amount` | 负债：`-amount`（还款），支出：不更新 |
| deposit | 负债：`-amount`（取现） | 资产：`+amount`，负债：`+amount`（借款） |
| transfer | 资产：`-amount`，负债：`+amount` | 资产：`+amount`，负债：`-amount` |

**注意**：`TransactionService` 和 `TransactionBulkService` 各自有一份 `calculateBalanceChanges`，逻辑必须保持一致。

### 账户验证 (`validateTransaction`)

`validateTransaction` 验证交易类型和账户归属，返回源账户和目标账户模型（含类型信息）。返回值 `(*model.Account, *model.Account, error)` 供 `calculateBalanceChanges` 使用账户类型计算余额变更。

### 余额更新

使用 SQL 表达式 `current_balance + change` 原子更新，避免先读后写竞态。所有余额更新操作都在数据库事务中执行。

### 事务一致性

所有写操作（Create/Update/Delete/Split/MergeSplits/ConvertType/Clone/BulkDelete）都在 `db.Transaction()` 中执行，确保交易记录和余额更新的原子性。Repository 层提供 `WithDB` 后缀方法（`CreateWithDB`/`DeleteWithDB`/`UpdateWithTagsAndDB`/`CreateBatchWithDB`/`AddTagsWithDB`/`AttachTagsWithDB`）接受外部事务对象。

### 删除策略

使用**硬删除**（`Unscoped().Delete()`），因为余额已在同一事务中硬回滚。如果使用软删除，恢复交易时余额不会自动恢复，导致数据不一致。`DeleteWithDB`、`DeleteBatch`、`BulkDelete` 中均使用 `Unscoped()`。

### 交易拆分

- 父交易：`ParentID = nil`，保留原始金额，余额变更只在父交易创建时发生一次
- 子交易：`ParentID = &parentID`，各子交易有独立金额/分类/标签，继承父交易的 Type/Date/SourceID/DestinationID
- 子交易金额之和必须等于父交易金额
- 列表查询通过 `parent_id IS NULL` 过滤只返回顶级交易
- **Update 同步**：更新父交易的 Type/SourceID/DestinationID 时，同步更新所有子交易的对应字段

### 标签关联

- **全量替换**（`UpdateWithTags`/`AttachTags`）：先删除旧标签，再插入新标签。用于单笔交易更新。
- **追加标签**（`AddTags`）：仅添加不存在的标签，不删除已有标签。用于批量编辑（`BulkEdit`）。

### 标签过滤

`applyFilter` 中 `TagID` 和 `TagIDs` 合并为互斥分支，同时指定时合并 ID 列表后单次 JOIN，避免产生两个 JOIN 导致结果不正确。

### Webhook 触发时机

Webhook 在事务成功提交后触发，避免事务回滚时 Webhook 已发出导致通知与实际状态不一致。

### 请求/响应 DTO

**CreateTransactionReq**：`type`(必填)、`date`(必填)、`description`(必填)、`amount`(必填)、`source_id`(必填)、`destination_id`(可选)、`category_id`、`notes`、`tags`([]uint64)、`splits`([]CreateSplitReq)

**UpdateTransactionReq**：与 Create 相同但**不含 `splits`** 字段。拆分请使用 Split/MergeSplits API。

**BulkEditReq**：`ids`(必填)、`category_id`、`notes`、`tag_ids`（追加标签，非替换）

**TransactionResp**：`id`、`type`、`date`、`description`、`amount`(string, StringFixed(4))、`source_id`、`source`(AccountResp)、`destination_id`、`destination`(*AccountResp)、`category_id`、`category`(*CategoryResp)、`notes`、`tags`([]TagResp)、`splits`([]SplitResp)、`bill_id`、`created_at`、`updated_at`

### 前端实现

- **类型** (`types/transaction.ts`)：`Transaction` 接口与后端 `TransactionResp` 字段对齐，`amount` 为 string 类型避免浮点精度问题
- **API** (`api/transaction.ts`)：`list`/`getTransaction`/`create`/`update`/`remove`/`search`

### 仪表盘净资产

`DashboardService.Get` 计算净资产（`total_balance`）= 所有资产账户余额之和 - 所有负债账户余额之和。同时返回 `total_assets`（总资产）和 `total_liabilities`（总负债）。

**DashboardResp**：`month_income`、`month_expense`、`net_income`、`total_balance`（净资产）、`total_assets`（总资产）、`total_liabilities`（总负债）、`budget_alerts`、`recurring_reminders`、`recent_txns`

前端仪表盘页面展示三张卡片：净资产、总资产、总负债。

## 账单管理

### 核心概念

账单（Bill）用于管理**周期性固定支出**，如房租、水电费、订阅服务等。与交易的关系是"计划 vs 实际"：账单是预期要付的钱，交易是实际付的钱。

- 账单设定重复规则（日/周/月/年）和下次到期日
- 系统每天 00:00 自动扫描到期账单，创建对应的支出交易并推进下次到期日
- 交易通过 `bill_id` 字段关联到账单，实现可追溯
- 仪表盘展示 7 天内即将到期的账单提醒

### 数据模型 (`model/bill.go`)

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:100) | 账单名称 |
| Amount | decimal(19,4) | 账单金额 |
| RepeatRule | string (size:20) | 重复规则：daily/weekly/monthly/yearly |
| NextDue | time.Time | 下次到期日期 |
| SourceID | *uint64 (indexed) | 支出账户ID（可选） |
| CategoryID | *uint64 (indexed) | 分类ID（可选） |
| Notes | text | 备注 |
| DeletedAt | gorm.DeletedAt | 软删除 |

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/bills` | POST | 创建账单 |
| `/api/v1/bills` | GET | 列表（按到期日升序） |
| `/api/v1/bills/:id` | GET | 详情 |
| `/api/v1/bills/:id` | PUT | 更新 |
| `/api/v1/bills/:id` | DELETE | 删除（同时清除关联交易的 bill_id） |

### 请求/响应 DTO

**CreateBillReq**：`name`(必填)、`amount`(必填，必须>0)、`repeat_rule`(必填，oneof=daily/weekly/monthly/yearly)、`next_due`(必填，YYYY-MM-DD)、`source_id`(*uint64)、`category_id`(*uint64)、`notes`

**UpdateBillReq**：`name`、`amount`、`repeat_rule`(omitempty oneof)、`next_due`、`source_id`(*uint64)、`category_id`(*uint64)、`notes`、`clear_source_id`(bool)、`clear_category_id`(bool)、`clear_notes`(bool)

**BillResp**：`id`、`name`、`amount`(string, StringFixed(4))、`repeat_rule`、`next_due`(time.Time)、`source_id`、`category_id`、`notes`、`created_at`、`updated_at`

### 业务规则

1. **金额校验**：创建和更新时金额必须大于零（`ErrBillAmountInvalid`，70001）
2. **到期处理**：`CreateTransactionFromBill` 在单个数据库事务中完成：创建支出交易（含 `bill_id` 关联）+ 推进到期日，确保原子性
3. **无 source_id 的账单**：到期时仅推进到期日，不创建交易（仅提醒模式）
4. **删除级联**：删除账单时在同一事务中清除所有关联交易的 `bill_id`（设为 NULL），再软删除账单
5. **空字符串语义**：`UpdateBillReq` 中 `name`/`amount`/`repeat_rule`/`next_due`/`notes` 为空字符串表示"不修改"
6. **清空可选字段**：`source_id`/`category_id` 使用 `*uint64` 指针，传值表示设置，但无法区分"未传"和"清空为null"。通过 `clear_source_id`/`clear_category_id`(bool) 字段显式清空。`notes` 通过 `clear_notes`(bool) 清空

### 到期处理流程 (`CreateTransactionFromBill`)

1. 获取账单 → 检查是否已到期（`NextDue > now` 则跳过）
2. 在 `db.Transaction` 中执行：
   - 若 `SourceID != nil`：调用 `txnService.CreateWithDB` 创建 withdrawal 交易（设置 `BillID` 关联）
   - 根据 `RepeatRule` 计算下次到期日（`calculateNextDue`）
   - 调用 `billRepo.UpdateWithDB` 更新到期日
3. 事务成功后：触发规则引擎和 Webhook 通知（`txnService.TriggerPostCreate`）

**下次到期日计算** (`calculateNextDue`)：
- daily：+1天
- weekly：+7天
- monthly：+1月
- yearly：+1年
- default：+1月

### 定时任务

`CronService` 每天 00:00 调用 `billRepo.GetAllDueBills()` 获取所有到期账单，逐条调用 `billService.CreateTransactionFromBill` 处理。不再有独立的 `createTransactionFromBill` 方法，消除代码重复。

### 仪表盘集成

`DashboardService` 调用 `billRepo.GetUpcoming(userID, 7)` 获取 7 天内到期的账单，返回 `BillReminderResp`（含 BillID、BillName、Amount、NextDue）。

### 前端实现

- **类型** (`types/bill.ts`)：`Bill` 接口与后端 `BillResp` 字段对齐，`RepeatRule` 枚举定义4种规则
- **API** (`api/bill.ts`)：`list`/`getBill`/`create`/`update`/`remove`
- **列表页** (`pages/bills/BillListPage.vue`)：表格展示账单列表，含逾期状态标签（前端根据 `next_due` 与当前日期比较计算）、创建/编辑对话框含账户和分类选择器

## 储蓄罐管理

### 核心概念

储蓄罐用于设定储蓄目标并跟踪进度。每个储蓄罐关联一个资产账户，记录目标金额、当前已存金额和可选的目标日期。

**储蓄罐是虚拟记账，存取不扣减/增加关联账户余额，不创建交易记录。** 但存入时校验关联账户余额是否足够（防止存入金额超过账户实际持有资金），取出时不校验账户余额（取出是释放虚拟锁定，永远合理）。

### 数据模型 (`model/piggy_bank.go`)

**PiggyBank**：

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint64 | 主键自增 |
| UserID | uint64 | 所属用户，索引，数据隔离 |
| Name | string (size:100) | 储蓄罐名称 |
| TargetAmount | decimal(19,4) | 目标金额，必须为正数 |
| CurrentAmount | decimal(19,4) | 当前已存金额，默认0 |
| AccountID | uint64 | 关联资产账户ID（存入时校验账户余额） |
| Order | int (default:0) | 排序权重 |
| TargetDate | *time.Time (indexed) | 目标完成日期（可选） |
| Notes | text | 备注 |
| Account | Account | 关联账户（Belongs To） |
| Events | []PiggyEvent | 存取事件记录（Has Many） |
| DeletedAt | gorm.DeletedAt | 软删除 |

**PiggyEvent**：

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint64 | 主键自增 |
| PiggyBankID | uint64 | 关联储蓄罐ID，索引 |
| Amount | decimal(19,4) | 正数=存入，负数=取出 |
| TransactionID | *uint64 | 关联交易ID（可选，当前未使用） |
| Note | string (text) | 操作备注 |
| CreatedAt | time.Time | 创建时间 |

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/piggy-banks` | POST | 创建储蓄罐 |
| `/api/v1/piggy-banks` | GET | 列表 |
| `/api/v1/piggy-banks/:id` | GET | 详情 |
| `/api/v1/piggy-banks/:id` | PUT | 更新（名称/目标金额/目标日期/备注） |
| `/api/v1/piggy-banks/:id` | DELETE | 删除（软删除） |
| `/api/v1/piggy-banks/:id/add` | POST | 存入金额 |
| `/api/v1/piggy-banks/:id/remove` | POST | 取出金额 |
| `/api/v1/piggy-banks/:id/events` | GET | 事件记录列表 |
| `/api/v1/piggy-banks/reorder` | PUT | 批量排序 |
| `/api/v1/piggy-banks/:id/reset` | POST | 重置历史（清空事件+金额归零） |

### 业务规则

1. **创建**：目标金额必须为正数；关联账户必填且必须是当前用户的 asset 类型账户（`ErrPiggyBankAccountInvalid`，170005）
2. **更新目标金额**：新目标金额不能小于当前已存金额（`ErrPiggyBankTargetTooSmall`，170004）
3. **存入**：双重校验：
   - 目标金额限制：存入后当前金额不能超过目标金额（`ErrPiggyBankOverDeposit`，170002）
   - 账户余额限制：关联账户余额 - 同账户所有储蓄罐已存总额 >= 本次存入金额（`ErrPiggyBankAccountInsufficient`，170007）
4. **取出**：金额必须为正数；取出金额不能超过当前已存金额（`ErrPiggyBankOverWithdraw`，170003）。取出时不校验账户余额（取出是释放虚拟锁定）
5. **完成百分比**：`CurrentAmount / TargetAmount * 100`，上限100%
6. **重置**：清空所有事件记录，当前金额归零，事务保证原子性

### 账户余额校验（存入时）

存入时校验关联账户余额，防止存入金额超过账户实际持有资金。多个储蓄罐可绑定同一账户，校验时需考虑同账户下所有储蓄罐的已存总额。

**校验逻辑**（`AddAmount`）：
1. 查询关联账户的 `current_balance`
2. 查询该账户下所有储蓄罐的 `current_amount` 之和（`SumCurrentAmountByAccount`）
3. 计算：`账户余额 - 同账户储蓄罐已存总额 >= 本次存入金额`
4. 不满足则返回 `ErrPiggyBankAccountInsufficient`

**可存入金额**（`available_deposit`，响应中实时计算）：
```
account_available = max(account.current_balance - SUM(同账户储蓄罐.current_amount), 0)
target_remaining = target_amount - current_amount
available_deposit = min(account_available, target_remaining)
```

**超额场景处理**：如果账户余额因支出减少，导致储蓄罐的 `current_amount` 超出了合理范围，系统不自动修正。下次存入时校验会阻止继续存入，用户只能先取出多存的金额，才能再存入。前端通过 `available_deposit` 字段实时展示可存入上限。

**举例**：
- 账户余额 10000，储蓄罐A（目标5000，已存3000），储蓄罐B（目标8000，已存4000）
- 已锁定总额 = 3000 + 4000 = 7000
- 账户剩余可分配 = 10000 - 7000 = 3000
- 储蓄罐A 再存入：min(3000, 5000-3000) = min(3000, 2000) = **2000**
- 如果账户支出 5000，余额变 5000：账户剩余可分配 = 5000 - 7000 = **-2000**（已超额）
- 储蓄罐A 存入被拒绝（`available_deposit = 0`），只能先取出

### 事务一致性

- **存入/取出**：在 `db.Transaction()` 中执行，金额更新（SQL 原子表达式）和事件记录创建在同一事务中
- **原子更新**：`AddAmountWithDB` 使用 `UPDATE ... SET current_amount = current_amount + ? WHERE current_amount + ? <= target_amount`，通过 `RowsAffected` 判断是否超额，避免先读后写竞态
- **取出原子更新**：`RemoveAmountWithDB` 使用 `UPDATE ... SET current_amount = current_amount - ? WHERE current_amount >= ?`
- **重置**：`ResetAmountWithDB` + `DeleteEventsWithDB` 在同一事务中执行
- **Repository WithDB 方法**：`UpdateWithDB`/`CreateEventWithDB`/`DeleteEventsWithDB`/`AddAmountWithDB`/`RemoveAmountWithDB`/`ResetAmountWithDB` 接受外部事务对象

### 错误码（17xxxx）

| 错误码 | Code | 说明 |
|--------|------|------|
| `ErrPiggyBankAmountInvalid` | 170001 | 目标金额必须为正数 |
| `ErrPiggyBankOverDeposit` | 170002 | 存入后超过目标金额 |
| `ErrPiggyBankOverWithdraw` | 170003 | 取出超过当前金额 |
| `ErrPiggyBankTargetTooSmall` | 170004 | 目标金额不能小于当前已存金额 |
| `ErrPiggyBankAccountInvalid` | 170005 | 关联账户无效（不存在/非 asset 类型/非本人） |
| `ErrPiggyBankNotFound` | 170006 | 储蓄罐不存在 |
| `ErrPiggyBankAccountInsufficient` | 170007 | 关联账户余额不足以存入该金额 |

### 请求/响应 DTO

**CreatePiggyBankReq**：`name`(必填)、`target_amount`(必填，正数)、`account_id`(必填)、`target_date`(可选)、`notes`(可选)

**UpdatePiggyBankReq**：`name`(可选)、`target_amount`(*string，可选，正数)、`target_date`(*string，可选)、`notes`(可选)、`clear_notes`(bool，可选，清空备注)

**AddAmountReq**：`amount`(必填，正数字符串)、`note`(可选)

**RemoveAmountReq**：`amount`(必填，正数字符串)、`note`(可选)

**PiggyBankResp**：`id`、`name`、`target_amount`(string)、`current_amount`(string)、`account_id`、`account`(*AccountResp)、`target_date`、`notes`、`percentage`(float64)、`available_deposit`(string，可存入金额)、`created_at`、`updated_at`

**PiggyEventResp**：`id`、`piggy_bank_id`、`amount`(string)、`transaction_id`、`note`、`created_at`

### 前端实现

- **类型** (`types/piggyBank.ts`)：`PiggyBank` 接口含 `available_deposit` 和 `account?` 字段，与后端 `PiggyBankResp` 完全对齐
- **API** (`api/piggyBank.ts`)：`list`/`getPiggyBank`/`create`/`update`/`remove`/`addAmount`/`removeAmount`/`getEvents`/`reorder`/`resetHistory`；`addAmount`/`removeAmount` 返回 `PiggyBank`（更新后的数据）
- **无独立 Store**：页面直接调用 API
- **列表页** (`pages/piggyBanks/PiggyBankListPage.vue`)：表格展示储蓄罐列表，含进度条和状态颜色；账户选择器过滤仅 asset 类型；存入对话框显示 `available_deposit` 作为可存入上限；行点击跳转详情页；金额输入添加前端校验
- **详情页** (`pages/piggyBanks/PiggyBankDetailPage.vue`)：三卡片（目标/当前/剩余）；进度条；基本信息（关联账户、目标日期、备注）；存取记录列表（区分存入绿色+号/取出红色-号）；存入/取出/编辑/删除/重置历史操作
