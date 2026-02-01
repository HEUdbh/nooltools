<script setup>
import Sidebar from './components/Sidebar.vue'
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
