/**
 * 存钱罐相关类型定义
 * 存钱罐用于设定储蓄目标，跟踪存款进度（如旅行基金、应急储备金等）
 */

/**
 * 存钱罐信息接口
 * 表示一个储蓄目标的完整信息
 */
export interface PiggyBank {
  id: number                    // 存钱罐唯一标识
  name: string                  // 存钱罐名称，如"旅行基金"
  account_id: number            // 关联的资产账户ID（存取款从该账户操作）
  target_amount: string        // 目标金额
  current_amount: string       // 当前已存金额
  start_date: string           // 开始存款日期
  target_date: string          // 目标完成日期
  order: number                // 排序序号
  notes: string                // 备注说明
  percentage: number           // 完成百分比（0-100）
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 存钱罐事件接口
 * 记录存钱罐的每次存入或取出操作
 */
export interface PiggyEvent {
  id: number                    // 事件唯一标识
  piggy_bank_id: number         // 关联的存钱罐ID
  amount: string               // 操作金额
  type: string                 // 操作类型，如"add"（存入）、"remove"（取出）
  note: string                 // 操作备注
  created_at: string           // 操作时间
}

/**
 * 创建存钱罐请求接口
 */
export interface CreatePiggyBankReq {
  name: string                 // 存钱罐名称（必填）
  account_id: number            // 关联账户ID（必填）
  target_amount: string        // 目标金额（必填）
  target_date?: string         // 目标完成日期（可选）
  notes?: string               // 备注说明（可选）
}

/**
 * 更新存钱罐请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdatePiggyBankReq {
  name?: string                // 存钱罐名称
  target_amount?: string       // 目标金额
  target_date?: string         // 目标完成日期
  notes?: string               // 备注说明
}

/**
 * 存取款请求接口
 * 存入操作使用此接口
 */
export interface AddAmountReq {
  amount: string               // 存入金额（必填）
  note?: string                // 操作备注（可选）
}

/**
 * 取款请求接口
 * 取出操作使用此接口
 */
export interface RemoveAmountReq {
  amount: string               // 取出金额（必填）
  note?: string                // 操作备注（可选）
}

/**
 * 存钱罐列表查询请求接口
 */
export interface PiggyBankListReq {
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}
