<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const STORAGE_CHILD_DIR = 'nooltools_data'

const activeTab = ref('about')
const loadingSettings = ref(false)
const selectingDirectory = ref(false)
const migrating = ref(false)
const restarting = ref(false)
const selectedParentDir = ref('')
const currentSettings = ref({
  current_data_dir: '',
  default_data_dir: '',
  is_custom: false
})
const migrationResult = ref(null)
const errorMessage = ref('')
const successMessage = ref('')

const targetDirPreview = computed(() => {
  if (!selectedParentDir.value) {
    return ''
  }
  const parent = selectedParentDir.value
  if (parent.endsWith('\\') || parent.endsWith('/')) {
    return `${parent}${STORAGE_CHILD_DIR}`
  }
  const separator = parent.includes('\\') ? '\\' : '/'
  return `${parent}${separator}${STORAGE_CHILD_DIR}`
})

watch(
  () => props.visible,
  (visible) => {
    if (!visible) {
      return
    }
    resetState()
    loadStorageSettings()
  }
)

function resetState() {
  activeTab.value = 'about'
  selectedParentDir.value = ''
  migrationResult.value = null
  errorMessage.value = ''
  successMessage.value = ''
}

function handleClose() {
  if (migrating.value || restarting.value) {
    return
  }
  emit('close')
}

async function loadStorageSettings() {
  loadingSettings.value = true
  errorMessage.value = ''
  try {
    const result = await window.go.main.app.GetStorageSettings()
    currentSettings.value = result || currentSettings.value
  } catch (error) {
    errorMessage.value = `读取存储设置失败: ${error?.message || String(error)}`
  } finally {
    loadingSettings.value = false
  }
}

async function handleSelectDirectory() {
  selectingDirectory.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const selectedDir = await window.go.main.app.SelectStorageParentDirectory()
    if (selectedDir) {
      selectedParentDir.value = selectedDir
    }
  } catch (error) {
    errorMessage.value = `选择目录失败: ${error?.message || String(error)}`
  } finally {
    selectingDirectory.value = false
  }
}

async function handleMigrateStorage() {
  if (!selectedParentDir.value) {
    errorMessage.value = '请先选择存储目录'
    return
  }

  migrating.value = true
  errorMessage.value = ''
  successMessage.value = ''
  migrationResult.value = null
  try {
    const result = await window.go.main.app.MigrateStorageDirectory(selectedParentDir.value)
    migrationResult.value = result
    await loadStorageSettings()
    successMessage.value = `迁移完成，新目录: ${result?.to_dir || ''}`
  } catch (error) {
    errorMessage.value = `迁移失败: ${error?.message || String(error)}`
  } finally {
    migrating.value = false
  }
}

async function handleRestartNow() {
  restarting.value = true
  errorMessage.value = ''
  try {
    await window.go.main.app.RestartApplication()
  } catch (error) {
    restarting.value = false
    errorMessage.value = `重启失败: ${error?.message || String(error)}`
  }
}

function handleRestartLater() {
  if (migrationResult.value) {
    migrationResult.value.restart_recommended = false
  }
  successMessage.value = '已应用新目录，你可以稍后手动重启应用。'
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click="handleClose">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <h3>设置</h3>
            <button class="close-btn" :disabled="migrating || restarting" @click="handleClose">x</button>
          </div>

          <div class="modal-body">
            <div class="tabs">
              <button
                :class="['tab-btn', { active: activeTab === 'about' }]"
                @click="activeTab = 'about'"
              >
                关于作者
              </button>
              <button
                :class="['tab-btn', { active: activeTab === 'storage' }]"
                @click="activeTab = 'storage'"
              >
                存储设置
              </button>
            </div>

            <div v-if="activeTab === 'about'" class="tab-content">
              <div class="contact-card">
                <div class="contact-item">
                  <div class="contact-icon">📧</div>
                  <div class="contact-info">
                    <h4>邮箱</h4>
                    <p>heuxry@outlook.com</p>
                  </div>
                </div>
                <div class="contact-item">
                  <div class="contact-icon">💬</div>
                  <div class="contact-info">
                    <h4>反馈建议</h4>
                    <p>欢迎提出宝贵意见和建议</p>
                  </div>
                </div>
                <div class="contact-item">
                  <div class="contact-icon">⭐</div>
                  <div class="contact-info">
                    <h4>支持我们</h4>
                    <p>感谢您的使用和支持，关注微信公众号"夜天炫安全"获取最新动态和资讯</p>
                    <div class="image-container">
                      <img src="/wechat.png" alt="微信公众号" style="width: 320px;">
                      <img src="/get.jpg" alt="赞助码" style="width: 220px;">
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="tab-content">
              <div class="storage-row">
                <span class="label">当前存储目录</span>
                <span class="value">{{ currentSettings.current_data_dir || '-' }}</span>
              </div>
              <div class="storage-row">
                <span class="label">默认存储目录</span>
                <span class="value">{{ currentSettings.default_data_dir || '-' }}</span>
              </div>
              <div class="storage-row">
                <span class="label">当前模式</span>
                <span class="value">{{ currentSettings.is_custom ? '自定义目录' : '默认目录' }}</span>
              </div>

              <div class="storage-actions">
                <button class="btn btn-select" :disabled="selectingDirectory || migrating" @click="handleSelectDirectory">
                  {{ selectingDirectory ? '选择中...' : '选择新目录' }}
                </button>
                <button
                  class="btn btn-migrate"
                  :disabled="!selectedParentDir || migrating"
                  @click="handleMigrateStorage"
                >
                  {{ migrating ? '迁移中...' : '迁移并应用' }}
                </button>
              </div>

              <div class="storage-row">
                <span class="label">已选父目录</span>
                <span class="value">{{ selectedParentDir || '-' }}</span>
              </div>
              <div class="storage-row">
                <span class="label">目标目录</span>
                <span class="value">{{ targetDirPreview || '-' }}</span>
              </div>

              <p class="tip">迁移会复制整个数据目录，旧目录文件会保留不删除。</p>

              <p v-if="errorMessage" class="message error">{{ errorMessage }}</p>
              <p v-if="successMessage" class="message success">{{ successMessage }}</p>

              <div v-if="migrationResult?.backed_up_conflicts?.length" class="conflict-box">
                <p class="conflict-title">已自动备份冲突项:</p>
                <ul>
                  <li v-for="item in migrationResult.backed_up_conflicts" :key="item">{{ item }}</li>
                </ul>
              </div>

              <div v-if="migrationResult?.restart_recommended" class="restart-box">
                <p>迁移成功，建议重启应用以完成环境刷新。</p>
                <div class="restart-actions">
                  <button class="btn btn-restart" :disabled="restarting" @click="handleRestartNow">
                    {{ restarting ? '重启中...' : '立即重启' }}
                  </button>
                  <button class="btn btn-later" :disabled="restarting" @click="handleRestartLater">稍后重启</button>
                </div>
              </div>

              <p v-if="loadingSettings" class="tip">正在读取设置...</p>
            </div>
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
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 2200;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--app-overlay-mask);
  padding: 20px;
}

.modal-container {
  width: 100%;
  max-width: 860px;
  max-height: calc(100vh - 40px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background-color: var(--app-surface);
  border-radius: 12px;
  box-shadow: var(--app-shadow-xl);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.modal-header h3 {
  color: var(--app-text-primary);
  font-size: 18px;
}

.close-btn {
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: transparent;
  color: var(--app-text-muted);
  cursor: pointer;
}

.close-btn:hover {
  background-color: var(--app-item-bg);
  color: var(--app-text-primary);
}

.close-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.modal-body {
  padding: 18px 20px;
  overflow: auto;
}

.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.tab-btn {
  border: none;
  border-radius: 8px;
  padding: 8px 14px;
  cursor: pointer;
  background-color: var(--app-item-bg);
  color: var(--app-text-secondary);
}

.tab-btn.active {
  background-color: var(--app-active-bg);
  color: var(--app-accent);
  font-weight: 500;
}

.tab-content {
  border: 1px solid var(--app-border);
  border-radius: 10px;
  background-color: var(--app-surface-2);
  padding: 16px;
}

.contact-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contact-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  background-color: var(--app-surface);
  border: 1px solid var(--app-border);
  border-radius: 8px;
  padding: 14px;
}

.contact-icon {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--app-item-bg);
  font-size: 24px;
}

.contact-info h4 {
  color: var(--app-text-primary);
  font-size: 16px;
}

.contact-info p {
  color: var(--app-text-secondary);
  margin-top: 4px;
}

.image-container {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.storage-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.label {
  min-width: 92px;
  color: var(--app-text-muted);
}

.value {
  color: var(--app-text-primary);
  word-break: break-all;
}

.storage-actions {
  display: flex;
  gap: 10px;
  margin: 12px 0;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 9px 14px;
  cursor: pointer;
  font-size: 14px;
}

.btn-select,
.btn-later {
  background-color: var(--app-item-bg);
  color: var(--app-text-secondary);
}

.btn-select:hover,
.btn-later:hover {
  background-color: var(--app-hover-bg);
}

.btn-migrate,
.btn-restart {
  background-color: var(--app-info);
  color: var(--app-text-inverse);
}

.btn-migrate:hover,
.btn-restart:hover {
  background-color: var(--app-info-hover);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.tip {
  color: var(--app-text-muted);
  margin-top: 10px;
}

.message {
  margin-top: 10px;
  padding: 8px 10px;
  border-radius: 6px;
}

.message.error {
  color: var(--app-danger-text);
  background-color: var(--app-danger-soft-bg);
  border: 1px solid var(--app-danger-soft-border);
}

.message.success {
  color: var(--app-success);
  background-color: var(--app-success-soft-bg);
  border: 1px solid var(--app-success-soft-border);
}

.conflict-box {
  margin-top: 12px;
  background-color: var(--app-surface);
  border: 1px solid var(--app-border);
  border-radius: 8px;
  padding: 10px;
}

.conflict-title {
  color: var(--app-text-primary);
  margin-bottom: 6px;
}

.conflict-box ul {
  list-style-position: inside;
  max-height: 140px;
  overflow: auto;
  color: var(--app-text-secondary);
}

.restart-box {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--app-warning-soft-border);
  background-color: var(--app-warning-soft-bg);
  border-radius: 8px;
}

.restart-box p {
  color: var(--app-warning-text);
}

.restart-actions {
  margin-top: 8px;
  display: flex;
  gap: 10px;
}
</style>
