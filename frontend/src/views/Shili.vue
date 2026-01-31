<template>
  <div class="shili-page">
    <ShiliSidebar
      :shilis="shilis"
      :active-shili-id="activeShiliId"
      @select="handleSelectShili"
      @add="handleAddShili"
      @delete="handleDeleteShili"
    />
    <div class="main-content">
      <div v-if="activeShili" class="shili-detail">
        <!-- 势力基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>基本信息</h2>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <label>名称</label>
              <div class="info-value">{{ activeShili.name }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'level' }" @dblclick="startEditBasicInfo('level')">
              <label>等级</label>
              <div v-if="editingBasicInfo.field !== 'level'" class="info-value">Lv.{{ activeShili.level }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.level"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('level')"
                @keyup.enter="saveBasicInfo('level')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'founder' }" @dblclick="startEditBasicInfo('founder')">
              <label>创始人</label>
              <div v-if="editingBasicInfo.field !== 'founder'" class="info-value">{{ activeShili.founder }}</div>
              <input
                v-else
                v-model="editingBasicInfo.founder"
                class="edit-input"
                @blur="saveBasicInfo('founder')"
                @keyup.enter="saveBasicInfo('founder')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'wealth' }" @dblclick="startEditBasicInfo('wealth')">
              <label>财富</label>
              <div v-if="editingBasicInfo.field !== 'wealth'" class="info-value">{{ activeShili.wealth }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.wealth"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('wealth')"
                @keyup.enter="saveBasicInfo('wealth')"
              />
            </div>
            <div class="info-item">
              <label>现有人数</label>
              <div class="info-value">{{ activeShili.memberCount }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'maxMembers' }" @dblclick="startEditBasicInfo('maxMembers')">
              <label>容纳人数</label>
              <div v-if="editingBasicInfo.field !== 'maxMembers'" class="info-value">{{ activeShili.maxMembers }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.maxMembers"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('maxMembers')"
                @keyup.enter="saveBasicInfo('maxMembers')"
              />
            </div>
          </div>
        </section>

        <!-- 职务部分 -->
        <section class="positions-section">
          <div class="section-header">
            <h2>职务</h2>
          </div>
          <div class="positions-list">
            <div
              v-for="position in activeShili.positions"
              :key="position.id"
              class="position-item"
            >
              <div v-if="editingPosition.id !== position.id" class="position-content" @dblclick="startEditPosition(position)">
                <div class="position-name">{{ position.positionName }}</div>
                <div class="position-person">{{ position.personName }}</div>
                <div v-if="position.description" class="position-description">{{ position.description }}</div>
              </div>
              <div v-else class="position-content editing">
                <input
                  v-model="editingPosition.positionName"
                  class="edit-input"
                  placeholder="职务名称"
                />
                <input
                  v-model="editingPosition.personName"
                  class="edit-input"
                  placeholder="姓名"
                />
                <input
                  v-model="editingPosition.description"
                  class="edit-input"
                  placeholder="描述"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="savePosition">保存</button>
                  <button class="cancel-btn" @click="cancelEditPosition">取消</button>
                </div>
              </div>
              <button v-if="editingPosition.id !== position.id" class="delete-btn" @click="handleDeletePosition(position.id)" title="删除职务">
                ×
              </button>
            </div>
            <div v-if="activeShili.positions.length === 0" class="empty-positions">
              <p>暂无职务</p>
              <p class="tip">点击下方按钮添加</p>
            </div>
          </div>
          <button class="add-btn" @click="showPositionModal = true">
            <span>+</span> 添加职务
          </button>
        </section>

        <!-- 其他属性部分 -->
        <section class="attributes-section">
          <div class="section-header">
            <h2>其他属性</h2>
          </div>
          <div class="attributes-list">
            <div
              v-for="attr in activeShili.attributes"
              :key="attr.id"
              class="attribute-item"
            >
              <div v-if="editingAttribute.id !== attr.id" class="attribute-content" @dblclick="startEditAttribute(attr)">
                <div class="attribute-name">{{ attr.name }}</div>
                <div v-if="attr.description" class="attribute-description">{{ attr.description }}</div>
                <div class="attribute-value">{{ attr.value }}</div>
              </div>
              <div v-else class="attribute-content editing">
                <input
                  v-model="editingAttribute.name"
                  class="edit-input"
                  placeholder="属性名称"
                />
                <input
                  v-model="editingAttribute.description"
                  class="edit-input"
                  placeholder="属性描述"
                />
                <input
                  v-model.number="editingAttribute.value"
                  type="number"
                  class="edit-input"
                  placeholder="属性值"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveAttribute">保存</button>
                  <button class="cancel-btn" @click="cancelEditAttribute">取消</button>
                </div>
              </div>
              <button v-if="editingAttribute.id !== attr.id" class="delete-btn" @click="handleDeleteAttribute(attr.id)" title="删除属性">
                ×
              </button>
            </div>
            <div v-if="activeShili.attributes.length === 0" class="empty-attributes">
              <p>暂无属性</p>
              <p class="tip">点击下方按钮添加</p>
            </div>
          </div>
          <button class="add-btn" @click="showAttributeModal = true">
            <span>+</span> 添加属性
          </button>
        </section>
      </div>
      <div v-else class="empty-state">
        <p>请选择或添加一个势力</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddShiliModal
      v-model:visible="showShiliModal"
      @confirm="handleAddShiliConfirm"
    />
    <AddPositionModal
      v-model:visible="showPositionModal"
      @confirm="handleAddPosition"
    />
    <AddAttributeModal
      v-model:visible="showAttributeModal"
      @confirm="handleAddAttribute"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import ShiliSidebar from '../components/ShiliSidebar.vue'
import AddShiliModal from '../components/AddShiliModal.vue'
import AddPositionModal from '../components/AddPositionModal.vue'
import AddAttributeModal from '../components/AddAttributeModal.vue'

// 势力数据
const shilis = ref([])
const activeShiliId = ref(null)
const showShiliModal = ref(false)
const showPositionModal = ref(false)
const showAttributeModal = ref(false)
const loading = ref(false)

// 基本信息编辑状态
const editingBasicInfo = ref({
  field: '', // 当前正在编辑的字段
  level: 0,
  founder: '',
  wealth: 0,
  maxMembers: 0
})

// 职务编辑状态
const editingPosition = ref({
  id: null,
  positionName: '',
  personName: '',
  description: ''
})

// 属性编辑状态
const editingAttribute = ref({
  id: null,
  name: '',
  description: '',
  value: 0
})

const activeShili = computed(() => {
  return shilis.value.find(s => s.id === activeShiliId.value)
})

// 组件挂载时加载势力列表
onMounted(async () => {
  await loadShilis()
})

// 从数据库加载势力列表
async function loadShilis() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllShili()
    shilis.value = result.map(shili => ({
      id: parseInt(shili.id),
      name: shili.name,
      level: parseInt(shili.level),
      founder: shili.founder || '',
      wealth: parseInt(shili.wealth),
      memberCount: parseInt(shili.memberCount),
      maxMembers: parseInt(shili.maxMembers),
      positions: [],
      attributes: []
    }))
    
    // 如果有数据，默认选中第一个
    if (shilis.value.length > 0 && !activeShiliId.value) {
      activeShiliId.value = shilis.value[0].id
      await loadShiliDetail(activeShiliId.value)
    }
  } catch (error) {
    console.error('加载势力列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 从数据库加载势力详细信息
async function loadShiliDetail(shiliId) {
  try {
    const result = await window.go.main.app.GetShiliInfo(parseInt(shiliId))
    const shili = shilis.value.find(s => s.id === parseInt(shiliId))
    if (shili) {
      // 更新基本信息
      shili.name = result.name
      shili.level = parseInt(result.level)
      shili.founder = result.founder || ''
      shili.wealth = parseInt(result.wealth)
      shili.memberCount = parseInt(result.memberCount)
      shili.maxMembers = parseInt(result.maxMembers)
      // 更新职务和属性
      shili.positions = result.positions || []
      shili.attributes = result.attributes || []
    }
  } catch (error) {
    console.error('加载势力详情失败:', error)
  }
}

// 选择势力
async function handleSelectShili(id) {
  activeShiliId.value = id
  await loadShiliDetail(id)
}

// 添加势力 - 显示弹窗
function handleAddShili() {
  showShiliModal.value = true
}

// 确认添加势力
async function handleAddShiliConfirm(data) {
  console.log('开始添加势力，数据:', data)
  try {
    console.log('调用 CreateShili 接口，参数:', data.name, data.founder, data.level, data.wealth, data.maxMembers)
    const newId = await window.go.main.app.CreateShili(
      data.name,
      data.founder,
      parseInt(data.level),
      parseInt(data.wealth),
      parseInt(data.maxMembers)
    )
    console.log('创建势力成功，ID:', newId)
    
    // 重新加载列表
    await loadShilis()
    activeShiliId.value = parseInt(newId)
    await loadShiliDetail(newId)
  } catch (error) {
    console.error('添加势力失败:', error)
    alert('添加势力失败: ' + error.message)
  }
}

// 删除势力
async function handleDeleteShili(shiliId) {
  if (!confirm('确定要删除这个势力吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteShili(shiliId)
    
    // 如果删除的是当前选中的势力，清空选中状态
    if (activeShiliId.value === shiliId) {
      activeShiliId.value = null
    }
    
    // 重新加载列表
    await loadShilis()
  } catch (error) {
    console.error('删除势力失败:', error)
    alert('删除势力失败: ' + error.message)
  }
}

// 添加职务
async function handleAddPosition(data) {
  if (!activeShili.value) return
  try {
    await window.go.main.app.AddShiliPosition(
      activeShili.value.id,
      data.positionName,
      data.personName,
      data.description
    )
    // 重新加载职务列表
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('添加职务失败:', error)
  }
}

// 删除职务
async function handleDeletePosition(id) {
  if (!activeShili.value) return
  try {
    await window.go.main.app.DeleteShiliPosition(id)
    // 重新加载职务列表
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('删除职务失败:', error)
  }
}

// 添加属性
async function handleAddAttribute(data) {
  if (!activeShili.value) return
  try {
    await window.go.main.app.AddShiliAttribute(
      activeShili.value.id,
      data.name,
      data.description,
      parseInt(data.value)
    )
    // 重新加载属性列表
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('添加属性失败:', error)
  }
}

// 删除属性
async function handleDeleteAttribute(id) {
  if (!activeShili.value) return
  try {
    await window.go.main.app.DeleteShiliAttribute(id)
    // 重新加载属性列表
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('删除属性失败:', error)
  }
}

// 开始编辑基本信息
function startEditBasicInfo(field) {
  if (!activeShili.value) return
  
  // 设置当前编辑的字段
  editingBasicInfo.value.field = field
  
  // 复制当前值到编辑状态
  editingBasicInfo.value[field] = activeShili.value[field]
  
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
  if (!activeShili.value) return
  
  const value = editingBasicInfo.value[field]
  
  try {
    // 调用后端接口更新基本信息
    await window.go.main.app.UpdateShiliBasicInfo(
      parseInt(activeShili.value.id),
      activeShili.value.name,
      field === 'level' ? parseInt(value) : parseInt(activeShili.value.level),
      field === 'founder' ? value : activeShili.value.founder,
      field === 'wealth' ? parseInt(value) : parseInt(activeShili.value.wealth),
      field === 'maxMembers' ? parseInt(value) : parseInt(activeShili.value.maxMembers)
    )
    
    // 重新加载势力信息
    await loadShiliDetail(activeShili.value.id)
    
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

// 开始编辑职务
function startEditPosition(position) {
  editingPosition.value = {
    id: position.id,
    positionName: position.positionName,
    personName: position.personName,
    description: position.description || ''
  }
  nextTick(() => {
    const inputs = document.querySelectorAll('.position-item.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存职务
async function savePosition() {
  try {
    await window.go.main.app.UpdateShiliPosition(
      editingPosition.value.id,
      editingPosition.value.positionName,
      editingPosition.value.personName,
      editingPosition.value.description
    )
    editingPosition.value = { id: null, positionName: '', personName: '', description: '' }
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('保存职务失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑职务
function cancelEditPosition() {
  editingPosition.value = { id: null, positionName: '', personName: '', description: '' }
}

// 开始编辑属性
function startEditAttribute(attr) {
  editingAttribute.value = {
    id: attr.id,
    name: attr.name,
    description: attr.description || '',
    value: attr.value
  }
  nextTick(() => {
    const inputs = document.querySelectorAll('.attribute-item.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存属性
async function saveAttribute() {
  try {
    await window.go.main.app.UpdateShiliAttribute(
      editingAttribute.value.id,
      editingAttribute.value.name,
      editingAttribute.value.description,
      editingAttribute.value.value
    )
    editingAttribute.value = { id: null, name: '', description: '', value: 0 }
    await loadShiliDetail(activeShili.value.id)
  } catch (error) {
    console.error('保存属性失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑属性
function cancelEditAttribute() {
  editingAttribute.value = { id: null, name: '', description: '', value: 0 }
}
</script>

<style scoped>
.shili-page {
  display: flex;
  height: 100%;
  background-color: #f5f7fa;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px;
}

.shili-detail {
  width: 100%;
}

.info-section,
.positions-section,
.attributes-section {
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

/* 职务列表 */
.positions-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.position-item {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e8eaed;
  transition: all 0.2s ease;
}

.position-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.position-content {
  flex: 1;
  min-width: 0;
}

.position-name {
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.position-person {
  font-size: 14px;
  color: #595959;
  margin-bottom: 8px;
}

.position-description {
  font-size: 13px;
  color: #8c92a0;
  line-height: 1.4;
}

/* 属性列表 */
.attributes-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.attribute-item {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e8eaed;
  transition: all 0.2s ease;
}

.attribute-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.attribute-content {
  flex: 1;
  min-width: 0;
}

.attribute-name {
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.attribute-description {
  font-size: 13px;
  color: #8c92a0;
  margin-bottom: 10px;
  line-height: 1.4;
}

.attribute-value {
  font-size: 18px;
  font-weight: 600;
  color: #1890ff;
}

/* 空状态 */
.empty-positions,
.empty-attributes {
  grid-column: 1 / -1;
  text-align: center;
  padding: 48px 20px;
  color: #8c92a0;
}

.empty-positions p,
.empty-attributes p {
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
.position-item.editing,
.attribute-item.editing {
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
  .shili-detail {
    max-width: 900px;
  }
  
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
  }
  
  .positions-list,
  .attributes-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 24px;
  }
  
  .info-grid,
  .positions-list,
  .attributes-list {
    grid-template-columns: 1fr;
  }
  
  .info-section,
  .positions-section,
  .attributes-section {
    padding: 24px;
  }
}
</style>