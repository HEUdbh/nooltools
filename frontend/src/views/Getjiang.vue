<template>
  <div class="getjiang-page">
    <div class="main-content">
      <!-- 抽奖区域 -->
      <section class="draw-section">
        <div class="section-header">
          <h2>🎰 抽奖</h2>
          <div class="draw-buttons">
            <button class="draw-btn once" @click="handleDrawOnce" :disabled="drawing">
              {{ drawing ? '抽奖中...' : '单抽' }}
            </button>
            <button class="draw-btn ten" @click="handleDrawTen" :disabled="drawing">
              {{ drawing ? '抽奖中...' : '十连抽' }}
            </button>
          </div>
        </div>

        <!-- 抽奖结果展示 -->
        <div v-if="drawResults.length > 0" class="draw-results">
          <div v-for="(result, index) in drawResults" :key="index" class="draw-result-item" :class="result.variety">
            <div class="result-name">{{ result.name }}</div>
            <div class="result-info">
              <span class="result-variety">{{ result.variety }}</span>
              <span class="result-rate">爆率: {{ (result.rate * 100).toFixed(2) }}%</span>
            </div>
            <div v-if="result.description" class="result-description">{{ result.description }}</div>
          </div>
        </div>
      </section>

      <!-- 奖池管理区域 -->
      <section class="prize-pool-section">
        <div class="section-header">
          <h2>🎁 奖池管理</h2>
          <button class="add-btn" @click="handleAddPrize">添加奖品</button>
        </div>

        <!-- 奖品列表 -->
        <div class="prize-list">
          <div v-for="prize in prizes" :key="prize.id" class="prize-item" :class="prize.variety">
            <div class="prize-info">
              <div class="prize-name">{{ prize.name }}</div>
              <div class="prize-details">
                <span class="prize-variety">{{ prize.variety }}</span>
                <span class="prize-rate">爆率: {{ (prize.rate * 100).toFixed(2) }}%</span>
              </div>
              <div v-if="prize.description" class="prize-description">{{ prize.description }}</div>
            </div>
            <div class="prize-actions">
              <button class="edit-btn" @click="handleEditPrize(prize)" title="编辑">✏️</button>
              <button class="delete-btn" @click="handleDeletePrize(prize.id)" title="删除">🗑️</button>
            </div>
          </div>
          <div v-if="prizes.length === 0" class="empty-state">
            <p>暂无奖品，请添加奖品</p>
          </div>
        </div>
      </section>

      <!-- 抽奖历史区域 -->
      <section class="history-section">
        <div class="section-header">
          <h2>📜 抽奖历史</h2>
          <button class="clear-btn" @click="handleClearHistory" :disabled="drawHistory.length === 0">
            清空历史
          </button>
        </div>

        <div class="history-list">
          <div v-for="(item, index) in drawHistory" :key="item.id" class="history-item" :class="item.variety">
            <div class="history-index">#{{ drawHistory.length - index }}</div>
            <div class="history-name">{{ item.name }}</div>
            <div class="history-time">{{ formatTime(item.drawn_at) }}</div>
          </div>
          <div v-if="drawHistory.length === 0" class="empty-state">
            <p>暂无抽奖历史</p>
          </div>
        </div>
      </section>
    </div>

    <!-- 弹窗组件 -->
    <PrizeModal
      v-model:visible="showPrizeModal"
      :prize="editingPrize"
      @confirm="handlePrizeConfirm"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PrizeModal from '../components/PrizeModal.vue'

// 数据
const prizes = ref([])
const drawResults = ref([])
const drawHistory = ref([])
const drawing = ref(false)
const showPrizeModal = ref(false)
const editingPrize = ref(null)

// 组件挂载时加载数据
onMounted(async () => {
  await loadPrizes()
  await loadDrawHistory()
})

// 加载奖品列表
async function loadPrizes() {
  try {
    const result = await window.go.main.app.GetAllPrizes()
    prizes.value = result.map(p => ({
      id: parseInt(p.id),
      name: p.name,
      rate: parseFloat(p.rate),
      description: p.description || '',
      variety: p.variety || '普通'
    }))
  } catch (error) {
    console.error('加载奖品列表失败:', error)
  }
}

// 加载抽奖历史
async function loadDrawHistory() {
  try {
    const result = await window.go.main.app.GetDrawHistory()
    drawHistory.value = result.map(h => ({
      id: parseInt(h.id),
      name: h.name,
      variety: h.variety || '普通',
      rate: parseFloat(h.rate),
      description: h.description || '',
      drawn_at: h.drawn_at
    }))
  } catch (error) {
    console.error('加载抽奖历史失败:', error)
  }
}

// 单抽
async function handleDrawOnce() {
  if (drawing.value) return
  if (prizes.value.length === 0) {
    alert('奖池中没有奖品，请先添加奖品')
    return
  }

  drawing.value = true
  try {
    const result = await window.go.main.app.DrawOnce()
    drawResults.value = [{
      id: parseInt(result.id),
      name: result.name,
      rate: parseFloat(result.rate),
      description: result.description || '',
      variety: result.variety || '普通'
    }]

    // 重新加载抽奖历史
    await loadDrawHistory()
  } catch (error) {
    console.error('单抽失败:', error)
    alert('抽奖失败: ' + error.message)
  } finally {
    drawing.value = false
  }
}

// 十连抽
async function handleDrawTen() {
  if (drawing.value) return
  if (prizes.value.length === 0) {
    alert('奖池中没有奖品，请先添加奖品')
    return
  }

  drawing.value = true
  try {
    const results = await window.go.main.app.DrawTen()
    drawResults.value = results.map(r => ({
      id: parseInt(r.id),
      name: r.name,
      rate: parseFloat(r.rate),
      description: r.description || '',
      variety: r.variety || '普通'
    }))

    // 重新加载抽奖历史
    await loadDrawHistory()
  } catch (error) {
    console.error('十连抽失败:', error)
    alert('抽奖失败: ' + error.message)
  } finally {
    drawing.value = false
  }
}

// 添加奖品
function handleAddPrize() {
  editingPrize.value = null
  showPrizeModal.value = true
}

// 编辑奖品
function handleEditPrize(prize) {
  editingPrize.value = { ...prize }
  showPrizeModal.value = true
}

// 确认添加/编辑奖品
async function handlePrizeConfirm(data) {
  try {
    if (editingPrize.value) {
      // 编辑奖品
      await window.go.main.app.UpdatePrize(
        editingPrize.value.id,
        data.name,
        data.rate,
        data.description,
        data.variety
      )
    } else {
      // 添加奖品
      await window.go.main.app.CreatePrize(
        data.name,
        data.rate,
        data.description,
        data.variety
      )
    }

    // 重新加载奖品列表
    await loadPrizes()
  } catch (error) {
    console.error('保存奖品失败:', error)
    alert('保存奖品失败: ' + error.message)
  }
}

// 删除奖品
async function handleDeletePrize(prizeId) {
  if (!confirm('确定要删除这个奖品吗？删除后无法恢复。')) {
    return
  }

  try {
    await window.go.main.app.DeletePrize(prizeId)
    await loadPrizes()
  } catch (error) {
    console.error('删除奖品失败:', error)
    alert('删除奖品失败: ' + error.message)
  }
}

// 清空抽奖历史
async function handleClearHistory() {
  if (!confirm('确定要清空所有抽奖历史吗？')) {
    return
  }

  try {
    await window.go.main.app.ClearDrawHistory()
    await loadDrawHistory()
  } catch (error) {
    console.error('清空历史失败:', error)
    alert('清空历史失败: ' + error.message)
  }
}

// 格式化时间
function formatTime(timeStr) {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
</script>

<style scoped>
.getjiang-page {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #f5f5f5;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.section-header h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

/* 抽奖区域 */
.draw-section {
  margin-bottom: 30px;
}

.draw-buttons {
  display: flex;
  gap: 15px;
}

.draw-btn {
  padding: 12px 30px;
  font-size: 16px;
  font-weight: bold;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.draw-btn.once {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.draw-btn.ten {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.draw-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.draw-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.draw-results {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.draw-result-item {
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.draw-result-item.传说 {
  border: 3px solid #ffd700;
  background: linear-gradient(135deg, #fff9c4 0%, #fff59d 100%);
}

.draw-result-item.史诗 {
  border: 3px solid #ff9800;
  background: linear-gradient(135deg, #ffe0b2 0%, #ffcc80 100%);
}

.draw-result-item.稀有 {
  border: 3px solid #9c27b0;
  background: linear-gradient(135deg, #e1bee7 0%, #ce93d8 100%);
}

.draw-result-item.普通 {
  border: 2px solid #9e9e9e;
}

.result-name {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
}

.result-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
  color: #666;
}

.result-variety {
  font-weight: bold;
}

.result-description {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}

/* 奖池管理区域 */
.prize-pool-section {
  margin-bottom: 30px;
}

.add-btn {
  padding: 8px 16px;
  background: #4caf50;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.add-btn:hover {
  background: #45a049;
  transform: translateY(-2px);
}

.prize-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 15px;
}

.prize-item {
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  transition: all 0.3s;
}

.prize-item:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.prize-item.传说 {
  border-left: 4px solid #ffd700;
}

.prize-item.史诗 {
  border-left: 4px solid #ff9800;
}

.prize-item.稀有 {
  border-left: 4px solid #9c27b0;
}

.prize-item.普通 {
  border-left: 4px solid #9e9e9e;
}

.prize-info {
  flex: 1;
}

.prize-name {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
}

.prize-details {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
  font-size: 14px;
  color: #666;
}

.prize-variety {
  font-weight: bold;
}

.prize-description {
  font-size: 13px;
  color: #999;
}

.prize-actions {
  display: flex;
  gap: 8px;
}

.edit-btn,
.delete-btn {
  padding: 6px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.edit-btn {
  background: #2196f3;
  color: white;
}

.edit-btn:hover {
  background: #1976d2;
}

.delete-btn {
  background: #f44336;
  color: white;
}

.delete-btn:hover {
  background: #d32f2f;
}

/* 抽奖历史区域 */
.history-section {
  margin-bottom: 30px;
}

.clear-btn {
  padding: 8px 16px;
  background: #f44336;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.clear-btn:hover:not(:disabled) {
  background: #d32f2f;
  transform: translateY(-2px);
}

.clear-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  background: white;
  border-radius: 6px;
  margin-bottom: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.history-item.传说 {
  border-left: 4px solid #ffd700;
}

.history-item.史诗 {
  border-left: 4px solid #ff9800;
}

.history-item.稀有 {
  border-left: 4px solid #9c27b0;
}

.history-item.普通 {
  border-left: 4px solid #9e9e9e;
}

.history-index {
  font-weight: bold;
  color: #666;
  min-width: 50px;
}

.history-name {
  flex: 1;
  font-weight: bold;
  color: #333;
}

.history-time {
  font-size: 12px;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}
</style>