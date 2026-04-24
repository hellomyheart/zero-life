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
 * 交易拆分信息
 * 用于一笔交易拆分为多笔子交易
 */
export interface Split {
  id: string              // 拆分ID
  amount: string          // 拆分金额
  category_id: string     // 分类ID
  category_name: string   // 分类名称
  tag_ids: string[]       // 标签ID列表
  tag_names: string[]     // 标签名称列表
  description: string     // 拆分描述
}

/**
 * 交易信息
 * 记录收入、支出、转账等财务交易
 */
export interface Transaction {
  id: string                      // 交易ID
  type: TransactionType           // 交易类型
  date: string                    // 交易日期
  description: string             // 交易描述
  amount: string                  // 交易金额
  source_account_id: string       // 源账户ID
  source_account_name: string     // 源账户名称
  destination_account_id: string  // 目标账户ID
  destination_account_name: string // 目标账户名称
  category_id: string             // 分类ID
  category_name: string           // 分类名称
  tag_ids: string[]               // 标签ID列表
  tag_names: string[]             // 标签名称列表
  splits: Split[]                 // 拆分交易列表
  created_at: string              // 创建时间
  updated_at: string              // 更新时间
}

/**
 * 创建交易拆分请求
 */
export interface CreateSplitReq {
  amount: string          // 拆分金额
  category_id: string     // 分类ID
  tag_ids: string[]       // 标签ID列表
  description: string     // 拆分描述
}

/**
 * 创建交易请求
 */
export interface CreateTransactionReq {
  type: TransactionType           // 交易类型
  date: string                    // 交易日期
  description: string             // 交易描述
  amount: string                  // 交易金额
  source_account_id: string       // 源账户ID
  destination_account_id?: string // 目标账户ID（转账时需要）
  category_id?: string            // 分类ID
  tag_ids?: string[]              // 标签ID列表
  splits?: CreateSplitReq[]       // 拆分交易列表
}

/**
 * 更新交易请求
 */
export interface UpdateTransactionReq {
  type?: TransactionType          // 交易类型
  date?: string                   // 交易日期
  description?: string            // 交易描述
  amount?: string                 // 交易金额
  source_account_id?: string      // 源账户ID
  destination_account_id?: string // 目标账户ID
  category_id?: string            // 分类ID
  tag_ids?: string[]              // 标签ID列表
  splits?: CreateSplitReq[]       // 拆分交易列表
}

/**
 * 交易列表查询请求
 */
export interface TransactionListReq {
  type?: TransactionType          // 交易类型筛选
  start_date?: string             // 开始日期
  end_date?: string               // 结束日期
  source_account_id?: string      // 源账户ID筛选
  destination_account_id?: string // 目标账户ID筛选
  category_id?: string            // 分类ID筛选
  tag_ids?: string[]              // 标签ID筛选
  keyword?: string                // 关键词搜索
  page?: number                   // 页码
  page_size?: number              // 每页数量
}
