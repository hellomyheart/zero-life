# REST API 接口

## 概述

Zero-Life 提供完整的 RESTful API，基于 JWT 实现 API 认证。API 版本为 1.0.0，基础路径为 `/api/v1/`。

## 认证

### JWT 认证

| 模式 | 说明 |
|------|------|
| Access Token | 登录获取，有效期 24 小时 |
| Refresh Token | 用于刷新 Access Token，有效期 7 天 |
| Personal Access Token | 用户在设置中生成的长期令牌 |

所有 API 请求需在 Header 中携带：`Authorization: Bearer <token>`

## API 端点总览

### 账户（Accounts）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /accounts | 列出账户 |
| POST | /accounts | 创建账户 |
| GET | /accounts/{id} | 获取账户详情 |
| PUT | /accounts/{id} | 更新账户 |
| DELETE | /accounts/{id} | 删除账户 |
| GET | /accounts/{id}/piggy-banks | 获取关联储蓄罐 |
| GET | /accounts/{id}/transactions | 获取关联交易 |
| GET | /accounts/{id}/attachments | 获取关联附件 |

### 附件（Attachments）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /attachments | 列出附件 |
| POST | /attachments | 创建附件 |
| GET | /attachments/{id} | 获取附件详情 |
| PUT | /attachments/{id} | 更新附件 |
| DELETE | /attachments/{id} | 删除附件 |
| GET | /attachments/{id}/download | 下载附件文件 |
| POST | /attachments/{id}/upload | 上传附件文件 |

### 账单/订阅（Bills / Subscriptions）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /bills | 列出账单 |
| POST | /bills | 创建账单 |
| GET | /bills/{id} | 获取账单详情 |
| PUT | /bills/{id} | 更新账单 |
| DELETE | /bills/{id} | 删除账单 |
| GET | /bills/{id}/attachments | 获取关联附件 |
| GET | /bills/{id}/rules | 获取关联规则 |
| GET | /bills/{id}/transactions | 获取关联交易 |

### 预算（Budgets）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /budgets | 列出预算 |
| POST | /budgets | 创建预算 |
| GET | /budgets/{id} | 获取预算详情 |
| PUT | /budgets/{id} | 更新预算 |
| DELETE | /budgets/{id} | 删除预算 |
| GET | /budgets/{id}/limits | 获取预算限额 |
| POST | /budgets/{id}/limits | 创建预算限额 |
| GET | /budgets/{id}/limits/{limitId} | 获取限额详情 |
| PUT | /budgets/{id}/limits/{limitId} | 更新限额 |
| DELETE | /budgets/{id}/limits/{limitId} | 删除限额 |
| GET | /budgets/{id}/transactions | 获取关联交易 |

### 可用预算（Available Budgets）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /available-budgets | 列出可用预算 |
| POST | /available-budgets | 创建可用预算 |
| GET | /available-budgets/{id} | 获取详情 |
| PUT | /available-budgets/{id} | 更新 |
| DELETE | /available-budgets/{id} | 删除 |

### 分类（Categories）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /categories | 列出分类 |
| POST | /categories | 创建分类 |
| GET | /categories/{id} | 获取分类详情 |
| PUT | /categories/{id} | 更新分类 |
| DELETE | /categories/{id} | 删除分类 |
| GET | /categories/{id}/transactions | 获取关联交易 |
| GET | /categories/{id}/attachments | 获取关联附件 |

### 币种（Currencies）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /currencies | 列出币种 |
| POST | /currencies | 创建币种 |
| GET | /currencies/{code} | 获取币种详情 |
| PUT | /currencies/{code} | 更新币种 |
| DELETE | /currencies/{code} | 删除币种 |
| POST | /currencies/{code}/enable | 启用币种 |
| POST | /currencies/{code}/disable | 禁用币种 |
| POST | /currencies/{code}/default | 设为默认 |
| GET | /currencies/{code}/accounts | 获取使用该币种的账户 |
| GET | /currencies/{code}/available-budgets | 获取可用预算 |
| GET | /currencies/{code}/bills | 获取账单 |
| GET | /currencies/{code}/budget-limits | 获取预算限额 |
| GET | /currencies/{code}/cer | 获取汇率 |
| GET | /currencies/{code}/recurrences | 获取定期交易 |
| GET | /currencies/{code}/rules | 获取规则 |
| GET | /currencies/{code}/transactions | 获取交易 |

### 汇率（Exchange Rates）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /exchange-rates | 列出汇率 |
| POST | /exchange-rates | 创建汇率 |
| GET | /exchange-rates/{id} | 获取汇率详情 |
| PUT | /exchange-rates/{id} | 更新汇率 |
| DELETE | /exchange-rates/{id} | 删除汇率 |

### 储蓄罐（Piggy Banks）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /piggy-banks | 列出储蓄罐 |
| POST | /piggy-banks | 创建储蓄罐 |
| GET | /piggy-banks/{id} | 获取储蓄罐详情 |
| PUT | /piggy-banks/{id} | 更新储蓄罐 |
| DELETE | /piggy-banks/{id} | 删除储蓄罐 |
| POST | /piggy-banks/{id}/add | 向储蓄罐存入金额 |
| POST | /piggy-banks/{id}/remove | 从储蓄罐取出金额 |
| GET | /piggy-banks/{id}/events | 获取储蓄事件 |
| GET | /piggy-banks/{id}/attachments | 获取关联附件 |
| GET | /piggy-banks/{id}/accounts | 获取关联账户 |

### 定期交易（Recurrences）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /recurrences | 列出定期交易 |
| POST | /recurrences | 创建定期交易 |
| GET | /recurrences/{id} | 获取定期交易详情 |
| PUT | /recurrences/{id} | 更新定期交易 |
| DELETE | /recurrences/{id} | 删除定期交易 |
| POST | /recurrences/{id}/trigger | 手动触发 |
| GET | /recurrences/{id}/transactions | 获取生成的交易 |

### 规则（Rules）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /rules | 列出规则 |
| POST | /rules | 创建规则 |
| GET | /rules/{id} | 获取规则详情 |
| PUT | /rules/{id} | 更新规则 |
| DELETE | /rules/{id} | 删除规则 |
| POST | /rules/{id}/test | 测试规则 |
| POST | /rules/{id}/trigger | 触发规则 |
| GET | /rules/{id}/validate-expression | 验证表达式 |

### 规则组（Rule Groups）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /rule-groups | 列出规则组 |
| POST | /rule-groups | 创建规则组 |
| GET | /rule-groups/{id} | 获取规则组详情 |
| PUT | /rule-groups/{id} | 更新规则组 |
| DELETE | /rule-groups/{id} | 删除规则组 |
| POST | /rule-groups/{id}/test | 测试规则组 |
| POST | /rule-groups/{id}/trigger | 触发规则组 |
| GET | /rule-groups/{id}/rules | 获取组内规则 |

### 标签（Tags）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /tags | 列出标签 |
| POST | /tags | 创建标签 |
| GET | /tags/{id} | 获取标签详情 |
| PUT | /tags/{id} | 更新标签 |
| DELETE | /tags/{id} | 删除标签 |
| GET | /tags/{id}/transactions | 获取关联交易 |
| GET | /tags/{id}/attachments | 获取关联附件 |

### 交易（Transactions）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /transactions | 列出交易 |
| POST | /transactions | 创建交易 |
| GET | /transactions/{id} | 获取交易详情 |
| PUT | /transactions/{id} | 更新交易 |
| DELETE | /transactions/{id} | 删除交易 |
| GET | /transactions/{id}/attachments | 获取关联附件 |
| GET | /transactions/{id}/piggy-bank-events | 获取储蓄事件 |

### 交易日志（Transaction Journals）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /transaction-journals/{id} | 获取交易日志详情 |
| DELETE | /transaction-journals/{id} | 删除交易日志 |
| GET | /transaction-journals/{id}/links | 获取关联链接 |

### 交易链接（Transaction Links / Link Types）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /transaction-links | 列出交易链接 |
| POST | /transaction-links | 创建交易链接 |
| GET | /transaction-links/{id} | 获取详情 |
| PUT | /transaction-links/{id} | 更新 |
| DELETE | /transaction-links/{id} | 删除 |
| GET | /link-types | 列出链接类型 |
| POST | /link-types | 创建链接类型 |
| GET | /link-types/{id} | 获取详情 |
| PUT | /link-types/{id} | 更新 |
| DELETE | /link-types/{id} | 删除 |

### Webhook

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /webhooks | 列出 Webhook |
| POST | /webhooks | 创建 Webhook |
| GET | /webhooks/{id} | 获取详情 |
| PUT | /webhooks/{id} | 更新 |
| DELETE | /webhooks/{id} | 删除 |
| POST | /webhooks/{id}/submit | 手动提交 |
| POST | /webhooks/{id}/trigger-transaction | 触发交易 |
| GET | /webhooks/{id}/messages | 获取消息列表 |
| GET | /webhooks/{id}/attempts | 获取尝试列表 |

### 偏好设置（Preferences）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /preferences | 列出偏好 |
| POST | /preferences | 创建或更新偏好（upsert 语义） |
| GET | /preferences/{name} | 获取指定偏好 |
| PUT | /preferences/{name} | 更新指定偏好 |
| DELETE | /preferences/{name} | 删除指定偏好 |

### 用户管理（Users，仅管理员）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /users | 列出用户 |
| POST | /users | 创建用户 |
| GET | /users/{id} | 获取用户详情 |
| PUT | /users/{id} | 更新用户 |
| DELETE | /users/{id} | 删除用户 |

### 对象分组（Object Groups）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /object-groups | 列出对象分组 |
| POST | /object-groups | 创建对象分组 |
| GET | /object-groups/{id} | 获取详情 |
| PUT | /object-groups/{id} | 更新 |
| DELETE | /object-groups/{id} | 删除 |

## 功能性端点

### 系统信息

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /about | 系统基本信息 |
| GET | /configuration | 系统配置（管理员） |
| PUT | /configuration | 更新配置（管理员） |

### 搜索

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /search/transactions | 搜索交易 |
| GET | /search/accounts | 搜索账户 |

### 自动补全

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /autocomplete/accounts | 账户补全 |
| GET | /autocomplete/categories | 分类补全 |
| GET | /autocomplete/tags | 标签补全 |
| GET | /autocomplete/budgets | 预算补全 |
| GET | /autocomplete/bills | 账单补全 |
| GET | /autocomplete/currencies | 币种补全 |
| GET | /autocomplete/piggy-banks | 储蓄罐补全 |
| GET | /autocomplete/rules | 规则补全 |
| GET | /autocomplete/rule-groups | 规则组补全 |
| GET | /autocomplete/recurrences | 定期交易补全 |
| GET | /autocomplete/webhooks | Webhook 补全 |
| GET | /autocomplete/transaction-types | 交易类型补全 |

### 图表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /chart/balance | 资产余额趋势 |
| GET | /chart/account/overview | 账户收支概览 |
| GET | /chart/budget/overview | 预算使用概览 |
| GET | /chart/category/overview | 分类收支概览 |

### 报表与洞察

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /summary/basic | 基本财务摘要 |
| GET | /insight/expense | 支出洞察 |
| GET | /insight/income | 收入洞察 |
| GET | /insight/transfer | 转账洞察 |

### 定时任务

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /cron/{token} | 触发定时任务 |

## 认证接口详细说明

### POST /auth/register

注册新用户。

**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}
```

**响应：** `201 Created`
```json
{
  "data": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

### POST /auth/login

登录获取 Token。

**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应：** `200 OK`
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

### POST /auth/refresh

刷新 Token。

**请求头：** `Authorization: Bearer <refresh_token>`

**响应：** `200 OK`
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

### POST /auth/logout

注销当前 Token（加入黑名单）。

### POST /auth/2fa/enable

启用双因素认证，返回 TOTP Secret 和 QR Code URL。

**响应：** `200 OK`
```json
{
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code_url": "otpauth://totp/ZeroLife:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=ZeroLife"
  }
}
```

### POST /auth/2fa/verify

验证双因素认证码，返回完整访问 Token。

**请求体：**
```json
{
  "code": "123456"
}
```

## 通用请求/响应规范

### 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| per_page | int | 50 | 每页数量（最大 500） |

### 日期范围参数

| 参数 | 类型 | 格式 | 说明 |
|------|------|------|------|
| start | string | YYYY-MM-DD | 开始日期 |
| end | string | YYYY-MM-DD | 结束日期 |

### 统一成功响应

```json
{
  "data": { ... }
}
```

### 统一分页响应

```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 100,
    "total_pages": 2
  },
  "links": {
    "self": "/api/v1/accounts?page=1",
    "first": "/api/v1/accounts?page=1",
    "last": "/api/v1/accounts?page=2",
    "next": "/api/v1/accounts?page=2"
  }
}
```

### 统一错误响应

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数校验失败",
    "details": [
      { "field": "name", "message": "名称不能为空" }
    ]
  }
}
```

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无返回体） |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 冲突（如重复创建） |
| 422 | 业务校验失败 |
| 429 | 请求限流 |
| 500 | 服务器内部错误 |
