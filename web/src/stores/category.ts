/**
 * 分类状态管理（Pinia Store）
 *
 * 组件功能：
 * - 缓存分类列表数据（树形结构），避免重复请求
 *
 * 数据流：
 * - categories: 从后端API获取，供交易表单、预算表单等场景的下拉选择使用
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category } from '@/types/category'
import { list as listCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  /** 分类列表数据（树形结构，包含子分类） */
  const categories = ref<Category[]>([])

  /**
   * 获取分类列表
   * 从后端API获取完整的分类树形结构并缓存
   */
  async function fetchCategories() {
    try {
      categories.value = await listCategories() as unknown as Category[]
    } catch {
      categories.value = []
    }
  }

  return {
    categories,
    fetchCategories,
  }
})
