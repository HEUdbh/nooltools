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
  background-color: var(--app-surface);
  border-right: 1px solid var(--app-border);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--app-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--app-text-primary);
}

.add-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: var(--app-success);
  color: var(--app-text-inverse);
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
  background-color: var(--app-success-hover);
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
  background-color: var(--app-bg);
  border: 1px solid var(--app-border);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.shopping-item:hover {
  background-color: var(--app-divider-soft);
  border-color: var(--app-border);
}

.shopping-item.active {
  background-color: var(--app-success-soft-bg);
  border-color: var(--app-success);
}

.shopping-name {
  font-weight: 500;
  color: var(--app-text-primary);
  margin-bottom: 4px;
  font-size: 14px;
}

.shopping-value {
  font-size: 12px;
  color: var(--app-text-secondary);
}

.delete-item-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  border: none;
  background-color: var(--app-danger);
  color: var(--app-text-inverse);
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
  background-color: var(--app-danger-hover);
  transform: scale(1.1);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--app-text-muted);
}

.empty-state p {
  margin: 5px 0;
}

.empty-state .tip {
  font-size: 12px;
  color: var(--app-text-muted);
}
</style>
