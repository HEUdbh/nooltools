<script setup>
import { computed } from 'vue'

const props = defineProps({
  daoju: {
    type: Array,
    default: () => []
  },
  activeDaojuId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete'])

const activeDaoju = computed(() => {
  return props.daoju.find(d => d.id === props.activeDaojuId)
})

function selectDaoju(d) {
  emit('select', d.id)
}

function addDaoju() {
  emit('add')
}

function deleteDaoju(daojuId, event) {
  event.stopPropagation()
  emit('delete', daojuId)
}
</script>

<template>
  <aside class="daoju-sidebar">
    <div class="sidebar-header">
      <h2>道具列表</h2>
      <button class="add-btn" @click="addDaoju" title="添加道具">
        <span>+</span>
      </button>
    </div>
    <div class="daoju-list">
      <div
        v-for="d in daoju"
        :key="d.id"
        :class="['daoju-item', { active: d.id === activeDaojuId }]"
        @click="selectDaoju(d)"
      >
        <div class="daoju-avatar">
          {{ d.name.charAt(0) }}
        </div>
        <div class="daoju-info">
          <div class="daoju-name">{{ d.name }}</div>
          <div class="daoju-level">Lv.{{ d.level }}</div>
        </div>
        <button
          class="delete-daoju-btn"
          @click="deleteDaoju(d.id, $event)"
          title="删除道具"
        >
          ×
        </button>
      </div>
      <div v-if="daoju.length === 0" class="empty-tip">
        <p>暂无道具</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.daoju-sidebar {
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

.daoju-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.daoju-item {
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

.daoju-item:hover {
  background-color: var(--app-hover-bg);
}

.daoju-item:hover .delete-daoju-btn {
  opacity: 1;
}

.daoju-item.active {
  background-color: var(--app-active-bg);
  border: 1px solid var(--app-accent);
}

.daoju-avatar {
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

.daoju-info {
  flex: 1;
  min-width: 0;
}

.daoju-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--app-text-primary);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.daoju-level {
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

.delete-daoju-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  color: var(--app-text-muted);
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

.delete-daoju-btn:hover {
  background-color: var(--app-danger-soft-bg);
  color: var(--app-danger);
}
</style>