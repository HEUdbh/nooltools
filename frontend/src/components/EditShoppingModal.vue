<template>
  <div v-if="visible" class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>修改商品</h3>
        <button class="close-btn" @click="handleCancel">×</button>
      </div>
      
      <div class="modal-body">
        <div class="form-group">
          <label>商品名称</label>
          <input
            v-model="formData.name"
            type="text"
            placeholder="请输入商品名称"
            @keyup.enter="handleConfirm"
          />
        </div>
        
        <div class="form-group">
          <label>商品价值</label>
          <input
            v-model.number="formData.value"
            type="number"
            placeholder="请输入商品价值"
            @keyup.enter="handleConfirm"
          />
        </div>
        
        <div class="form-group">
          <label>商品描述</label>
          <textarea
            v-model="formData.description"
            placeholder="请输入商品描述"
            rows="3"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label>购买条件</label>
          <textarea
            v-model="formData.condition"
            placeholder="请输入购买条件"
            rows="2"
          ></textarea>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="cancel-btn" @click="handleCancel">取消</button>
        <button class="confirm-btn" @click="handleConfirm">确认</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  shopping: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const formData = ref({
  id: null,
  name: '',
  value: 0,
  description: '',
  condition: ''
})

// 监听 shopping 变化，更新表单数据
watch(() => props.shopping, (newVal) => {
  if (newVal) {
    formData.value = {
      id: newVal.id,
      name: newVal.name || '',
      value: newVal.value || 0,
      description: newVal.description || '',
      condition: newVal.condition || ''
    }
  }
}, { immediate: true })

function handleOverlayClick() {
  handleCancel()
}

function handleCancel() {
  emit('update:visible', false)
}

function handleConfirm() {
  if (!formData.value.name) {
    alert('请输入商品名称')
    return
  }
  
  if (formData.value.value === undefined || formData.value.value === null || formData.value.value === '') {
    alert('请输入商品价值')
    return
  }
  
  emit('confirm', { ...formData.value })
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
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  animation: slideIn 0.2s;
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
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background-color: transparent;
  color: #999;
  border-radius: 50%;
  cursor: pointer;
  font-size: 24px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: #f5f5f5;
  color: #666;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #4CAF50;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.modal-footer {
  padding: 20px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.cancel-btn,
.confirm-btn {
  padding: 10px 24px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.cancel-btn {
  background-color: #f5f5f5;
  color: #666;
}

.cancel-btn:hover {
  background-color: #e0e0e0;
}

.confirm-btn {
  background-color: #4CAF50;
  color: white;
}

.confirm-btn:hover {
  background-color: #45a049;
}
</style>