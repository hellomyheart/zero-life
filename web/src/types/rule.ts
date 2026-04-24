/**
 * 规则相关类型定义
 * 规则用于自动化交易处理，由"条件"和"动作"组成
 * 当交易满足所有条件时，系统自动执行指定的动作（如自动分类、添加标签等）
 */

/**
 * 规则条件接口
 * 定义触发规则需要满足的条件
 */
export interface RuleCondition {
  field: string                // 要匹配的字段，如"description"、"amount"、"source_account_id"
  operator: string             // 比较运算符，如"equals"、"contains"、"gt"、"lt"
  value: string                // 匹配值
}

/**
 * 规则动作接口
 * 定义条件满足时要执行的操作
 */
export interface RuleAction {
  type: string                 // 动作类型，如"set_category"、"add_tag"、"set_description"
  field: string                // 要修改的字段
  value: string                // 要设置的值
}

/**
 * 规则信息接口
 * 表示一个完整的自动化规则
 */
export interface Rule {
  id: string                   // 规则唯一标识
  name: string                 // 规则名称，如"自动分类餐饮消费"
  conditions: RuleCondition[]  // 触发条件列表（所有条件需同时满足）
  actions: RuleAction[]        // 执行动作列表（条件满足时依次执行）
  is_enabled: boolean          // 是否启用
  priority: number             // 优先级（数字越小优先级越高）
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建规则时的条件请求接口
 */
export interface RuleConditionReq {
  field: string                // 要匹配的字段（必填）
  operator: string             // 比较运算符（必填）
  value: string                // 匹配值（必填）
}

/**
 * 创建规则时的动作请求接口
 */
export interface RuleActionReq {
  type: string                 // 动作类型（必填）
  field: string                // 要修改的字段（必填）
  value: string                // 要设置的值（必填）
}

/**
 * 创建规则请求接口
 */
export interface CreateRuleReq {
  name: string                 // 规则名称（必填）
  conditions: RuleConditionReq[] // 触发条件列表（必填）
  actions: RuleActionReq[]     // 执行动作列表（必填）
  is_enabled?: boolean         // 是否启用（可选，默认true）
  priority?: number            // 优先级（可选，默认0）
}

/**
 * 更新规则请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateRuleReq {
  name?: string                // 规则名称
  conditions?: RuleConditionReq[] // 触发条件列表
  actions?: RuleActionReq[]    // 执行动作列表
  is_enabled?: boolean         // 是否启用
  priority?: number            // 优先级
}

/**
 * 执行规则请求接口
 * 手动对指定交易应用规则
 */
export interface ExecuteRuleReq {
  transaction_ids: string[]    // 要应用规则的交易ID列表
}
