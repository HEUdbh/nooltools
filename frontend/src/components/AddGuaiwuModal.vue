<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const name = ref('')
const type = ref('')
const level = ref('1')
const health = ref('100')
const attack = ref('10')
const defense = ref('10')
const rewards = ref('')

watch(() => props.visible, (newVal) => {
  if (!newVal) {
    // 关闭时清空表单
    name.value = ''
    type.value = ''
    level.value = '1'
    health.value = '100'
    attack.value = '10'
    defense.value = '10'
    rewards.value = ''
  }
})

function handleConfirm() {
  if (!name.value.trim()) {
    alert('请输入怪物名称')
    return
  }
  if (!level.value || level.value <= 0) {
    alert('请输入有效的怪物等级')
    return
  }

  const levelNum = parseInt(level.value) || 1
  const healthNum = parseInt(health.value) || 100
  const attackNum = parseInt(attack.value) || 10
  const defenseNum = parseInt(defense.value) || 10

  emit('confirm', {
    name: name.value.trim(),
    type: type.value.trim(),
    level: levelNum,
    health: healthNum,
    attack: attackNum,
    defense: defenseNum,
    rewards: rewards.value.trim()
  })

  emit('update:visible', false)
}

function handleCancel() {
  emit('update:visible', false)
}

function handleClose() {
  emit('update:visible', false)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click="handleClose">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <h3>添加怪物</h3>
            <button class="close-btn" @click="handleClose">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>怪物名称 <span class="required">*</span></label>
              <input
                v-model="name"
                type="text"
                placeholder="请输入怪物名称"
                maxlength="50"
              />
            </div>
            <div class="form-group">
              <label>怪物类型</label>
              <input
                v-model="type"
                type="text"
                placeholder="请输入怪物类型（可选）"
                maxlength="50"
              />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>等级 <span class="required">*</span></label>
                <input
                  v-model="level"
                  type="number"
                  placeholder="1"
                  min="1"
                />
              </div>
              <div class="form-group">
                <label>生命值</label>
                <input
                  v-model="health"
                  type="number"
                  placeholder="100"
                  min="1"
                />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>攻击</label>
                <input
                  v-model="attack"
                  type="number"
                  placeholder="10"
                  min="0"
                />
              </div>
              <div class="form-group">
                <label>防御</label>
                <input
                  v-model="defense"
                  type="number"
                  placeholder="10"
                  min="0"
                />
              </div>
            </div>
            <div class="form-group">
              <label>奖励</label>
              <input
                v-model="rewards"
                type="text"
                placeholder="请输入奖励内容（可选）"
                maxlength="200"
              />
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-cancel" @click="handleCancel">取消</button>
            <button class="btn btn-confirm" @click="handleConfirm">创建</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--app-overlay-mask);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-container {
  background-color: var(--app-surface);
  border-radius: 12px;
  box-shadow: var(--app-shadow-lg);
  width: 100%;
  max-width: 480px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  color: var(--app-text-muted);
  font-size: 28px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
}

.close-btn:hover {
  background-color: var(--app-item-bg);
  color: var(--app-text-primary);
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--app-text-primary);
  margin-bottom: 8px;
}

.required {
  color: var(--app-danger);
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--app-border);
  border-radius: 6px;
  font-size: 14px;
  color: var(--app-text-primary);
  transition: all 0.2s ease;
  box-sizing: border-box;
  font-family: inherit;
}

.form-group input:focus {
  outline: none;
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-focus);
}

.form-group input::placeholder {
  color: var(--app-text-muted);
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--app-divider-soft);
  background-color: var(--app-bg);
}

.btn {
  padding: 8px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel {
  background-color: var(--app-surface);
  color: var(--app-text-primary);
  border: 1px solid var(--app-border);
}

.btn-cancel:hover {
  background-color: var(--app-item-bg);
  border-color: var(--app-text-muted);
}

.btn-confirm {
  background-color: var(--app-accent);
  color: var(--app-text-inverse);
}

.btn-confirm:hover {
  background-color: var(--app-accent-hover);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.2s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.95);
}
</style>