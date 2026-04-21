<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { list, remove, toggleStatus, execute } from '@/api/rule'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Rule, CreateRuleReq, RuleConditionReq, RuleActionReq } from '@/types/rule'

const { t } = useI18n()

const rules = ref<Rule[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref<string | null>(null)

const form = ref<CreateRuleReq>({
  name: '',
  conditions: [{ field: '', operator: '', value: '' }],
  actions: [{ type: '', field: '', value: '' }],
  is_enabled: true,
  priority: 0,
})

async function fetchRules() {
  loading.value = true
  try {
    rules.value = await list() as unknown as Rule[]
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogTitle.value = t('rule.create')
  editingId.value = null
  form.value = {
    name: '',
    conditions: [{ field: '', operator: '', value: '' }],
    actions: [{ type: '', field: '', value: '' }],
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
    conditions: rule.conditions.map((c) => ({ field: c.field, operator: c.operator, value: c.value })),
    actions: rule.actions.map((a) => ({ type: a.type, field: a.field, value: a.value })),
    is_enabled: rule.is_enabled,
    priority: rule.priority,
  }
  dialogVisible.value = true
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm(t('rule.deleteConfirm'), t('common.confirm'), { type: 'warning' })
    await remove(id)
    ElMessage.success(t('common.success'))
    await fetchRules()
  } catch {
    // cancelled or error
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
    await ElMessageBox.confirm(`Execute rule "${rule.name}"?`, t('common.confirm'), { type: 'info' })
    await execute(rule.id, { transaction_ids: [] })
    ElMessage.success(t('common.success'))
  } catch {
    // cancelled or error
  }
}

function addCondition() {
  form.value.conditions.push({ field: '', operator: '', value: '' })
}

function removeCondition(index: number) {
  form.value.conditions.splice(index, 1)
}

function addAction() {
  form.value.actions.push({ type: '', field: '', value: '' })
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
  return rule.actions.map((a) => `${a.type}: ${a.field} = ${a.value}`).join(', ')
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('rule.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('rule.priority')">
          <el-input-number v-model="form.priority" :min="0" />
        </el-form-item>

        <el-divider>{{ t('rule.conditions') }}</el-divider>
        <div v-for="(cond, index) in form.conditions" :key="'c' + index" class="rule-item">
          <el-row :gutter="8">
            <el-col :span="6"><el-input v-model="cond.field" :placeholder="t('rule.field')" /></el-col>
            <el-col :span="6"><el-input v-model="cond.operator" :placeholder="t('rule.operator')" /></el-col>
            <el-col :span="8"><el-input v-model="cond.value" :placeholder="t('rule.value')" /></el-col>
            <el-col :span="4"><el-button type="danger" link @click="removeCondition(index)">{{ t('common.delete') }}</el-button></el-col>
          </el-row>
        </div>
        <el-button type="primary" link @click="addCondition">+ Condition</el-button>

        <el-divider>{{ t('rule.actions') }}</el-divider>
        <div v-for="(act, index) in form.actions" :key="'a' + index" class="rule-item">
          <el-row :gutter="8">
            <el-col :span="5"><el-input v-model="act.type" :placeholder="t('rule.actionType')" /></el-col>
            <el-col :span="6"><el-input v-model="act.field" :placeholder="t('rule.field')" /></el-col>
            <el-col :span="8"><el-input v-model="act.value" :placeholder="t('rule.value')" /></el-col>
            <el-col :span="5"><el-button type="danger" link @click="removeAction(index)">{{ t('common.delete') }}</el-button></el-col>
          </el-row>
        </div>
        <el-button type="primary" link @click="addAction">+ Action</el-button>
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
