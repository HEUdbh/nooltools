<script setup>
import { computed } from 'vue'

const props = defineProps({
  shilis: {
    type: Array,
    default: () => []
  },
  activeShiliId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete'])

const activeShili = computed(() => {
  return props.shilis.find(s => s.id === props.activeShiliId)
})

function selectShili(shili) {
  emit('select', shili.id)
}

function addShili() {
  emit('add')
}

function deleteShili(shiliId, event) {
  event.stopPropagation()
  emit('delete', shiliId)
}
</script>

<template>
  <aside class="shili-sidebar">
    <div class="sidebar-header">
      <h2>势力列表</h2>
      <button class="add-btn" @click="addShili" title="添加势力">
        <span>+</span>
      </button>
    </div>
    <div class="shili-list">
      <div
        v-for="shili in shilis"
        :key="shili.id"
        :class="['shili-item', { active: shili.id === activeShiliId }]"
        @click="selectShili(shili)"
      >
        <div class="shili-avatar">
          {{ shili.name.charAt(0) }}
        </div>
        <div class="shili-info">
          <div class="shili-name">{{ shili.name }}</div>
          <div class="shili-level">Lv.{{ shili.level }}</div>
        </div>
        <button
          class="delete-shili-btn"
          @click="deleteShili(shili.id, $event)"
          title="删除势力"
        >
          ×
        </button>
      </div>
      <div v-if="shilis.length === 0" class="empty-tip">
        <p>暂无势力</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.shili-sidebar {
  width: 200px;
  background-color: var(--app-surface);
  border-right: 1px solid var(--app-border);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--app-divider-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin: 0;
}

.add-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: var(--app-accent);
  color: var(--app-text-inverse);
  border-radius: 50%;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background-color: var(--app-accent-hover);
  transform: scale(1.1);
}

.shili-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.shili-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: var(--app-item-bg);
  position: relative;
}

.shili-item:hover {
  background-color: var(--app-hover-bg);
}

.shili-item:hover .delete-shili-btn {
  opacity: 1;
}

.shili-item.active {
  background-color: var(--app-active-bg);
  border: 1px solid var(--app-accent);
}

.shili-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--app-gradient-primary);
  color: var(--app-text-inverse);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
}

.shili-info {
  flex: 1;
  min-width: 0;
}

.shili-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--app-text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.shili-level {
  font-size: 12px;
  color: var(--app-text-muted);
}

.empty-tip {
  text-align: center;
  padding: 40px 20px;
  color: var(--app-text-muted);
}

.empty-tip p {
  margin: 8px 0;
}

.tip {
  font-size: 12px;
  color: var(--app-text-muted);
}

.delete-shili-btn {
  width: 24px;
  height: 24px;
  border: none;
  background-color: var(--app-danger);
  color: var(--app-text-inverse);
  border-radius: 50%;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: all 0.2s ease;
  position: absolute;
  right: 8px;
}

.delete-shili-btn:hover {
  background-color: var(--app-danger-hover);
  transform: scale(1.1);
}
</style>