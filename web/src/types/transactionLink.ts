/**
 * 交易关联相关类型定义
 * 对应后端 model/transaction_link.go 和 DTO
 */

/**
 * 交易关联类型枚举
 */
export type TransactionLinkType = 'rolled_back' | 'reconciled' | 'linked';

/**
 * 交易关联数据类型
 * 用于记录两笔交易之间的关联关系
 */
export interface TransactionLink {
  /** 关联记录唯一标识 */
  id: number;
  /** 源交易 ID */
  transaction_id: number;
  /** 关联类型 */
  link_type: TransactionLinkType;
  /** 目标交易日志 ID */
  linked_journal_id: number;
  /** 创建时间 */
  created_at: string;
}

/**
 * 创建交易关联请求参数
 */
export interface CreateTransactionLinkRequest {
  /** 源交易 ID */
  transaction_id: number;
  /** 关联类型 */
  link_type: TransactionLinkType;
  /** 目标交易日志 ID */
  linked_journal_id: number;
}

export interface UpdateTransactionLinkRequest {
  /** 源交易 ID */
  transaction_id?: number;
  /** 关联类型 */
  link_type?: TransactionLinkType;
  /** 目标交易日志 ID */
  linked_journal_id?: number;
}

/**
 * 交易关联列表查询参数
 */
export interface TransactionLinkListParams {
  /** 按交易 ID 筛选 */
  transaction_id?: number;
  /** 页码 */
  page?: number;
  /** 每页条数 */
  page_size?: number;
}
