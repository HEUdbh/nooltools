<template>
  <div class="shopping-sidebar">
    <div class="sidebar-header">
      <h3>商城商品</h3>
      <button class="add-btn" @click="handleAdd" title="添加商品">
        <span>+</span>
      </button>
    </div>
    
    <div class="shopping-list">
      <div
        v-for="item in shopping"
        :key="item.id"
        :class="['shopping-item', { active: item.id === activeShoppingId }]"
        @click="handleSelect(item.id)"
      >
        <div class="shopping-name">{{ item.name }}</div>
        <div class="shopping-value">价值: {{ item.value }}</div>
        <button
          v-if="item.id === activeShoppingId"
          class="delete-item-btn"
          @click.stop="handleDelete(item.id)"
          title="删除商品"
        >
          ×
        </button>
      </div>
      
      <div v-if="shopping.length === 0" class="empty-state">
        <p>暂无商品</p>
        <p class="tip">点击上方按钮添加</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  shopping: {
    type: Array,
    default: () => []
  },
  activeShoppingId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete'])

function handleSelect(id) {
  emit('select', id)
}

function handleAdd() {
  emit('add')
}

function handleDelete(id) {
  emit('delete', id)
}
</script>

<style scoped>
.shopping-sidebar {
  width: 250px;
  background-color: white;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.add-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: #4CAF50;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 24px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.add-btn:hover {
  background-color: #45a049;
  transform: scale(1.1);
}

.shopping-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.shopping-item {
  padding: 12px 15px;
  margin-bottom: 8px;
  background-color: #fafafa;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.shopping-item:hover {
  background-color: #f0f0f0;
  border-color: #d0d0d0;
}

.shopping-item.active {
  background-color: #e8f5e9;
  border-color: #4CAF50;
}

.shopping-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  font-size: 14px;
}

.shopping-value {
  font-size: 12px;
  color: #666;
}

.delete-item-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  border: none;
  background-color: #ff5252;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.delete-item-btn:hover {
  background-color: #ff1744;
  transform: scale(1.1);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-state p {
  margin: 5px 0;
}

.empty-state .tip {
  font-size: 12px;
  color: #bbb;
}
</style>