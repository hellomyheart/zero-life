/**
 * 循环交易相关类型定义
 * 循环交易用于记录定期重复发生的交易（如每月工资、每周通勤费等）
 * 系统会在到期时自动根据循环交易创建正式的交易记录
 */

/**
 * 重复类型枚举
 * 定义循环交易的重复频率
 * - Daily: 每天重复
 * - Weekly: 每周重复
 * - Monthly: 每月重复
 * - Yearly: 每年重复
 */
export enum RecurrenceType {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Yearly = 'yearly',
}

/**
 * 循环交易信息接口
 * 表示一个定期重复的交易模板
 */
export interface RecurringTransaction {
  id: string                   // 循环交易唯一标识
  title: string                // 交易标题，如"每月房租"
  type: string                 // 交易类型，如"deposit"、"withdrawal"、"transfer"
  amount: string               // 交易金额
  source_account_id: string    // 源账户ID（资金流出的账户）
  destination_account_id: string // 目标账户ID（资金流入的账户，转账时使用）
  category_id: string          // 关联分类ID
  recurrence_type: RecurrenceType // 重复类型（每天/每周/每月/每年）
  repeat_interval: number      // 重复间隔（如 recurrence_type=Monthly 且 repeat_interval=2 表示每两月一次）
  start_date: string           // 首次执行日期
  next_date: string            // 下次执行日期（系统自动计算）
  end_date: string             // 结束日期（超过此日期不再重复）
  active: boolean              // 是否活跃（暂停时设为false）
  max_repetitions: number      // 最大重复次数（0表示无限重复）
  description: string          // 交易描述
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建循环交易请求接口
 */
export interface CreateRecurringTransactionReq {
  title: string                // 交易标题（必填）
  type: string                 // 交易类型（必填）
  amount: string               // 交易金额（必填）
  source_account_id: string    // 源账户ID（必填）
  destination_account_id?: string // 目标账户ID（转账时必填）
  category_id?: string         // 关联分类ID（可选）
  recurrence_type: RecurrenceType // 重复类型（必填）
  repeat_interval: number      // 重复间隔（必填）
  start_date: string           // 首次执行日期（必填）
  end_date?: string            // 结束日期（可选，不填则无限重复）
  max_repetitions?: number     // 最大重复次数（可选，不填则无限重复）
  description?: string         // 交易描述（可选）
}

/**
 * 更新循环交易请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateRecurringTransactionReq {
  title?: string               // 交易标题
  amount?: string              // 交易金额
  recurrence_type?: RecurrenceType // 重复类型
  repeat_interval?: number     // 重复间隔
  start_date?: string          // 首次执行日期
  end_date?: string            // 结束日期
  active?: boolean             // 是否活跃
  max_repetitions?: number     // 最大重复次数
  description?: string         // 交易描述
}

/**
 * 循环交易列表查询请求接口
 */
export interface RecurringTransactionListReq {
  active?: boolean             // 是否仅显示活跃的循环交易
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}
