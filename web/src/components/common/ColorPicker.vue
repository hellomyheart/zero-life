<script setup lang="ts">
/**
 * ColorPicker 颜色选择器组件
 * 功能：
 * - 提供预设调色板快速选择常用颜色
 * - 支持通过 el-color-picker 自定义任意颜色
 * - v-model 双向绑定颜色值（十六进制格式）
 */
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

/** 预设颜色调色板，覆盖常用标签颜色 */
const presetColors = [
  '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
  '#00B5AD', '#11A0F8', '#6C5CE7', '#FD79A8', '#FDCB6E',
  '#E17055', '#00CEC9', '#0984E3', '#6C5CE7', '#B2BEC3',
  '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7',
]

/** 当前选中的颜色值 */
const colorValue = computed({
  get: () => props.modelValue,
  set: (val: string) => emit('update:modelValue', val),
})

/**
 * 选择预设颜色
 * @param color 颜色值
 */
function selectPreset(color: string) {
  colorValue.value = color
}
</script>

<template>
  <div class="color-picker">
    <div class="preset-colors">
      <div
        v-for="color in presetColors"
        :key="color"
        class="color-swatch"
        :class="{ active: modelValue === color }"
        :style="{ backgroundColor: color }"
        :title="color"
        @click="selectPreset(color)"
      />
    </div>
    <div class="custom-color">
      <el-color-picker v-model="colorValue" size="small" />
      <span class="custom-label">{{ t('tag.color') }}</span>
    </div>
  </div>
</template>

<style scoped>
.color-picker {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preset-colors {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.color-swatch {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.color-swatch:hover {
  transform: scale(1.15);
}

.color-swatch.active {
  border-color: #303133;
  box-shadow: 0 0 0 2px rgba(48, 49, 51, 0.2);
}

.custom-color {
  display: flex;
  align-items: center;
  gap: 8px;
}

.custom-label {
  font-size: 13px;
  color: #606266;
}
</style>
