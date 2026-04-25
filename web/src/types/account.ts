/**
 * 账户相关类型定义
 * 账户是记账系统的核心实体，用于记录资金的来源和去向
 * 账户分为四种类型：资产(asset)、支出(expense)、收入(income)、负债(liability)
 * 注意：字段名必须与后端 JSON 响应完全对应
 */

/**
 * 账户类型枚举
 * - Asset: 资产账户，如银行卡、现金、投资账户等（资金存放处）
 * - Expense: 支出账户，如餐饮、交通等（资金流出方向）
 * - Revenue: 收入账户，如工资、利息等（资金流入方向，后端字段值为 revenue）
 * - Liability: 负债账户，如信用卡、贷款等（欠款）
 */
export enum AccountType {
  Asset = 'asset',
  Expense = 'expense',
  Revenue = 'revenue',
  Liability = 'liability',
}

/**
 * 货币信息接口（账户响应中嵌套的货币对象）
 */
export interface AccountCurrencyResp {
  id: number                    // 货币ID
  code: string                  // 货币代码，如"CNY"、"USD"
  name: string                  // 货币名称
  symbol: string                // 货币符号，如"¥"、"$"
  decimal_places: number        // 小数位数
}

/**
 * 账户信息接口
 * 表示系统中的一个完整账户，字段与后端 AccountResp 完全对应
 */
export interface Account {
  id: number                    // 账户唯一标识（后端返回 uint64）
  name: string                  // 账户名称，如"招商银行储蓄卡"
  type: AccountType             // 账户类型
  currency_id: number           // 货币ID（后端返回 uint64）
  currency: AccountCurrencyResp // 关联货币信息对象
  initial_balance: string       // 初始余额（字符串避免浮点精度问题）
  current_balance: string       // 当前余额（由系统根据交易自动计算）
  is_virtual: boolean           // 是否为虚拟账户（虚拟账户不参与实际资金计算）
  notes: string                 // 备注
  created_at: string            // 创建时间
  updated_at: string            // 最后更新时间
}

/**
 * 创建账户请求接口
 * 创建新账户时需要提供的参数
 * 注意：字段名必须与后端 CreateAccountReq JSON tag 完全对应
 */
export interface CreateAccountReq {
  name: string                  // 账户名称（必填）
  type: AccountType             // 账户类型（必填，asset/expense/revenue/liability）
  currency_id: number           // 货币ID（必填，后端要求 uint64 的货币ID，不是货币代码）
  initial_balance: string       // 初始余额（必填）
  is_virtual?: boolean          // 是否为虚拟账户（可选，默认为false）
}

/**
 * 更新账户请求接口
 * 仅允许修改名称、备注和虚拟属性
 * 账户类型、货币、初始余额创建后不可修改
 */
export interface UpdateAccountReq {
  name?: string                 // 账户名称
  notes?: string                // 备注
  is_virtual?: boolean          // 是否为虚拟账户
}

/**
 * 账户列表查询请求接口
 * 用于筛选和分页查询账户列表
 */
export interface AccountListReq {
  type?: AccountType            // 按账户类型筛选
  search?: string               // 按名称搜索
  page?: number                 // 页码（从1开始）
  page_size?: number            // 每页数量
  sort?: string                 // 排序字段
}
