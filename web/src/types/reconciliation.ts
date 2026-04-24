/**
 * 对账相关类型定义
 * 对账用于核对账户的账面余额与实际余额是否一致，帮助发现记账错误或遗漏
 */

/**
 * 对账记录接口
 * 表示一次对账操作的完整信息
 */
export interface Reconciliation {
  id: string                   // 对账记录唯一标识
  account_id: string           // 对账的账户ID
  start_date: string           // 对账期间起始日期
  end_date: string             // 对账期间结束日期
  start_balance: string        // 期初余额（对账期间开始时的账面余额）
  end_balance: string          // 期末余额（对账期间结束时的实际余额）
  book_balance: string         // 账面余额（系统根据交易计算出的余额）
  difference: string           // 差额（实际余额与账面余额的差异，0表示一致）
  status: string               // 对账状态，如"reconciled"（已对平）、"discrepancy"（有差异）
  notes: string                // 对账备注说明
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建对账记录请求接口
 */
export interface CreateReconciliationReq {
  account_id: string           // 对账的账户ID（必填）
  start_date: string           // 起始日期（必填）
  end_date: string             // 结束日期（必填）
  start_balance: string        // 期初余额（必填）
  end_balance: string          // 期末余额（必填）
  notes?: string               // 备注说明（可选）
}

/**
 * 更新对账记录请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateReconciliationReq {
  end_balance?: string         // 期末余额
  status?: string              // 对账状态
  notes?: string               // 备注说明
}

/**
 * 对账记录列表查询请求接口
 */
export interface ReconciliationListReq {
  account_id?: string          // 按账户ID筛选
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}
