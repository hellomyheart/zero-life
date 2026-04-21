import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Tag } from '@/types/tag'
import { list as listTags } from '@/api/tag'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])

  async function fetchTags() {
    tags.value = await listTags() as unknown as Tag[]
  }

  return {
    tags,
    fetchTags,
  }
})
