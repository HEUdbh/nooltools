<template>
  <div class="shiqing-sidebar">
    <div class="sidebar-header">
      <h3>任务列表</h3>
      <button class="add-btn" @click="$emit('add')" title="添加任务">
        <span>+</span>
      </button>
    </div>
    <div class="sidebar-content">
      <div
        v-for="item in shiqing"
        :key="item.id"
        class="shiqing-item"
        :class="{ active: item.id === activeShiqingId }"
        @click="$emit('select', item.id)"
      >
        <div class="shiqing-name">{{ item.name }}</div>
        <button
          class="delete-btn"
          @click.stop="$emit('delete', item.id)"
          title="删除任务"
        >
          ×
        </button>
      </div>
      <div v-if="shiqing.length === 0" class="empty-state">
        <p>暂无任务</p>
        <p class="tip">点击 + 添加新任务</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  shiqing: {
    type: Array,
    default: () => []
  },
  activeShiqingId: {
    type: Number,
    default: null
  }
})

defineEmits(['select', 'add', 'delete'])
</script>

<style scoped>
.shiqing-sidebar {
  width: 280px;
  background-color: white;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 24px;
  border-bottom: 1px solid #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.add-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: #1890ff;
  color: white;
  border-radius: 6px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
}

.add-btn:hover {
  background-color: #40a9ff;
  transform: scale(1.05);
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.shiqing-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  margin-bottom: 8px;
  background-color: #f8f9fa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.shiqing-item:hover {
  background-color: #f0f7ff;
  border-color: #1890ff;
}

.shiqing-item.active {
  background-color: #e6f7ff;
  border-color: #1890ff;
}

.shiqing-name {
  flex: 1;
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  color: #8c92a0;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
  margin-left: 12px;
  flex-shrink: 0;
}

.delete-btn:hover {
  background-color: #fff1f0;
  color: #ff4d4f;
}

.empty-state {
  text-align: center;
  padding: 48px 20px;
  color: #8c92a0;
}

.empty-state p {
  margin: 10px 0;
}

.tip {
  font-size: 13px;
  color: #b0b8c3;
}

/* 滚动条样式 */
.sidebar-content::-webkit-scrollbar {
  width: 6px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: #b0b8c3;
}
</style>