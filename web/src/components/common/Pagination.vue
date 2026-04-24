<script setup lang="ts">
// 分页组件 - 封装Element Plus分页，支持页码和每页数量切换
const props = defineProps<{
  total: number // 总记录数
  page: number // 当前页码
  pageSize: number // 每页数量
}>()

const emit = defineEmits<{
  'update:page': [value: number] // 页码变化
  'update:pageSize': [value: number] // 每页数量变化
}>()

// 切换页码
function handleCurrentChange(val: number) {
  emit('update:page', val)
}

// 切换每页数量时，重置到第1页
function handleSizeChange(val: number) {
  emit('update:pageSize', val)
  emit('update:page', 1)
}
</script>

<template>
  <div class="pagination-wrapper">
    <!-- layout: 总数、每页数量、上/下一页、页码、跳转 -->
    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
