import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category } from '@/types/category'
import { list as listCategories } from '@/api/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])

  async function fetchCategories() {
    categories.value = await listCategories() as unknown as Category[]
  }

  return {
    categories,
    fetchCategories,
  }
})
