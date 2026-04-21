export interface RuleCondition {
  field: string
  operator: string
  value: string
}

export interface RuleAction {
  type: string
  field: string
  value: string
}

export interface Rule {
  id: string
  name: string
  conditions: RuleCondition[]
  actions: RuleAction[]
  is_enabled: boolean
  priority: number
  created_at: string
  updated_at: string
}

export interface RuleConditionReq {
  field: string
  operator: string
  value: string
}

export interface RuleActionReq {
  type: string
  field: string
  value: string
}

export interface CreateRuleReq {
  name: string
  conditions: RuleConditionReq[]
  actions: RuleActionReq[]
  is_enabled?: boolean
  priority?: number
}

export interface UpdateRuleReq {
  name?: string
  conditions?: RuleConditionReq[]
  actions?: RuleActionReq[]
  is_enabled?: boolean
  priority?: number
}

export interface ExecuteRuleReq {
  transaction_ids: string[]
}
