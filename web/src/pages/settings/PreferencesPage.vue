<!--
  用户偏好设置页面
  功能：管理用户个性化配置（键值对形式）
  对应后端 API：/api/v1/preferences
-->
<template>
  <div class="preferences-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('preference.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            {{ t('common.add') }}
          </el-button>
        </div>
      </template>

      <!-- 偏好列表 -->
      <el-table :data="preferences" v-loading="loading" style="width: 100%">
        <el-table-column prop="key" :label="t('preference.key')" />
        <el-table-column prop="value" :label="t('preference.value')">
          <template #default="{ row }">
            <el-input
              v-if="row.editing"
              v-model="row.value"
              size="small"
              @keyup.enter="handleUpdate(row)"
            />
            <span v-else>{{ row.value }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.editing"
              link
              type="success"
              @click="handleUpdate(row)"
            >
              {{ t('common.save') }}
            </el-button>
            <el-button
              v-else
              link
              type="primary"
              @click="handleEdit(row)"
            >
              {{ t('common.edit') }}
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加偏好对话框 -->
    <el-dialog v-model="dialogVisible" :title="t('common.add')" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('preference.key')" prop="key">
          <el-input v-model="form.key" :placeholder="t('preference.keyPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('preference.value')" prop="value">
          <el-input
            v-model="form.value"
            type="textarea"
            :rows="3"
            :placeholder="t('preference.valuePlaceholder')"
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
 * 用户偏好设置页面
 * 提供用户偏好的查询、添加、编辑、删除功能
 */
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import { useI18n } from 'vue-i18n';
import * as preferenceApi from '@/api/preference';
import type { Preference, SetPreferenceRequest } from '@/types/preference';

const { t } = useI18n();

/** 加载状态 */
const loading = ref(false);
/** 提交状态 */
const submitting = ref(false);
/** 对话框显示状态 */
const dialogVisible = ref(false);

/** 偏好列表数据 */
const preferences = ref<Array<Preference & { editing?: boolean }>>([]);

/** 表单引用 */
const formRef = ref<FormInstance>();

/** 表单数据 */
const form = reactive<SetPreferenceRequest>({
  key: '',
  value: ''
});

/** 表单验证规则 */
const rules: FormRules = {
  key: [
    { required: true, message: t('common.required'), trigger: 'blur' },
    { pattern: /^[a-z0-9_.]+$/, message: '键名只能包含小写字母、数字、下划线和点', trigger: 'blur' }
  ],
  value: [{ required: true, message: t('common.required'), trigger: 'blur' }]
};

/**
 * 获取偏好列表数据
 */
async function fetchData() {
  loading.value = true;
  try {
    const res = await preferenceApi.list() as unknown as Preference[];
    preferences.value = res.map(item => ({ ...item, editing: false }));
  } catch (error) {
    ElMessage.error(t('common.fetchError'));
  } finally {
    loading.value = false;
  }
}

/**
 * 打开添加对话框
 */
function handleAdd() {
  Object.assign(form, { key: '', value: '' });
  dialogVisible.value = true;
}

/**
 * 编辑偏好值
 * @param row 选中的偏好数据
 */
function handleEdit(row: Preference & { editing?: boolean }) {
  row.editing = true;
}

/**
 * 更新偏好值
 * @param row 选中的偏好数据
 */
async function handleUpdate(row: Preference & { editing?: boolean }) {
  try {
    await preferenceApi.set({ key: row.key, value: row.value });
    ElMessage.success(t('common.success'));
    row.editing = false;
    fetchData();
  } catch (error) {
    ElMessage.error(t('common.submitError'));
  }
}

/**
 * 删除偏好
 * @param row 选中的偏好数据
 */
function handleDelete(row: Preference) {
  ElMessageBox.confirm(
    t('common.confirmDelete'),
    t('common.warning'),
    { type: 'warning' }
  ).then(async () => {
    try {
      await preferenceApi.remove(row.key);
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
 * 提交添加表单
 */
async function handleSubmit() {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    
    submitting.value = true;
    try {
      await preferenceApi.set(form);
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
.preferences-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
