<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleCancel">
    <div class="modal-container">
      <div class="modal-header">
        <h3>添加职务</h3>
        <button class="close-btn" @click="handleCancel">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label for="positionName">职务名称 *</label>
          <input
            id="positionName"
            v-model="formData.positionName"
            type="text"
            placeholder="请输入职务名称"
            @keyup.enter="handleConfirm"
          />
        </div>
        <div class="form-group">
          <label for="personName">姓名 *</label>
          <input
            id="personName"
            v-model="formData.personName"
            type="text"
            placeholder="请输入姓名"
            @keyup.enter="handleConfirm"
          />
        </div>
        <div class="form-group">
          <label for="description">描述</label>
          <textarea
            id="description"
            v-model="formData.description"
            placeholder="请输入职务描述（可选）"
            rows="3"
          ></textarea>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-cancel" @click="handleCancel">取消</button>
        <button class="btn btn-confirm" @click="handleConfirm" :disabled="!isFormValid">
          确定
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const formData = ref({
  positionName: '',
  personName: '',
  description: ''
})

// 表单验证
const isFormValid = computed(() => {
  return formData.value.positionName.trim() !== '' && formData.value.personName.trim() !== ''
})

// 监听visible变化，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    formData.value = {
      positionName: '',
      personName: '',
      description: ''
    }
  }
})

// 处理确认
function handleConfirm() {
  if (!isFormValid.value) return
  
  emit('confirm', {
    positionName: formData.value.positionName.trim(),
    personName: formData.value.personName.trim(),
    description: formData.value.description.trim()
  })
  
  emit('update:visible', false)
}

// 处理取消
function handleCancel() {
  emit('update:visible', false)
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-container {
  background-color: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  font-size: 24px;
  color: #8c92a0;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
}

.close-btn:hover {
  background-color: #f5f7fa;
  color: #2c3e50;
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

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  color: #2c3e50;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: #b0b8c3;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f2f5;
}

.btn {
  padding: 8px 20px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.btn-cancel {
  background-color: white;
  border-color: #d9d9d9;
  color: #595959;
}

.btn-cancel:hover {
  border-color: #8c92a0;
  color: #2c3e50;
}

.btn-confirm {
  background-color: #1890ff;
  color: white;
}

.btn-confirm:hover:not(:disabled) {
  background-color: #40a9ff;
}

.btn-confirm:disabled {
  background-color: #d9d9d9;
  color: #8c92a0;
  cursor: not-allowed;
}
</style>