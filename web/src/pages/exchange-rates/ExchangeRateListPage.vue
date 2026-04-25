<!-- 汇率管理页面 - 展示汇率列表，支持筛选、新增、编辑和删除汇率记录 -->
<script setup lang="ts">
// 导入 Vue 响应式和生命周期 API
import { ref, reactive, computed, onMounted } from 'vue'
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n'
// 导入汇率相关 API（列表、创建、更新、删除）
import { list, create, update, remove } from '@/api/exchange-rate'
// 导入 Element Plus 消息提示和确认框组件
import { ElMessage, ElMessageBox } from 'element-plus'
// 导入货币状态管理
import { useCurrencyStore } from '@/stores/currency'
// 导入汇率相关类型定义
import type { ExchangeRate, ExchangeRateListParams, CreateExchangeRateReq, UpdateExchangeRateReq } from '@/types/exchange-rate'
// 导入货币类型定义
import type { Currency } from '@/types/currency'
// 导入分页组件
import Pagination from '@/components/common/Pagination.vue'

// 国际化翻译函数
const { t } = useI18n()
// 货币状态管理实例
const currencyStore = useCurrencyStore()

// 已启用的货币列表（计算属性，只包含 is_enabled 为 true 的货币）
const enabledCurrencies = computed(() => currencyStore.currencies.filter((c: Currency) => c.is_enabled))

// 汇率列表数据
const exchangeRates = ref<ExchangeRate[]>([])
// 数据总条数，用于分页
const total = ref(0)
// 列表加载状态
const loading = ref(false)

// 新增/编辑对话框是否可见
const dialogVisible = ref(false)
// 对话框标题（新增 / 编辑）
const dialogTitle = ref('')
// 当前正在编辑的汇率 ID，为 null 时表示新增模式
const editingId = ref<string | null>(null)

// 新增/编辑表单数据
const form = ref({
  from_currency_id: '' as string | number,
  to_currency_id: '' as string | number,
  date: '',
  rate: '',
})

// 筛选条件
const filter = reactive({
  from_currency_id: '' as string | number,
  to_currency_id: '' as string | number,
  start_date: '',
  end_date: '',
  page: 1,
  page_size: 20,
})

// 获取汇率列表数据
async function fetchExchangeRates() {
  loading.value = true
  try {
    // 构建查询参数，空值不传给后端
    const params: ExchangeRateListParams = {
      page: filter.page,
      page_size: filter.page_size,
      from_currency_id: filter.from_currency_id ? Number(filter.from_currency_id) : undefined,
      to_currency_id: filter.to_currency_id ? Number(filter.to_currency_id) : undefined,
      start_date: filter.start_date || undefined,
      end_date: filter.end_date || undefined,
    }

    const res = await list(params)
    const data = res as unknown as { items: ExchangeRate[]; total: number }
    exchangeRates.value = data.items || []
    total.value = data.total || 0
  } catch {
    ElMessage.error(t('common.fetchError') || 'Failed to load data')
  } finally {
    loading.value = false
  }
}

// 点击筛选按钮，重置页码后查询
function handleFilter() {
  filter.page = 1
  fetchExchangeRates()
}

// 重置筛选条件并重新查询
function resetFilter() {
  filter.from_currency_id = ''
  filter.to_currency_id = ''
  filter.start_date = ''
  filter.end_date = ''
  filter.page = 1
  fetchExchangeRates()
}

// 点击新增按钮，打开新增对话框
function handleCreate() {
  dialogTitle.value = t('common.create')
  editingId.value = null
  form.value = { from_currency_id: '', to_currency_id: '', date: '', rate: '' } as typeof form.value
  dialogVisible.value = true
}

// 点击编辑按钮，填充表单数据并打开编辑对话框
function handleEdit(row: ExchangeRate) {
  dialogTitle.value = t('common.edit')
  editingId.value = String(row.id)
  form.value = {
    from_currency_id: row.from_currency_id,
    to_currency_id: row.to_currency_id,
    date: row.date,
    rate: row.rate,
  }
  dialogVisible.value = true
}

// 提交新增/编辑表单
async function handleSubmit() {
  // 基本字段校验
  if (!form.value.from_currency_id || !form.value.to_currency_id || !form.value.date || !form.value.rate) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    if (editingId.value) {
      // 编辑模式：调用更新 API
      await update(String(editingId.value), {
        from_currency_id: Number(form.value.from_currency_id),
        to_currency_id: Number(form.value.to_currency_id),
        date: form.value.date,
        rate: form.value.rate,
      } as UpdateExchangeRateReq)
    } else {
      // 新增模式：调用创建 API
      await create({
        from_currency_id: Number(form.value.from_currency_id),
        to_currency_id: Number(form.value.to_currency_id),
        date: form.value.date,
        rate: form.value.rate,
      } as CreateExchangeRateReq)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    // 提交成功后刷新列表
    await fetchExchangeRates()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

// 删除汇率记录（带确认弹窗）
async function handleDelete(id: string) {
  try {
    // 弹出确认框，用户确认后才执行删除
    await ElMessageBox.confirm(t('exchangeRate.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchExchangeRates()
  } catch (err) {
    // 用户点击取消时不提示错误
    if ((err as string) !== 'cancel') ElMessage.error(t('common.fetchError') || 'Failed to load data')
  }
}

// 分页页码变化
function handlePageChange(page: number) {
  filter.page = page
  fetchExchangeRates()
}

// 每页条数变化，重置到第一页
function handleSizeChange(size: number) {
  filter.page_size = size
  filter.page = 1
  fetchExchangeRates()
}

// 页面挂载时先加载货币列表，再加载汇率数据
onMounted(async () => {
  await currencyStore.fetchCurrencies()
  fetchExchangeRates()
})
</script>

<template>
  <div class="exchange-rate-list-page">
    <el-card>
      <!-- 页面标题和新增按钮 -->
      <template #header>
        <div class="header">
          <h2>{{ t('exchangeRate.title') }}</h2>
          <el-button type="primary" @click="handleCreate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <!-- 筛选条件表单 -->
      <el-form :inline="true" :model="filter" class="filter-form">
        <el-form-item :label="t('exchangeRate.fromCurrency')">
          <el-select v-model="filter.from_currency_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('exchangeRate.toCurrency')">
          <el-select v-model="filter.to_currency_id" :placeholder="t('common.selectPlaceholder')" filterable clearable>
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.startDate')">
          <el-date-picker v-model="filter.start_date" type="date" value-format="YYYY-MM-DD" clearable />
        </el-form-item>
        <el-form-item :label="t('common.endDate')">
          <el-date-picker v-model="filter.end_date" type="date" value-format="YYYY-MM-DD" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">{{ t('common.search') }}</el-button>
          <el-button @click="resetFilter">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 汇率数据表格 -->
      <el-table :data="exchangeRates" v-loading="loading" stripe>
        <el-table-column prop="from_currency.code" :label="t('exchangeRate.fromCurrency')" />
        <el-table-column prop="to_currency.code" :label="t('exchangeRate.toCurrency')" />
        <el-table-column prop="date" :label="t('exchangeRate.date')" />
        <el-table-column prop="rate" :label="t('exchangeRate.rate')" />
        <!-- 操作列：编辑和删除 -->
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

      <!-- 分页组件 -->
      <Pagination
        :total="total"
        :page="filter.page"
        :page-size="filter.page_size"
        @update:page="handlePageChange"
        @update:page-size="handleSizeChange"
      />
    </el-card>

    <!-- 新增/编辑汇率对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="t('exchangeRate.fromCurrency')">
          <el-select v-model="form.from_currency_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('exchangeRate.toCurrency')">
          <el-select v-model="form.to_currency_id" :placeholder="t('common.selectPlaceholder')" filterable clearable style="width: 100%">
            <el-option v-for="c in enabledCurrencies" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('exchangeRate.date')">
          <el-date-picker v-model="form.date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="t('exchangeRate.rate')">
          <el-input v-model="form.rate" />
        </el-form-item>
      </el-form>
      <!-- 对话框底部按钮 -->
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

.filter-form {
  margin-bottom: 16px;
}
</style>
