<script setup lang="ts">
// 饼图组件 - 基于ECharts封装，用于展示分类占比
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 注册ECharts必需的组件
use([TitleComponent, TooltipComponent, LegendComponent, PieChart, CanvasRenderer])

const props = defineProps<{
  data: Record<string, unknown>[] // 数据源数组
  nameField: string // 名称字段名（饼图扇区名称）
  valueField: string // 数值字段名（饼图扇区大小）
  title?: string // 图表标题
}>()

// 根据props动态计算ECharts配置项
const option = computed(() => ({
  title: {
    text: props.title || '',
    left: 'center',
  },
  tooltip: {
    trigger: 'item' as const,
    formatter: '{b}: {c} ({d}%)', // 格式：名称: 值 (百分比%)
  },
  legend: {
    orient: 'vertical' as const,
    left: 'left',
    bottom: 0,
  },
  series: [
    {
      type: 'pie' as const,
      radius: '60%',
      center: ['50%', '50%'],
      data: props.data.map((item) => ({
        name: item[props.nameField],
        value: item[props.valueField],
      })),
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)',
        },
      },
    },
  ],
}))
</script>

<template>
  <VChart :option="option" autoresize style="height: 400px; width: 100%" />
</template>
