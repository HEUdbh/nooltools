<template>
  <div class="shopping-page">
    <ShoppingSidebar
      :shopping="shopping"
      :active-shopping-id="activeShoppingId"
      @select="handleSelectShopping"
      @add="handleAddShopping"
      @delete="handleDeleteShopping"
    />
    <div class="main-content">
      <div v-if="activeShopping" class="shopping-detail">
        <!-- 商品基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>{{ activeShopping.name }}</h2>
          </div>
        </section>

        <!-- 商品详细信息 -->
        <section class="detail-section">
          <div v-if="editingShopping.id !== activeShopping.id" class="shopping-content" @dblclick="startEditShopping">
            <div class="info-row">
              <span class="label">价值：</span>
              <span class="value">{{ activeShopping.value }}</span>
            </div>
            <div class="info-row">
              <span class="label">描述：</span>
              <span class="value">{{ activeShopping.description || '暂无描述' }}</span>
            </div>
            <div class="info-row">
              <span class="label">购买条件：</span>
              <span class="value">{{ activeShopping.condition || '无限制' }}</span>
            </div>
          </div>
          <div v-else class="shopping-content editing">
            <div class="info-row">
              <span class="label">名称：</span>
              <input
                v-model="editingShopping.name"
                class="edit-input"
                placeholder="商品名称"
              />
            </div>
            <div class="info-row">
              <span class="label">价值：</span>
              <input
                v-model.number="editingShopping.value"
                type="number"
                class="edit-input"
                placeholder="商品价值"
              />
            </div>
            <div class="info-row">
              <span class="label">描述：</span>
              <textarea
                v-model="editingShopping.description"
                class="edit-input textarea"
                placeholder="商品描述"
                rows="3"
              ></textarea>
            </div>
            <div class="info-row">
              <span class="label">购买条件：</span>
              <textarea
                v-model="editingShopping.condition"
                class="edit-input textarea"
                placeholder="购买条件"
                rows="2"
              ></textarea>
            </div>
            <div class="edit-actions">
              <button class="save-btn" @click="saveShopping">保存</button>
              <button class="cancel-btn" @click="cancelEditShopping">取消</button>
            </div>
          </div>
          <button v-if="editingShopping.id !== activeShopping.id" class="delete-btn" @click="handleDeleteShopping(activeShopping.id)" title="删除商品">
            ×
          </button>
        </section>
      </div>
      <div v-else class="empty-state">
        <p>请选择或添加一个商品</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddShoppingModal
      v-model:visible="showShoppingModal"
      @confirm="handleAddShoppingConfirm"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ShoppingSidebar from '../components/ShoppingSidebar.vue'
import AddShoppingModal from '../components/AddShoppingModal.vue'

// 商城数据
const shopping = ref([])
const activeShoppingId = ref(null)
const showShoppingModal = ref(false)
const loading = ref(false)

// 商品编辑状态
const editingShopping = ref({
  id: null,
  name: '',
  value: 0,
  description: '',
  condition: ''
})

const activeShopping = computed(() => {
  return shopping.value.find(s => s.id === activeShoppingId.value)
})

// 组件挂载时加载商城列表
onMounted(async () => {
  await loadShopping()
})

// 从数据库加载商城列表
async function loadShopping() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllShopping()
    shopping.value = result.map(s => ({
      id: parseInt(s.id),
      name: s.name,
      value: parseInt(s.value),
      description: s.description || '',
      condition: s.condition || ''
    }))
    
    // 如果有数据，默认选中第一个
    if (shopping.value.length > 0 && !activeShoppingId.value) {
      activeShoppingId.value = shopping.value[0].id
    }
  } catch (error) {
    console.error('加载商城列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 选择商品
function handleSelectShopping(id) {
  activeShoppingId.value = id
}

// 添加商品 - 显示弹窗
function handleAddShopping() {
  showShoppingModal.value = true
}

// 确认添加商品
async function handleAddShoppingConfirm(data) {
  try {
    await window.go.main.app.CreateShopping(data.name, data.value, data.description, data.condition)
    
    // 重新加载列表
    await loadShopping()
  } catch (error) {
    console.error('添加商品失败:', error)
    alert('添加商品失败: ' + error.message)
  }
}

// 删除商品
async function handleDeleteShopping(shoppingId) {
  if (!confirm('确定要删除这个商品吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteShopping(shoppingId)
    
    // 如果删除的是当前选中的商品，清空选中状态
    if (activeShoppingId.value === shoppingId) {
      activeShoppingId.value = null
    }
    
    // 重新加载列表
    await loadShopping()
  } catch (error) {
    console.error('删除商品失败:', error)
    alert('删除商品失败: ' + error.message)
  }
}

// 开始编辑商品
function startEditShopping() {
  if (!activeShopping.value) return
  
  editingShopping.value = {
    id: activeShopping.value.id,
    name: activeShopping.value.name,
    value: activeShopping.value.value,
    description: activeShopping.value.description,
    condition: activeShopping.value.condition
  }
}

// 取消编辑商品
function cancelEditShopping() {
  editingShopping.value = {
    id: null,
    name: '',
    value: 0,
    description: '',
    condition: ''
  }
}

// 保存商品
async function saveShopping() {
  if (!editingShopping.value.name) {
    alert('请输入商品名称')
    return
  }
  
  if (editingShopping.value.value === undefined || editingShopping.value.value === null) {
    alert('请输入商品价值')
    return
  }
  
  try {
    await window.go.main.app.UpdateShopping(
      editingShopping.value.id,
      editingShopping.value.name,
      editingShopping.value.value,
      editingShopping.value.description,
      editingShopping.value.condition
    )
    
    // 重新加载列表
    await loadShopping()
    
    // 清除编辑状态
    cancelEditShopping()
  } catch (error) {
    console.error('更新商品失败:', error)
    alert('更新商品失败: ' + error.message)
  }
}
</script>

<style scoped>
.shopping-page {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f5f5;
}

.shopping-detail {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-header {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #e0e0e0;
}

.section-header h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.info-section {
  margin-bottom: 20px;
}

.detail-section {
  position: relative;
}

.shopping-content {
  padding: 15px;
  background-color: #fafafa;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
  cursor: pointer;
  transition: all 0.2s;
}

.shopping-content:hover {
  background-color: #f0f0f0;
  border-color: #d0d0d0;
}

.shopping-content.editing {
  cursor: default;
  background-color: white;
  border: 2px solid #4CAF50;
}

.info-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 15px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  min-width: 100px;
  font-weight: bold;
  color: #666;
  padding-top: 5px;
}

.value {
  flex: 1;
  color: #333;
  line-height: 1.6;
}

.edit-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.edit-input:focus {
  outline: none;
  border-color: #4CAF50;
}

.edit-input.textarea {
  resize: vertical;
  min-height: 60px;
  font-family: inherit;
}

.edit-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #e0e0e0;
}

.save-btn,
.cancel-btn {
  padding: 8px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.save-btn {
  background-color: #4CAF50;
  color: white;
}

.save-btn:hover {
  background-color: #45a049;
}

.cancel-btn {
  background-color: #f5f5f5;
  color: #666;
}

.cancel-btn:hover {
  background-color: #e0e0e0;
}

.delete-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 30px;
  height: 30px;
  border: none;
  background-color: #ff5252;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 20px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.delete-btn:hover {
  background-color: #ff1744;
  transform: scale(1.1);
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
  font-size: 16px;
}
</style>