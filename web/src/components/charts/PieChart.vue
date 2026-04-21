<script setup lang="ts">
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

use([TitleComponent, TooltipComponent, LegendComponent, PieChart, CanvasRenderer])

const props = defineProps<{
  data: Record<string, unknown>[]
  nameField: string
  valueField: string
  title?: string
}>()

const option = computed(() => ({
  title: {
    text: props.title || '',
    left: 'center',
  },
  tooltip: {
    trigger: 'item' as const,
    formatter: '{b}: {c} ({d}%)',
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
