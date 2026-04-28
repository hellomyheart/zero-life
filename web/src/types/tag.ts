/**
 * 标签相关类型定义
 * 标签用于对交易进行标记和分组（如"出差"、"日常"等）
 * 与分类不同，标签是一种更灵活的横向归类方式，一笔交易可以有多个标签
 * 支持最多5级树形结构
 */

/**
 * 标签信息接口
 * 表示系统中的一个标签，支持树形结构
 */
export interface Tag {
  id: number                       // 标签唯一标识（后端 uint64）
  name: string                     // 标签名称，如"出差"、"日常"
  color: string                    // 标签颜色（十六进制色值），如"#FF5733"
  parent_id: number | null         // 父标签ID，为null时表示顶级标签
  transaction_count: number        // 关联的交易数量
  children?: Tag[]                 // 子标签列表（树形结构时使用）
  created_at: string               // 创建时间
  updated_at: string               // 最后更新时间
}

/**
 * 创建标签请求接口
 */
export interface CreateTagReq {
  name: string                     // 标签名称（必填）
  color: string                    // 标签颜色（必填，十六进制色值）
  parent_id?: number | null        // 父标签ID，为null时表示顶级标签
}

/**
 * 更新标签请求接口
 * 所有字段均为可选，只传需要修改的字段
 */
export interface UpdateTagReq {
  name?: string                    // 标签名称
  color?: string                   // 标签颜色
  parent_id?: number | null        // 父标签ID
}
