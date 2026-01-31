<template>
  <div class="page-container">
    <div class="page-header">
      <h1>武器管理</h1>
      <p>管理游戏中的武器信息</p>
    </div>
    <div class="page-content">
      <div class="empty-state">
        <p>武器页面内容开发中...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
// 武器页面组件
</script>

<style scoped>
.page-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 32px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.page-header p {
  font-size: 14px;
  color: #8c92a0;
}

.page-content {
  flex: 1;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #8c92a0;
  font-size: 16px;
}
</style>