<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([TitleComponent, TooltipComponent, LegendComponent, GridComponent, LineChart, CanvasRenderer])

const props = defineProps<{
  data: Record<string, unknown>[]
  xField: string
  yFields: { field: string; name: string }[]
  title?: string
}>()

const option = computed(() => ({
  title: {
    text: props.title || '',
    left: 'center',
  },
  tooltip: {
    trigger: 'axis',
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
    type: 'line' as const,
    data: props.data.map((item) => item[y.field]),
    smooth: true,
  })),
}))
</script>

<template>
  <VChart :option="option" autoresize style="height: 400px; width: 100%" />
</template>
