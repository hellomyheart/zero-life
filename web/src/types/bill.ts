/**
 * 账单相关类型定义
 * 账单用于管理定期重复的固定支出（如房租、水电费、订阅服务等）
 */

/**
 * 重复规则枚举
 * 定义账单的重复频率
 * - Daily: 每天重复
 * - Weekly: 每周重复
 * - Monthly: 每月重复
 * - Quarterly: 每季度重复
 * - Yearly: 每年重复
 */
export enum RepeatRule {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Quarterly = 'quarterly',
  Yearly = 'yearly',
}

/**
 * 账单信息接口
 * 表示一个定期账单的完整信息
 */
export interface Bill {
  id: string                   // 账单唯一标识
  name: string                 // 账单名称，如"月租费"
  amount: string               // 账单金额
  account_id: string           // 关联的支出账户ID
  account_name: string         // 关联的支出账户名称（冗余字段，方便展示）
  category_id: string          // 关联的分类ID
  category_name: string        // 关联的分类名称（冗余字段，方便展示）
  repeat_rule: RepeatRule      // 重复规则
  next_due_date: string        // 下次到期日期
  is_overdue: boolean          // 是否已逾期（当前日期超过到期日但未支付）
  description: string          // 账单描述说明
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建账单请求接口
 */
export interface CreateBillReq {
  name: string                 // 账单名称（必填）
  amount: string               // 账单金额（必填）
  account_id: string           // 支出账户ID（必填）
  category_id: string          // 分类ID（必填）
  repeat_rule: RepeatRule      // 重复规则（必填）
  next_due_date: string        // 首次到期日期（必填）
  description?: string         // 账单描述（可选）
}

/**
 * 更新账单请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateBillReq {
  name?: string                // 账单名称
  amount?: string              // 账单金额
  account_id?: string          // 支出账户ID
  category_id?: string         // 分类ID
  repeat_rule?: RepeatRule     // 重复规则
  next_due_date?: string       // 下次到期日期
  description?: string         // 账单描述
}
