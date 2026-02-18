import { createRouter, createWebHashHistory } from 'vue-router'

// 页面组件（懒加载）
const MarkdownEditor = () => import('../views/MarkdownEditor.vue')
const Renwu = () => import('../views/Renwu.vue')
const Wuqi = () => import('../views/Wuqi.vue')
const Daoju = () => import('../views/Daoju.vue')
const Chongwu = () => import('../views/Chongwu.vue')
const Shiqing = () => import('../views/Shiqing.vue')
const Shili = () => import('../views/Shili.vue')
const Guaiwu = () => import('../views/Guaiwu.vue')
const Beibao = () => import('../views/Beibao.vue')
const Shopping = () => import('../views/Shopping.vue')
const Contact = () => import('../views/Contact.vue')

// 路由配置
const routes = [
  {
    path: '/',
    redirect: '/markdown'
  },
  {
    path: '/markdown',
    name: 'MarkdownEditor',
    component: MarkdownEditor,
    meta: { title: 'Markdown编辑器', icon: '📝' }
  },
  {
    path: '/renwu',
    name: 'Renwu',
    component: Renwu,
    meta: { title: '人物', icon: '👤' }
  },
  {
    path: '/wuqi',
    name: 'Wuqi',
    component: Wuqi,
    meta: { title: '武器', icon: '⚔️' }
  },
  {
    path: '/daoju',
    name: 'Daoju',
    component: Daoju,
    meta: { title: '道具', icon: '🎒' }
  },
  {
    path: '/chongwu',
    name: 'Chongwu',
    component: Chongwu,
    meta: { title: '宠物', icon: '🐾' }
  },
  {
    path: '/shiqing',
    name: 'Shiqing',
    component: Shiqing,
    meta: { title: '任务', icon: '📋' }
  },
  {
    path: '/shili',
    name: 'Shili',
    component: Shili,
    meta: { title: '势力', icon: '🏰' }
  },
  {
    path: '/guaiwu',
    name: 'Guaiwu',
    component: Guaiwu,
    meta: { title: '怪物', icon: '👹' }
  },
  {
    path: '/beibao',
    name: 'Beibao',
    component: Beibao,
    meta: { title: '背包', icon: '🎒' }
  },
  {
    path: '/shopping',
    name: 'Shopping',
    component: Shopping,
    meta: { title: '商城', icon: '🛒' }
  },
  {
    path: '/contact',
    name: 'Contact',
    component: Contact,
    meta: { title: '联系作者', icon: '📧' }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || 'NoolTools'} - NoolTools`
  next()
})

export default router