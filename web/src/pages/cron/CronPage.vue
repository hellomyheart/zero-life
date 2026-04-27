<script setup lang="ts">
import { ref } from 'vue'
import { listTasks, runTask, type CronTask, type CronTaskResult } from '@/api/cron'
import { ElMessage, ElMessageBox } from 'element-plus'

const tasks = ref<CronTask[]>([])
const loading = ref(false)
const running = ref<string | null>(null)
const results = ref<Map<string, CronTaskResult>>(new Map())

async function fetchTasks() {
  loading.value = true
  try {
    const res = await listTasks()
    tasks.value = (res as unknown as CronTask[]) || []
  } catch {
    ElMessage.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

async function handleRun(task: CronTask) {
  try {
    await ElMessageBox.confirm(`确定要手动执行「${task.name}」吗？`, '确认执行', { type: 'warning' })
  } catch {
    return
  }

  running.value = task.id
  try {
    const res = await runTask(task.id)
    const result = res as unknown as CronTaskResult
    results.value.set(task.id, result)
    if (result?.success) {
      ElMessage.success(`${task.name} 执行成功`)
    } else {
      ElMessage.error(`${task.name} 执行失败: ${result?.message}`)
    }
  } catch {
    ElMessage.error(`${task.name} 执行异常`)
  } finally {
    running.value = null
  }
}

fetchTasks()
</script>

<template>
  <div class="cron-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>定时任务管理</h2>
          <el-button :loading="loading" @click="fetchTasks">刷新</el-button>
        </div>
      </template>

      <el-table :data="tasks" v-loading="loading" stripe>
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="250" />
        <el-table-column prop="schedule" label="执行周期" min-width="120" />
        <el-table-column label="最近执行结果" min-width="150">
          <template #default="{ row }">
            <template v-if="results.has(row.id)">
              <el-tag :type="results.get(row.id)!.success ? 'success' : 'danger'" size="small">
                {{ results.get(row.id)!.message }}
              </el-tag>
            </template>
            <template v-else>
              <span style="color: var(--app-text-secondary)">未执行</span>
            </template>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" :loading="running === row.id" @click="handleRun(row)">
              手动执行
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.cron-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>