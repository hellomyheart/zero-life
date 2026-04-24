<script setup lang="ts">
/**
 * 汇率列表页面
 * 功能：
 * - 展示汇率列表，支持分页
 * - 支持按货币对、日期范围筛选
 * - 支持创建、编辑、删除汇率（使用内联对话框）
 * - 支持货币转换计算
 */
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, create, update, remove } from '@/api/exchange-rate'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { ExchangeRate } from '@/types/exchange-rate'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

// 响应式数据
const exchangeRates = ref<ExchangeRate[]>([])
const total = ref(0)
const loading = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

// 表单数据 - 创建/编辑汇率时使用
const form = ref({
  from_currency_id: '',
  to_currency_id: '',
  date: '',
  rate: '',
})

// 筛选条件
const filter = reactive({
  from_currency_id: '',
  to_currency_id: '',
  start_date: '',
  end_date: '',
  page: 1,
  page_size: 20,
})

/**
 * 获取汇率列表
 */
async function fetchExchangeRates() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: filter.page,
      page_size: filter.page_size
    }
    if (filter.from_currency_id) params.from_currency_id = filter.from_currency_id
    if (filter.to_currency_id) params.to_currency_id = filter.to_currency_id
    if (filter.start_date) params.start_date = filter.start_date
    if (filter.end_date) params.end_date = filter.end_date

    const res = await list(params as unknown as typeof filter)
    const data = res as unknown as { items: ExchangeRate[]; total: number }
    exchangeRates.value = data.items || []
    total.value = data.total || 0
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

/**
 * 打开创建汇率对话框
 */
function handleCreate() {
  dialogTitle.value = t('common.create')
  editingId.value = null
  form.value = { from_currency_id: '', to_currency_id: '', date: '', rate: '' }
  dialogVisible.value = true
}

/**
 * 打开编辑汇率对话框 - 将当前行数据填充到表单
 */
function handleEdit(row: ExchangeRate) {
  dialogTitle.value = t('common.edit')
  editingId.value = String(row.id)
  form.value = {
    from_currency_id: String(row.from_currency_id),
    to_currency_id: String(row.to_currency_id),
    date: row.date,
    rate: row.rate,
  }
  dialogVisible.value = true
}

/**
 * 提交表单 - 根据editingId判断是创建还是编辑
 */
async function handleSubmit() {
  if (!form.value.from_currency_id || !form.value.to_currency_id || !form.value.date || !form.value.rate) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      await update(editingId.value, form.value)
    } else {
      await create(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchExchangeRates()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

/**
 * 删除汇率
 */
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('exchangeRate.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchExchangeRates()
  } catch {
    // cancelled or error
  }
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  filter.page = page
  fetchExchangeRates()
}

/**
 * 每页记录数变化
 */
function handleSizeChange(size: number) {
  filter.page_size = size
  filter.page = 1
  fetchExchangeRates()
}

onMounted(() => {
  fetchExchangeRates()
})
</script>

<template>
  <div class="exchange-rate-list-page">
    <el-card>
      <template #header>
        <div class="header">
          <h2>{{ t('exchangeRate.title') }}</h2>
          <el-button type="primary" @click="handleCreate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <el-table :data="exchangeRates" v-loading="loading" stripe>
        <el-table-column prop="from_currency.code" :label="t('exchangeRate.fromCurrency')" />
        <el-table-column prop="to_currency.code" :label="t('exchangeRate.toCurrency')" />
        <el-table-column prop="date" :label="t('exchangeRate.date')" />
        <el-table-column prop="rate" :label="t('exchangeRate.rate')" />
        <el-table-column :label="t('common.actions')" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <Pagination
        :total="total"
        :page="filter.page"
        :page-size="filter.page_size"
        @page-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </el-card>

    <!-- 创建/编辑汇率对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="t('exchangeRate.fromCurrency')">
          <el-input v-model="form.from_currency_id" />
        </el-form-item>
        <el-form-item :label="t('exchangeRate.toCurrency')">
          <el-input v-model="form.to_currency_id" />
        </el-form-item>
        <el-form-item :label="t('exchangeRate.date')">
          <el-input v-model="form.date" type="date" />
        </el-form-item>
        <el-form-item :label="t('exchangeRate.rate')">
          <el-input v-model="form.rate" />
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
.exchange-rate-list-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
}
</style>
