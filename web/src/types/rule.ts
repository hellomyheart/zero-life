/**
 * 规则相关类型定义
 * 字段名与后端 JSON tag 完全对应
 */

/**
 * 规则条件接口（对应后端 RuleConditionResp）
 */
export interface RuleCondition {
  id: number
  field: string
  operator: string
  value: string
}

/**
 * 规则动作接口（对应后端 RuleActionResp）
 * 后端仅有 type 和 value，无 field 字段
 */
export interface RuleAction {
  id: number
  type: string
  value: string
}

/**
 * 规则信息接口（对应后端 RuleResp）
 */
export interface Rule {
  id: number
  name: string
  priority: number
  is_enabled: boolean
  logic_type: string
  trigger: string
  conditions: RuleCondition[]
  actions: RuleAction[]
  created_at: string
  updated_at: string
}

/**
 * 创建规则时的条件请求接口（对应后端 RuleConditionReq）
 */
export interface RuleConditionReq {
  field: string
  operator: string
  value: string
}

/**
 * 创建规则时的动作请求接口（对应后端 RuleActionReq）
 * 后端仅有 type 和 value，无 field 字段
 */
export interface RuleActionReq {
  type: string
  value: string
}

/**
 * 创建规则请求接口（对应后端 CreateRuleReq）
 * logic_type 和 trigger 为必填字段
 */
export interface CreateRuleReq {
  name: string
  priority?: number
  is_enabled?: boolean
  logic_type: string
  trigger: string
  conditions: RuleConditionReq[]
  actions: RuleActionReq[]
}

/**
 * 更新规则请求接口（对应后端 UpdateRuleReq）
 */
export interface UpdateRuleReq {
  name?: string
  priority?: number | null
  is_enabled?: boolean | null
  logic_type?: string
  trigger?: string
  conditions?: RuleConditionReq[]
  actions?: RuleActionReq[]
}

/**
 * 执行规则请求接口（对应后端 ExecuteRuleReq）
 * 后端使用 start_date/end_date，非 transaction_ids
 */
export interface ExecuteRuleReq {
  start_date: string
  end_date: string
}

/**
 * 规则执行结果响应接口（对应后端 RuleExecuteResultResp）
 */
export interface RuleExecuteResultResp {
  matched_count: number
  success_count: number
  fail_count: number
  errors?: string[]
}
