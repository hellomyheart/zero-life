/**
 * 规则相关API接口
 * 提供规则的增删改查、启用/禁用切换和手动执行功能
 * 规则由"条件"和"动作"组成，当交易满足条件时自动执行动作（如自动分类、添加标签等），
 * 实现交易的自动化处理
 */
import { get, post, put, del } from '@/utils/request'
import type { Rule, CreateRuleReq, UpdateRuleReq, ExecuteRuleReq } from '@/types/rule'

/**
 * 获取所有规则列表
 * @returns 规则数组
 * @endpoint GET /rules
 */
export function list() {
  return get<Rule[]>('/rules')
}

/**
 * 获取单个规则详情
 * @param id - 规则ID
 * @returns 规则详细信息（含条件和动作列表）
 * @endpoint GET /rules/:id
 */
export function getRule(id: number) {
  return get<Rule>(`/rules/${id}`)
}

/**
 * 创建新规则
 * @param data - 创建请求参数（名称、条件列表、动作列表、是否启用、优先级）
 * @returns 创建成功的规则信息
 * @endpoint POST /rules
 */
export function create(data: CreateRuleReq) {
  return post<Rule>('/rules', data)
}

/**
 * 更新规则信息
 * @param id - 规则ID
 * @param data - 更新内容（名称、条件、动作、启用状态、优先级，均为可选）
 * @returns 更新后的规则信息
 * @endpoint PUT /rules/:id
 */
export function update(id: number, data: UpdateRuleReq) {
  return put<Rule>(`/rules/${id}`, data)
}

/**
 * 删除规则
 * @param id - 规则ID
 * @endpoint DELETE /rules/:id
 */
export function remove(id: number) {
  return del<void>(`/rules/${id}`)
}

/**
 * 切换规则的启用/禁用状态
 * 如果当前启用则变为禁用，反之亦然
 * @param id - 规则ID
 * @returns 更新后的规则信息
 * @endpoint PUT /rules/:id/toggle
 */
export function toggleStatus(id: number) {
  return put<Rule>(`/rules/${id}/toggle`)
}

/**
 * 手动执行规则
 * 对指定的交易列表应用该规则的条件和动作
 * @param id - 规则ID
 * @param data - 执行请求参数（要应用规则的交易ID列表）
 * @endpoint POST /rules/:id/execute
 */
export function execute(id: number, data: ExecuteRuleReq) {
  return post<void>(`/rules/${id}/execute`, data)
}
