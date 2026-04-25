// 国际化配置文件 - 配置多语言支持（中文/英文），用于整个应用的文本翻译
import { createI18n } from 'vue-i18n'
// 导入中文语言包
import zhCN from './locales/zh-CN.json'
// 导入英文语言包
import enUS from './locales/en-US.json'

// 创建 i18n 实例
const i18n = createI18n({
  // 使用 Composition API 模式（而非 Options API 的 legacy 模式）
  legacy: false,
  // 默认语言：优先从 localStorage 读取用户上次选择的语言，否则使用中文
  locale: localStorage.getItem('locale') || 'zh-CN',
  // 回退语言：当当前语言的翻译缺失时，使用中文作为兜底
  fallbackLocale: 'zh-CN',
  // 语言消息映射，键为语言代码，值为对应的翻译内容
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

export default i18n
