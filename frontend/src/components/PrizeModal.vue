<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleClose">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ isEdit ? '编辑奖品' : '添加奖品' }}</h3>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>

      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label for="name">奖品名称 *</label>
          <input
            id="name"
            v-model="formData.name"
            type="text"
            placeholder="请输入奖品名称"
            required
          />
        </div>

        <div class="form-group">
          <label for="rate">爆率 *</label>
          <input
            id="rate"
            v-model.number="formData.rate"
            type="number"
            step="0.0001"
            min="0"
            max="1"
            placeholder="请输入爆率（0-1之间的小数）"
            required
          />
          <div class="form-hint">例如：0.01 表示 1% 的爆率</div>
        </div>

        <div class="form-group">
          <label for="variety">品种 *</label>
          <select id="variety" v-model="formData.variety" required>
            <option value="普通">普通</option>
            <option value="稀有">稀有</option>
            <option value="史诗">史诗</option>
            <option value="传说">传说</option>
          </select>
        </div>

        <div class="form-group">
          <label for="description">描述</label>
          <textarea
            id="description"
            v-model="formData.description"
            rows="3"
            placeholder="请输入奖品描述（可选）"
          ></textarea>
        </div>

        <div class="modal-footer">
          <button type="button" class="cancel-btn" @click="handleClose">取消</button>
          <button type="submit" class="confirm-btn">确认</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  prize: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'confirm'])

const formData = ref({
  name: '',
  rate: 0.01,
  variety: '普通',
  description: ''
})

const isEdit = computed(() => !!props.prize)

// 监听弹窗显示状态，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.prize) {
      // 编辑模式，填充数据
      formData.value = {
        name: props.prize.name,
        rate: props.prize.rate,
        variety: props.prize.variety,
        description: props.prize.description || ''
      }
    } else {
      // 添加模式，重置表单
      formData.value = {
        name: '',
        rate: 0.01,
        variety: '普通',
        description: ''
      }
    }
  }
})

// 关闭弹窗
function handleClose() {
  emit('update:visible', false)
}

// 提交表单
function handleSubmit() {
  // 验证爆率
  if (formData.value.rate < 0 || formData.value.rate > 1) {
    alert('爆率必须在 0 到 1 之间')
    return
  }

  // 验证奖品名称
  if (!formData.value.name.trim()) {
    alert('请输入奖品名称')
    return
  }

  emit('confirm', { ...formData.value })
  handleClose()
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
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
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 500px;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(30px);
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
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group textarea {
  resize: vertical;
  font-family: inherit;
}

.form-hint {
  font-size: 12px;
  color: #999;
  margin-top: 6px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #e0e0e0;
}

.cancel-btn,
.confirm-btn {
  padding: 10px 24px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.cancel-btn:hover {
  background: #e0e0e0;
}

.confirm-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}
</style>