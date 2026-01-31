<script setup>
import { computed } from 'vue'

const props = defineProps({
  guaiwus: {
    type: Array,
    default: () => []
  },
  activeGuaiwuId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete'])

const activeGuaiwu = computed(() => {
  return props.guaiwus.find(g => g.id === props.activeGuaiwuId)
})

function selectGuaiwu(guaiwu) {
  emit('select', guaiwu.id)
}

function addGuaiwu() {
  emit('add')
}

function deleteGuaiwu(guaiwuId, event) {
  event.stopPropagation()
  emit('delete', guaiwuId)
}
</script>

<template>
  <aside class="guaiwu-sidebar">
    <div class="sidebar-header">
      <h2>怪物列表</h2>
      <button class="add-btn" @click="addGuaiwu" title="添加怪物">
        <span>+</span>
      </button>
    </div>
    <div class="guaiwu-list">
      <div
        v-for="guaiwu in guaiwus"
        :key="guaiwu.id"
        :class="['guaiwu-item', { active: guaiwu.id === activeGuaiwuId }]"
        @click="selectGuaiwu(guaiwu)"
      >
        <div class="guaiwu-avatar">
          {{ guaiwu.name.charAt(0) }}
        </div>
        <div class="guaiwu-info">
          <div class="guaiwu-name">{{ guaiwu.name }}</div>
          <div class="guaiwu-level">Lv.{{ guaiwu.level }}</div>
        </div>
        <button
          class="delete-guaiwu-btn"
          @click="deleteGuaiwu(guaiwu.id, $event)"
          title="删除怪物"
        >
          ×
        </button>
      </div>
      <div v-if="guaiwus.length === 0" class="empty-tip">
        <p>暂无怪物</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.guaiwu-sidebar {
  width: 200px;
  background-color: #ffffff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-header h2 {
  font-size: 18px;
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
  border-radius: 50%;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background-color: #40a9ff;
  transform: scale(1.1);
}

.guaiwu-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.guaiwu-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #f8f9fa;
  position: relative;
}

.guaiwu-item:hover {
  background-color: #f0f7ff;
}

.guaiwu-item:hover .delete-guaiwu-btn {
  opacity: 1;
}

.guaiwu-item.active {
  background-color: #e6f7ff;
  border: 1px solid #1890ff;
}

.guaiwu-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
}

.guaiwu-info {
  flex: 1;
  min-width: 0;
}

.guaiwu-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.guaiwu-level {
  font-size: 12px;
  color: #8c92a0;
}

.empty-tip {
  text-align: center;
  padding: 40px 20px;
  color: #8c92a0;
}

.empty-tip p {
  margin: 8px 0;
}

.tip {
  font-size: 12px;
  color: #b0b8c3;
}

.delete-guaiwu-btn {
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
  opacity: 0;
  flex-shrink: 0;
}

.delete-guaiwu-btn:hover {
  background-color: #fff1f0;
  color: #ff4d4f;
}
</style>