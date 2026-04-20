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
| 时区 | 请求日期默认 UTC，可通过偏好设置指定 |

### 统一响应格式

**单资源成功响应：**

```json
{
  "data": {
    "id": 1,
    "name": "招商银行储蓄卡",
    "created_at": "2024-03-15T10:30:00Z",
    "updated_at": "2024-03-15T10:30:00Z"
  }
}
```

**列表分页响应：**

```json
{
  "data": [
    { "id": 1, "name": "..." },
    { "id": 2, "name": "..." }
  ],
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
    "next": "/api/v1/accounts?page=2",
    "prev": null
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
      { "field": "name", "message": "名称不能为空" },
      { "field": "amount", "message": "金额必须为正数" }
    ]
  }
}
```

### HTTP 状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|----------|
| 200 | 成功 | GET/PUT 成功 |
| 201 | 创建成功 | POST 创建资源成功 |
| 204 | 删除成功 | DELETE 成功（无返回体） |
| 400 | 请求参数错误 | 参数格式、类型不合法 |
| 401 | 未认证 | Token 缺失、过期、无效 |
| 403 | 无权限 | 角色权限不足 |
| 404 | 资源不存在 | ID 对应资源不存在 |
| 409 | 冲突 | 重复创建（如邮箱已注册） |
| 422 | 业务校验失败 | 业务规则不满足（如账户类型不匹配） |
| 429 | 请求限流 | 超出频率限制 |
| 500 | 服务器内部错误 | 未预期的服务端异常 |

### 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码，从 1 开始 |
| per_page | int | 50 | 每页数量，范围 1-500 |

### 通用日期范围参数

| 参数 | 类型 | 格式 | 说明 |
|------|------|------|------|
| start | string | YYYY-MM-DD | 开始日期（含） |
| end | string | YYYY-MM-DD | 结束日期（含） |

### 金额字段约定

- 所有金额以**字符串**形式传输，避免浮点精度丢失
- 正数表示贷方（资金流入），负数表示借方（资金流出）
- 前端展示时取绝对值，根据交易类型决定正负号显示

---

## 认证接口

### GET /api/v1/auth/me

获取当前登录用户信息。

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "default_group_id": 1,
    "two_factor_enabled": false,
    "blocked": false,
    "groups": [
      { "id": 1, "title": "默认组", "user_role": "OWNER" }
    ]
  }
}
```

### POST /api/v1/auth/register

注册新用户，同时自动创建默认用户组。

**请求体：**

| 字段 | 类型 | 必填 | 校验 | 说明 |
|------|------|------|------|------|
| email | string | 是 | email, max=255 | 邮箱地址，全局唯一 |
| password | string | 是 | min=8, max=72 | 密码 |
| password_confirmation | string | 是 | 必须与 password 一致 | 确认密码 |

**响应：** `201 Created`

```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "default_group_id": 1,
    "created_at": "2024-03-15T10:00:00Z"
  }
}
```

### POST /api/v1/auth/login

登录获取 Token。如果用户启用了 2FA，返回临时 Token 需要先验证。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码 |

**响应（无 2FA）：** `200 OK`

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
  }
}
```

**响应（已启用 2FA）：** `200 OK`

```json
{
  "data": {
    "two_factor_required": true,
    "temp_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 300
  }
}
```

### POST /api/v1/auth/refresh

使用 refresh_token 刷新 access_token。

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

### POST /api/v1/auth/logout

注销当前 Token，加入 Redis 黑名单。

**请求头：** `Authorization: Bearer <access_token>`

**响应：** `204 No Content`

### POST /api/v1/auth/2fa/enable

启用双因素认证。返回 TOTP Secret 和 QR Code URL，用户需扫码后调用 2fa/verify 完成启用。

**请求头：** `Authorization: Bearer <access_token>`

**响应：** `200 OK`

```json
{
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code_url": "otpauth://totp/ZeroLife:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=ZeroLife"
  }
}
```

### POST /api/v1/auth/2fa/verify

验证双因素认证码。首次验证完成 2FA 启用，后续验证返回完整访问 Token。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| code | string | 是 | 6 位 TOTP 验证码 |

**响应：** `200 OK`

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
  }
}
```

---

## 账户接口

### GET /api/v1/accounts

列出账户，支持按类型、活跃状态过滤和搜索。

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type | string | - | 账户类型：asset/expense/revenue/debt/loan/mortgage |
| active | bool | - | 活跃状态过滤：true/false |
| search | string | - | 名称模糊搜索 |
| page | int | 1 | 页码 |
| per_page | int | 50 | 每页数量 |

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "name": "招商银行储蓄卡",
      "account_type": "asset",
      "account_role": "defaultAsset",
      "currency_id": 1,
      "currency_code": "CNY",
      "currency_symbol": "¥",
      "virtual_balance": "0",
      "current_balance": "50000.00",
      "is_active": true,
      "is_virtual": false,
      "object_group_id": 1,
      "order": 0,
      "notes": "主储蓄卡",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-03-15T10:00:00Z"
    }
  ],
  "meta": { "page": 1, "per_page": 50, "total": 5, "total_pages": 1 }
}
```

### POST /api/v1/accounts

创建账户。如果是资产类账户，可同时设置期初余额。

**请求体：**

| 字段 | 类型 | 必填 | 校验 | 说明 |
|------|------|------|------|------|
| name | string | 是 | max=255 | 账户名称 |
| type | string | 是 | 枚举 | 账户类型 |
| account_role | string | 否 | 枚举，仅资产类 | 账户角色：defaultAsset/sharedAsset/savingAsset/ccAsset/cashWalletAsset |
| currency_id | uint64 | 是 | 存在的币种 | 默认币种 |
| virtual_balance | string | 否 | decimal | 虚拟余额（显示用，不参与计算） |
| opening_balance | string | 否 | decimal | 期初余额（仅资产类） |
| opening_balance_date | string | 否 | date | 期初余额日期 |
| is_active | bool | 否 | 默认 true | 是否活跃 |
| is_virtual | bool | 否 | 默认 false | 是否虚拟账户 |
| object_group_id | uint64 | 否 | - | 对象分组 ID |
| notes | string | 否 | - | 备注 |
| interest | string | 否 | decimal | 利率（仅负债类） |
| interest_period | string | 否 | 枚举 | 利息周期：daily/weekly/monthly/quarterly/half-year/yearly |

**响应：** `201 Created`

### GET /api/v1/accounts/:id

获取账户详情，含当前余额计算。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| date | string | 计算余额的截止日期 (YYYY-MM-DD)，默认当天 |

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "name": "招商银行储蓄卡",
    "account_type": "asset",
    "account_role": "defaultAsset",
    "currency_id": 1,
    "current_balance": "50000.00",
    "virtual_balance": "0",
    "is_active": true,
    "notes": "主储蓄卡",
    "meta": [
      { "name": "account_role", "data": "defaultAsset" }
    ],
    "location": { "latitude": 0, "longitude": 0, "zoom_level": 6 },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-03-15T10:00:00Z"
  }
}
```

### PUT /api/v1/accounts/:id

更新账户。请求体与 POST 相同，所有字段可选。

**响应：** `200 OK`

### DELETE /api/v1/accounts/:id

删除账户（软删除）。如果账户下有交易，返回 422 错误。

**响应：** `204 No Content`

### GET /api/v1/accounts/:id/transactions

获取账户关联交易，支持分页和日期范围。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | 交易类型过滤：withdrawal/deposit/transfer |
| start | string | 开始日期 |
| end | string | 结束日期 |
| page | int | 页码 |
| per_page | int | 每页数量 |

### GET /api/v1/accounts/:id/piggy-banks

获取账户关联的储蓄罐列表。

### GET /api/v1/accounts/:id/attachments

获取账户关联附件列表。

---

## 交易接口

### GET /api/v1/transactions

列出交易组，支持多维度过滤。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | 交易类型：withdrawal/deposit/transfer |
| start | string | 开始日期 (YYYY-MM-DD) |
| end | string | 结束日期 (YYYY-MM-DD) |
| category_id | uint64 | 分类 ID 过滤 |
| budget_id | uint64 | 预算 ID 过滤 |
| bill_id | uint64 | 账单 ID 过滤 |
| tag_id | uint64 | 标签 ID 过滤 |
| account_id | uint64 | 账户 ID 过滤 |
| page | int | 页码 |
| per_page | int | 每页数量 |

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "title": null,
      "journals": [
        {
          "id": 1,
          "transaction_type_id": "withdrawal",
          "description": "超市购物",
          "date": "2024-03-15T00:00:00Z",
          "currency_id": 1,
          "category_id": 3,
          "budget_id": 2,
          "bill_id": null,
          "tags": ["日常", "食品"],
          "transactions": [
            {
              "id": 1,
              "account_id": 1,
              "amount": "-156.80",
              "currency_id": 1
            },
            {
              "id": 2,
              "account_id": 5,
              "amount": "156.80",
              "currency_id": 1
            }
          ]
        }
      ],
      "created_at": "2024-03-15T10:00:00Z",
      "updated_at": "2024-03-15T10:00:00Z"
    }
  ],
  "meta": { "page": 1, "per_page": 50, "total": 30, "total_pages": 1 }
}
```

### POST /api/v1/transactions

创建交易（支持拆分交易）。创建后异步触发规则引擎和 Webhook。

**请求体：**

| 字段 | 类型 | 必填 | 校验 | 说明 |
|------|------|------|------|------|
| type | string | 是 | withdrawal/deposit/transfer | 交易类型 |
| date | string | 是 | date | 交易日期 |
| description | string | 是 | max=255 | 交易描述 |
| title | string | 否 | max=255 | 交易组标题（拆分交易时建议填写） |
| transactions | array | 是 | min=1 | 交易明细数组 |

**transactions 数组每项：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source_id | uint64 | 是 | 来源账户 ID |
| destination_id | uint64 | 是 | 目标账户 ID（不能与 source_id 相同） |
| amount | string | 是 | 金额（正数字符串） |
| currency_id | uint64 | 是 | 交易币种 ID |
| foreign_amount | string | 否 | 外币金额 |
| foreign_currency_id | uint64 | 否 | 外币币种 ID |
| category_id | uint64 | 否 | 分类 ID |
| budget_id | uint64 | 否 | 预算 ID（仅支出类型） |
| bill_id | uint64 | 否 | 账单 ID |
| tags | string[] | 否 | 标签文本数组 |
| description | string | 否 | 拆分项描述（覆盖主描述） |
| notes | string | 否 | 备注 |
| reconciled | bool | 否 | 默认 false，是否已对账 |

**请求体示例（单笔支出）：**

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
      "notes": "周末采购"
    }
  ]
}
```

**请求体示例（拆分交易）：**

```json
{
  "type": "withdrawal",
  "date": "2024-03-15",
  "description": "超市购物",
  "title": "超市购物",
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

**请求体示例（跨币种交易）：**

```json
{
  "type": "withdrawal",
  "date": "2024-03-15",
  "description": "Amazon购物",
  "transactions": [
    {
      "source_id": 1,
      "destination_id": 5,
      "amount": "50.00",
      "currency_id": 2,
      "foreign_amount": "350.00",
      "foreign_currency_id": 1,
      "category_id": 3
    }
  ]
}
```

**响应：** `201 Created`

### GET /api/v1/transactions/:id

获取交易组详情，包含所有日志和交易记录。

**响应：** `200 OK`（同列表中的单条结构）

### PUT /api/v1/transactions/:id

更新交易组。请求体与 POST 相同，会整体替换交易内容。

**响应：** `200 OK`

### DELETE /api/v1/transactions/:id

删除交易组（软删除），同时删除所有关联的交易日志和交易记录。

**响应：** `204 No Content`

### GET /api/v1/transactions/:id/attachments

获取交易关联附件列表。

### GET /api/v1/transactions/:id/piggy-bank-events

获取交易关联的储蓄罐事件列表。

---

## 预算接口

### GET /api/v1/budgets

列出预算，含当前周期的限额和已用金额。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 限额开始日期 (YYYY-MM-DD) |
| end | string | 限额结束日期 (YYYY-MM-DD) |
| page | int | 页码 |
| per_page | int | 每页数量 |

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "name": "餐饮",
      "auto_budget_type": 1,
      "auto_budget_amount": "2000.00",
      "auto_budget_period": "monthly",
      "is_active": true,
      "order": 0,
      "current_limit": {
        "id": 10,
        "amount": "2000.00",
        "start_date": "2024-03-01",
        "end_date": "2024-03-31",
        "spent": "1800.00",
        "left": "200.00"
      },
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /api/v1/budgets

创建预算。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 预算名称，max=255 |
| auto_budget_type | int | 否 | 自动预算类型：1=RESET, 2=ROLLOVER, 3=ADJUSTED |
| auto_budget_amount | string | 否 | 自动预算金额（decimal） |
| auto_budget_period | string | 否 | 周期：daily/weekly/monthly/quarterly/half-year/yearly |
| is_active | bool | 否 | 默认 true |
| order | int | 否 | 排序顺序 |

**响应：** `201 Created`

### GET /api/v1/budgets/:id

获取预算详情，含当前周期限额和已用金额。

### PUT /api/v1/budgets/:id

更新预算。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/budgets/:id

删除预算。

### GET /api/v1/budgets/:id/limits

获取预算的所有限额记录。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期过滤 |
| end | string | 结束日期过滤 |

### POST /api/v1/budgets/:id/limits

创建预算限额。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| currency_id | uint64 | 是 | 币种 ID |
| start_date | string | 是 | 限额开始日期 |
| end_date | string | 是 | 限额结束日期 |
| amount | string | 是 | 限额金额（decimal） |

**响应：** `201 Created`

### PUT /api/v1/budgets/:id/limits/:limitId

更新预算限额。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/budgets/:id/limits/:limitId

删除预算限额。

### GET /api/v1/budgets/:id/transactions

获取预算关联交易，支持分页和日期范围。

---

## 账单接口

### GET /api/v1/bills

列出账单，含下次预期日期和已支付状态。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| active | bool | 活跃状态过滤 |
| page | int | 页码 |
| per_page | int | 每页数量 |

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "name": "房租",
      "amount_min": "3000.00",
      "amount_max": "3500.00",
      "date": "2024-04-01",
      "repeat_freq": "monthly",
      "is_active": true,
      "paid": false,
      "next_expected_match": "2024-04-01",
      "currency_id": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /api/v1/bills

创建账单。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 账单名称，max=255 |
| match | string | 是 | 匹配规则，多个用逗号分隔 |
| amount_min | string | 是 | 最低金额（decimal） |
| amount_max | string | 是 | 最高金额（decimal） |
| date | string | 是 | 下次预期日期 |
| end_date | string | 否 | 结束日期 |
| extension_date | string | 否 | 延期日期 |
| repeat_freq | string | 是 | 重复频率：daily/weekly/monthly/quarterly/half-year/yearly |
| skip | int | 否 | 跳过次数，默认 0 |
| automatch | bool | 否 | 是否自动匹配，默认 true |
| is_active | bool | 否 | 默认 true |
| currency_id | uint64 | 是 | 币种 ID |
| object_group_id | uint64 | 否 | 对象分组 ID |
| notes | string | 否 | 备注 |

**响应：** `201 Created`

### GET /api/v1/bills/:id

获取账单详情。

### PUT /api/v1/bills/:id

更新账单。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/bills/:id

删除账单。

### GET /api/v1/bills/:id/transactions

获取账单关联交易，支持分页。

### GET /api/v1/bills/:id/attachments

获取账单关联附件。

### GET /api/v1/bills/:id/rules

获取与账单关联的规则列表。

---

## 储蓄罐接口

### GET /api/v1/piggy-banks

列出储蓄罐，含进度百分比。

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "name": "旅行基金",
      "target_amount": "20000.00",
      "current_amount": "8000.00",
      "progress_percent": 40.0,
      "start_date": "2024-01-01",
      "target_date": "2024-12-31",
      "account_id": 1,
      "currency_id": 1,
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /api/v1/piggy-banks

创建储蓄罐。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 储蓄罐名称，max=255 |
| account_id | uint64 | 是 | 关联账户 ID |
| target_amount | string | 是 | 目标金额（decimal） |
| start_date | string | 否 | 起始日期 |
| target_date | string | 否 | 目标日期 |
| currency_id | uint64 | 是 | 币种 ID |
| object_group_id | uint64 | 否 | 对象分组 ID |
| notes | string | 否 | 备注 |

**响应：** `201 Created`

### GET /api/v1/piggy-banks/:id

获取储蓄罐详情，含进度百分比。

### PUT /api/v1/piggy-banks/:id

更新储蓄罐。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/piggy-banks/:id

删除储蓄罐。

### POST /api/v1/piggy-banks/:id/add

向储蓄罐存入金额。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | string | 是 | 存入金额（正数） |
| transaction_journal_id | uint64 | 否 | 关联的交易日志 ID |

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "current_amount": "8500.00",
    "progress_percent": 42.5
  }
}
```

### POST /api/v1/piggy-banks/:id/remove

从储蓄罐取出金额。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | string | 是 | 取出金额（正数） |
| transaction_journal_id | uint64 | 否 | 关联的交易日志 ID |

### GET /api/v1/piggy-banks/:id/events

获取储蓄事件列表。

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "piggy_bank_id": 1,
      "journal_id": 100,
      "amount": "500.00",
      "created_at": "2024-03-15T10:00:00Z"
    }
  ]
}
```

### GET /api/v1/piggy-banks/:id/attachments

获取储蓄罐关联附件。

### GET /api/v1/piggy-banks/:id/accounts

获取储蓄罐关联账户。

---

## 分类接口

### GET /api/v1/categories

列出分类。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| per_page | int | 每页数量 |

### POST /api/v1/categories

创建分类。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 分类名称，max=255，同用户组内唯一 |
| notes | string | 否 | 备注 |

**响应：** `201 Created`

### GET /api/v1/categories/:id

获取分类详情，含当月支出/收入统计。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 统计开始日期 |
| end | string | 统计结束日期 |

**响应：** `200 OK`

```json
{
  "data": {
    "id": 3,
    "name": "食品",
    "spent": "3200.00",
    "earned": "0.00",
    "notes": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### PUT /api/v1/categories/:id

更新分类。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/categories/:id

删除分类。分类下的交易不会被删除，但分类关联会被清除。

### GET /api/v1/categories/:id/transactions

获取分类关联交易，支持分页和日期范围。

### GET /api/v1/categories/:id/attachments

获取分类关联附件。

---

## 标签接口

### GET /api/v1/tags

列出标签。

### POST /api/v1/tags

创建标签。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag | string | 是 | 标签文本，max=255，同用户组内唯一 |
| date | string | 否 | 标签日期 |
| description | string | 否 | 标签描述 |
| tag_mode | string | 否 | 标签模式：nothing/balancing/expenseTransfer/incomeTransfer |

**响应：** `201 Created`

### GET /api/v1/tags/:id

获取标签详情，含当月统计。

### PUT /api/v1/tags/:id

更新标签。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/tags/:id

删除标签。标签与交易的关联会被清除。

### GET /api/v1/tags/:id/transactions

获取标签关联交易，支持分页和日期范围。

### GET /api/v1/tags/:id/attachments

获取标签关联附件。

---

## 币种接口

### GET /api/v1/currencies

列出所有币种。

### POST /api/v1/currencies

创建币种。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 货币全名，max=128 |
| code | string | 是 | 货币代码，max=16，全局唯一（如 CNY） |
| symbol | string | 是 | 货币符号，max=16（如 ¥） |
| decimal_places | int | 否 | 小数位数，默认 2 |
| enabled | bool | 否 | 是否启用，默认 true |
| is_default | bool | 否 | 是否设为默认，默认 false |

**响应：** `201 Created`

### GET /api/v1/currencies/:code

获取币种详情（按代码查询）。

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "name": "人民币",
    "code": "CNY",
    "symbol": "¥",
    "decimal_places": 2,
    "enabled": true,
    "is_default": true
  }
}
```

### PUT /api/v1/currencies/:code

更新币种。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/currencies/:code

删除币种。如果有交易使用该币种，返回 422 错误。

### POST /api/v1/currencies/:code/enable

为当前用户启用该币种。

**响应：** `204 No Content`

### POST /api/v1/currencies/:code/disable

为当前用户禁用该币种。如果有交易使用该币种，返回 422 错误。

### POST /api/v1/currencies/:code/default

将该币种设为当前用户的默认币种。

### GET /api/v1/currencies/:code/accounts

获取使用该币种的账户列表。

### GET /api/v1/currencies/:code/available-budgets

获取该币种的可用预算列表。

### GET /api/v1/currencies/:code/bills

获取该币种的账单列表。

### GET /api/v1/currencies/:code/budget-limits

获取该币种的预算限额列表。

### GET /api/v1/currencies/:code/cer

获取该币种的汇率记录。

### GET /api/v1/currencies/:code/recurrences

获取该币种的定期交易列表。

### GET /api/v1/currencies/:code/rules

获取该币种的规则列表。

### GET /api/v1/currencies/:code/transactions

获取该币种的交易列表。

---

## 汇率接口

### GET /api/v1/exchange-rates

列出汇率记录。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| from | string | 源币种代码 |
| to | string | 目标币种代码 |
| date | string | 汇率日期 (YYYY-MM-DD) |
| page | int | 页码 |
| per_page | int | 每页数量 |

### POST /api/v1/exchange-rates

创建汇率记录。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| from_currency_id | uint64 | 是 | 源币种 ID |
| to_currency_id | uint64 | 是 | 目标币种 ID |
| date | string | 是 | 汇率日期 |
| rate | string | 是 | 汇率值（decimal，如 7.8500 表示 1 USD = 7.85 CNY） |

**响应：** `201 Created`

### GET /api/v1/exchange-rates/:id

获取汇率详情。

### PUT /api/v1/exchange-rates/:id

更新汇率。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/exchange-rates/:id

删除汇率记录。

---

## 规则接口

### GET /api/v1/rules

列出规则，含触发器和动作。

### POST /api/v1/rules

创建规则。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| rule_group_id | uint64 | 是 | 所属规则组 ID |
| title | string | 是 | 规则名称，max=255 |
| description | string | 否 | 规则描述 |
| strict | bool | 否 | 严格模式（AND），默认 true |
| stop_processing | bool | 否 | 执行后停止后续规则，默认 false |
| order | int | 否 | 组内执行顺序 |
| is_active | bool | 否 | 默认 true |
| triggers | array | 是 | 触发器数组，min=1 |
| actions | array | 是 | 动作数组，min=1 |

**triggers 数组每项：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 触发器类型（见下方枚举） |
| value | string | 否 | 触发器值 |
| order | int | 否 | 检查顺序 |
| active | bool | 否 | 默认 true |
| stop_processing | bool | 否 | 匹配后停止检查后续触发器 |

**触发器类型枚举：**

| 类型 | 说明 | 需要值 |
|------|------|--------|
| description_contains / description_starts / description_ends / description_is | 描述匹配 | 是 |
| amount_less / amount_more / amount_is | 金额比较 | 是（数字） |
| foreign_amount_less / foreign_amount_more / foreign_amount_is | 外币金额比较 | 是 |
| from_account_starts / from_account_ends / from_account_is / from_account_contains | 来源账户匹配 | 是 |
| to_account_starts / to_account_ends / to_account_is / to_account_contains | 目标账户匹配 | 是 |
| account_is / account_contains | 任意账户匹配 | 是 |
| category_is | 分类匹配 | 是 |
| tag_is | 标签匹配 | 是 |
| has_category / has_no_category / has_tag / has_no_tag | 存在性检查 | 否 |
| date_before / date_after / date_is | 日期匹配 | 是（YYYY-MM-DD） |
| transaction_type | 交易类型匹配 | 是（withdrawal/deposit/transfer） |
| notes_contain / notes_is / notes_empty / notes_not_empty | 备注匹配 | 视情况 |
| budget_is / bill_is | 预算/账单匹配 | 是 |
| has_attachments | 有附件 | 否 |
| reconciled | 已对账 | 否 |
| user_action | 用户手动触发（始终为 true） | 否 |

**actions 数组每项：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 动作类型（见下方枚举） |
| value | string | 否 | 动作值（可使用表达式语言 `=amount*0.1`） |
| order | int | 否 | 执行顺序 |
| active | bool | 否 | 默认 true |
| stop_processing | bool | 否 | 执行后停止后续动作 |

**动作类型枚举：**

| 类型 | 说明 | 需要值 |
|------|------|--------|
| set_category / clear_category | 设置/清除分类 | 分类名称 / 无 |
| set_budget / clear_budget | 设置/清除预算 | 预算名称 / 无 |
| add_tag / remove_tag / remove_all_tags | 标签操作 | 标签文本 / 无 |
| set_description / append_description | 描述操作 | 文本 |
| set_notes / clear_notes | 备注操作 | 文本 / 无 |
| set_source_account / set_destination_account | 设置账户 | 账户名称 |
| set_source_to_cash / set_destination_to_cash | 设为现金账户 | 无 |
| set_amount | 设置金额 | 新金额 |
| link_to_bill | 关联账单 | 账单名称 |
| convert_withdrawal / convert_deposit / convert_transfer | 转换交易类型 | 无 |
| switch_accounts | 交换来源和目标账户 | 无 |
| update_piggy | 更新储蓄罐 | 储蓄罐名称\|金额 |
| delete_transaction | 删除交易 | 无 |

**请求体示例：**

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

**响应：** `201 Created`

### GET /api/v1/rules/:id

获取规则详情，含完整触发器和动作列表。

### PUT /api/v1/rules/:id

更新规则。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/rules/:id

删除规则。

### POST /api/v1/rules/:id/test

测试规则（预览匹配结果，不修改数据）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| transaction_ids | uint64[] | 否 | 指定测试的交易 ID 列表，为空则测试所有交易 |

**响应：** `200 OK`

```json
{
  "data": {
    "matched_count": 3,
    "matched_transactions": [
      { "id": 10, "description": "超市购物", "amount": "156.80" }
    ],
    "proposed_actions": [
      { "transaction_id": 10, "action": "set_category", "value": "食品" }
    ]
  }
}
```

### POST /api/v1/rules/:id/trigger

触发规则（实际执行，修改数据）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| transaction_ids | uint64[] | 否 | 指定触发的交易 ID 列表 |

**响应：** `200 OK`

```json
{
  "data": {
    "executed_count": 3,
    "actions_applied": 6
  }
}
```

### GET /api/v1/rules/:id/validate-expression

验证规则动作中的表达式语言。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| expression | string | 表达式字符串 |

**响应：** `200 OK`

```json
{
  "data": {
    "valid": true,
    "result": "15.68"
  }
}
```

---

## 规则组接口

### GET /api/v1/rule-groups

列出规则组，含组内规则数量。

### POST /api/v1/rule-groups

创建规则组。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 规则组名称，max=255 |
| description | string | 否 | 规则组描述 |
| order | int | 否 | 执行顺序 |
| is_active | bool | 否 | 默认 true |
| stop_processing | bool | 否 | 执行完本组后停止，默认 false |

**响应：** `201 Created`

### GET /api/v1/rule-groups/:id

获取规则组详情，含组内规则列表。

### PUT /api/v1/rule-groups/:id

更新规则组。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/rule-groups/:id

删除规则组及组内所有规则。

### POST /api/v1/rule-groups/:id/test

测试规则组（预览所有组内规则效果）。

**请求体：** 同规则测试

### POST /api/v1/rule-groups/:id/trigger

触发规则组（实际执行组内所有规则）。

**请求体：** 同规则触发

### GET /api/v1/rule-groups/:id/rules

获取组内规则列表。

---

## 定期交易接口

### GET /api/v1/recurrences

列出定期交易。

### POST /api/v1/recurrences

创建定期交易。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题，max=255 |
| description | string | 否 | 描述 |
| type | string | 是 | 交易类型：withdrawal/deposit/transfer |
| first_date | string | 是 | 首次执行日期 |
| repeat_until | string | 否 | 重复截止日期 |
| apply_rules | bool | 否 | 创建后是否触发规则，默认 false |
| is_active | bool | 否 | 默认 true |
| repetitions | array | 是 | 重复模式数组 |
| transactions | array | 是 | 交易模板数组 |

**repetitions 数组每项：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 重复类型：daily/weekly/monthly/ndom/yearly |
| moment | string | 否 | 重复时刻（weekly:1-7, monthly:1-31, ndom:第N个星期D） |
| skip | int | 否 | 跳过间隔，默认 0 |
| weekend | int | 否 | 周末处理：1=保持/2=前周五/3=后周一/4=前周五不跨月 |

**transactions 数组每项：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source_id | uint64 | 是 | 来源账户 ID |
| destination_id | uint64 | 是 | 目标账户 ID |
| amount | string | 是 | 金额 |
| currency_id | uint64 | 是 | 币种 ID |
| foreign_amount | string | 否 | 外币金额 |
| foreign_currency_id | uint64 | 否 | 外币币种 ID |
| description | string | 否 | 描述 |
| category_id | uint64 | 否 | 分类 ID |
| budget_id | uint64 | 否 | 预算 ID |
| bill_id | uint64 | 否 | 账单 ID |
| tags | string[] | 否 | 标签数组 |
| notes | string | 否 | 备注 |

**请求体示例：**

```json
{
  "title": "每月房租",
  "type": "withdrawal",
  "first_date": "2024-01-01",
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

**响应：** `201 Created`

### GET /api/v1/recurrences/:id

获取定期交易详情。

### PUT /api/v1/recurrences/:id

更新定期交易。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/recurrences/:id

删除定期交易。

### POST /api/v1/recurrences/:id/trigger

手动触发定期交易，立即创建一笔交易。

**响应：** `200 OK`

```json
{
  "data": {
    "transaction_group_id": 100,
    "message": "定期交易已触发，交易组 #100 已创建"
  }
}
```

### GET /api/v1/recurrences/:id/transactions

获取该定期交易已生成的交易列表。

---

## Webhook 接口

### GET /api/v1/webhooks

列出 Webhook。

### POST /api/v1/webhooks

创建 Webhook。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | Webhook 名称，max=255 |
| url | string | 是 | 回调 URL，max=512，必须为 https |
| secret | string | 否 | 签名密钥，max=128 |
| trigger | string | 是 | 触发条件（见枚举） |
| response | string | 否 | 响应内容，默认 RELEVANT |
| delivery | string | 否 | 交付方式，默认 JSON |
| active | bool | 否 | 默认 true |

**trigger 枚举：**

| 值 | 说明 |
|----|------|
| ANY | 任意事件（匹配所有触发条件） |
| STORE_TRANSACTION | 创建交易时 |
| UPDATE_TRANSACTION | 更新交易时 |
| DESTROY_TRANSACTION | 删除交易时 |
| STORE_BUDGET | 创建预算时 |
| UPDATE_BUDGET | 更新预算时 |
| DESTROY_BUDGET | 删除预算时 |
| STORE_UPDATE_BUDGET_LIMIT | 创建或更新预算限额时 |

**response 枚举：** TRANSACTIONS / ACCOUNTS / BUDGET / RELEVANT / NONE

**响应：** `201 Created`

### GET /api/v1/webhooks/:id

获取 Webhook 详情。

### PUT /api/v1/webhooks/:id

更新 Webhook。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/webhooks/:id

删除 Webhook。

### POST /api/v1/webhooks/:id/submit

手动提交未发送的消息（重试失败的发送）。

**响应：** `200 OK`

```json
{
  "data": {
    "resubmitted_count": 2,
    "message": "2 条消息已重新提交"
  }
}
```

### POST /api/v1/webhooks/:id/trigger-transaction

手动触发指定交易的 Webhook。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| transaction_id | uint64 | 是 | 交易 ID |

**响应：** `200 OK`

```json
{
  "data": {
    "triggered": true,
    "message_id": 15,
    "message": "Webhook 已触发"
  }
}
```

### GET /api/v1/webhooks/:id/messages

获取 Webhook 消息列表，含发送状态。

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "sent": true,
      "errored": false,
      "created_at": "2024-03-15T10:00:00Z"
    }
  ]
}
```

### GET /api/v1/webhooks/:id/attempts

获取发送尝试列表，含 HTTP 状态码和响应。

---

## 附件接口

### GET /api/v1/attachments

列出附件，支持按关联实体过滤。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| attachable_type | string | 关联实体类型：account/bill/budget/category/piggy_bank/tag/journal |
| attachable_id | uint64 | 关联实体 ID |
| page | int | 页码 |
| per_page | int | 每页数量 |

### POST /api/v1/attachments

创建附件记录（仅创建元数据，需再调用 upload 上传文件）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| attachable_type | string | 是 | 关联实体类型 |
| attachable_id | uint64 | 是 | 关联实体 ID |
| filename | string | 是 | 文件名，max=255 |
| title | string | 否 | 标题 |
| description | string | 否 | 描述 |

**响应：** `201 Created`

### POST /api/v1/attachments/:id/upload

上传附件文件内容。

**请求：** `multipart/form-data`

| 字段 | 类型 | 说明 |
|------|------|------|
| file | file | 文件内容，最大 10MB |

**响应：** `200 OK`

### GET /api/v1/attachments/:id/download

下载附件文件内容。

**响应：** 文件流，Content-Type 为附件的 MIME 类型

### PUT /api/v1/attachments/:id

更新附件元数据。

### DELETE /api/v1/attachments/:id

删除附件（含文件）。

---

## 搜索接口

### GET /api/v1/search/transactions

搜索交易，支持操作符语法。

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| query | string | 是 | 搜索查询（操作符语法） |
| page | int | 否 | 页码 |
| per_page | int | 否 | 每页数量 |

**查询语法：** `key:value key2:"value with spaces"`

**操作符示例：**

| 操作符 | 示例 | 说明 |
|--------|------|------|
| description_contains | `description_contains:"超市"` | 描述包含 |
| description_is | `description_is:"工资"` | 描述精确匹配 |
| amount_is | `amount_is:156.80` | 金额精确匹配 |
| amount_min | `amount_min:100` | 金额最小值 |
| amount_max | `amount_max:500` | 金额最大值 |
| date_after | `date_after:"2024-03-01"` | 日期晚于 |
| date_before | `date_before:"2024-03-31"` | 日期早于 |
| source_is | `source_is:"招商银行"` | 来源账户名 |
| destination_is | `destination_is:"超市"` | 目标账户名 |
| category_is | `category_is:"食品"` | 分类名 |
| tag_is | `tag_is:"日常"` | 标签名 |
| has_any_category | `has_any_category:true` | 有分类 |
| has_any_tag | `has_any_tag:true` | 有标签 |
| on_budget | `on_budget:true` | 有预算 |
| reconciled | `reconciled:true` | 已对账 |

**组合示例：**
```
description_contains:"超市" amount_max:500 date_after:"2024-03-01" category_is:"食品"
```

### GET /api/v1/search/accounts

搜索账户。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| query | string | 搜索关键词 |
| type | string | 账户类型过滤 |
| page | int | 页码 |
| per_page | int | 每页数量 |

---

## 自动补全接口

所有自动补全接口返回简化的 `{id, name}` 列表，用于前端下拉选择。

### GET /api/v1/autocomplete/accounts

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| query | string | 搜索关键词，min=1 |
| type | string | 账户类型过滤 |
| limit | int | 返回数量，默认 10，最大 50 |

**响应：** `200 OK`

```json
{
  "data": [
    { "id": 1, "name": "招商银行储蓄卡", "type": "asset" },
    { "id": 2, "name": "支付宝", "type": "asset" }
  ]
}
```

### GET /api/v1/autocomplete/categories

**查询参数：** query, limit

### GET /api/v1/autocomplete/tags

**查询参数：** query, limit

### GET /api/v1/autocomplete/budgets

**查询参数：** query, limit

### GET /api/v1/autocomplete/bills

**查询参数：** query, limit

### GET /api/v1/autocomplete/currencies

**查询参数：** query, limit

### GET /api/v1/autocomplete/piggy-banks

**查询参数：** query, limit

### GET /api/v1/autocomplete/rules

**查询参数：** query, limit

### GET /api/v1/autocomplete/rule-groups

**查询参数：** query, limit

### GET /api/v1/autocomplete/recurrences

**查询参数：** query, limit

### GET /api/v1/autocomplete/webhooks

**查询参数：** query, limit

### GET /api/v1/autocomplete/transaction-types

返回所有交易类型。

**响应：** `200 OK`

```json
{
  "data": [
    { "id": "withdrawal", "name": "支出" },
    { "id": "deposit", "name": "收入" },
    { "id": "transfer", "name": "转账" }
  ]
}
```

---

## 报表与洞察接口

### GET /api/v1/summary/basic

基本财务摘要，返回指定日期范围内的汇总数据。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 (YYYY-MM-DD) |
| end | string | 结束日期 (YYYY-MM-DD) |
| currency_id | uint64 | 币种过滤（可选） |

**响应：** `200 OK`

```json
{
  "data": {
    "total_assets": "150000.00",
    "total_liabilities": "800000.00",
    "net_worth": "-650000.00",
    "total_spent": "8500.00",
    "total_earned": "25000.00",
    "net_worth_change": "16500.00"
  }
}
```

### GET /api/v1/insight/expense

支出洞察，按指定维度分组统计支出。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |
| group_by | string | 分组维度：category/account/tag（默认 category） |
| currency_id | uint64 | 币种过滤 |

**响应：** `200 OK`

```json
{
  "data": [
    { "id": 3, "name": "食品", "amount": "3200.00", "percent": 37.6 },
    { "id": 4, "name": "交通", "amount": "1500.00", "percent": 17.6 },
    { "id": 5, "name": "娱乐", "amount": "1200.00", "percent": 14.1 }
  ]
}
```

### GET /api/v1/insight/income

收入洞察。参数和响应格式同 expense。

### GET /api/v1/insight/transfer

转账洞察。参数和响应格式同 expense。

---

## 图表接口

### GET /api/v1/chart/balance

资产余额趋势图数据，返回时间序列数据。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |
| account_id | uint64 | 账户 ID（可选，不传则汇总所有资产账户） |
| currency_id | uint64 | 币种过滤 |

**响应：** `200 OK`

```json
{
  "data": {
    "labels": ["2024-03-01", "2024-03-02", "2024-03-03"],
    "datasets": [
      {
        "label": "招商银行",
        "data": ["50000.00", "49850.00", "50350.00"]
      },
      {
        "label": "支付宝",
        "data": ["3200.50", "3100.50", "3500.50"]
      }
    ]
  }
}
```

### GET /api/v1/chart/account/overview

账户收支概览图，返回指定账户在日期范围内的收入/支出汇总。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |
| account_id | uint64 | 账户 ID |

**响应：** `200 OK`

```json
{
  "data": {
    "labels": ["收入", "支出", "转账流入", "转账流出"],
    "datasets": [
      {
        "label": "招商银行",
        "data": ["25000.00", "8500.00", "3000.00", "1500.00"]
      }
    ]
  }
}
```

### GET /api/v1/chart/budget/overview

预算使用概览图，返回各预算的限额和已用金额。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |

**响应：** `200 OK`

```json
{
  "data": {
    "labels": ["餐饮", "交通", "娱乐"],
    "datasets": [
      { "label": "预算限额", "data": ["2000.00", "1000.00", "1500.00"] },
      { "label": "已支出", "data": ["1800.00", "500.00", "300.00"] }
    ]
  }
}
```

### GET /api/v1/chart/category/overview

分类收支概览图。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |

**响应：** `200 OK`

```json
{
  "data": {
    "labels": ["食品", "交通", "娱乐", "日用品"],
    "datasets": [
      { "label": "支出", "data": ["3200.00", "1500.00", "1200.00", "800.00"] },
      { "label": "收入", "data": ["0.00", "0.00", "0.00", "0.00"] }
    ]
  }
}
```

---

## 偏好设置接口

### GET /api/v1/preferences

列出当前用户的所有偏好设置。

**响应：** `200 OK`

```json
{
  "data": [
    { "name": "currencyPreference", "data": "CNY" },
    { "name": "languagePreference", "data": "zh-CN" },
    { "name": "viewRangePreference", "data": "1M" },
    { "name": "fiscalYearStart", "data": "1" }
  ]
}
```

### POST /api/v1/preferences

创建或更新偏好设置（upsert 语义）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 偏好名称，max=128 |
| data | string | 是 | 偏好数据 |

**响应：** `200 OK`

### GET /api/v1/preferences/:name

获取指定偏好。

**响应：** `200 OK`

```json
{
  "data": {
    "name": "currencyPreference",
    "data": "CNY"
  }
}
```

### PUT /api/v1/preferences/:name

更新指定偏好（upsert 语义，不存在则创建）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| data | any | 是 | 偏好值 |

**响应：** `200 OK`

### DELETE /api/v1/preferences/:name

删除指定偏好，恢复为默认值。

---

## 用户管理接口（管理员）

### GET /api/v1/users

列出所有用户（需要 OWNER 全局角色）。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| per_page | int | 每页数量 |

### POST /api/v1/users

创建用户（管理员创建）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱 |
| password | string | 是 | 密码，min=8 |
| blocked | bool | 否 | 是否封禁，默认 false |

**响应：** `201 Created`

### GET /api/v1/users/:id

获取用户详情，含用户组和角色信息。

### PUT /api/v1/users/:id

更新用户。可修改 email、password、blocked 状态。

### DELETE /api/v1/users/:id

删除用户（软删除）。

---

## 对象分组接口

### GET /api/v1/object-groups

列出对象分组。

### POST /api/v1/object-groups

创建对象分组。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 分组标题，max=255 |
| order | int | 否 | 排序顺序 |

### GET /api/v1/object-groups/:id

获取对象分组详情。

### PUT /api/v1/object-groups/:id

更新对象分组。

### DELETE /api/v1/object-groups/:id

删除对象分组。

---

## 交易链接接口

### GET /api/v1/transaction-links

列出交易链接。

### POST /api/v1/transaction-links

创建交易链接，关联两个交易日志。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| link_type_id | uint64 | 是 | 链接类型 ID |
| source_id | uint64 | 是 | 源交易日志 ID |
| destination_id | uint64 | 是 | 目标交易日志 ID |

### GET /api/v1/transaction-links/:id

获取交易链接详情。

### PUT /api/v1/transaction-links/:id

更新交易链接。

### DELETE /api/v1/transaction-links/:id

删除交易链接。

### GET /api/v1/link-types

列出链接类型。

### POST /api/v1/link-types

创建链接类型。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 名称 |
| outward | string | 是 | 正向描述 |
| inward | string | 是 | 反向描述 |
| editable | bool | 否 | 是否可编辑，默认 true |

### GET /api/v1/link-types/:id

获取链接类型详情。

### PUT /api/v1/link-types/:id

更新链接类型。

### DELETE /api/v1/link-types/:id

删除链接类型（仅可编辑的）。

---

## 交易日志接口

### GET /api/v1/transaction-journals/:id

获取交易日志详情，含关联的分类、预算、标签、交易记录。

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "transaction_type_id": "withdrawal",
    "description": "超市购物",
    "date": "2024-03-15T00:00:00Z",
    "currency_id": 1,
    "category_id": 3,
    "budget_id": 2,
    "bill_id": null,
    "tags": ["日常", "食品"],
    "transactions": [
      { "id": 1, "account_id": 1, "amount": "-156.80", "currency_id": 1 },
      { "id": 2, "account_id": 5, "amount": "156.80", "currency_id": 1 }
    ],
    "notes": "周末采购",
    "created_at": "2024-03-15T10:00:00Z"
  }
}
```

### DELETE /api/v1/transaction-journals/:id

删除交易日志（软删除），同时删除关联的交易记录。

**响应：** `204 No Content`

### GET /api/v1/transaction-journals/:id/links

获取交易日志的关联链接列表。

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "link_type_id": 1,
      "link_type_name": "Related",
      "source_id": 1,
      "destination_id": 5,
      "inward": "相关",
      "outward": "相关"
    }
  ]
}
```

---

## 可用预算接口

### GET /api/v1/available-budgets

列出可用预算，支持按币种和日期范围过滤。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| currency_id | uint64 | 币种 ID 过滤 |
| start | string | 开始日期 (YYYY-MM-DD) |
| end | string | 结束日期 (YYYY-MM-DD) |
| page | int | 页码 |
| per_page | int | 每页数量 |

### POST /api/v1/available-budgets

创建可用预算。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| currency_id | uint64 | 是 | 币种 ID |
| amount | string | 是 | 可用金额（decimal） |
| start_date | string | 是 | 开始日期 |
| end_date | string | 是 | 结束日期 |

**响应：** `201 Created`

### GET /api/v1/available-budgets/:id

获取可用预算详情。

### PUT /api/v1/available-budgets/:id

更新可用预算。请求体与 POST 相同，所有字段可选。

### DELETE /api/v1/available-budgets/:id

删除可用预算。

---

## 用户组与成员管理接口

### GET /api/v1/user-groups

列出用户可见的用户组。

### POST /api/v1/user-groups

创建用户组。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 组名称，max=255 |

**响应：** `201 Created`

### GET /api/v1/user-groups/:id

获取用户组详情，含成员列表。

### PUT /api/v1/user-groups/:id

更新用户组。

### DELETE /api/v1/user-groups/:id

删除用户组（仅 OWNER 可操作）。

### GET /api/v1/user-groups/:id/memberships

获取用户组成员列表。

**响应：** `200 OK`

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "user_email": "admin@example.com",
      "user_role": "OWNER",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /api/v1/user-groups/:id/memberships

添加成员到用户组。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | uint64 | 是 | 用户 ID |
| user_role | string | 是 | 组内角色（见 UserRoleEnum） |

### PUT /api/v1/user-groups/:id/memberships/:membershipId

更新成员角色。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_role | string | 是 | 新的组内角色 |

### DELETE /api/v1/user-groups/:id/memberships/:membershipId

移除成员。

---

## 邀请接口

### POST /api/v1/invitations

创建邀请。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_group_id | uint64 | 是 | 邀请加入的用户组 ID |
| email | string | 是 | 被邀请者邮箱 |

**响应：** `201 Created`

```json
{
  "data": {
    "id": 1,
    "email": "newuser@example.com",
    "invite_code": "abc123def456",
    "expires_at": "2024-04-20T00:00:00Z"
  }
}
```

### POST /api/v1/invitations/redeem

兑换邀请码。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| invite_code | string | 是 | 邀请码 |

**响应：** `200 OK`

### GET /api/v1/invitations

列出当前用户组内的邀请（需 OWNER 或 FULL 角色）。

### DELETE /api/v1/invitations/:id

撤销邀请。

---

## 全局角色接口

### GET /api/v1/roles

列出所有全局角色。

**响应：** `200 OK`

```json
{
  "data": [
    { "id": 1, "name": "owner", "description": "站点所有者" },
    { "id": 2, "name": "demo", "description": "演示用户" }
  ]
}
```

### POST /api/v1/users/:id/roles

为用户分配全局角色（需站点 OWNER）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | uint64 | 是 | 角色 ID |

### DELETE /api/v1/users/:id/roles/:roleId

移除用户的全局角色（需站点 OWNER）。

---

## 数据导出接口

### POST /api/v1/export/accounts

导出账户数据为 CSV。

**响应：** CSV 文件流，Content-Type: text/csv

### POST /api/v1/export/bills

导出账单数据为 CSV。

### POST /api/v1/export/budgets

导出预算数据为 CSV。

### POST /api/v1/export/categories

导出分类数据为 CSV。

### POST /api/v1/export/piggy-banks

导出储蓄罐数据为 CSV。

### POST /api/v1/export/recurrences

导出定期交易数据为 CSV。

### POST /api/v1/export/rules

导出规则数据为 CSV。

### POST /api/v1/export/tags

导出标签数据为 CSV。

### POST /api/v1/export/transactions

导出交易数据为 CSV。

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| start | string | 开始日期 |
| end | string | 结束日期 |
| type | string | 交易类型过滤 |

---

## 数据导入接口

### POST /api/v1/imports

创建导入任务。

**请求体：** `multipart/form-data`

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 导入文件（CSV） |
| type | string | 是 | 导入类型：csv/spectre/nordigen/bunq |

**响应：** `201 Created`

```json
{
  "data": {
    "id": 1,
    "type": "csv",
    "status": "uploaded",
    "created_at": "2024-03-15T10:00:00Z"
  }
}
```

### GET /api/v1/imports/:id

获取导入任务状态。

**响应：** `200 OK`

```json
{
  "data": {
    "id": 1,
    "type": "csv",
    "status": "completed",
    "imported_count": 150,
    "skipped_count": 5,
    "error_count": 2,
    "created_at": "2024-03-15T10:00:00Z"
  }
}
```

### POST /api/v1/imports/:id/configure

配置导入列映射（CSV 导入第二步）。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| column_map | object | 是 | 列映射：系统字段 -> CSV 列索引 |
| date_format | string | 否 | 日期格式，默认 Y-m-d |
| delimiter | string | 否 | 分隔符，默认 , |
| currency_id | uint64 | 否 | 默认币种 ID |
| duplicate_check | bool | 否 | 是否启用重复检测，默认 true |

### POST /api/v1/imports/:id/preview

预览导入数据（不实际写入）。

**响应：** `200 OK`

```json
{
  "data": {
    "total_rows": 200,
    "preview_rows": [
      { "date": "2024-03-15", "description": "超市购物", "amount": "156.80" }
    ],
    "duplicates_count": 5
  }
}
```

### POST /api/v1/imports/:id/execute

执行导入。

**响应：** `200 OK`

```json
{
  "data": {
    "imported_count": 150,
    "skipped_count": 5,
    "error_count": 2,
    "errors": [
      { "row": 10, "message": "账户不存在" }
    ]
  }
}
```

### DELETE /api/v1/imports/:id

删除导入任务。

---

## 数据销毁接口

### POST /api/v1/data/destroy

选择性销毁指定类型的数据。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| types | string[] | 是 | 要销毁的数据类型：accounts/bills/budgets/categories/piggy_banks/recurrences/rules/tags/transactions |

**响应：** `200 OK`

```json
{
  "data": {
    "destroyed": {
      "transactions": 500,
      "accounts": 20,
      "categories": 15
    },
    "message": "数据已销毁"
  }
}
```

### POST /api/v1/data/purge

彻底清除当前用户组的所有财务数据。

**响应：** `200 OK`

```json
{
  "data": {
    "message": "所有数据已清除"
  }
}
```

> **警告：** 此操作不可逆，将删除所有交易、账户、预算、分类、标签、规则等财务数据。

---

## 批量操作接口

### POST /api/v1/transactions/bulk

批量更新交易。

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| transaction_ids | uint64[] | 是 | 要更新的交易 ID 列表 |
| category_id | uint64 | 否 | 批量设置分类 |
| budget_id | uint64 | 否 | 批量设置预算 |
| add_tags | string[] | 否 | 批量添加标签 |
| remove_tags | string[] | 否 | 批量移除标签 |
| notes | string | 否 | 批量设置备注 |

**响应：** `200 OK`

```json
{
  "data": {
    "updated_count": 25,
    "message": "25 笔交易已更新"
  }
}
```

---

## 系统接口

### GET /api/v1/health

健康检查（无需认证）。

**响应：** `200 OK`

```json
{
  "status": "ok"
}
```

### GET /api/v1/about

获取系统基本信息（无需认证）。

**响应：** `200 OK`

```json
{
  "data": {
    "version": "1.0.0",
    "api_version": "1.0.0",
    "go_version": "go1.22",
    "environment": "production"
  }
}
```

### GET /api/v1/configuration

获取系统配置（需要 OWNER 全局角色）。

### PUT /api/v1/configuration

更新系统配置（需要 OWNER 全局角色）。

**请求体：**

```json
{
  "data": {
    "is_demo_site": false,
    "single_user_mode": false,
    "must_confirm_account": true
  }
}
```

### GET /api/v1/cron/:token

定时任务触发端点。token 为系统配置的 cron_token，用于防止未授权调用。

**执行的任务：**
1. 检查并创建到期的定期交易
2. 检查到期账单并发送提醒
3. 从外部源下载最新汇率
4. 检查是否有新版本可用
5. 清理过期的 session 和日志

**响应：** `200 OK`

```json
{
  "data": {
    "recurrences_created": 2,
    "bills_reminded": 1,
    "rates_updated": 28,
    "version_checked": true,
    "message": "Cron job executed successfully"
  }
}
```
