/**
 * 对象分组相关类型定义
 * 对应后端 model/object_group.go 和 DTO
 */

/**
 * 对象分组数据类型
 * 用于将相关实体（交易、账单、预算等）分组管理
 */
export interface ObjectGroup {
  /** 分组唯一标识 */
  id: number;
  /** 分组名称 */
  name: string;
  /** 关联对象类型（多态） */
  groupable_type: 'transaction' | 'bill' | 'budget' | 'piggy_bank';
  /** 关联对象 ID */
  groupable_id: number;
  /** 创建时间 */
  created_at: string;
  /** 更新时间 */
  updated_at: string;
}

/**
 * 创建对象分组请求参数
 */
export interface CreateObjectGroupRequest {
  /** 分组名称 */
  name: string;
  /** 关联对象类型 */
  groupable_type: 'transaction' | 'bill' | 'budget' | 'piggy_bank';
  /** 关联对象 ID */
  groupable_id: number;
}

/**
 * 更新对象分组请求参数
 */
export interface UpdateObjectGroupRequest {
  /** 分组名称 */
  name: string;
}

/**
 * 对象分组列表查询参数
 */
export interface ObjectGroupListParams {
  /** 按关联对象类型筛选 */
  groupable_type?: string;
  /** 页码 */
  page?: number;
  /** 每页条数 */
  page_size?: number;
}
