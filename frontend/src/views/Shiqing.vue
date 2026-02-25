<script setup>
import { ref, onMounted } from 'vue'
import ShiqingSidebar from '../components/ShiqingSidebar.vue'
import AddShiqingModal from '../components/AddShiqingModal.vue'
import AddDetailModal from '../components/AddDetailModal.vue'

const shiqing = ref([])
const activeShiqingId = ref(null)
const currentShiqing = ref(null)
const details = ref([])

// 弹窗控制
const showAddShiqingModal = ref(false)
const showDetailModal = ref(false)
const editingDetail = ref(null)

// 基本信息编辑
const editingLocation = ref(false)
const editingTime = ref(false)
const tempLocation = ref('')
const tempTime = ref('')

// 加载所有任务
async function loadShiqing() {
  try {
    const result = await window.go.main.app.GetAllShiqing()
    shiqing.value = result || []
  } catch (error) {
    console.error('加载任务失败:', error)
    alert('加载任务失败')
  }
}

// 加载任务详情
async function loadShiqingDetail() {
  if (!activeShiqingId.value) {
    details.value = []
    return
  }
  
  try {
    const result = await window.go.main.app.GetShiqingDetails(activeShiqingId.value)
    details.value = result || []
  } catch (error) {
    console.error('加载任务详情失败:', error)
    alert('加载任务详情失败')
  }
}

// 选择任务
async function handleSelectShiqing(id) {
  activeShiqingId.value = id
  currentShiqing.value = shiqing.value.find(s => s.id === id)
  await loadShiqingDetail()
}

// 添加任务
async function handleAddShiqing(data) {
  try {
    const newId = await window.go.main.app.CreateShiqing(data.name, data.location, data.time)
    console.log('创建任务成功，ID:', newId)
    
    // 重新加载列表
    await loadShiqing()
    // 自动选择新创建的任务
    if (newId) {
      await handleSelectShiqing(newId)
    }
  } catch (error) {
    console.error('创建任务失败:', error)
    alert('创建任务失败: ' + error.message)
  }
}

// 删除任务
async function handleDeleteShiqing(id) {
  if (!confirm('确定要删除这个任务吗？')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteShiqing(id)
    
    // 如果删除的是当前选中的任务，清空选中状态
    if (activeShiqingId.value === id) {
      activeShiqingId.value = null
      currentShiqing.value = null
      details.value = []
    }
    
    // 重新加载列表
    await loadShiqing()
  } catch (error) {
    console.error('删除任务失败:', error)
    alert('删除任务失败: ' + error.message)
  }
}

// 添加任务描述
async function handleAddDetail(data) {
  if (!activeShiqingId.value) {
    alert('请先选择一个任务')
    return
  }
  
  try {
    await window.go.main.app.AddShiqingDetail(activeShiqingId.value, data.description)
    // 重新加载详情列表
    await loadShiqingDetail()
  } catch (error) {
    console.error('添加任务描述失败:', error)
    alert('添加任务描述失败: ' + error.message)
  }
}

// 编辑任务描述
async function handleEditDetail(detail) {
  editingDetail.value = detail
  showDetailModal.value = true
}

// 保存编辑的任务描述
async function handleSaveDetail(data) {
  if (!editingDetail.value || !activeShiqingId.value) {
    return
  }
  
  try {
    await window.go.main.app.UpdateShiqingDetail(
      editingDetail.value.id,
      data.description
    )
    // 重新加载详情列表
    await loadShiqingDetail()
  } catch (error) {
    console.error('更新任务描述失败:', error)
    alert('更新任务描述失败: ' + error.message)
  }
}

// 删除任务描述
async function handleDeleteDetail(detailId) {
  if (!confirm('确定要删除这个任务描述吗？')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteShiqingDetail(activeShiqingId.value, detailId)
    // 重新加载详情列表
    await loadShiqingDetail()
  } catch (error) {
    console.error('删除任务描述失败:', error)
    alert('删除任务描述失败: ' + error.message)
  }
}

// 打开详情弹窗
function openDetailModal() {
  editingDetail.value = null
  showDetailModal.value = true
}

// 处理详情弹窗确认
function handleDetailModalConfirm(data) {
  if (editingDetail.value) {
    handleSaveDetail(data)
  } else {
    handleAddDetail(data)
  }
}

// 编辑任务地点
function startEditLocation() {
  tempLocation.value = currentShiqing.value?.location || ''
  editingLocation.value = true
}

// 保存任务地点
async function saveLocation() {
  if (!activeShiqingId.value) {
    return
  }
  
  try {
    await window.go.main.app.UpdateShiqingBasicInfo(
      activeShiqingId.value,
      tempLocation.value,
      currentShiqing.value?.time || ''
    )
    // 重新加载任务信息
    await loadShiqing()
    currentShiqing.value = shiqing.value.find(s => s.id === activeShiqingId.value)
    editingLocation.value = false
  } catch (error) {
    console.error('更新任务地点失败:', error)
    alert('更新任务地点失败: ' + error.message)
  }
}

// 取消编辑任务地点
function cancelEditLocation() {
  editingLocation.value = false
  tempLocation.value = ''
}

// 编辑任务时间
function startEditTime() {
  tempTime.value = currentShiqing.value?.time || ''
  editingTime.value = true
}

// 保存任务时间
async function saveTime() {
  if (!activeShiqingId.value) {
    return
  }
  
  try {
    await window.go.main.app.UpdateShiqingBasicInfo(
      activeShiqingId.value,
      currentShiqing.value?.location || '',
      tempTime.value
    )
    // 重新加载任务信息
    await loadShiqing()
    currentShiqing.value = shiqing.value.find(s => s.id === activeShiqingId.value)
    editingTime.value = false
  } catch (error) {
    console.error('更新任务时间失败:', error)
    alert('更新任务时间失败: ' + error.message)
  }
}

// 取消编辑任务时间
function cancelEditTime() {
  editingTime.value = false
  tempTime.value = ''
}

onMounted(() => {
  loadShiqing()
})
</script>

<template>
  <div class="shiqing-page">
    <ShiqingSidebar
      :shiqing="shiqing"
      :active-shiqing-id="activeShiqingId"
      @select="handleSelectShiqing"
      @add="showAddShiqingModal = true"
      @delete="handleDeleteShiqing"
    />
    
    <div class="main-content">
      <div v-if="!currentShiqing" class="empty-state">
        <p>请选择一个任务查看详情</p>
        <p class="tip">或点击左侧 + 添加新任务</p>
      </div>
      
      <div v-else class="shiqing-detail">
        <!-- 基本信息 -->
        <div class="info-section">
          <h2>基本信息</h2>
          <div class="info-grid">
            <div class="info-item">
              <label>任务名称</label>
              <div class="info-value">{{ currentShiqing.name }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingLocation }">
              <label>任务地点</label>
              <div v-if="!editingLocation" class="info-value" @click="startEditLocation">
                {{ currentShiqing.location || '未设置' }}
              </div>
              <div v-else class="edit-wrapper">
                <input
                  v-model="tempLocation"
                  class="edit-input"
                  placeholder="请输入任务地点"
                  @keyup.enter="saveLocation"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveLocation">保存</button>
                  <button class="cancel-btn" @click="cancelEditLocation">取消</button>
                </div>
              </div>
            </div>
            <div class="info-item" :class="{ editing: editingTime }">
              <label>任务时间</label>
              <div v-if="!editingTime" class="info-value" @click="startEditTime">
                {{ currentShiqing.time || '未设置' }}
              </div>
              <div v-else class="edit-wrapper">
                <input
                  v-model="tempTime"
                  class="edit-input"
                  placeholder="请输入任务时间"
                  @keyup.enter="saveTime"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveTime">保存</button>
                  <button class="cancel-btn" @click="cancelEditTime">取消</button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 任务详情 -->
        <div class="details-section">
          <h2>任务详情</h2>
          <div class="details-list">
            <div
              v-for="detail in details"
              :key="detail.id"
              class="detail-item"
            >
              <div class="detail-content">
                <p class="detail-description">{{ detail.description }}</p>
              </div>
              <div class="detail-actions">
                <button class="edit-btn" @click="handleEditDetail(detail)" title="编辑">✎</button>
                <button class="delete-btn" @click="handleDeleteDetail(detail.id)" title="删除">×</button>
              </div>
            </div>
            <div v-if="details.length === 0" class="empty-details">
              <p>暂无任务描述</p>
              <p class="tip">点击下方按钮添加任务描述</p>
            </div>
          </div>
          <button class="add-btn" @click="openDetailModal">
            <span>+</span> 添加任务描述
          </button>
        </div>
      </div>
    </div>
    
    <!-- 添加任务弹窗 -->
    <AddShiqingModal
      v-model:visible="showAddShiqingModal"
      @confirm="handleAddShiqing"
    />
    
    <!-- 添加/编辑任务描述弹窗 -->
    <AddDetailModal
      v-model:visible="showDetailModal"
      :detail="editingDetail"
      @confirm="handleDetailModalConfirm"
    />
  </div>
</template>

<style scoped>
.shiqing-page {
  display: flex;
  height: 100%;
  background-color: var(--app-bg);
}

.main-content {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 500px;
  color: var(--app-text-muted);
  font-size: 16px;
}

.empty-state p {
  margin: 10px 0;
}

.tip {
  font-size: 13px;
  color: var(--app-text-muted);
}

.shiqing-detail {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-section,
.details-section {
  background-color: var(--app-surface);
  border-radius: 12px;
  padding: 32px;
  box-shadow: var(--app-shadow-md);
}

.info-section {
  flex-shrink: 0;
}

.details-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.details-section .details-list {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 24px;
}

.info-section h2,
.details-section h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin: 0 0 24px 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-item label {
  font-size: 13px;
  color: var(--app-text-muted);
  margin-bottom: 12px;
  font-weight: 500;
}

.info-value {
  font-size: 16px;
  color: var(--app-text-primary);
  font-weight: 500;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.info-value:hover {
  background-color: var(--app-bg);
}

.info-value:empty::before {
  content: '未设置';
  color: var(--app-text-muted);
}

.edit-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.edit-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--app-border);
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.edit-input:focus {
  outline: none;
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-focus);
}

.edit-actions {
  display: flex;
  gap: 8px;
}

.save-btn {
  flex: 1;
  padding: 6px 12px;
  border: none;
  background-color: var(--app-accent);
  color: var(--app-text-inverse);
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.save-btn:hover {
  background-color: var(--app-accent-hover);
}

.cancel-btn {
  flex: 1;
  padding: 6px 12px;
  border: 1px solid var(--app-border);
  background-color: var(--app-surface);
  color: var(--app-text-secondary);
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  border-color: var(--app-text-muted);
  color: var(--app-text-primary);
}

.details-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.detail-item {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  background-color: var(--app-item-bg);
  border-radius: 8px;
  border: 1px solid var(--app-border);
  transition: all 0.2s ease;
}

.detail-item:hover {
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-accent-sm);
}

.detail-content {
  flex: 1;
  min-width: 0;
}

.detail-description {
  font-size: 14px;
  color: var(--app-text-primary);
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.detail-actions {
  display: flex;
  gap: 8px;
  margin-left: 16px;
  flex-shrink: 0;
}

.edit-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--app-text-muted);
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  padding: 0;
}

.edit-btn:hover {
  background-color: var(--app-hover-bg);
  color: var(--app-accent);
}

.delete-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--app-text-muted);
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  padding: 0;
  line-height: 1;
}

.delete-btn:hover {
  background-color: var(--app-danger-soft-bg);
  color: var(--app-danger);
}

.empty-details {
  text-align: center;
  padding: 48px 20px;
  color: var(--app-text-muted);
}

.empty-details p {
  margin: 10px 0;
}

.add-btn {
  width: 100%;
  padding: 14px 20px;
  border: 2px dashed var(--app-border);
  background: none;
  border-radius: 8px;
  font-size: 14px;
  color: var(--app-text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.add-btn:hover {
  border-color: var(--app-accent);
  color: var(--app-accent);
  background-color: var(--app-hover-bg);
}

.add-btn span {
  font-size: 20px;
  font-weight: 600;
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .shiqing-detail {
    max-width: 900px;
  }
  
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 24px;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .info-section,
  .details-section {
    padding: 24px;
  }
}
</style>