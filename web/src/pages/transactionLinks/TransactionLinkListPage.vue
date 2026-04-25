<!-- 交易关联页面 - 管理交易之间的关联关系（冲销、对账、关联），支持筛选、新增、编辑和删除 -->
<template>
  <div class="transaction-links-page">
    <el-card>
      <!-- 页面标题和新增按钮 -->
      <template #header>
        <div class="card-header">
          <span>{{ t('transactionLink.title') }}</span>
          <el-button type="primary" @click="handleCreate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <!-- 筛选条件表单 -->
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item :label="t('transactionLink.transaction')">
          <el-input-number
            v-model="filters.transaction_id"
            :placeholder="t('transactionLink.transactionPlaceholder')"
            :min="1"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">{{ t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 交易关联数据表格 -->
      <el-table :data="links" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" :label="t('common.id')" width="80" />
        <el-table-column prop="transaction_id" :label="t('transactionLink.transaction')" width="120" />
        <!-- 关联类型列，使用标签展示不同颜色 -->
        <el-table-column prop="link_type" :label="t('transactionLink.type')" width="150">
          <template #default="{ row }">
            <el-tag :type="getLinkTypeTag(row.link_type)">
              {{ getLinkTypeLabel(row.link_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="linked_journal_id" :label="t('transactionLink.linkedJournal')" width="150" />
        <el-table-column prop="created_at" :label="t('common.createdAt')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <!-- 操作列：编辑和删除 -->
        <el-table-column :label="t('common.actions')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页组件 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="fetchData"
        @size-change="fetchData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 新增/编辑交易关联对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="140px">
        <!-- 关联的交易 ID -->
        <el-form-item :label="t('transactionLink.transaction')" prop="transaction_id">
          <el-input-number
            v-model="form.transaction_id"
            :min="1"
            :placeholder="t('transactionLink.transactionPlaceholder')"
            style="width: 100%"
          />
        </el-form-item>
        <!-- 关联类型选择 -->
        <el-form-item :label="t('transactionLink.type')" prop="link_type">
          <el-select v-model="form.link_type" :placeholder="t('common.select')" style="width: 100%">
            <el-option
              v-for="item in linkTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <!-- 关联的日记账 ID -->
        <el-form-item :label="t('transactionLink.linkedJournal')" prop="linked_journal_id">
          <el-input-number
            v-model="form.linked_journal_id"
            :min="1"
            :placeholder="t('transactionLink.linkedJournalPlaceholder')"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <!-- 对话框底部按钮 -->
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
// 导入 Vue 响应式、计算属性和生命周期 API
import { ref, reactive, computed, onMounted } from 'vue';
// 导入 Element Plus 消息提示和确认框组件
import { ElMessage, ElMessageBox } from 'element-plus';
// 导入表单类型定义
import type { FormInstance, FormRules } from 'element-plus';
// 导入国际化钩子函数
import { useI18n } from 'vue-i18n';
// 导入交易关联 API
import * as transactionLinkApi from '@/api/transactionLink';
// 导入交易关联类型定义
import type { TransactionLink, CreateTransactionLinkRequest, UpdateTransactionLinkRequest } from '@/types/transactionLink';
// 导入日期格式化工具
import { formatDate } from '@/utils/format';

// 国际化翻译函数
const { t } = useI18n();

// 关联类型下拉选项（冲销、对账、关联）
const linkTypeOptions = computed(() => [
  { label: t('transactionLink.typeRolledBack'), value: 'rolled_back' },
  { label: t('transactionLink.typeReconciled'), value: 'reconciled' },
  { label: t('transactionLink.typeLinked'), value: 'linked' }
]);

// 列表加载状态
const loading = ref(false);
// 提交加载状态（新增/编辑时使用）
const submitting = ref(false);
// 对话框是否可见
const dialogVisible = ref(false);
// 对话框标题
const dialogTitle = ref('');
// 当前编辑的关联 ID，null 表示新增模式
const editingId = ref<number | null>(null);

// 交易关联列表数据
const links = ref<TransactionLink[]>([]);

// 分页参数
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
});

// 筛选条件
const filters = reactive({
  transaction_id: undefined as number | undefined
});

// 表单引用，用于调用验证方法
const formRef = ref<FormInstance>();

// 新增/编辑表单数据
const form = reactive<CreateTransactionLinkRequest>({
  transaction_id: 1,
  link_type: 'linked',
  linked_journal_id: 1
});

// 表单验证规则
const rules: FormRules = {
  transaction_id: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  link_type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  linked_journal_id: [{ required: true, message: t('common.required'), trigger: 'blur' }]
};

// 根据关联类型返回对应的标签颜色类型
function getLinkTypeTag(type: string): 'danger' | 'success' | 'info' {
  const typeMap: Record<string, 'danger' | 'success' | 'info'> = {
    rolled_back: 'danger',    // 冲销 - 红色
    reconciled: 'success',    // 对账 - 绿色
    linked: 'info'            // 关联 - 蓝色
  };
  return typeMap[type] || 'info';
}

// 根据关联类型返回对应的显示文本
function getLinkTypeLabel(type: string): string {
  const typeMap: Record<string, string> = {
    rolled_back: t('transactionLink.typeRolledBack'),
    reconciled: t('transactionLink.typeReconciled'),
    linked: t('transactionLink.typeLinked')
  };
  return typeMap[type] || type;
}

// 获取交易关联列表数据
async function fetchData() {
  loading.value = true;
  try {
    // 构建查询参数，有筛选条件时才传
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    };
    if (filters.transaction_id) {
      params.transaction_id = filters.transaction_id;
    }

    const res = await transactionLinkApi.list(params);
    links.value = res.items;
    pagination.total = res.total;
  } catch (error) {
    ElMessage.error(t('common.fetchError'));
  } finally {
    loading.value = false;
  }
}

// 重置筛选条件并重新查询
function handleReset() {
  filters.transaction_id = undefined;
  pagination.page = 1;
  fetchData();
}

// 点击新增按钮，重置表单并打开对话框
function handleCreate() {
  dialogTitle.value = t('common.create');
  editingId.value = null;
  Object.assign(form, {
    transaction_id: 1,
    link_type: 'linked',
    linked_journal_id: 1
  });
  dialogVisible.value = true;
}

// 点击编辑按钮，填充表单数据并打开对话框
function handleEdit(row: TransactionLink) {
  dialogTitle.value = t('common.edit');
  editingId.value = row.id;
  Object.assign(form, {
    transaction_id: row.transaction_id,
    link_type: row.link_type,
    linked_journal_id: row.linked_journal_id
  });
  dialogVisible.value = true;
}

// 删除交易关联（带确认弹窗）
function handleDelete(row: TransactionLink) {
  ElMessageBox.confirm(
    t('common.confirmDelete'),
    t('common.warning'),
    { type: 'warning' }
  ).then(async () => {
    try {
      await transactionLinkApi.remove(row.id);
      ElMessage.success(t('common.success'));
      fetchData();
    } catch (error) {
      ElMessage.error(t('common.deleteError'));
    }
  }).catch(() => {});
}

// 提交新增/编辑表单
async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) return;

    submitting.value = true;
    try {
      if (editingId.value) {
        // 编辑模式：调用更新 API
        await transactionLinkApi.update(editingId.value, form as UpdateTransactionLinkRequest);
      } else {
        // 新增模式：调用创建 API
        await transactionLinkApi.create(form);
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

// 页面挂载时加载数据
onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.transaction-links-page {
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
