// 分类状态管理 - 管理分类列表
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category } from '@/types/category'
import { list as listCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([]) // 分类列表

  // fetchCategories 获取分类列表
  async function fetchCategories() {
    categories.value = await listCategories() as unknown as Category[]
  }

  return {
    categories,
    fetchCategories,
  }
})
