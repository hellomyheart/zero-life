<!--
  交易关联管理页面
  功能：管理交易之间的关联关系（如回滚、对账、关联等）
  对应后端 API：/api/v1/transaction-links
-->
<template>
  <div class="transaction-links-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('transactionLink.title') }}</span>
          <el-button type="primary" @click="handleCreate">
            {{ t('common.create') }}
          </el-button>
        </div>
      </template>

      <!-- 筛选条件 -->
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

      <!-- 关联列表 -->
      <el-table :data="links" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" :label="t('common.id')" width="80" />
        <el-table-column prop="transaction_id" :label="t('transactionLink.transaction')" width="120" />
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
        <el-table-column :label="t('common.actions')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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

    <!-- 创建对话框 -->
    <el-dialog v-model="dialogVisible" :title="t('common.create')" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="140px">
        <el-form-item :label="t('transactionLink.transaction')" prop="transaction_id">
          <el-input-number
            v-model="form.transaction_id"
            :min="1"
            :placeholder="t('transactionLink.transactionPlaceholder')"
            style="width: 100%"
          />
        </el-form-item>
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
        <el-form-item :label="t('transactionLink.linkedJournal')" prop="linked_journal_id">
          <el-input-number
            v-model="form.linked_journal_id"
            :min="1"
            :placeholder="t('transactionLink.linkedJournalPlaceholder')"
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
 * 交易关联管理页面
 * 提供交易关联的创建、查询、删除功能
 */
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import { useI18n } from 'vue-i18n';
import * as transactionLinkApi from '@/api/transactionLink';
import type { TransactionLink, CreateTransactionLinkRequest } from '@/types/transactionLink';
import { formatDate } from '@/utils/format';

const { t } = useI18n();

/** 关联类型选项 */
const linkTypeOptions = [
  { label: '回滚', value: 'rolled_back' },
  { label: '已对账', value: 'reconciled' },
  { label: '已关联', value: 'linked' }
];

/** 加载状态 */
const loading = ref(false);
/** 提交状态 */
const submitting = ref(false);
/** 对话框显示状态 */
const dialogVisible = ref(false);

/** 关联列表数据 */
const links = ref<TransactionLink[]>([]);

/** 分页信息 */
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
});

/** 筛选条件 */
const filters = reactive({
  transaction_id: undefined as number | undefined
});

/** 表单引用 */
const formRef = ref<FormInstance>();

/** 表单数据 */
const form = reactive<CreateTransactionLinkRequest>({
  transaction_id: 1,
  link_type: 'linked',
  linked_journal_id: 1
});

/** 表单验证规则 */
const rules: FormRules = {
  transaction_id: [{ required: true, message: t('common.required'), trigger: 'blur' }],
  link_type: [{ required: true, message: t('common.required'), trigger: 'change' }],
  linked_journal_id: [{ required: true, message: t('common.required'), trigger: 'blur' }]
};

/**
 * 获取关联类型标签颜色
 * @param type 关联类型
 * @returns 标签类型
 */
function getLinkTypeTag(type: string): 'danger' | 'success' | 'info' {
  const typeMap: Record<string, 'danger' | 'success' | 'info'> = {
    rolled_back: 'danger',
    reconciled: 'success',
    linked: 'info'
  };
  return typeMap[type] || 'info';
}

/**
 * 获取关联类型标签文本
 * @param type 关联类型
 * @returns 类型标签文本
 */
function getLinkTypeLabel(type: string): string {
  const typeMap: Record<string, string> = {
    rolled_back: '回滚',
    reconciled: '已对账',
    linked: '已关联'
  };
  return typeMap[type] || type;
}

/**
 * 获取关联列表数据
 */
async function fetchData() {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    };
    if (filters.transaction_id) {
      params.transaction_id = filters.transaction_id;
    }
    
    const res = await transactionLinkApi.list(params);
    links.value = res.list;
    pagination.total = res.total;
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
  filters.transaction_id = undefined;
  pagination.page = 1;
  fetchData();
}

/**
 * 打开创建对话框
 */
function handleCreate() {
  Object.assign(form, {
    transaction_id: 1,
    link_type: 'linked',
    linked_journal_id: 1
  });
  dialogVisible.value = true;
}

/**
 * 删除关联
 * @param row 选中的关联数据
 */
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
      await transactionLinkApi.create(form);
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
