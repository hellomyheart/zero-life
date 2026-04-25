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
  id: number                    // 分类唯一标识（后端 uint64）
  name: string                  // 分类名称，如"餐饮"、"交通"
  parent_id: number | null      // 父分类ID，顶级分类为null
  icon: string                  // 分类图标
  notes: string                 // 备注
  sort_order: number            // 排序序号
  children: Category[]          // 子分类列表（递归结构）
  created_at: string            // 创建时间
  updated_at: string            // 最后更新时间
}

/**
 * 创建分类请求接口
 */
export interface CreateCategoryReq {
  name: string                  // 分类名称（必填）
  parent_id?: number | null     // 父分类ID（可选，为null或不传表示顶级分类）
  icon?: string                 // 分类图标（可选）
  notes?: string                // 备注（可选）
}

export interface UpdateCategoryReq {
  name?: string                 // 分类名称
  icon?: string                 // 分类图标
  notes?: string                // 备注
  sort_order?: number           // 排序序号
}
