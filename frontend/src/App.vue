<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import Sidebar from './components/Sidebar.vue'
import UpdateNoticeModal from './components/UpdateNoticeModal.vue'
import { EventsOff, EventsOn } from '../wailsjs/runtime/runtime'

const showUpdateModal = ref(false)
const updateInfo = ref(null)
const isUpdating = ref(false)
const updateProgress = ref({
  stage: '',
  percent: 0,
  message: '',
  detail: ''
})
const theme = ref('light')
const THEME_STORAGE_KEY = 'nooltools-theme'

function handleUpdateProgress(event) {
  if (!event || typeof event !== 'object') {
    return
  }

  updateProgress.value = {
    stage: event.stage || '',
    percent: Number.isFinite(event.percent) ? event.percent : 0,
    message: event.message || '',
    detail: event.detail || ''
  }

  if (event.stage === 'error') {
    isUpdating.value = false
  }
}

function applyTheme(nextTheme) {
  const normalizedTheme = nextTheme === 'dark' ? 'dark' : 'light'
  theme.value = normalizedTheme
  document.body.classList.toggle('theme-dark', normalizedTheme === 'dark')
  document.body.classList.toggle('theme-light', normalizedTheme === 'light')
  localStorage.setItem(THEME_STORAGE_KEY, normalizedTheme)
}

function toggleTheme() {
  applyTheme(theme.value === 'dark' ? 'light' : 'dark')
}

onMounted(async () => {
  EventsOn('update:progress', handleUpdateProgress)
  applyTheme(localStorage.getItem(THEME_STORAGE_KEY) || 'light')

  try {
    const result = await window.go.main.app.CheckReleaseUpdate()
    if (result?.has_update) {
      updateInfo.value = result
      showUpdateModal.value = true
    }
  } catch (error) {
    console.error('检查更新失败:', error)
  }
})

onBeforeUnmount(() => {
  EventsOff('update:progress')
})

function handleCloseUpdateModal() {
  if (isUpdating.value) {
    return
  }
  showUpdateModal.value = false
}

async function handleStartAutoUpdate() {
  if (isUpdating.value) {
    return
  }

  isUpdating.value = true
  updateProgress.value = {
    stage: 'preparing',
    percent: 0,
    message: '正在准备更新',
    detail: ''
  }

  try {
    await window.go.main.app.StartAutoUpdate()
  } catch (error) {
    isUpdating.value = false
    updateProgress.value = {
      stage: 'error',
      percent: 0,
      message: '自动更新失败',
      detail: error?.message || String(error)
    }
  }
}
</script>

<template>
  <div class="app-container">
    <!-- 左侧侧边栏 -->
    <Sidebar :theme="theme" @toggle-theme="toggleTheme" />

    <!-- 右侧内容区域 -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
  <UpdateNoticeModal
    :visible="showUpdateModal"
    :update-info="updateInfo"
    :is-updating="isUpdating"
    :progress="updateProgress"
    @close="handleCloseUpdateModal"
    @start-update="handleStartAutoUpdate"
  />
</template>

<style>
.app-container {
  display: flex;
  height: 100vh;
  background-color: var(--app-bg);
}

/* 主内容区域样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
</style>
