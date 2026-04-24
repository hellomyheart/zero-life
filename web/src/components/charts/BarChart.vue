<script setup lang="ts">
// 柱状图组件 - 基于ECharts封装，支持多系列柱状图
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

// 注册ECharts必需的组件
use([TitleComponent, TooltipComponent, LegendComponent, GridComponent, BarChart, CanvasRenderer])

const props = defineProps<{
  data: Record<string, unknown>[] // 数据源数组
  xField: string // X轴字段名
  yFields: { field: string; name: string }[] // Y轴字段配置，支持多组柱子
  title?: string // 图表标题
}>()

// 根据props动态计算ECharts配置项
const option = computed(() => ({
  title: {
    text: props.title || '',
    left: 'center',
  },
  tooltip: {
    trigger: 'axis' as const,
  },
  legend: {
    data: props.yFields.map((y) => y.name),
    bottom: 0,
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '12%',
    containLabel: true,
  },
  xAxis: {
    type: 'category' as const,
    data: props.data.map((item) => item[props.xField]),
  },
  yAxis: {
    type: 'value' as const,
  },
  series: props.yFields.map((y) => ({
    name: y.name,
    type: 'bar' as const,
    data: props.data.map((item) => item[y.field]),
  })),
}))
</script>

<template>
  <VChart :option="option" autoresize style="height: 400px; width: 100%" />
</template>
