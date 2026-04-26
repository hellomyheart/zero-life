// 应用入口文件 - 创建Vue实例并注册全局插件
import { createApp } from 'vue'
import { createPinia } from 'pinia' // Pinia状态管理
import ElementPlus from 'element-plus' // Element Plus UI组件库
import 'element-plus/dist/index.css' // Element Plus样式
import * as ElementPlusIconsVue from '@element-plus/icons-vue' // Element Plus图标
import i18n from './i18n' // 国际化
import router from './router' // 路由
import App from './App.vue' // 根组件
import './styles/index.css' // 全局样式（暗黑模式、响应式）

const app = createApp(App)

// 注册全局插件
app.use(createPinia()) // 状态管理
app.use(router) // 路由
app.use(ElementPlus) // Element Plus（locale由el-config-provider动态控制）
app.use(i18n) // 国际化

// 全局注册所有Element Plus图标组件，可在模板中直接使用图标名
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 挂载应用到#app元素
app.mount('#app')
