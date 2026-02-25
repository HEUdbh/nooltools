<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const props = defineProps({
  theme: {
    type: String,
    default: 'light'
  }
})
const emit = defineEmits(['toggle-theme'])

// 导航菜单项
const menuItems = [
  { id: 'markdown', name: '小说编辑器', icon: '📝', path: '/markdown' },
  { id: 'renwu', name: '人物', icon: '👤', path: '/renwu' },
  { id: 'wuqi', name: '武器', icon: '⚔️', path: '/wuqi' },
  { id: 'daoju', name: '道具', icon: '🎒', path: '/daoju' },
  { id: 'beibao', name: '背包', icon: '🎒', path: '/beibao' },
  { id: 'shopping', name: '商城', icon: '🛒', path: '/shopping' },
  { id: 'getjiang', name: '抽奖', icon: '🎰', path: '/getjiang' },
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

function handleToggleTheme() {
  emit('toggle-theme')
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
    <div class="sidebar-footer">
      <button class="theme-toggle-btn" @click="handleToggleTheme">
        {{ props.theme === 'dark' ? '☀️ 白天模式' : '🌙 黑夜模式' }}
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 220px;
  background-color: var(--app-surface);
  border-right: 1px solid var(--app-border);
  display: flex;
  flex-direction: column;
  box-shadow: var(--app-shadow-md);
  height: 100%;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.sidebar-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: var(--app-text-primary);
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
  color: var(--app-text-secondary);
}

.nav-button:hover {
  background-color: var(--app-hover-bg);
  color: var(--app-accent);
}

.nav-button.active {
  background-color: var(--app-active-bg);
  color: var(--app-accent);
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

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--app-divider-soft);
}

.theme-toggle-btn {
  width: 100%;
  padding: 10px 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  background-color: var(--app-muted-bg);
  color: var(--app-text-secondary);
  transition: all 0.2s ease;
}

.theme-toggle-btn:hover {
  background-color: var(--app-hover-bg);
  color: var(--app-accent);
}
</style>
