/**
 * 通用类型定义
 * 提供API响应、分页结果和分页参数等基础类型，被其他模块复用
 */

/**
 * API统一响应接口
 * 后端所有API返回的统一格式，包含状态码、消息和数据
 * @template T - data字段的类型
 */
export interface ApiResponse<T> {
  code: number                 // 业务状态码（0通常表示成功）
  message: string              // 响应消息（成功或错误描述）
  data: T                      // 响应数据，具体类型由泛型参数决定
}

/**
 * 分页结果接口
 * 后端返回的分页数据统一格式
 * @template T - 列表项的类型
 */
export interface PageResult<T> {
  items: T[]                   // 当前页的数据列表
  total: number                // 符合条件的总记录数
  page: number                 // 当前页码
  page_size: number            // 每页数量
  total_pages: number          // 总页数
}

/**
 * 分页参数接口
 * 请求分页数据时使用的通用分页参数
 */
export interface PageParams {
  page?: number                // 页码（从1开始）
  page_size?: number           // 每页数量
}
