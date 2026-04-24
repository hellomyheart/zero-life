/**
 * 分类相关类型定义
 * 分类用于对交易进行归类（如餐饮、交通、工资等）
 * 分类支持树形结构，即一个分类可以有子分类（如"餐饮"下有"外卖"、"堂食"）
 */

/**
 * 分类信息接口
 * 表示一个分类节点，包含子分类形成树形结构
 */
export interface Category {
  id: string                   // 分类唯一标识
  name: string                 // 分类名称，如"餐饮"、"交通"
  parent_id: string | null     // 父分类ID，顶级分类为null
  children: Category[]         // 子分类列表（递归结构）
  created_at: string           // 创建时间
  updated_at: string           // 最后更新时间
}

/**
 * 创建分类请求接口
 */
export interface CreateCategoryReq {
  name: string                 // 分类名称（必填）
  parent_id?: string | null    // 父分类ID（可选，为null或不传表示顶级分类）
}

/**
 * 更新分类请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateCategoryReq {
  name?: string                // 分类名称
  parent_id?: string | null    // 父分类ID（传null可将其变为顶级分类）
}
