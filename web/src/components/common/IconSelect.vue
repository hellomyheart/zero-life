<script setup lang="ts">
import { ref, computed } from 'vue'
import * as Icons from '@element-plus/icons-vue'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const iconNames = Object.keys(Icons)
const search = ref('')
const pickerVisible = ref(false)

const filteredIcons = computed(() => {
  if (!search.value) return iconNames
  return iconNames.filter(name => name.toLowerCase().includes(search.value.toLowerCase()))
})

function handleSelect(name: string) {
  emit('update:modelValue', name)
  pickerVisible.value = false
}

function handleClear() {
  emit('update:modelValue', '')
}
</script>

<template>
  <div class="icon-select">
    <div class="icon-select-trigger" @click="pickerVisible = !pickerVisible">
      <el-icon v-if="modelValue" :size="18">
        <component :is="(Icons as any)[modelValue]" />
      </el-icon>
      <span v-else class="icon-placeholder">选择图标</span>
      <el-icon v-if="modelValue" class="icon-clear" @click.stop="handleClear" :size="14">
        <Icons.CircleClose />
      </el-icon>
    </div>
    <el-popover
      v-model:visible="pickerVisible"
      placement="bottom-start"
      :width="360"
      trigger="click"
    >
      <div class="icon-picker">
        <el-input
          v-model="search"
          placeholder="搜索图标..."
          clearable
          size="small"
          class="icon-search"
        />
        <div class="icon-grid">
          <div
            v-for="name in filteredIcons"
            :key="name"
            class="icon-item"
            :class="{ active: modelValue === name }"
            :title="name"
            @click="handleSelect(name)"
          >
            <el-icon :size="20">
              <component :is="(Icons as any)[name]" />
            </el-icon>
          </div>
        </div>
      </div>
    </el-popover>
  </div>
</template>

<style scoped>
.icon-select-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 11px;
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-border-radius-base);
  cursor: pointer;
  min-width: 120px;
  position: relative;
}

.icon-select-trigger:hover {
  border-color: var(--el-border-color-hover);
}

.icon-placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 14px;
}

.icon-clear {
  position: absolute;
  right: 8px;
  color: var(--el-text-color-secondary);
  cursor: pointer;
}

.icon-clear:hover {
  color: var(--el-text-color-primary);
}

.icon-picker {
  max-height: 320px;
  overflow-y: auto;
}

.icon-search {
  margin-bottom: 8px;
}

.icon-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 4px;
}

.icon-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-item:hover {
  background-color: var(--el-fill-color-light);
}

.icon-item.active {
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}
</style>
