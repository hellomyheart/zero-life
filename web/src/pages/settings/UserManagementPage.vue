<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, update, remove, lock, unlock } from '@/api/user'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { User, UpdateUserReq } from '@/types/user'

const { t } = useI18n()

const users = ref<User[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref<string | null>(null)
const editForm = ref<UpdateUserReq>({})

async function fetchUsers() {
  loading.value = true
  try {
    const res = await list({}) as unknown as { items: User[] }
    users.value = res.items || []
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleEdit(row: User) {
  editingId.value = row.id
  editForm.value = { nickname: row.nickname, role: row.role, language: row.language }
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!editingId.value) return
  try {
    await update(editingId.value, editForm.value)
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchUsers()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('user.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchUsers()
  } catch {
    // cancelled or error
  }
}

async function handleLock(id: string) {
  try {
    await lock(id)
    ElMessage.success(t('common.success'))
    await fetchUsers()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleUnlock(id: string) {
  try {
    await unlock(id)
    ElMessage.success(t('common.success'))
    await fetchUsers()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div class="user-management-page">
    <h2>{{ t('user.title') }}</h2>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="email" :label="t('auth.email')" />
      <el-table-column prop="nickname" :label="t('auth.nickname')" width="120" />
      <el-table-column prop="role" :label="t('user.role')" width="100" />
      <el-table-column prop="language" :label="t('profile.language')" width="100" />
      <el-table-column :label="t('user.locked')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.locked ? 'danger' : 'success'" size="small">{{ row.locked ? t('common.yes') : t('common.no') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('attachment.createdAt')" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="260" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link :type="row.locked ? 'success' : 'warning'" @click="row.locked ? handleUnlock(row.id) : handleLock(row.id)">{{ row.locked ? t('user.unlock') : t('user.lock') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="t('common.edit')" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item :label="t('auth.nickname')">
          <el-input v-model="editForm.nickname" />
        </el-form-item>
        <el-form-item :label="t('user.role')">
          <el-select v-model="editForm.role">
            <el-option label="User" value="user" />
            <el-option label="Admin" value="admin" />
            <el-option label="Owner" value="owner" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('profile.language')">
          <el-select v-model="editForm.language">
            <el-option label="中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
h2 {
  margin: 0 0 16px 0;
}
</style>
