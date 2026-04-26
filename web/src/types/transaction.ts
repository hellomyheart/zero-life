/**
 * 交易相关类型定义
 * 注意：字段名和类型必须与后端 JSON 完全对应
 */

/**
 * 交易类型枚举
 * - Deposit: 存款/收入，资金流入账户
 * - Withdrawal: 取款/支出，资金流出账户
 * - Transfer: 转账，账户间资金转移
 */
export enum TransactionType {
  Deposit = 'deposit',
  Withdrawal = 'withdrawal',
  Transfer = 'transfer',
}

/**
 * 交易中的账户信息（后端 AccountResp 嵌套）
 */
export interface TransactionAccountResp {
  id: number                    // 账户ID（后端 uint64）
  name: string                  // 账户名称
  type: string                  // 账户类型
  currency_id: number           // 货币ID
  currency: {                   // 货币信息
    id: number
    code: string
    name: string
    symbol: string
  }
  current_balance: string       // 当前余额
  is_virtual: boolean           // 是否虚拟账户
}

/**
 * 交易中的分类信息（后端 CategoryResp 嵌套）
 */
export interface TransactionCategoryResp {
  id: number                    // 分类ID（后端 uint64）
  name: string                  // 分类名称
}

/**
 * 交易中的标签信息（后端 TagResp 嵌套）
 */
export interface TransactionTagResp {
  id: number                    // 标签ID（后端 uint64）
  name: string                  // 标签名称
  color: string                 // 标签颜色
}

/**
 * 交易拆分项响应（后端 SplitResp）
 */
export interface SplitResp {
  id: number                    // 拆分ID（后端 uint64）
  amount: string                // 拆分金额
  category_id: number | null    // 分类ID
  category: TransactionCategoryResp | null // 分类信息
  tags: TransactionTagResp[]    // 标签列表
  notes: string                 // 备注
}

/**
 * 交易信息（后端 TransactionResp）
 * 字段名与后端 JSON tag 完全对应
 */
export interface Transaction {
  id: number                          // 交易ID（后端 uint64）
  type: TransactionType               // 交易类型
  date: string                        // 交易日期（ISO 8601 格式）
  description: string                 // 交易描述
  amount: string                      // 交易金额
  source_id: number                   // 源账户ID（后端 uint64）
  source: TransactionAccountResp      // 源账户信息
  destination_id: number | null       // 目标账户ID（转账时有值）
  destination: TransactionAccountResp | null // 目标账户信息
  category_id: number | null          // 分类ID
  category: TransactionCategoryResp | null // 分类信息
  notes: string                       // 备注
  tags: TransactionTagResp[]          // 标签列表（对象数组，不是ID数组）
  splits: SplitResp[]                 // 拆分交易列表
  bill_id: number | null              // 关联账单ID
  created_at: string                  // 创建时间
  updated_at: string                  // 更新时间
}

/**
 * 创建交易拆分请求
 * 字段名与后端 CreateSplitReq JSON tag 对应
 */
export interface CreateSplitReq {
  amount: string                // 拆分金额（必填）
  category_id?: number          // 分类ID（后端 *uint64）
  tags?: number[]               // 标签ID列表（后端 []uint64）
  description?: string          // 拆分描述
  notes?: string                // 备注
}

/**
 * 创建交易请求
 * 字段名与后端 CreateTransactionReq JSON tag 完全对应
 */
export interface CreateTransactionReq {
  type: TransactionType         // 交易类型（必填）
  date: string                  // 交易日期（必填，格式 YYYY-MM-DD HH:mm）
  description: string           // 交易描述（必填）
  amount: string                // 交易金额（必填，必须大于0）
  source_id: number | undefined    // 源账户ID（必填，后端 uint64）
  destination_id?: number       // 目标账户ID（转账时必填，后端 *uint64）
  category_id?: number          // 分类ID（后端 *uint64）
  notes?: string                // 备注
  tags?: number[]               // 标签ID列表（后端 []uint64）
  splits?: CreateSplitReq[]     // 拆分交易列表
}

/**
 * 更新交易请求
 * 字段名与后端 UpdateTransactionReq JSON tag 对应
 */
export interface UpdateTransactionReq {
  type: TransactionType         // 交易类型（必填）
  date: string                  // 交易日期（必填）
  description: string           // 交易描述（必填）
  amount: string                // 交易金额（必填）
  source_id: number             // 源账户ID（必填）
  destination_id?: number       // 目标账户ID
  category_id?: number          // 分类ID
  notes?: string                // 备注
  tags?: number[]               // 标签ID列表
  splits?: CreateSplitReq[]     // 拆分交易列表
}

/**
 * 交易列表查询请求（对应后端 TransactionListReq）
 */
export interface TransactionListReq {
  type?: TransactionType        // 交易类型筛选
  start_date?: string           // 开始日期
  end_date?: string             // 结束日期
  account_id?: number           // 账户ID筛选
  category_id?: number          // 分类ID筛选
  tag_id?: number               // 标签ID筛选
  sort?: string                 // 排序字段，默认 -date
  page?: number                 // 页码
  page_size?: number            // 每页数量
}

/**
 * 交易搜索请求（对应后端 TransactionSearchReq）
 */
export interface TransactionSearchReq {
  keyword?: string              // 关键词搜索（描述、备注）
  type?: TransactionType        // 交易类型筛选
  start_date?: string           // 开始日期
  end_date?: string             // 结束日期
  min_amount?: string           // 最小金额
  max_amount?: string           // 最大金额
  account_id?: number           // 账户ID筛选
  category_id?: number          // 分类ID筛选
  tag_id?: number               // 标签ID筛选
  sort?: string                 // 排序字段
  page?: number                 // 页码
  page_size?: number            // 每页数量
}
