<script setup>
import { computed } from 'vue'

const props = defineProps({
  characters: {
    type: Array,
    default: () => []
  },
  activeCharacterId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add'])

const activeCharacter = computed(() => {
  return props.characters.find(c => c.id === props.activeCharacterId)
})

function selectCharacter(character) {
  emit('select', character.id)
}

function addCharacter() {
  emit('add')
}
</script>

<template>
  <aside class="character-sidebar">
    <div class="sidebar-header">
      <h2>人物列表</h2>
      <button class="add-btn" @click="addCharacter" title="添加人物">
        <span>+</span>
      </button>
    </div>
    <div class="character-list">
      <div
        v-for="character in characters"
        :key="character.id"
        :class="['character-item', { active: character.id === activeCharacterId }]"
        @click="selectCharacter(character)"
      >
        <div class="character-avatar">
          {{ character.name.charAt(0) }}
        </div>
        <div class="character-info">
          <div class="character-name">{{ character.name }}</div>
          <div class="character-level">Lv.{{ character.level }}</div>
        </div>
      </div>
      <div v-if="characters.length === 0" class="empty-tip">
        <p>暂无人物</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.character-sidebar {
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

.character-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.character-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #f8f9fa;
}

.character-item:hover {
  background-color: #f0f7ff;
}

.character-item.active {
  background-color: #e6f7ff;
  border: 1px solid #1890ff;
}

.character-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  margin-right: 12px;
  flex-shrink: 0;
}

.character-info {
  flex: 1;
  min-width: 0;
}

.character-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.character-level {
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
</style>