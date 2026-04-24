# Zero Life 代码注释规范

## 目的
本文档旨在为 Go 后端和 Vue 前端代码提供统一的注释规范，帮助开发者（尤其是 Go 和 Vue 初学者）快速理解代码结构和业务逻辑。

## Go 后端注释规范

### 1. 文件头注释
每个 Go 文件开头应包含包说明：
```go
// Package model 定义了系统的数据模型，对应数据库表结构
package model
```

### 2. 类型/结构体注释
每个类型定义前应有详细说明：
```go
// Account 账户模型，对应 accounts 表
// 管理用户的各类财务账户，包括资产、负债、收入、支出账户
// 
// 主要功能：
// - 支持多种账户类型（asset/expense/revenue/liability）
// - 自动计算当前余额
// - 支持多货币
// - 支持软删除
type Account struct {
    // ID 账户唯一标识，主键自增
    ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
    
    // UserID 所属用户 ID
    // 每个账户都属于一个特定用户，用户间数据隔离
    UserID uint64 `gorm:"not null;index" json:"user_id"`
    
    // Type 账户类型
    // 可选值：asset(资产), expense(支出), revenue(收入), liability(负债)
    // 资产账户：银行账户、现金账户等
    // 支出账户：用于记录支出分类（如：餐饮、交通）
    // 收入账户：用于记录收入分类（如：工资、奖金）
    // 负债账户：信用卡、贷款等
    Type AccountType `gorm:"not null;size:20" json:"type"`
}
```

### 3. 函数注释
每个公开函数应有完整注释，包括参数和返回值说明：
```go
// Create 创建账户
// 
// 业务逻辑：
// 1. 验证账户名称不能为空
// 2. 验证账户类型是否合法
// 3. 检查货币是否存在
// 4. 设置初始余额（如果未提供则默认为 0）
// 5. 计算当前余额（初始余额 + 所有交易）
// 6. 保存到数据库
//
// 参数：
//   - userID: 用户 ID
//   - req: 创建账户请求参数
//
// 返回：
//   - *Account: 创建的账户对象
//   - error: 错误信息，成功时为 nil
//
// 示例：
//   account, err := service.Create(1, &CreateAccountReq{
//       Name: "招商银行卡",
//       Type: "asset",
//       CurrencyID: 1,
//       InitialBalance: decimal.NewFromFloat(1000),
//   })
func (s *AccountService) Create(userID uint64, req *request.CreateAccountReq) (*model.Account, error) {
    // 实现代码...
}
```

### 4. 复杂业务逻辑注释
对于复杂的业务逻辑，应在关键步骤添加行内注释：
```go
// CalculateBalance 计算账户余额
func (s *AccountService) CalculateBalance(accountID uint64) (decimal.Decimal, error) {
    // 1. 获取账户初始余额
    account, err := s.repo.GetByID(accountID)
    if err != nil {
        return decimal.Zero, err
    }
    
    // 2. 查询所有与该账户相关的交易
    // 包括：作为源账户的支出/转账，作为目标账户的收入/转账
    transactions, err := s.transactionRepo.GetByAccountID(accountID)
    if err != nil {
        return decimal.Zero, err
    }
    
    // 3. 从初始余额开始，累加所有交易金额
    balance := account.InitialBalance
    for _, tx := range transactions {
        switch tx.Type {
        case model.TransactionTypeWithdrawal:
            // 支出：减少余额
            if tx.SourceID == accountID {
                balance = balance.Sub(tx.Amount)
            }
        case model.TransactionTypeDeposit:
            // 收入：增加余额
            if tx.DestinationID != nil && *tx.DestinationID == accountID {
                balance = balance.Add(tx.Amount)
            }
        case model.TransactionTypeTransfer:
            // 转账：根据账户角色增减
            if tx.SourceID == accountID {
                balance = balance.Sub(tx.Amount)
            } else if tx.DestinationID != nil && *tx.DestinationID == accountID {
                balance = balance.Add(tx.Amount)
            }
        }
    }
    
    return balance, nil
}
```

### 5. 枚举类型注释
枚举值应有清晰的说明：
```go
// TransactionType 交易类型枚举
// 定义三种基本交易类型，决定资金流向和账户角色
type TransactionType string

const (
    // TransactionTypeDeposit 存款/收入
    // 资金流入：从收入账户（如工资）到资产账户（如银行卡）
    // 特点：增加资产账户余额
    TransactionTypeDeposit TransactionType = "deposit"
    
    // TransactionTypeWithdrawal 取款/支出
    // 资金流出：从资产账户（如银行卡）到支出账户（如餐饮）
    // 特点：减少资产账户余额
    TransactionTypeWithdrawal TransactionType = "withdrawal"
    
    // TransactionTypeTransfer 转账
    // 账户间转移：从一个资产账户到另一个资产账户
    // 特点：不改变总资产，只是资金在不同账户间移动
    TransactionTypeTransfer TransactionType = "transfer"
)
```

## Vue 前端注释规范

### 1. 文件头注释
每个组件文件开头应说明组件功能：
```vue
<!--
  DashboardPage.vue - 仪表盘页面
  功能：展示用户财务概况，包括月度收支、预算预警、账单提醒等
  
  主要模块：
  - 收支概览卡片：本月收入、支出、净收入
  - 资产总览：所有资产账户余额总和
  - 预算预警：显示超支或即将超支的预算
  - 账单提醒：显示即将到期的账单
  - 最近交易：显示最新的 5 笔交易
  
  数据来源：/api/v1/dashboard 接口
-->
```

### 2. Script 块注释
```vue
<script setup lang="ts">
/**
 * 仪表盘页面组件
 * 
 * 功能说明：
 * - 页面加载时自动获取仪表盘数据
 * - 支持数据刷新
 * - 响应式布局，适配不同屏幕尺寸
 * 
 * 数据流：
 * 1. onMounted 钩子调用 getDashboard() API
 * 2. 数据保存到 dashboard ref
 * 3. 模板根据数据渲染各个卡片和表格
 * 
 * @module pages/dashboard
 */

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getDashboard } from '@/api/dashboard'
import { formatAmount, formatDate } from '@/utils/format'
import type { DashboardResp } from '@/types/dashboard'

const { t } = useI18n()

/** 仪表盘数据，包含收支、预算、账单等信息 */
const dashboard = ref<DashboardResp | null>(null)

/** 加载状态，true 表示正在请求数据 */
const loading = ref(true)

/**
 * 获取仪表盘数据
 * 调用 API 获取用户的财务概览数据
 * 
 * 错误处理：
 * - API 调用失败时静默处理，显示空状态
 * - 最终都会将 loading 设为 false
 */
onMounted(async () => {
  try {
    dashboard.value = await getDashboard() as unknown as DashboardResp
  } catch (error) {
    // TODO: 添加错误提示
    console.error('Failed to fetch dashboard data:', error)
  } finally {
    loading.value = false
  }
})
</script>
```

### 3. 模板注释
对于复杂的模板结构，添加 HTML 注释：
```vue
<template>
  <div class="dashboard-page" v-loading="loading">
    <!-- 收支概览卡片 -->
    <!-- 三列布局：收入、支出、净收入 -->
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>{{ t('dashboard.income') }}</template>
          <div class="amount income">{{ formatAmount(dashboard.total_income) }}</div>
        </el-card>
      </el-col>
      <!-- ... 其他卡片 ... -->
    </el-row>

    <!-- 预算预警表格 -->
    <!-- 显示使用率超过 80% 的预算，超支时显示红色警告 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>{{ t('dashboard.budgetAlerts') }}</template>
          <el-table :data="dashboard.budget_alerts" size="small">
            <!-- 使用率进度条：绿色 (<80%) -> 黄色 (80-100%) -> 红色 (>100%) -->
            <el-table-column prop="usage_rate" :label="t('budget.usageRate')">
              <template #default="{ row }">
                <el-progress 
                  :percentage="Math.round(row.usage_rate * 100)" 
                  :status="row.status === 'exceeded' ? 'exception' : row.status === 'warning' ? 'warning' : undefined" 
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
```

### 4. 函数注释
```typescript
/**
 * 格式化金额显示
 * 将原始金额数据转换为带货币符号的字符串
 * 
 * @param amount - 原始金额数字
 * @param currency - 货币代码（可选，默认为 CNY）
 * @returns 格式化后的金额字符串，如 "¥1,234.56"
 * 
 * @example
 * formatAmount(1234.56) // => "¥1,234.56"
 * formatAmount(-500, 'USD') // => "-$500.00"
 */
function formatAmount(amount: number, currency?: string): string {
  // 实现代码...
}
```

### 5. 复杂逻辑注释
```typescript
/**
 * 处理预算预警状态
 * 根据预算使用率确定预警级别
 * 
 * 业务规则：
 * - 使用率 < 80%: 正常 (normal)
 * - 使用率 80%-100%: 预警 (warning)
 * - 使用率 > 100%: 超支 (exceeded)
 * 
 * @param usageRate - 预算使用率（0-1 之间的小数）
 * @returns 预警状态字符串
 */
function getBudgetStatus(usageRate: number): 'normal' | 'warning' | 'exceeded' {
  if (usageRate > 1) {
    return 'exceeded' // 超支
  } else if (usageRate > 0.8) {
    return 'warning' // 预警
  } else {
    return 'normal' // 正常
  }
}
```

## 特殊场景注释

### 1. 重要业务规则
```go
// ⚠️ 重要业务规则：
// 账户余额计算必须考虑所有相关交易，包括：
// 1. 作为源账户的支出交易（减少余额）
// 2. 作为目标账户的收入交易（增加余额）
// 3. 作为源账户或目标账户的转账交易（减少或增加）
// 4. 已对账的交易（确保准确性）
// 
// 注意：计算时必须过滤已软删除的交易（DeletedAt != nil）
```

### 2. 性能优化说明
```go
// 性能优化：
// 使用批量查询而非 N+1 查询
// 先一次性获取所有账户 ID，然后批量查询余额
// 避免在循环中执行 SQL 查询
// 
// 优化前：O(N) 次查询（N 为账户数量）
// 优化后：O(1) 次查询
```

### 3. TODO/FIXME 注释
```go
// TODO: 添加缓存机制，减少数据库查询
// FIXME: 当账户数量>1000 时性能下降，需要分页处理
// XXX: 这里有个临时方案，等 v2 API 上线后移除
```

## 注释质量检查清单

- [ ] 每个公开函数都有完整的注释
- [ ] 复杂业务逻辑有详细的步骤说明
- [ ] 枚举值有清晰的含义说明
- [ ] 关键算法有复杂度分析或性能说明
- [ ] 重要业务规则有⚠️标记
- [ ] TODO/FIXME 注释标注了优先级
- [ ] 注释使用中文（团队约定）
- [ ] 注释简洁明了，避免废话

## 示例文件

完整示例请参考：
- Go: `server/internal/model/rule.go` (规则引擎模型)
- Go: `server/internal/service/transaction_service.go` (交易服务)
- Vue: `web/src/pages/transactions/TransactionListPage.vue` (交易列表页)
- Vue: `web/src/pages/objectGroups/ObjectGroupListPage.vue` (对象分组页)
