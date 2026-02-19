<script setup>
import { computed } from 'vue'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  updateInfo: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close'])

const formattedPublishedAt = computed(() => {
  const value = props.updateInfo?.published_at
  if (!value) {
    return '-'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
})

const releaseNotes = computed(() => {
  const notes = props.updateInfo?.release_notes || ''
  return notes.trim() || 'No release notes.'
})

function handleClose() {
  emit('close')
}

function handleOpenRelease() {
  const releaseURL = props.updateInfo?.release_url
  if (releaseURL) {
    BrowserOpenURL(releaseURL)
  }
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click="handleClose">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <h3>发现新版本</h3>
            <button class="close-btn" @click="handleClose">x</button>
          </div>

          <div class="modal-body">
            <p class="row">
              <span class="label">当前版本:</span>
              <span class="value">{{ updateInfo?.current_version || '-' }}</span>
            </p>
            <p class="row">
              <span class="label">最新版本:</span>
              <span class="value">{{ updateInfo?.latest_version || '-' }}</span>
            </p>
            <p class="row">
              <span class="label">发布时间:</span>
              <span class="value">{{ formattedPublishedAt }}</span>
            </p>
            <p class="row">
              <span class="label">发布名称:</span>
              <span class="value">{{ updateInfo?.release_name || '-' }}</span>
            </p>

            <div class="notes-wrapper">
              <div class="label">更新说明:</div>
              <pre class="notes">{{ releaseNotes }}</pre>
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn btn-cancel" @click="handleClose">稍后提醒</button>
            <button class="btn btn-confirm" @click="handleOpenRelease">立即更新</button>
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
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.modal-container {
  width: 100%;
  max-width: 640px;
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.modal-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 18px;
}

.close-btn {
  border: none;
  background: none;
  cursor: pointer;
  color: #8c92a0;
  font-size: 16px;
  width: 28px;
  height: 28px;
  border-radius: 6px;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #2c3e50;
}

.modal-body {
  padding: 18px 20px;
}

.row {
  margin-bottom: 8px;
  display: flex;
  gap: 8px;
}

.label {
  color: #8c92a0;
  min-width: 72px;
}

.value {
  color: #2c3e50;
  font-weight: 500;
  word-break: break-all;
}

.notes-wrapper {
  margin-top: 12px;
}

.notes {
  margin-top: 6px;
  padding: 10px 12px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e8eaed;
  max-height: 220px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
  color: #34495e;
  line-height: 1.5;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid #f0f2f5;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 8px 14px;
  font-size: 14px;
  cursor: pointer;
}

.btn-cancel {
  background-color: #f3f4f6;
  color: #475569;
}

.btn-cancel:hover {
  background-color: #e5e7eb;
}

.btn-confirm {
  background-color: #1677ff;
  color: #ffffff;
}

.btn-confirm:hover {
  background-color: #0958d9;
}
</style>
