/**
 * 标签状态管理（Pinia Store）
 *
 * 组件功能：
 * - 缓存标签列表数据，避免重复请求
 *
 * 数据流：
 * - tags: 从后端API获取所有标签数据，供交易表单等场景的多选下拉使用
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Tag } from '@/types/tag'
import { list as listTags } from '@/api/tag'

export const useTagStore = defineStore('tag', () => {
  /** 标签列表数据，包含名称、颜色和关联交易数量 */
  const tags = ref<Tag[]>([])

  /**
   * 获取标签列表
   * 从后端API获取所有标签数据并缓存
   */
  async function fetchTags() {
    tags.value = await listTags() as unknown as Tag[]
  }

  return {
    tags,
    fetchTags,
  }
})
