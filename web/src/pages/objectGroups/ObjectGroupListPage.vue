<!--
  对象分组管理页面
  功能：管理各类实体（交易、账单、预算、储蓄罐）的分组和排序
  对应后端 API：/api/v1/object-groups
-->
<template>
  <div class="object-groups-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('objectGroup.title') }}</span>
          <el-button type="primary" @click="handleCreate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item :label="t('objectGroup.type')">
          <el-select v-model="filters.groupable_type" clearable :placeholder="t('common.all')">
            <el-option :label="t('objectGroup.typeTransaction')" value="transaction" />
            <el-option :label="t('objectGroup.typeRecurringTransaction')" value="recurring_transaction" />
            <el-option :label="t('objectGroup.typeBudget')" value="budget" />
            <el-option :label="t('objectGroup.typePiggyBank')" value="piggy_bank" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 分组列表 -->
      <el-table :data="groups" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" :label="t('common.id')" width="80" />
        <el-table-column prop="name" :label="t('objectGroup.name')" />
        <el-table-column prop="groupable_type" :label="t('objectGroup.type')">
          <template #default="{ row }">
            <el-tag>{{ getTypeLabel(row.groupable_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="groupable_id" :label="t('objectGroup.entityId')" width="120" />
        <el-table-column prop="created_at" :label="t('common.createdAt')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('common.edit') : t('common.create')"
      width="500px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('objectGroup.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('objectGroup.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('objectGroup.type')" prop="groupable_type">
          <el-select v-model="form.groupable_type" :placeholder="t('common.select')" style="width: 100%">
            <el-option :label="t('objectGroup.typeTransaction')" value="transaction" />
            <el-option :label="t('objectGroup.typeRecurringTransaction')" value="recurring_transaction" />
            <el-option :label="t('objectGroup.typeBudget')" value="budget" />
            <el-option :label="t('objectGroup.typePiggyBank')" value="piggy_bank" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('objectGroup.entityId')" prop="groupable_id">
          <el-input-number
            v-model="form.groupable_id"
            :min="1"
            :placeholder="t('objectGroup.entityIdPlaceholder')"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
/**
 * 对象分组管理页面
 * 提供对象分组的增删改查功能
 */
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import { useI18n } from 'vue-i18n';
import * as objectGroupApi from '@/api/objectGroup';
import type { ObjectGroup, CreateObjectGroupRequest } from '@/types/objectGroup';
import { formatDate } from '@/utils/format';

const { t } = useI18n();

/** 加载状态 */
const loading = ref(false);
/** 提交状态 */
const submitting = ref(false);
/** 对话框显示状态 */
const dialogVisible = ref(false);
/** 是否编辑模式 */
const isEdit = ref(false);
/** 当前编辑的分组 ID */
const currentId = ref<number>();

/** 分组列表数据 */
const groups = ref<ObjectGroup[]>([]);

/** 筛选条件 */
const filters = reactive({
  groupable_type: ''
});

/** 表单引用 */
const formRef = ref<FormInstance>();

/** 表单数据 */
const form = reactive<CreateObjectGroupRequest>({
  name: '',
  groupable_type: 'transaction',
  groupable_id: 1
});

/** 表单验证规则 */
const rules: FormRules = {
  name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  groupable_type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  groupable_id: [{ required: true, message: t('common.required'), trigger: 'blur' }]
};

/**
 * 获取类型标签文本 - 使用i18n国际化
 * @param type 类型值
 * @returns 类型标签文本
 */
function getTypeLabel(type: string): string {
  const typeMap: Record<string, string> = {
    transaction: t('objectGroup.typeTransaction'),
    recurring_transaction: t('objectGroup.typeRecurringTransaction'),
    budget: t('objectGroup.typeBudget'),
    piggy_bank: t('objectGroup.typePiggyBank')
  };
  return typeMap[type] || type;
}

/**
 * 获取分组列表数据
 */
async function fetchData() {
  loading.value = true;
  try {
    const params = filters.groupable_type ? { groupable_type: filters.groupable_type } : undefined;
    const res = await objectGroupApi.list(params);
    groups.value = res;
  } catch (error) {
    ElMessage.error(t('common.fetchError'));
  } finally {
    loading.value = false;
  }
}

/**
 * 重置筛选条件
 */
function handleReset() {
  filters.groupable_type = '';
  fetchData();
}

/**
 * 打开创建对话框
 */
function handleCreate() {
  isEdit.value = false;
  currentId.value = undefined;
  Object.assign(form, {
    name: '',
    groupable_type: 'transaction',
    groupable_id: 1
  });
  dialogVisible.value = true;
}

/**
 * 打开编辑对话框
 * @param row 选中的分组数据
 */
function handleEdit(row: ObjectGroup) {
  isEdit.value = true;
  currentId.value = row.id;
  Object.assign(form, {
    name: row.name,
    groupable_type: row.groupable_type,
    groupable_id: row.groupable_id
  });
  dialogVisible.value = true;
}

/**
 * 删除分组
 * @param row 选中的分组数据
 */
function handleDelete(row: ObjectGroup) {
  ElMessageBox.confirm(
    t('common.confirmDelete'),
    t('common.warning'),
    { type: 'warning' }
  ).then(async () => {
    try {
      await objectGroupApi.remove(row.id);
      ElMessage.success(t('common.success'));
      fetchData();
    } catch (error) {
      ElMessage.error(t('common.deleteError'));
    }
  }).catch(() => {
    // 用户取消
  });
}

/**
 * 提交表单数据
 */
async function handleSubmit() {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    
    submitting.value = true;
    try {
      if (isEdit.value && currentId.value) {
        await objectGroupApi.update(currentId.value, { name: form.name });
      } else {
        await objectGroupApi.create(form);
      }
      ElMessage.success(t('common.success'));
      dialogVisible.value = false;
      fetchData();
    } catch (error) {
      ElMessage.error(t('common.submitError'));
    } finally {
      submitting.value = false;
    }
  });
}

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.object-groups-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 20px;
}
</style>
