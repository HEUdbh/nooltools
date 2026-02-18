<template>
  <div class="beibao-page">
    <BeibaoSidebar
      :beibao="beibao"
      :active-beibao-id="activeBeibaoId"
      @select="handleSelectBeibao"
      @add="handleAddBeibao"
      @delete="handleDeleteBeibao"
      @update="handleUpdateBeibao"
    />
    <div class="main-content">
      <div v-if="activeBeibao" class="beibao-detail">
        <!-- 背包基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>{{ activeBeibao.name }}</h2>
          </div>
        </section>

        <!-- 物品列表部分 -->
        <section class="items-section">
          <div class="section-header">
            <h2>物品列表</h2>
          </div>
          <div class="items-list">
            <div
              v-for="item in activeBeibao.items"
              :key="item.id"
              class="item-card"
            >
              <div v-if="editingItem.id !== item.id" class="item-content" @dblclick="startEditItem(item)">
                <div class="item-header">
                  <div class="item-name">{{ item.name }}</div>
                  <div class="item-quantity">x{{ item.quantity }}</div>
                </div>
                <div v-if="item.description" class="item-description">{{ item.description }}</div>
              </div>
              <div v-else class="item-content editing">
                <div class="item-header">
                  <input
                    v-model="editingItem.name"
                    class="edit-input"
                    placeholder="物品名称"
                  />
                  <input
                    v-model.number="editingItem.quantity"
                    type="number"
                    class="edit-input quantity-input"
                    placeholder="数量"
                  />
                </div>
                <input
                  v-model="editingItem.description"
                  class="edit-input"
                  placeholder="物品说明"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveItem">保存</button>
                  <button class="cancel-btn" @click="cancelEditItem">取消</button>
                </div>
              </div>
              <button v-if="editingItem.id !== item.id" class="delete-btn" @click="handleDeleteItem(item.id)" title="删除物品">
                ×
              </button>
            </div>
            <div v-if="activeBeibao.items.length === 0" class="empty-items">
              <p>暂无物品</p>
              <p class="tip">点击下方按钮添加</p>
            </div>
          </div>
          <button class="add-btn" @click="showItemModal = true">
            <span>+</span> 添加物品
          </button>
        </section>
      </div>
      <div v-else class="empty-state">
        <p>请选择或添加一个背包</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddBeibaoModal
      v-model:visible="showBeibaoModal"
      @confirm="handleAddBeibaoConfirm"
    />
    <AddItemModal
      v-model:visible="showItemModal"
      @confirm="handleAddItem"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import BeibaoSidebar from '../components/BeibaoSidebar.vue'
import AddItemModal from '../components/AddItemModal.vue'
import AddBeibaoModal from '../components/AddBeibaoModal.vue'

// 背包数据
const beibao = ref([])
const activeBeibaoId = ref(null)
const showItemModal = ref(false)
const showBeibaoModal = ref(false)
const loading = ref(false)

// 物品编辑状态
const editingItem = ref({
  id: null,
  name: '',
  quantity: 1,
  description: ''
})

const activeBeibao = computed(() => {
  return beibao.value.find(b => b.id === activeBeibaoId.value)
})

// 组件挂载时加载背包列表
onMounted(async () => {
  await loadBeibao()
})

// 从数据库加载背包列表
async function loadBeibao() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllBeibao()
    beibao.value = result.map(b => ({
      id: parseInt(b.id),
      name: b.name,
      items: []
    }))
    
    // 如果有数据，默认选中第一个
    if (beibao.value.length > 0 && !activeBeibaoId.value) {
      activeBeibaoId.value = beibao.value[0].id
      await loadBeibaoDetail(activeBeibaoId.value)
    }
  } catch (error) {
    console.error('加载背包列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 从数据库加载背包详细信息
async function loadBeibaoDetail(beibaoId) {
  try {
    const result = await window.go.main.app.GetBeibaoInfo(parseInt(beibaoId))
    const b = beibao.value.find(b => b.id === parseInt(beibaoId))
    if (b) {
      // 更新基本信息
      b.name = result.name
      // 更新物品
      b.items = result.items || []
    }
  } catch (error) {
    console.error('加载背包详情失败:', error)
  }
}

// 选择背包
async function handleSelectBeibao(id) {
  activeBeibaoId.value = id
  await loadBeibaoDetail(id)
}

// 添加背包 - 显示弹窗
function handleAddBeibao() {
  showBeibaoModal.value = true
}

// 确认添加背包
async function handleAddBeibaoConfirm(data) {
  console.log('开始添加背包，数据:', data)
  try {
    console.log('调用 CreateBeibao 接口，参数:', data.name)
    const newId = await window.go.main.app.CreateBeibao(data.name)
    console.log('创建背包成功，ID:', newId)
    
    // 重新加载列表
    await loadBeibao()
    activeBeibaoId.value = parseInt(newId)
    await loadBeibaoDetail(newId)
  } catch (error) {
    console.error('添加背包失败:', error)
    alert('添加背包失败: ' + error.message)
  }
}

// 删除背包
async function handleDeleteBeibao(beibaoId) {
  if (!confirm('确定要删除这个背包吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteBeibao(beibaoId)
    
    // 如果删除的是当前选中的背包，清空选中状态
    if (activeBeibaoId.value === beibaoId) {
      activeBeibaoId.value = null
    }
    
    // 重新加载列表
    await loadBeibao()
  } catch (error) {
    console.error('删除背包失败:', error)
    alert('删除背包失败: ' + error.message)
  }
}

// 更新背包名称
async function handleUpdateBeibao(beibaoId, name) {
  try {
    await window.go.main.app.UpdateBeibao(beibaoId, name)
    // 重新加载背包信息
    await loadBeibaoDetail(beibaoId)
  } catch (error) {
    console.error('更新背包失败:', error)
    alert('更新失败: ' + error.message)
  }
}

// 添加物品
async function handleAddItem(data) {
  if (!activeBeibao.value) return
  try {
    await window.go.main.app.AddBeibaoItem(
      activeBeibao.value.id,
      data.name,
      parseInt(data.quantity),
      data.description
    )
    // 重新加载物品列表
    await loadBeibaoDetail(activeBeibao.value.id)
  } catch (error) {
    console.error('添加物品失败:', error)
  }
}

// 删除物品
async function handleDeleteItem(id) {
  if (!activeBeibao.value) return
  try {
    await window.go.main.app.DeleteBeibaoItem(id)
    // 重新加载物品列表
    await loadBeibaoDetail(activeBeibao.value.id)
  } catch (error) {
    console.error('删除物品失败:', error)
  }
}

// 开始编辑物品
function startEditItem(item) {
  editingItem.value = {
    id: item.id,
    name: item.name,
    quantity: item.quantity,
    description: item.description || ''
  }
  nextTick(() => {
    const inputs = document.querySelectorAll('.item-card.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存物品
async function saveItem() {
  try {
    await window.go.main.app.UpdateBeibaoItem(
      editingItem.value.id,
      editingItem.value.name,
      parseInt(editingItem.value.quantity),
      editingItem.value.description
    )
    editingItem.value = { id: null, name: '', quantity: 1, description: '' }
    await loadBeibaoDetail(activeBeibao.value.id)
  } catch (error) {
    console.error('保存物品失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑物品
function cancelEditItem() {
  editingItem.value = { id: null, name: '', quantity: 1, description: '' }
}
</script>

<style scoped>
.beibao-page {
  display: flex;
  height: 100%;
  background-color: #f5f7fa;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px;
}

.beibao-detail {
  width: 100%;
}

.info-section,
.items-section {
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

/* 物品列表 */
.items-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.item-card {
  position: relative;
  background-color: #fafbfc;
  border: 1px solid #e8eaed;
  border-radius: 8px;
  padding: 20px;
  transition: all 0.2s ease;
}

.item-card:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.item-card.editing {
  border-color: #1890ff;
  background-color: #fff;
}

.item-content {
  cursor: pointer;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.item-name {
  font-size: 16px;
  font-weight: 500;
  color: #2c3e50;
  flex: 1;
}

.item-quantity {
  font-size: 14px;
  color: #1890ff;
  font-weight: 600;
  background-color: #e6f7ff;
  padding: 4px 12px;
  border-radius: 12px;
}

.item-description {
  font-size: 14px;
  color: #5a6c7d;
  line-height: 1.6;
  padding-top: 8px;
  border-top: 1px solid #e8eaed;
}

.edit-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.2s;
  margin-bottom: 8px;
}

.edit-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.quantity-input {
  width: 80px !important;
  margin-left: 8px;
}

.edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.save-btn,
.cancel-btn {
  flex: 1;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.save-btn {
  background-color: #1890ff;
  color: white;
}

.save-btn:hover {
  background-color: #40a9ff;
}

.cancel-btn {
  background-color: #f5f5f5;
  color: #5a6c7d;
}

.cancel-btn:hover {
  background-color: #e8e8e8;
}

.delete-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 24px;
  height: 24px;
  border: none;
  background-color: #ff4d4f;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  opacity: 0;
  transition: all 0.2s;
}

.item-card:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background-color: #ff7875;
  transform: scale(1.1);
}

.empty-items {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  color: #8c8c8c;
}

.empty-items p {
  margin: 8px 0;
}

.empty-items .tip {
  font-size: 14px;
  color: #bfbfbf;
}

.add-btn {
  width: 100%;
  padding: 14px 24px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.add-btn:hover {
  background-color: #40a9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.add-btn span {
  font-size: 20px;
  font-weight: 600;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #8c8c8c;
  font-size: 16px;
}
</style>