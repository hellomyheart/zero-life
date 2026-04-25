// Element Plus 语言包映射
// 根据 i18n 语言代码获取对应的 Element Plus locale
// 用于 el-config-provider 动态切换组件语言（分页、日期选择等）
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'

export const elementLocales: Record<string, any> = {
  'zh-CN': zhCn,
  'en-US': en,
}
