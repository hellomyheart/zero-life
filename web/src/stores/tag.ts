// 标签状态管理 - 管理标签列表
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Tag } from '@/types/tag'
import { list as listTags } from '@/api/tag'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([]) // 标签列表

  // fetchTags 获取标签列表
  async function fetchTags() {
    tags.value = await listTags() as unknown as Tag[]
  }

  return {
    tags,
    fetchTags,
  }
})
