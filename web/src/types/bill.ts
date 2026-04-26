/**
 * 账单相关类型定义
 * 账单用于管理定期重复的固定支出（如房租、水电费、订阅服务等）
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 重复规则枚举
 * - Daily: 每天重复
 * - Weekly: 每周重复
 * - Monthly: 每月重复
 * - Yearly: 每年重复
 */
export enum RepeatRule {
  Daily = 'daily',
  Weekly = 'weekly',
  Monthly = 'monthly',
  Yearly = 'yearly',
}

/**
 * 账单信息接口（对应后端 BillResp）
 */
export interface Bill {
  id: number                    // 账单唯一标识（后端 uint64）
  name: string                  // 账单名称
  amount: string                // 账单金额
  repeat_rule: RepeatRule       // 重复规则
  next_due: string              // 下次到期日
  source_id: number | null      // 支出账户ID
  category_id: number | null    // 分类ID
  notes: string                 // 备注
  created_at: string            // 创建时间
  updated_at: string            // 最后更新时间
}

/**
 * 创建账单请求接口（对应后端 CreateBillReq）
 */
export interface CreateBillReq {
  name: string                  // 账单名称（必填）
  amount: string                // 账单金额（必填，必须大于0）
  repeat_rule: RepeatRule       // 重复规则（必填）
  next_due: string              // 下次到期日期（必填，YYYY-MM-DD）
  source_id?: number | null     // 支出账户ID（可选）
  category_id?: number | null   // 分类ID（可选）
  notes?: string                // 备注（可选）
}

/**
 * 更新账单请求接口（对应后端 UpdateBillReq）
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateBillReq {
  name?: string                 // 账单名称
  amount?: string               // 账单金额
  repeat_rule?: RepeatRule      // 重复规则
  next_due?: string             // 下次到期日期
  source_id?: number | null     // 支出账户ID
  category_id?: number | null   // 分类ID
  notes?: string                // 备注
}
