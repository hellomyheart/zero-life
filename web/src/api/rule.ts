// 规则API接口 - 规则增删改查、启用/禁用和执行
import { get, post, put, del } from '@/utils/request'
import type { Rule, CreateRuleReq, UpdateRuleReq, ExecuteRuleReq } from '@/types/rule'

export function list() {
  return get<Rule[]>('/rules')
}

export function getRule(id: string) {
  return get<Rule>(`/rules/${id}`)
}

export function create(data: CreateRuleReq) {
  return post<Rule>('/rules', data)
}

export function update(id: string, data: UpdateRuleReq) {
  return put<Rule>(`/rules/${id}`, data)
}

export function remove(id: string) {
  return del<void>(`/rules/${id}`)
}

export function toggleStatus(id: string) {
  return put<Rule>(`/rules/${id}/toggle`)
}

export function execute(id: string, data: ExecuteRuleReq) {
  return post<void>(`/rules/${id}/execute`, data)
}
