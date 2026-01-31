<template>
  <div class="daoju-page">
    <DaojuSidebar
      :daoju="daoju"
      :active-daoju-id="activeDaojuId"
      @select="handleSelectDaoju"
      @add="handleAddDaoju"
      @delete="handleDeleteDaoju"
    />
    <div class="main-content">
      <div v-if="activeDaoju" class="daoju-detail">
        <!-- 道具基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>基本信息</h2>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <label>名称</label>
              <div class="info-value">{{ activeDaoju.name }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'level' }" @dblclick="startEditBasicInfo('level')">
              <label>等级</label>
              <div v-if="editingBasicInfo.field !== 'level'" class="info-value">Lv.{{ activeDaoju.level }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.level"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('level')"
                @keyup.enter="saveBasicInfo('level')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'holder' }" @dblclick="startEditBasicInfo('holder')">
              <label>所有人</label>
              <div v-if="editingBasicInfo.field !== 'holder'" class="info-value">{{ activeDaoju.holder || '无' }}</div>
              <input
                v-else
                v-model="editingBasicInfo.holder"
                class="edit-input"
                @blur="saveBasicInfo('holder')"
                @keyup.enter="saveBasicInfo('holder')"
              />
            </div>
          </div>
        </section>

        <!-- 功能描述部分 -->
        <section class="functions-section">
          <div class="section-header">
            <h2>功能描述</h2>
          </div>
          <div class="functions-list">
            <div
              v-for="func in activeDaoju.functions"
              :key="func.id"
              class="function-item"
            >
              <div v-if="editingFunction.id !== func.id" class="function-content" @dblclick="startEditFunction(func)">
                <div class="function-name">{{ func.name }}</div>
                <div v-if="func.description" class="function-description">{{ func.description }}</div>
              </div>
              <div v-else class="function-content editing">
                <input
                  v-model="editingFunction.name"
                  class="edit-input"
                  placeholder="功能名称"
                />
                <input
                  v-model="editingFunction.description"
                  class="edit-input"
                  placeholder="功能表述"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveFunction">保存</button>
                  <button class="cancel-btn" @click="cancelEditFunction">取消</button>
                </div>
              </div>
              <button v-if="editingFunction.id !== func.id" class="delete-btn" @click="handleDeleteFunction(func.id)" title="删除功能">
                ×
              </button>
            </div>
            <div v-if="activeDaoju.functions.length === 0" class="empty-functions">
              <p>暂无功能描述</p>
              <p class="tip">点击下方按钮添加</p>
            </div>
          </div>
          <button class="add-btn" @click="showFunctionModal = true">
            <span>+</span> 添加功能
          </button>
        </section>
      </div>
      <div v-else class="empty-state">
        <p>请选择或添加一个道具</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddDaojuModal
      v-model:visible="showDaojuModal"
      @confirm="handleAddDaojuConfirm"
    />
    <AddFunctionModal
      v-model:visible="showFunctionModal"
      @confirm="handleAddFunction"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import DaojuSidebar from '../components/DaojuSidebar.vue'
import AddFunctionModal from '../components/AddFunctionModal.vue'
import AddDaojuModal from '../components/AddDaojuModal.vue'

// 道具数据
const daoju = ref([])
const activeDaojuId = ref(null)
const showFunctionModal = ref(false)
const showDaojuModal = ref(false)
const loading = ref(false)

// 基本信息编辑状态
const editingBasicInfo = ref({
  field: '', // 当前正在编辑的字段
  level: 0,
  holder: ''
})

// 功能编辑状态
const editingFunction = ref({
  id: null,
  name: '',
  description: ''
})

const activeDaoju = computed(() => {
  return daoju.value.find(d => d.id === activeDaojuId.value)
})

// 组件挂载时加载道具列表
onMounted(async () => {
  await loadDaoju()
})

// 从数据库加载道具列表
async function loadDaoju() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllDaoju()
    daoju.value = result.map(d => ({
      id: parseInt(d.id),
      name: d.name,
      level: parseInt(d.level),
      holder: d.holder,
      functions: []
    }))
    
    // 如果有数据，默认选中第一个
    if (daoju.value.length > 0 && !activeDaojuId.value) {
      activeDaojuId.value = daoju.value[0].id
      await loadDaojuDetail(activeDaojuId.value)
    }
  } catch (error) {
    console.error('加载道具列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 从数据库加载道具详细信息
async function loadDaojuDetail(daojuId) {
  try {
    const result = await window.go.main.app.GetDaojuInfo(parseInt(daojuId))
    const d = daoju.value.find(d => d.id === parseInt(daojuId))
    if (d) {
      // 更新基本信息
      d.name = result.name
      d.level = parseInt(result.level)
      d.holder = result.holder
      // 更新功能
      d.functions = result.functions || []
    }
  } catch (error) {
    console.error('加载道具详情失败:', error)
  }
}

// 选择道具
async function handleSelectDaoju(id) {
  activeDaojuId.value = id
  await loadDaojuDetail(id)
}

// 添加道具 - 显示弹窗
function handleAddDaoju() {
  showDaojuModal.value = true
}

// 确认添加道具
async function handleAddDaojuConfirm(data) {
  console.log('开始添加道具，数据:', data)
  try {
    console.log('调用 CreateDaoju 接口，参数:', data.name, data.holder, data.level)
    const newId = await window.go.main.app.CreateDaoju(
      data.name,
      data.holder,
      parseInt(data.level)
    )
    console.log('创建道具成功，ID:', newId)
    
    // 重新加载列表
    await loadDaoju()
    activeDaojuId.value = parseInt(newId)
    await loadDaojuDetail(newId)
  } catch (error) {
    console.error('添加道具失败:', error)
    alert('添加道具失败: ' + error.message)
  }
}

// 删除道具
async function handleDeleteDaoju(daojuId) {
  if (!confirm('确定要删除这个道具吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteDaoju(daojuId)
    
    // 如果删除的是当前选中的道具，清空选中状态
    if (activeDaojuId.value === daojuId) {
      activeDaojuId.value = null
    }
    
    // 重新加载列表
    await loadDaoju()
  } catch (error) {
    console.error('删除道具失败:', error)
    alert('删除道具失败: ' + error.message)
  }
}

// 添加功能
async function handleAddFunction(data) {
  if (!activeDaoju.value) return
  try {
    await window.go.main.app.AddDaojuFunction(
      activeDaoju.value.id,
      data.name,
      data.description
    )
    // 重新加载功能列表
    await loadDaojuDetail(activeDaoju.value.id)
  } catch (error) {
    console.error('添加功能失败:', error)
  }
}

// 删除功能
async function handleDeleteFunction(id) {
  if (!activeDaoju.value) return
  try {
    await window.go.main.app.DeleteDaojuFunction(id)
    // 重新加载功能列表
    await loadDaojuDetail(activeDaoju.value.id)
  } catch (error) {
    console.error('删除功能失败:', error)
  }
}

// 开始编辑基本信息
function startEditBasicInfo(field) {
  if (!activeDaoju.value) return
  
  // 设置当前编辑的字段
  editingBasicInfo.value.field = field
  
  // 复制当前值到编辑状态
  editingBasicInfo.value[field] = activeDaoju.value[field]
  
  // 聚焦到对应的输入框
  nextTick(() => {
    const inputs = document.querySelectorAll('.info-item.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存基本信息
async function saveBasicInfo(field) {
  if (!activeDaoju.value) return
  
  const value = editingBasicInfo.value[field]
  
  try {
    // 调用后端接口更新基本信息
    await window.go.main.app.UpdateDaojuBasicInfo(
      parseInt(activeDaoju.value.id),
      field === 'holder' ? value : activeDaoju.value.holder,
      field === 'level' ? parseInt(value) : parseInt(activeDaoju.value.level)
    )
    
    // 重新加载道具信息
    await loadDaojuDetail(activeDaoju.value.id)
    
    // 清空编辑状态
    editingBasicInfo.value.field = ''
    editingBasicInfo.value[field] = ''
  } catch (error) {
    console.error('更新基本信息失败:', error)
    alert('更新失败: ' + error.message)
    // 清空编辑状态
    editingBasicInfo.value.field = ''
    editingBasicInfo.value[field] = ''
  }
}

// 开始编辑功能
function startEditFunction(func) {
  editingFunction.value = {
    id: func.id,
    name: func.name,
    description: func.description || ''
  }
  nextTick(() => {
    const inputs = document.querySelectorAll('.function-item.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存功能
async function saveFunction() {
  try {
    await window.go.main.app.UpdateDaojuFunction(
      editingFunction.value.id,
      editingFunction.value.name,
      editingFunction.value.description
    )
    editingFunction.value = { id: null, name: '', description: '' }
    await loadDaojuDetail(activeDaoju.value.id)
  } catch (error) {
    console.error('保存功能失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑功能
function cancelEditFunction() {
  editingFunction.value = { id: null, name: '', description: '' }
}
</script>

<style scoped>
.daoju-page {
  display: flex;
  height: 100%;
  background-color: #f5f7fa;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px;
}

.daoju-detail {
  width: 100%;
}

.info-section,
.functions-section {
  background-color: white;
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 28px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.section-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f2f5;
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

/* 基本信息 */
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
  color: #8c92a0;
  margin-bottom: 12px;
  font-weight: 500;
}

.info-value {
  font-size: 16px;
  color: #2c3e50;
  font-weight: 500;
}

/* 功能列表 */
.functions-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.function-item {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e8eaed;
  transition: all 0.2s ease;
}

.function-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.function-content {
  flex: 1;
  min-width: 0;
}

.function-name {
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.function-description {
  font-size: 13px;
  color: #8c92a0;
  line-height: 1.6;
}

.empty-functions {
  grid-column: 1 / -1;
  text-align: center;
  padding: 48px 20px;
  color: #8c92a0;
}

.empty-functions p {
  margin: 10px 0;
}

.tip {
  font-size: 13px;
  color: #b0b8c3;
}

/* 删除按钮 */
.delete-btn {
  width: 30px;
  height: 30px;
  border: none;
  background: none;
  color: #8c92a0;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

.delete-btn:hover {
  background-color: #fff1f0;
  color: #ff4d4f;
}

/* 添加按钮 */
.add-btn {
  width: 100%;
  padding: 14px 20px;
  border: 2px dashed #d9d9d9;
  background: none;
  border-radius: 8px;
  font-size: 14px;
  color: #8c92a0;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.add-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
  background-color: #f0f7ff;
}

.add-btn span {
  font-size: 20px;
  font-weight: 600;
}

/* 空状态 */
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 500px;
  color: #8c92a0;
  font-size: 16px;
}

/* 编辑输入框 */
.edit-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  margin-bottom: 8px;
  transition: all 0.2s ease;
}

.edit-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.edit-input::placeholder {
  color: #b0b8c3;
}

/* 编辑状态 */
.function-item.editing {
  border-color: #1890ff;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

/* 编辑按钮 */
.edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.save-btn {
  flex: 1;
  padding: 6px 12px;
  border: none;
  background-color: #1890ff;
  color: white;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.save-btn:hover {
  background-color: #40a9ff;
}

.cancel-btn {
  flex: 1;
  padding: 6px 12px;
  border: 1px solid #d9d9d9;
  background-color: white;
  color: #595959;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancel-btn:hover {
  border-color: #8c92a0;
  color: #2c3e50;
}

/* 基本信息编辑状态 */
.info-item.editing .info-value {
  display: none;
}

.info-item.editing .edit-input {
  display: block;
}

.info-item .edit-input {
  display: none;
}

/* 编辑状态下的info-item样式 */
.info-item[style*="cursor: pointer"] {
  cursor: pointer;
}

.info-item:hover {
  cursor: pointer;
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .daoju-detail {
    max-width: 900px;
  }
  
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
  }
  
  .functions-list {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 24px;
  }
  
  .info-grid,
  .functions-list {
    grid-template-columns: 1fr;
  }
  
  .info-section,
  .functions-section {
    padding: 24px;
  }
}
</style>