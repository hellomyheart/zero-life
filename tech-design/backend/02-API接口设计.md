# API 接口设计

## 总体规范

### 基础信息

| 项目 | 说明 |
|------|------|
| 基础路径 | `/api/v1` |
| 协议 | HTTPS |
| 认证 | Bearer Token (JWT) |
| 内容类型 | `application/json` |
| 字符编码 | UTF-8 |

### 统一响应格式

**成功响应：**

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 100,
    "total_pages": 2
  }
}
```

**列表响应：**

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

**错误响应：**

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

### 分页参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| page | 1 | 页码 |
| per_page | 50 | 每页数量（最大 500） |

### 通用过滤参数

| 参数 | 说明 |
|------|------|
| start | 开始日期 (YYYY-MM-DD) |
| end | 结束日期 (YYYY-MM-DD) |
| type | 类型过滤 |

## 认证接口

### POST /api/v1/auth/register

注册新用户。

**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}
```

### POST /api/v1/auth/login

登录获取 Token。

**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应：**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "refresh_token": "..."
  }
}
```

### POST /api/v1/auth/refresh

刷新 Token。

### POST /api/v1/auth/logout

注销 Token（加入黑名单）。

### POST /api/v1/auth/2fa/enable

启用双因素认证。

### POST /api/v1/auth/2fa/verify

验证双因素认证码。

## 账户接口

### GET /api/v1/accounts

列出账户。

**查询参数：**
| 参数 | 说明 |
|------|------|
| type | 账户类型过滤 (asset/expense/revenue/debt/loan/mortgage) |
| active | 活跃状态 (true/false) |
| search | 名称搜索 |

### POST /api/v1/accounts

创建账户。

**请求体：**
```json
{
  "name": "招商银行储蓄卡",
  "type": "asset",
  "account_role": "defaultAsset",
  "currency_id": 1,
  "virtual_balance": "0",
  "opening_balance": "10000.00",
  "opening_balance_date": "2024-01-01",
  "is_active": true,
  "notes": "主储蓄卡",
  "object_group_id": 1
}
```

### GET /api/v1/accounts/:id

获取账户详情（含当前余额）。

**查询参数：**
| 参数 | 说明 |
|------|------|
| date | 计算余额的截止日期 |

### PUT /api/v1/accounts/:id

更新账户。

### DELETE /api/v1/accounts/:id

删除账户。

### GET /api/v1/accounts/:id/transactions

获取账户关联交易。

### GET /api/v1/accounts/:id/piggy-banks

获取账户关联储蓄罐。

### GET /api/v1/accounts/:id/attachments

获取账户关联附件。

## 交易接口

### GET /api/v1/transactions

列出交易组。

**查询参数：**
| 参数 | 说明 |
|------|------|
| type | 交易类型 (withdrawal/deposit/transfer) |
| start | 开始日期 |
| end | 结束日期 |
| category_id | 分类过滤 |
| budget_id | 预算过滤 |
| bill_id | 账单过滤 |
| tag_id | 标签过滤 |
| account_id | 账户过滤 |

### POST /api/v1/transactions

创建交易（支持拆分交易）。

**请求体（单笔支出）：**
```json
{
  "type": "withdrawal",
  "date": "2024-03-15",
  "description": "超市购物",
  "transactions": [
    {
      "source_id": 1,
      "destination_id": 5,
      "amount": "156.80",
      "currency_id": 1,
      "category_id": 3,
      "budget_id": 2,
      "tags": ["日常", "食品"],
      "notes": "周末采购",
      "foreign_amount": null,
      "foreign_currency_id": null
    }
  ]
}
```

**请求体（拆分交易）：**
```json
{
  "type": "withdrawal",
  "date": "2024-03-15",
  "description": "超市购物",
  "transactions": [
    {
      "source_id": 1,
      "destination_id": 5,
      "amount": "100.00",
      "currency_id": 1,
      "category_id": 3,
      "description": "食品"
    },
    {
      "source_id": 1,
      "destination_id": 5,
      "amount": "56.80",
      "currency_id": 1,
      "category_id": 4,
      "description": "日用品"
    }
  ]
}
```

### GET /api/v1/transactions/:id

获取交易组详情。

### PUT /api/v1/transactions/:id

更新交易组。

### DELETE /api/v1/transactions/:id

删除交易组。

### GET /api/v1/transactions/:id/attachments

获取交易关联附件。

## 预算接口

### GET /api/v1/budgets

列出预算。

**查询参数：**
| 参数 | 说明 |
|------|------|
| start | 限额开始日期 |
| end | 限额结束日期 |

### POST /api/v1/budgets

创建预算。

**请求体：**
```json
{
  "name": "餐饮",
  "auto_budget_type": 1,
  "auto_budget_amount": "2000.00",
  "auto_budget_period": "monthly",
  "is_active": true
}
```

### GET /api/v1/budgets/:id

获取预算详情（含当前周期限额和已用金额）。

### PUT /api/v1/budgets/:id

更新预算。

### DELETE /api/v1/budgets/:id

删除预算。

### GET /api/v1/budgets/:id/limits

获取预算限额列表。

### POST /api/v1/budgets/:id/limits

创建/设置预算限额。

**请求体：**
```json
{
  "currency_id": 1,
  "start_date": "2024-03-01",
  "end_date": "2024-03-31",
  "amount": "2000.00"
}
```

### GET /api/v1/budgets/:id/transactions

获取预算关联交易。

## 账单接口

### GET /api/v1/bills

列出账单。

### POST /api/v1/bills

创建账单。

**请求体：**
```json
{
  "name": "房租",
  "match": "房租|rent",
  "amount_min": "3000.00",
  "amount_max": "3500.00",
  "date": "2024-03-01",
  "repeat_freq": "monthly",
  "currency_id": 1,
  "is_active": true
}
```

### GET /api/v1/bills/:id

获取账单详情。

### PUT /api/v1/bills/:id

更新账单。

### DELETE /api/v1/bills/:id

删除账单。

### GET /api/v1/bills/:id/transactions

获取账单关联交易。

## 储蓄罐接口

### GET /api/v1/piggy-banks

列出储蓄罐。

### POST /api/v1/piggy-banks

创建储蓄罐。

**请求体：**
```json
{
  "name": "旅行基金",
  "account_id": 1,
  "target_amount": "20000.00",
  "start_date": "2024-01-01",
  "target_date": "2024-12-31",
  "currency_id": 1
}
```

### GET /api/v1/piggy-banks/:id

获取储蓄罐详情（含进度百分比）。

### PUT /api/v1/piggy-banks/:id

更新储蓄罐。

### DELETE /api/v1/piggy-banks/:id

删除储蓄罐。

### POST /api/v1/piggy-banks/:id/add

向储蓄罐存入金额。

**请求体：**
```json
{
  "amount": "500.00",
  "transaction_journal_id": null
}
```

### POST /api/v1/piggy-banks/:id/remove

从储蓄罐取出金额。

### GET /api/v1/piggy-banks/:id/events

获取储蓄事件列表。

## 分类接口

### GET /api/v1/categories

列出分类。

### POST /api/v1/categories

创建分类。

### GET /api/v1/categories/:id

获取分类详情。

### PUT /api/v1/categories/:id

更新分类。

### DELETE /api/v1/categories/:id

删除分类。

### GET /api/v1/categories/:id/transactions

获取分类关联交易。

## 标签接口

### GET /api/v1/tags

列出标签。

### POST /api/v1/tags

创建标签。

### GET /api/v1/tags/:id

获取标签详情。

### PUT /api/v1/tags/:id

更新标签。

### DELETE /api/v1/tags/:id

删除标签。

### GET /api/v1/tags/:id/transactions

获取标签关联交易。

## 币种接口

### GET /api/v1/currencies

列出币种。

### POST /api/v1/currencies

创建币种。

### GET /api/v1/currencies/:code

获取币种详情。

### PUT /api/v1/currencies/:code

更新币种。

### POST /api/v1/currencies/:code/enable

启用币种。

### POST /api/v1/currencies/:code/disable

禁用币种。

### POST /api/v1/currencies/:code/default

设为默认币种。

## 汇率接口

### GET /api/v1/exchange-rates

列出汇率记录。

**查询参数：**
| 参数 | 说明 |
|------|------|
| from | 源币种代码 |
| to | 目标币种代码 |
| date | 汇率日期 |

### POST /api/v1/exchange-rates

创建/设置汇率。

**请求体：**
```json
{
  "from_currency_id": 1,
  "to_currency_id": 2,
  "date": "2024-03-15",
  "rate": "7.8500"
}
```

### PUT /api/v1/exchange-rates/:id

更新汇率。

### DELETE /api/v1/exchange-rates/:id

删除汇率。

## 规则接口

### GET /api/v1/rules

列出规则。

### POST /api/v1/rules

创建规则。

**请求体：**
```json
{
  "rule_group_id": 1,
  "title": "超市交易自动分类",
  "strict": true,
  "triggers": [
    { "type": "description_contains", "value": "超市", "order": 1 },
    { "type": "amount_less", "value": "500", "order": 2 }
  ],
  "actions": [
    { "type": "set_category", "value": "食品", "order": 1 },
    { "type": "add_tag", "value": "日常", "order": 2 }
  ]
}
```

### GET /api/v1/rules/:id

获取规则详情。

### PUT /api/v1/rules/:id

更新规则。

### DELETE /api/v1/rules/:id

删除规则。

### POST /api/v1/rules/:id/test

测试规则（预览匹配结果，不修改数据）。

**请求体：**
```json
{
  "transaction_ids": [1, 2, 3]
}
```

### POST /api/v1/rules/:id/trigger

触发规则（实际执行）。

## 规则组接口

### GET /api/v1/rule-groups

列出规则组。

### POST /api/v1/rule-groups

创建规则组。

### GET /api/v1/rule-groups/:id

获取规则组详情。

### PUT /api/v1/rule-groups/:id

更新规则组。

### DELETE /api/v1/rule-groups/:id

删除规则组。

### POST /api/v1/rule-groups/:id/test

测试规则组。

### POST /api/v1/rule-groups/:id/trigger

触发规则组。

### GET /api/v1/rule-groups/:id/rules

获取组内规则列表。

## 定期交易接口

### GET /api/v1/recurrences

列出定期交易。

### POST /api/v1/recurrences

创建定期交易。

**请求体：**
```json
{
  "title": "每月房租",
  "type": "withdrawal",
  "first_date": "2024-01-01",
  "repeat_until": null,
  "apply_rules": true,
  "repetitions": [
    { "type": "monthly", "moment": "1", "skip": 0, "weekend": 1 }
  ],
  "transactions": [
    {
      "source_id": 1,
      "destination_id": 5,
      "amount": "3200.00",
      "currency_id": 1,
      "description": "月租"
    }
  ]
}
```

### GET /api/v1/recurrences/:id

获取定期交易详情。

### PUT /api/v1/recurrences/:id

更新定期交易。

### DELETE /api/v1/recurrences/:id

删除定期交易。

### POST /api/v1/recurrences/:id/trigger

手动触发定期交易。

## Webhook 接口

### GET /api/v1/webhooks

列出 Webhook。

### POST /api/v1/webhooks

创建 Webhook。

**请求体：**
```json
{
  "title": "交易通知",
  "url": "https://example.com/webhook",
  "secret": "my-secret-key",
  "trigger": "STORE_TRANSACTION",
  "response": "TRANSACTIONS",
  "delivery": "JSON",
  "active": true
}
```

### GET /api/v1/webhooks/:id

获取 Webhook 详情。

### PUT /api/v1/webhooks/:id

更新 Webhook。

### DELETE /api/v1/webhooks/:id

删除 Webhook。

### POST /api/v1/webhooks/:id/submit

手动提交未发送的消息。

### GET /api/v1/webhooks/:id/messages

获取 Webhook 消息列表。

### GET /api/v1/webhooks/:id/attempts

获取发送尝试列表。

## 附件接口

### GET /api/v1/attachments

列出附件。

### POST /api/v1/attachments

创建附件记录。

### POST /api/v1/attachments/:id/upload

上传附件文件（multipart/form-data）。

### GET /api/v1/attachments/:id/download

下载附件文件。

### DELETE /api/v1/attachments/:id

删除附件。

## 搜索接口

### GET /api/v1/search/transactions

搜索交易。

**查询参数：**
| 参数 | 说明 |
|------|------|
| query | 搜索查询（支持操作符语法） |
| page | 页码 |
| per_page | 每页数量 |

**查询语法示例：**
```
description_contains:"超市" amount_max:500 date_after:"2024-03-01"
```

### GET /api/v1/search/accounts

搜索账户。

## 自动补全接口

### GET /api/v1/autocomplete/accounts

账户补全。

**查询参数：**
| 参数 | 说明 |
|------|------|
| query | 搜索关键词 |
| type | 账户类型过滤 |
| limit | 返回数量（默认 10） |

### GET /api/v1/autocomplete/categories

分类补全。

### GET /api/v1/autocomplete/tags

标签补全。

### GET /api/v1/autocomplete/budgets

预算补全。

### GET /api/v1/autocomplete/bills

账单补全。

### GET /api/v1/autocomplete/currencies

币种补全。

### GET /api/v1/autocomplete/piggy-banks

储蓄罐补全。

## 报表接口

### GET /api/v1/summary/basic

基本财务摘要。

**查询参数：**
| 参数 | 说明 |
|------|------|
| start | 开始日期 |
| end | 结束日期 |
| currency_id | 币种过滤 |

**响应：**
```json
{
  "data": {
    "total_assets": "150000.00",
    "total_spent": "8500.00",
    "total_earned": "25000.00",
    "net_worth_change": "16500.00"
  }
}
```

### GET /api/v1/insight/expense

支出洞察。

**查询参数：**
| 参数 | 说明 |
|------|------|
| start | 开始日期 |
| end | 结束日期 |
| group_by | 分组维度 (category/account/tag) |

### GET /api/v1/insight/income

收入洞察。

### GET /api/v1/insight/transfer

转账洞察。

## 图表接口

### GET /api/v1/chart/balance

资产余额趋势图数据。

**查询参数：**
| 参数 | 说明 |
|------|------|
| start | 开始日期 |
| end | 结束日期 |
| account_id | 账户 ID（可选，不传则汇总） |

**响应：**
```json
{
  "data": {
    "labels": ["2024-03-01", "2024-03-02", ...],
    "datasets": [
      {
        "label": "招商银行",
        "data": [10000.00, 9850.00, ...]
      }
    ]
  }
}
```

### GET /api/v1/chart/account/overview

账户收支概览图。

### GET /api/v1/chart/budget/overview

预算使用概览图。

### GET /api/v1/chart/category/overview

分类收支概览图。

## 偏好设置接口

### GET /api/v1/preferences

列出用户偏好。

### POST /api/v1/preferences

创建/更新偏好。

**请求体：**
```json
{
  "name": "currencyPreference",
  "data": "CNY"
}
```

### GET /api/v1/preferences/:name

获取指定偏好。

## 用户管理接口（管理员）

### GET /api/v1/users

列出用户。

### POST /api/v1/users

创建用户。

### GET /api/v1/users/:id

获取用户详情。

### PUT /api/v1/users/:id

更新用户。

### DELETE /api/v1/users/:id

删除用户。

## 系统接口

### GET /api/v1/about

系统基本信息。

### GET /api/v1/configuration

系统配置（管理员）。

### PUT /api/v1/configuration

更新系统配置（管理员）。

### GET /api/v1/cron/:token

定时任务触发端点。
