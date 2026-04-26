<script setup lang="ts">
/**
 * 规则列表页面
 * 字段名与后端 JSON tag 完全对应：
 * - RuleAction 仅有 type 和 value（无 field）
 * - CreateRuleReq 必须包含 logic_type 和 trigger
 * - ExecuteRuleReq 使用 start_date/end_date（非 transaction_ids）
 */
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, remove, toggleStatus, execute } from '@/api/rule'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Rule, CreateRuleReq, RuleConditionReq, RuleActionReq } from '@/types/rule'
import Pagination from '@/components/common/Pagination.vue'

const { t } = useI18n()

const rules = ref<Rule[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<number | null>(null)
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

/** 逻辑类型选项 */
const logicTypeOptions = [
  { value: 'and', label: t('rule.logicAnd') },
  { value: 'or', label: t('rule.logicOr') },
]

/** 触发时机选项 */
const triggerOptions = [
  { value: 'on_create', label: t('rule.triggerOnCreate') },
  { value: 'on_update', label: t('rule.triggerOnUpdate') },
]

/** 条件字段选项（对应后端 oneof 验证） */
const conditionFieldOptions = [
  { value: 'description', label: t('rule.fieldDescription') },
  { value: 'amount', label: t('rule.fieldAmount') },
  { value: 'source_account', label: t('rule.fieldSourceAccount') },
  { value: 'destination_account', label: t('rule.fieldDestinationAccount') },
  { value: 'category', label: t('rule.fieldCategory') },
  { value: 'tag', label: t('rule.fieldTag') },
  { value: 'transaction_type', label: t('rule.fieldTransactionType') },
  { value: 'budget', label: t('rule.fieldBudget') },
  { value: 'bill', label: t('rule.fieldBill') },
  { value: 'notes', label: t('rule.fieldNotes') },
  { value: 'date_after', label: t('rule.fieldDateAfter') },
  { value: 'date_before', label: t('rule.fieldDateBefore') },
]

/** 条件运算符选项（对应后端 oneof 验证） */
const conditionOperatorOptions = [
  { value: 'contains', label: t('rule.opContains') },
  { value: 'equals', label: t('rule.opEquals') },
  { value: 'starts_with', label: t('rule.opStartsWith') },
  { value: 'ends_with', label: t('rule.opEndsWith') },
  { value: 'not_contains', label: t('rule.opNotContains') },
  { value: 'not_equals', label: t('rule.opNotEquals') },
  { value: 'less', label: t('rule.opLess') },
  { value: 'more', label: t('rule.opMore') },
  { value: 'is_empty', label: t('rule.opIsEmpty') },
  { value: 'is_not_empty', label: t('rule.opIsNotEmpty') },
]

/** 动作类型选项（对应后端 oneof 验证） */
const actionTypeOptions = [
  { value: 'set_category', label: t('rule.actionSetCategory') },
  { value: 'add_tag', label: t('rule.actionAddTag') },
  { value: 'set_notes', label: t('rule.actionSetNotes') },
  { value: 'set_budget', label: t('rule.actionSetBudget') },
  { value: 'remove_tag', label: t('rule.actionRemoveTag') },
  { value: 'set_description', label: t('rule.actionSetDescription') },
  { value: 'clear_category', label: t('rule.actionClearCategory') },
  { value: 'clear_budget', label: t('rule.actionClearBudget') },
  { value: 'clear_notes', label: t('rule.actionClearNotes') },
  { value: 'append_notes', label: t('rule.actionAppendNotes') },
  { value: 'prepend_notes', label: t('rule.actionPrependNotes') },
]

const form = ref<CreateRuleReq>({
  name: '',
  logic_type: 'and',
  trigger: 'on_create',
  conditions: [{ field: '', operator: '', value: '' }],
  actions: [{ type: '', value: '' }],
  is_enabled: true,
  priority: 0,
})

async function fetchRules() {
  loading.value = true
  try {
    const res = await list()
    const data = res as unknown as Rule[]
    rules.value = data
    pagination.total = data.length
  } catch {
    ElMessage.error(t('common.fetchError'))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.page = page
  fetchRules()
}

function handleSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  fetchRules()
}

function handleCreate() {
  dialogTitle.value = t('rule.create')
  editingId.value = null
  form.value = {
    name: '',
    logic_type: 'and',
    trigger: 'on_create',
    conditions: [{ field: '', operator: '', value: '' }],
    actions: [{ type: '', value: '' }],
    is_enabled: true,
    priority: 0,
  }
  dialogVisible.value = true
}

function handleEdit(rule: Rule) {
  dialogTitle.value = t('rule.edit')
  editingId.value = rule.id
  form.value = {
    name: rule.name,
    logic_type: rule.logic_type,
    trigger: rule.trigger,
    conditions: rule.conditions.map((c): RuleConditionReq => ({ field: c.field, operator: c.operator, value: c.value })),
    actions: rule.actions.map((a): RuleActionReq => ({ type: a.type, value: a.value })),
    is_enabled: rule.is_enabled,
    priority: rule.priority,
  }
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t('rule.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchRules()
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

async function handleToggleStatus(rule: Rule) {
  try {
    await toggleStatus(rule.id)
    ElMessage.success(t('common.success'))
    await fetchRules()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

async function handleExecute(rule: Rule) {
  try {
    await ElMessageBox.confirm(t('rule.executeConfirm', { name: rule.name }), t('common.confirm'), { type: 'info' })
    await execute(rule.id, { start_date: '', end_date: '' })
    ElMessage.success(t('common.success'))
  } catch (err) {
    if ((err as string) !== 'cancel') ElMessage.error(t('common.failed'))
  }
}

function addCondition() {
  form.value.conditions.push({ field: '', operator: '', value: '' })
}

function removeCondition(index: number) {
  form.value.conditions.splice(index, 1)
}

function addAction() {
  form.value.actions.push({ type: '', value: '' })
}

function removeAction(index: number) {
  form.value.actions.splice(index, 1)
}

async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning(t('common.required'))
    return
  }
  try {
    const { create: createRule, update: updateRule } = await import('@/api/rule')
    if (editingId.value) {
      await updateRule(editingId.value, form.value)
    } else {
      await createRule(form.value)
    }
    ElMessage.success(t('common.success'))
    dialogVisible.value = false
    await fetchRules()
  } catch (err) {
    ElMessage.error((err as Error).message || t('common.failed'))
  }
}

function getConditionSummary(rule: Rule): string {
  return rule.conditions.map((c) => `${c.field} ${c.operator} ${c.value}`).join(', ')
}

function getActionSummary(rule: Rule): string {
  return rule.actions.map((a) => `${a.type}: ${a.value}`).join(', ')
}
</script>

<template>
  <div class="rule-list-page">
    <div class="page-header">
      <h2>{{ t('rule.title') }}</h2>
      <el-button type="primary" @click="handleCreate">{{ t('rule.create') }}</el-button>
    </div>

    <el-table :data="rules" v-loading="loading" stripe>
      <el-table-column prop="name" :label="t('rule.name')" />
      <el-table-column prop="logic_type" :label="t('rule.logicType')" width="80">
        <template #default="{ row }">{{ logicTypeOptions.find(o => o.value === row.logic_type)?.label || row.logic_type }}</template>
      </el-table-column>
      <el-table-column prop="trigger" :label="t('rule.trigger')" width="120">
        <template #default="{ row }">{{ triggerOptions.find(o => o.value === row.trigger)?.label || row.trigger }}</template>
      </el-table-column>
      <el-table-column :label="t('rule.conditions')" min-width="200">
        <template #default="{ row }">{{ getConditionSummary(row) }}</template>
      </el-table-column>
      <el-table-column :label="t('rule.actions')" min-width="200">
        <template #default="{ row }">{{ getActionSummary(row) }}</template>
      </el-table-column>
      <el-table-column :label="t('rule.enabled')" width="100">
        <template #default="{ row }">
          <el-switch :model-value="row.is_enabled" @change="handleToggleStatus(row)" />
        </template>
      </el-table-column>
      <el-table-column :label="t('common.edit')" width="240" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="primary" @click="handleExecute(row)">{{ t('rule.execute') }}</el-button>
          <el-button link type="danger" @click="handleDelete(row.id)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <Pagination
      :total="pagination.total"
      :page="pagination.page"
      :page-size="pagination.page_size"
      @update:page="handlePageChange"
      @update:page-size="handleSizeChange"
    />

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('rule.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('rule.logicType')">
          <el-select v-model="form.logic_type" style="width: 100%">
            <el-option v-for="opt in logicTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('rule.trigger')">
          <el-select v-model="form.trigger" style="width: 100%">
            <el-option v-for="opt in triggerOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('rule.priority')">
          <el-input-number v-model="form.priority" :min="0" />
        </el-form-item>

        <el-divider>{{ t('rule.conditions') }}</el-divider>
        <div v-for="(cond, index) in form.conditions" :key="'c' + index" class="rule-item">
          <el-row :gutter="8">
            <el-col :span="6"><el-select v-model="cond.field" :placeholder="t('rule.field')" filterable><el-option v-for="opt in conditionFieldOptions" :key="opt.value" :label="opt.label" :value="opt.value" /></el-select></el-col>
            <el-col :span="6"><el-select v-model="cond.operator" :placeholder="t('rule.operator')" filterable><el-option v-for="opt in conditionOperatorOptions" :key="opt.value" :label="opt.label" :value="opt.value" /></el-select></el-col>
            <el-col :span="8"><el-input v-model="cond.value" :placeholder="t('rule.value')" /></el-col>
            <el-col :span="4"><el-button type="danger" link @click="removeCondition(index)">{{ t('common.delete') }}</el-button></el-col>
          </el-row>
        </div>
        <el-button type="primary" link @click="addCondition">+ {{ t('rule.conditions') }}</el-button>

        <el-divider>{{ t('rule.actions') }}</el-divider>
        <div v-for="(act, index) in form.actions" :key="'a' + index" class="rule-item">
          <el-row :gutter="8">
            <el-col :span="6"><el-select v-model="act.type" :placeholder="t('rule.actionType')" filterable><el-option v-for="opt in actionTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" /></el-select></el-col>
            <el-col :span="10"><el-input v-model="act.value" :placeholder="t('rule.value')" /></el-col>
            <el-col :span="8"><el-button type="danger" link @click="removeAction(index)">{{ t('common.delete') }}</el-button></el-col>
          </el-row>
        </div>
        <el-button type="primary" link @click="addAction">+ {{ t('rule.actions') }}</el-button>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}

.rule-item {
  margin-bottom: 8px;
}
</style>