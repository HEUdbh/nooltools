<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  beibao: {
    type: Array,
    default: () => []
  },
  activeBeibaoId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete', 'update'])

const editingId = ref(null)
const editingName = ref('')

const activeBeibao = computed(() => {
  return props.beibao.find(b => b.id === props.activeBeibaoId)
})

function selectBeibao(b) {
  emit('select', b.id)
}

function addBeibao() {
  emit('add')
}

function deleteBeibao(beibaoId, event) {
  event.stopPropagation()
  emit('delete', beibaoId)
}

function startEdit(beibaoId, currentName, event) {
  event.stopPropagation()
  editingId.value = beibaoId
  editingName.value = currentName
}

function cancelEdit() {
  editingId.value = null
  editingName.value = ''
}

function saveEdit(beibaoId, event) {
  event.stopPropagation()
  if (editingName.value.trim()) {
    emit('update', beibaoId, editingName.value.trim())
  }
  editingId.value = null
  editingName.value = ''
}

function handleKeydown(event, beibaoId) {
  if (event.key === 'Enter') {
    saveEdit(beibaoId, event)
  } else if (event.key === 'Escape') {
    cancelEdit()
  }
}
</script>

<template>
  <aside class="beibao-sidebar">
    <div class="sidebar-header">
      <h2>背包列表</h2>
      <button class="add-btn" @click="addBeibao" title="添加背包">
        <span>+</span>
      </button>
    </div>
    <div class="beibao-list">
      <div
        v-for="b in beibao"
        :key="b.id"
        :class="['beibao-item', { active: b.id === activeBeibaoId }]"
        @click="selectBeibao(b)"
      >
        <div class="beibao-avatar">
          {{ b.name.charAt(0) }}
        </div>
        <div class="beibao-info">
          <div v-if="editingId !== b.id" class="beibao-name" @dblclick="startEdit(b.id, b.name, $event)">
            {{ b.name }}
          </div>
          <input
            v-else
            v-model="editingName"
            class="edit-input"
            @click.stop
            @blur="saveEdit(b.id, $event)"
            @keydown="handleKeydown($event, b.id)"
            ref="editInput"
          />
        </div>
        <button
          v-if="editingId !== b.id"
          class="delete-beibao-btn"
          @click="deleteBeibao(b.id, $event)"
          title="删除背包"
        >
          ×
        </button>
      </div>
      <div v-if="beibao.length === 0" class="empty-tip">
        <p>暂无背包</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.beibao-sidebar {
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

.beibao-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.beibao-item {
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

.beibao-item:hover {
  background-color: #f0f7ff;
}

.beibao-item:hover .delete-beibao-btn {
  opacity: 1;
}

.beibao-item.active {
  background-color: #e6f7ff;
  border: 1px solid #1890ff;
}

.beibao-avatar {
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

.beibao-info {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
}

.beibao-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  cursor: text;
}

.beibao-name:hover {
  color: #1890ff;
}

.edit-input {
  width: 100%;
  padding: 4px 8px;
  border: 1px solid #1890ff;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
  background-color: white;
}

.edit-input:focus {
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
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

.delete-beibao-btn {
  width: 24px;
  height: 24px;
  border: none;
  background-color: #ff4d4f;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  opacity: 0;
  transition: all 0.2s ease;
  margin-left: 8px;
  flex-shrink: 0;
}

.delete-beibao-btn:hover {
  background-color: #ff7875;
  transform: scale(1.1);
}
</style>