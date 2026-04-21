<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const authStore = useAuthStore()
const route = useRoute()

const isAuthPage = computed(() => {
  return ['/login', '/register', '/reset-password'].includes(route.path)
})
</script>

<template>
  <AppLayout v-if="authStore.isAuthenticated && !isAuthPage" />
  <router-view v-else />
</template>
