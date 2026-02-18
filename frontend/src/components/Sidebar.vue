<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// 导航菜单项
const menuItems = [
  { id: 'markdown', name: '小说编辑器', icon: '📝', path: '/markdown' },
  { id: 'renwu', name: '人物', icon: '👤', path: '/renwu' },
  { id: 'wuqi', name: '武器', icon: '⚔️', path: '/wuqi' },
  { id: 'daoju', name: '道具', icon: '🎒', path: '/daoju' },
  { id: 'beibao', name: '背包', icon: '🎒', path: '/beibao' },
  { id: 'chongwu', name: '宠物', icon: '🐾', path: '/chongwu' },
  { id: 'shiqing', name: '任务', icon: '📋', path: '/shiqing' },
  { id: 'shili', name: '势力', icon: '🏰', path: '/shili' },
  { id: 'guaiwu', name: '怪物', icon: '👹', path: '/guaiwu' },
  { id: 'contact', name: '联系作者', icon: '📧', path: '/contact' }
]

// 当前激活的菜单项
const activeItem = computed(() => {
  return route.path
})

// 导航跳转
function navigate(path) {
  router.push(path)
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <h1>小说信息管理</h1>
    </div>
    <nav class="sidebar-nav">
      <button
        v-for="item in menuItems"
        :key="item.id"
        :class="['nav-button', { active: activeItem === item.path }]"
        @click="navigate(item.path)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span class="nav-text">{{ item.name }}</span>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 220px;
  background-color: #ffffff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.04);
  height: 100%;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.sidebar-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  letter-spacing: 1px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}

.nav-button {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 8px;
  border: none;
  background-color: transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 15px;
  color: #5a6c7d;
}

.nav-button:hover {
  background-color: #f0f7ff;
  color: #1890ff;
}

.nav-button.active {
  background-color: #e6f7ff;
  color: #1890ff;
  font-weight: 500;
}

.nav-icon {
  font-size: 20px;
  margin-right: 12px;
  width: 24px;
  text-align: center;
}

.nav-text {
  flex: 1;
  text-align: left;
}
</style>