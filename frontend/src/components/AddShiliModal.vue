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
const founder = ref('')
const level = ref('1')
const wealth = ref('0')
const maxMembers = ref('10')

watch(() => props.visible, (newVal) => {
  if (!newVal) {
    // 关闭时清空表单
    name.value = ''
    founder.value = ''
    level.value = '1'
    wealth.value = '0'
    maxMembers.value = '10'
  }
})

function handleConfirm() {
  if (!name.value.trim()) {
    alert('请输入势力名称')
    return
  }
  if (!founder.value.trim()) {
    alert('请输入创始人')
    return
  }
  if (!level.value || level.value <= 0) {
    alert('请输入有效的势力等级')
    return
  }
  if (!maxMembers.value || maxMembers.value <= 0) {
    alert('请输入有效的容纳人数')
    return
  }

  const levelNum = parseInt(level.value) || 1
  const wealthNum = parseInt(wealth.value) || 0
  const maxMembersNum = parseInt(maxMembers.value) || 10

  emit('confirm', {
    name: name.value.trim(),
    founder: founder.value.trim(),
    level: levelNum,
    wealth: wealthNum,
    maxMembers: maxMembersNum
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
            <h3>添加势力</h3>
            <button class="close-btn" @click="handleClose">×</button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>势力名称 <span class="required">*</span></label>
              <input
                v-model="name"
                type="text"
                placeholder="请输入势力名称"
                maxlength="50"
              />
            </div>
            <div class="form-group">
              <label>创始人 <span class="required">*</span></label>
              <input
                v-model="founder"
                type="text"
                placeholder="请输入创始人"
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
                <label>财富</label>
                <input
                  v-model="wealth"
                  type="number"
                  placeholder="0"
                  min="0"
                />
              </div>
            </div>
            <div class="form-group">
              <label>容纳人数 <span class="required">*</span></label>
              <input
                v-model="maxMembers"
                type="number"
                placeholder="10"
                min="1"
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
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-container {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 480px;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: transparent;
  color: #8c92a0;
  font-size: 24px;
  cursor: pointer;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background-color: #f0f2f5;
  color: #2c3e50;
}

.modal-body {
  padding: 24px;
  max-height: 60vh;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 8px;
}

.required {
  color: #ff4d4f;
  margin-left: 4px;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.form-group input::placeholder {
  color: #bfbfbf;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.form-row .form-group {
  flex: 1;
  margin-bottom: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f2f5;
  background-color: #fafafa;
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
  background-color: white;
  color: #595959;
  border: 1px solid #d9d9d9;
}

.btn-cancel:hover {
  background-color: #fafafa;
  border-color: #bfbfbf;
}

.btn-confirm {
  background-color: #1890ff;
  color: white;
}

.btn-confirm:hover {
  background-color: #40a9ff;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9);
}
</style>