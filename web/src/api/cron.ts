import { get, post } from '@/utils/request'

export interface CronTask {
  id: string
  name: string
  description: string
  schedule: string
}

export interface CronTaskResult {
  task_id: string
  success: boolean
  message: string
  errors?: string[]
}

export function listTasks() {
  return get<CronTask[]>('/cron')
}

export function runTask(id: string) {
  return post<CronTaskResult>(`/cron/${id}/run`)
}
