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
  },
  isUpdating: {
    type: Boolean,
    default: false
  },
  progress: {
    type: Object,
    default: () => ({
      stage: '',
      percent: 0,
      message: '',
      detail: ''
    })
  }
})

const emit = defineEmits(['close', 'start-update'])
const TEXT = {
  title: '发现新版本',
  currentVersion: '当前版本:',
  latestVersion: '最新版本:',
  publishedAt: '发布时间:',
  releaseName: '发布名称:',
  assetName: '资产文件:',
  releaseNotes: '更新说明:',
  autoUpdateDisabled: '当前版本不支持自动替换，请去下载页手动更新。',
  progressDefault: '更新进行中...',
  remindLater: '稍后提醒',
  updateNow: '立即更新',
  updating: '更新中...',
  goDownload: '去下载页'
}

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

const canAutoUpdate = computed(() => {
  return Boolean(props.updateInfo?.can_auto_update)
})

const progressPercent = computed(() => {
  const value = Number(props.progress?.percent || 0)
  if (value < 0) {
    return 0
  }
  if (value > 100) {
    return 100
  }
  return value
})

const hasProgress = computed(() => {
  return Boolean(props.progress?.stage)
})

const isError = computed(() => {
  return props.progress?.stage === 'error'
})

function handleClose() {
  if (props.isUpdating) {
    return
  }
  emit('close')
}

function openReleasePage() {
  const releaseURL = props.updateInfo?.release_url
  if (releaseURL) {
    BrowserOpenURL(releaseURL)
  }
}

function handleFallbackDownload() {
  openReleasePage()
  emit('close')
}

function handleStartUpdate() {
  emit('start-update')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click="handleClose">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <h3>{{ TEXT.title }}</h3>
            <button class="close-btn" :disabled="isUpdating" @click="handleClose">x</button>
          </div>

          <div class="modal-body">
            <p class="row">
              <span class="label">{{ TEXT.currentVersion }}</span>
              <span class="value">{{ updateInfo?.current_version || '-' }}</span>
            </p>
            <p class="row">
              <span class="label">{{ TEXT.latestVersion }}</span>
              <span class="value">{{ updateInfo?.latest_version || '-' }}</span>
            </p>
            <p class="row">
              <span class="label">{{ TEXT.publishedAt }}</span>
              <span class="value">{{ formattedPublishedAt }}</span>
            </p>
            <p class="row">
              <span class="label">{{ TEXT.releaseName }}</span>
              <span class="value">{{ updateInfo?.release_name || '-' }}</span>
            </p>
            <p class="row">
              <span class="label">{{ TEXT.assetName }}</span>
              <span class="value">{{ updateInfo?.asset_name || '-' }}</span>
            </p>

            <div class="notes-wrapper">
              <div class="label">{{ TEXT.releaseNotes }}</div>
              <pre class="notes">{{ releaseNotes }}</pre>
            </div>

            <div v-if="!canAutoUpdate" class="auto-update-disabled">
              {{ updateInfo?.auto_update_reason || TEXT.autoUpdateDisabled }}
            </div>

            <div v-if="hasProgress" class="progress-wrapper">
              <div class="progress-header">
                <span>{{ progress?.message || TEXT.progressDefault }}</span>
                <span>{{ progressPercent }}%</span>
              </div>
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: `${progressPercent}%` }"></div>
              </div>
              <p v-if="progress?.detail" :class="['progress-detail', { error: isError }]">
                {{ progress.detail }}
              </p>
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn btn-cancel" :disabled="isUpdating" @click="handleClose">{{ TEXT.remindLater }}</button>
            <button
              v-if="canAutoUpdate"
              class="btn btn-confirm"
              :disabled="isUpdating"
              @click="handleStartUpdate"
            >
              {{ isUpdating ? TEXT.updating : TEXT.updateNow }}
            </button>
            <button v-else class="btn btn-confirm" @click="handleFallbackDownload">{{ TEXT.goDownload }}</button>
            <button
              v-if="isError && canAutoUpdate"
              class="btn btn-fallback"
              @click="openReleasePage"
            >
              {{ TEXT.goDownload }}
            </button>
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

.close-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.auto-update-disabled {
  margin-top: 12px;
  padding: 8px 10px;
  background: #fff7e6;
  border: 1px solid #ffe7ba;
  color: #ad6800;
  border-radius: 6px;
}

.progress-wrapper {
  margin-top: 14px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  color: #334155;
  margin-bottom: 6px;
}

.progress-track {
  width: 100%;
  height: 10px;
  background: #edf2f7;
  border-radius: 999px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #1677ff, #36cfc9);
  transition: width 0.2s ease;
}

.progress-detail {
  margin-top: 8px;
  color: #64748b;
  word-break: break-word;
}

.progress-detail.error {
  color: #d4380d;
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

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-fallback {
  background-color: #fff1f0;
  color: #cf1322;
}

.btn-fallback:hover {
  background-color: #ffccc7;
}
</style>
