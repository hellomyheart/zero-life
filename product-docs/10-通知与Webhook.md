# 通知与 Webhook

## 通知系统

Zero-Life 支持多种通知渠道，用于向用户或管理员推送重要事件通知。

### 通知渠道

| 渠道 | 说明 |
|------|------|
| Mail（邮件） | 通过 SMTP 发送邮件通知 |
| Slack | 通过 Slack Webhook 发送通知 |
| Pushover | 通过 Pushover 服务推送通知 |

### 通知类型

**用户通知：**
- 密码重置请求
- 新版本检查结果
- 账单到期提醒

**管理员通知：**
- 新用户注册
- 版本检查结果

### 邮件通知路由

邮件通知的发送目标按以下优先级确定：

1. 如果用户有自定义通知路由方法，使用自定义路由
2. 检查用户偏好 `remote_guard_alt_email`（备用邮箱）
3. 如果用户是 demo 角色，发送给站点所有者
4. 默认使用用户注册邮箱

### Slack 通知路由

Slack 通知的发送目标按以下逻辑确定：

1. 获取全局配置 `slack_webhook_url`
2. 如果通知类型为 `owner`，使用全局 Webhook URL
3. 如果是 `UserRegistration` 或 `VersionCheckResult` 通知，使用全局 Webhook URL
4. 其他通知使用用户个人偏好 `slack_webhook_url`

### Pushover 通知路由

Pushover 通知需要配置：
- `pushover_app_token` - 应用 Token（全局配置）
- `pushover_user_token` - 用户 Token（用户偏好配置）

## Webhook 系统

Webhook 允许在特定事件发生时向外部 URL 发送 HTTP 请求，实现与其他系统的集成。

### Webhook 属性

| 属性 | 说明 |
|------|------|
| title | Webhook 名称 |
| url | 回调 URL |
| secret | 签名密钥 |
| trigger | 触发条件 |
| response | 响应内容 |
| delivery | 交付方式 |
| active | 是否启用 |

### Webhook 触发器（Trigger）

| 触发器 | 说明 |
|--------|------|
| ANY | 任意事件（匹配所有触发条件） |
| STORE_TRANSACTION | 创建交易时 |
| UPDATE_TRANSACTION | 更新交易时 |
| DESTROY_TRANSACTION | 删除交易时 |
| STORE_BUDGET | 创建预算时 |
| UPDATE_BUDGET | 更新预算时 |
| DESTROY_BUDGET | 删除预算时 |
| STORE_UPDATE_BUDGET_LIMIT | 创建或更新预算限额时 |

### Webhook 响应内容（Response）

| 响应类型 | 说明 |
|----------|------|
| TRANSACTIONS | 返回相关交易数据 |
| ACCOUNTS | 返回相关账户数据 |
| BUDGET | 返回相关预算数据 |
| RELEVANT | 返回相关数据（由系统判断） |
| NONE | 不返回数据 |

### Webhook 交付方式（Delivery）

| 方式 | 说明 |
|------|------|
| JSON | 以 JSON 格式发送 HTTP POST 请求 |

### Webhook 消息（WebhookMessage）

每次 Webhook 触发会生成一条消息：

| 属性 | 说明 |
|------|------|
| sent | 是否已发送 |
| errored | 是否出错 |
| uuid | 消息唯一标识 |
| message | 消息内容（JSON） |
| logs | 发送日志（JSON） |

### Webhook 尝试（WebhookAttempt）

每次发送尝试记录：

| 属性 | 说明 |
|------|------|
| webhook_message_id | 所属消息 |
| status | HTTP 状态码 |
| logs | 尝试日志 |
| response | 响应内容 |

### Webhook 发送流程

```
1. 事件发生（如创建交易）
2. 检查匹配的 Webhook（trigger 匹配）
3. 生成 WebhookMessage
4. 尝试发送 HTTP POST 请求到 Webhook URL
5. 记录 WebhookAttempt（状态码、响应）
6. 如果成功，标记 sent = true
7. 如果失败，标记 errored = true，记录错误日志
```

### Webhook 签名

Webhook 请求包含签名头，用于验证请求来源：

- 使用 Webhook 的 secret 密钥
- 对消息体进行 HMAC 签名
- 签名放在 HTTP 请求头中

### Webhook 手动操作

| 操作 | 说明 |
|------|------|
| 手动提交 | 重新发送未成功的消息 |
| 触发交易 | 手动触发指定交易的 Webhook |

## 定时任务（Cron）

系统通过定时任务执行以下周期性操作：

| 任务 | 说明 |
|------|------|
| 定期交易生成 | 检查并创建到期的定期交易 |
| 账单提醒 | 检查到期账单并发送提醒 |
| 汇率更新 | 从外部源下载最新汇率 |
| 版本检查 | 检查是否有新版本可用 |
| 数据清理 | 清理过期的 session、日志等 |

### Cron 端点

通过 API 触发定时任务：`GET /cron/{cliToken}`

- `cliToken` 为系统配置的定时任务令牌
- 用于外部调度器（如系统 crontab）定期调用

## 报表与洞察

### 基本摘要

`GET /summary/basic` 返回当前周期的财务摘要：
- 总资产余额
- 总支出
- 总收入
- 净值变化

### 洞察端点

| 端点 | 说明 |
|------|------|
| /insight/expense | 支出洞察（按分类、账户、标签分组） |
| /insight/income | 收入洞察 |
| /insight/transfer | 转账洞察 |

洞察端点支持日期范围和多种分组维度，用于生成详细的财务分析报表。

### 图表端点

| 端点 | 说明 |
|------|------|
| /chart/balance | 资产余额趋势 |
| /chart/account/overview | 账户收支概览 |
| /chart/budget/overview | 预算使用概览 |
| /chart/category/overview | 分类收支概览 |
