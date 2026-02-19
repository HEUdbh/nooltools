<script setup>
import { onMounted, ref } from 'vue'
import Sidebar from './components/Sidebar.vue'
import UpdateNoticeModal from './components/UpdateNoticeModal.vue'

const showUpdateModal = ref(false)
const updateInfo = ref(null)

onMounted(async () => {
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

function handleCloseUpdateModal() {
  showUpdateModal.value = false
}
</script>

<template>
  <div class="app-container">
    <!-- 左侧侧边栏 -->
    <Sidebar />

    <!-- 右侧内容区域 -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
  <UpdateNoticeModal
    :visible="showUpdateModal"
    :update-info="updateInfo"
    @close="handleCloseUpdateModal"
  />
</template>

<style>
.app-container {
  display: flex;
  height: 100vh;
  background-color: #f5f7fa;
}

/* 主内容区域样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
</style>
