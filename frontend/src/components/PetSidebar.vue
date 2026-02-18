<script setup>
import { computed } from 'vue'

const props = defineProps({
  pets: {
    type: Array,
    default: () => []
  },
  activePetId: {
    type: [String, Number],
    default: null
  }
})

const emit = defineEmits(['select', 'add', 'delete'])

const activePet = computed(() => {
  return props.pets.find(p => p.id === props.activePetId)
})

function selectPet(pet) {
  emit('select', pet.id)
}

function addPet() {
  emit('add')
}

function deletePet(petId, event) {
  event.stopPropagation()
  emit('delete', petId)
}
</script>

<template>
  <aside class="pet-sidebar">
    <div class="sidebar-header">
      <h2>宠物列表</h2>
      <button class="add-btn" @click="addPet" title="添加宠物">
        <span>+</span>
      </button>
    </div>
    <div class="pet-list">
      <div
        v-for="pet in pets"
        :key="pet.id"
        :class="['pet-item', { active: pet.id === activePetId }]"
        @click="selectPet(pet)"
      >
        <div class="pet-avatar">
          {{ pet.name.charAt(0) }}
        </div>
        <div class="pet-info">
          <div class="pet-name">{{ pet.name }}</div>
          <div class="pet-level">Lv.{{ pet.level }}</div>
        </div>
        <button
          class="delete-pet-btn"
          @click="deletePet(pet.id, $event)"
          title="删除宠物"
        >
          ×
        </button>
      </div>
      <div v-if="pets.length === 0" class="empty-tip">
        <p>暂无宠物</p>
        <p class="tip">点击 + 添加</p>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.pet-sidebar {
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

.pet-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.pet-item {
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

.pet-item:hover {
  background-color: #f0f7ff;
}

.pet-item:hover .delete-pet-btn {
  opacity: 1;
}

.pet-item.active {
  background-color: #e6f7ff;
  border: 1px solid #1890ff;
}

.pet-avatar {
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

.pet-info {
  flex: 1;
  min-width: 0;
}

.pet-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pet-level {
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

.delete-pet-btn {
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

.delete-pet-btn:hover {
  background-color: #fff1f0;
  color: #ff4d4f;
}
</style>