<template>
  <div class="weapon-page">
    <WeaponSidebar
      :weapons="weapons"
      :active-weapon-id="activeWeaponId"
      @select="handleSelectWeapon"
      @add="handleAddWeapon"
      @delete="handleDeleteWeapon"
    />
    <div class="main-content">
      <div v-if="activeWeapon" class="weapon-detail">
        <!-- 武器基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>基本信息</h2>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <label>名称</label>
              <div class="info-value">{{ activeWeapon.name }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'holder' }" @dblclick="startEditBasicInfo('holder')">
              <label>所有人</label>
              <div v-if="editingBasicInfo.field !== 'holder'" class="info-value">{{ activeWeapon.holder || '无' }}</div>
              <input
                v-else
                v-model="editingBasicInfo.holder"
                class="edit-input"
                @blur="saveBasicInfo('holder')"
                @keyup.enter="saveBasicInfo('holder')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'level' }" @dblclick="startEditBasicInfo('level')">
              <label>等级</label>
              <div v-if="editingBasicInfo.field !== 'level'" class="info-value">Lv.{{ activeWeapon.level }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.level"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('level')"
                @keyup.enter="saveBasicInfo('level')"
              />
            </div>
          </div>
        </section>

        <!-- 武器属性 -->
        <section class="attributes-section">
          <div class="section-header">
            <h2>属性</h2>
          </div>
          <div class="attributes-list">
            <div
              v-for="attr in activeWeapon.attributes"
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
          </div>
          <button class="add-btn" @click="showAttributeModal = true">
            <span>+</span> 添加属性
          </button>
        </section>

        <!-- 技能部分 -->
        <section class="skills-section">
          <div class="section-header">
            <h2>技能</h2>
          </div>
          <div class="skills-list">
            <div
              v-for="skill in activeWeapon.skills"
              :key="skill.id"
              class="skill-item"
            >
              <div v-if="editingSkill.id !== skill.id" class="skill-content" @dblclick="startEditSkill(skill)">
                <div class="skill-name">{{ skill.name }}</div>
                <div v-if="skill.description" class="skill-description">{{ skill.description }}</div>
              </div>
              <div v-else class="skill-content editing">
                <input
                  v-model="editingSkill.name"
                  class="edit-input"
                  placeholder="技能名称"
                />
                <input
                  v-model="editingSkill.description"
                  class="edit-input"
                  placeholder="技能描述"
                />
                <div class="edit-actions">
                  <button class="save-btn" @click="saveSkill">保存</button>
                  <button class="cancel-btn" @click="cancelEditSkill">取消</button>
                </div>
              </div>
              <button v-if="editingSkill.id !== skill.id" class="delete-btn" @click="handleDeleteSkill(skill.id)" title="删除技能">
                ×
              </button>
            </div>
            <div v-if="activeWeapon.skills.length === 0" class="empty-skills">
              <p>暂无技能</p>
              <p class="tip">点击下方按钮添加</p>
            </div>
          </div>
          <button class="add-btn" @click="showSkillModal = true">
            <span>+</span> 添加技能
          </button>
        </section>
      </div>
      <div v-else class="empty-state">
        <p>请选择或添加一个武器</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddWeaponModal
      v-model:visible="showWeaponModal"
      @confirm="handleAddWeaponConfirm"
    />
    <AddAttributeModal
      v-model:visible="showAttributeModal"
      @confirm="handleAddAttribute"
    />
    <AddSkillModal
      v-model:visible="showSkillModal"
      @confirm="handleAddSkill"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import WeaponSidebar from '../components/WeaponSidebar.vue'
import AddAttributeModal from '../components/AddAttributeModal.vue'
import AddSkillModal from '../components/AddSkillModal.vue'
import AddWeaponModal from '../components/AddWeaponModal.vue'

// 武器数据
const weapons = ref([])
const activeWeaponId = ref(null)
const showAttributeModal = ref(false)
const showSkillModal = ref(false)
const showWeaponModal = ref(false)
const loading = ref(false)

// 基本信息编辑状态
const editingBasicInfo = ref({
  field: '', // 当前正在编辑的字段
  holder: '',
  level: 0
})

// 属性编辑状态
const editingAttribute = ref({
  id: null,
  name: '',
  description: '',
  value: 0
})

// 技能编辑状态
const editingSkill = ref({
  id: null,
  name: '',
  description: ''
})

const activeWeapon = computed(() => {
  return weapons.value.find(w => w.id === activeWeaponId.value)
})

// 组件挂载时加载武器列表
onMounted(async () => {
  await loadWeapons()
})

// 从数据库加载武器列表
async function loadWeapons() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllWeapons()
    weapons.value = result.map(weapon => ({
      id: parseInt(weapon.id),
      name: weapon.name,
      holder: weapon.holder,
      level: parseInt(weapon.level),
      attributes: [],
      skills: []
    }))
    
    // 如果有数据，默认选中第一个
    if (weapons.value.length > 0 && !activeWeaponId.value) {
      activeWeaponId.value = weapons.value[0].id
      await loadWeaponDetail(activeWeaponId.value)
    }
  } catch (error) {
    console.error('加载武器列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 从数据库加载武器详细信息
async function loadWeaponDetail(weaponId) {
  try {
    const result = await window.go.main.app.GetWeaponInfo(parseInt(weaponId))
    const weapon = weapons.value.find(w => w.id === parseInt(weaponId))
    if (weapon) {
      // 更新基本信息
      weapon.name = result.name
      weapon.holder = result.holder
      weapon.level = parseInt(result.level)
      // 更新属性和技能
      weapon.attributes = result.attributes || []
      weapon.skills = result.skills || []
    }
  } catch (error) {
    console.error('加载武器详情失败:', error)
  }
}

// 选择武器
async function handleSelectWeapon(id) {
  activeWeaponId.value = id
  await loadWeaponDetail(id)
}

// 添加武器 - 显示弹窗
function handleAddWeapon() {
  showWeaponModal.value = true
}

// 确认添加武器
async function handleAddWeaponConfirm(data) {
  console.log('开始添加武器，数据:', data)
  try {
    console.log('调用 CreateWeapon 接口，参数:', data.name, data.holder, data.level)
    const newId = await window.go.main.app.CreateWeapon(
      data.name,
      data.holder,
      parseInt(data.level)
    )
    console.log('创建武器成功，ID:', newId)
    
    // 重新加载列表
    await loadWeapons()
    activeWeaponId.value = parseInt(newId)
    await loadWeaponDetail(newId)
  } catch (error) {
    console.error('添加武器失败:', error)
    alert('添加武器失败: ' + error.message)
  }
}

// 删除武器
async function handleDeleteWeapon(weaponId) {
  if (!confirm('确定要删除这个武器吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteWeapon(weaponId)
    
    // 如果删除的是当前选中的武器，清空选中状态
    if (activeWeaponId.value === weaponId) {
      activeWeaponId.value = null
    }
    
    // 重新加载列表
    await loadWeapons()
  } catch (error) {
    console.error('删除武器失败:', error)
    alert('删除武器失败: ' + error.message)
  }
}

// 添加属性
async function handleAddAttribute(data) {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.AddWeaponAttribute(
      activeWeapon.value.id,
      data.name,
      data.description,
      parseInt(data.value)
    )
    // 重新加载属性列表
    await loadWeaponDetail(activeWeapon.value.id)
  } catch (error) {
    console.error('添加属性失败:', error)
  }
}

// 删除属性
async function handleDeleteAttribute(id) {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.DeleteWeaponAttribute(id)
    // 重新加载属性列表
    await loadWeaponDetail(activeWeapon.value.id)
  } catch (error) {
    console.error('删除属性失败:', error)
  }
}

// 开始编辑属性
function startEditAttribute(attr) {
  editingAttribute.value = {
    id: attr.id,
    name: attr.name,
    description: attr.description || '',
    value: attr.value
  }
}

// 保存属性
async function saveAttribute() {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.UpdateWeaponAttribute(
      editingAttribute.value.id,
      editingAttribute.value.name,
      editingAttribute.value.description,
      editingAttribute.value.value
    )
    // 重新加载属性列表
    await loadWeaponDetail(activeWeapon.value.id)
    // 清除编辑状态
    editingAttribute.value = {
      id: null,
      name: '',
      description: '',
      value: 0
    }
  } catch (error) {
    console.error('更新属性失败:', error)
  }
}

// 取消编辑属性
function cancelEditAttribute() {
  editingAttribute.value = {
    id: null,
    name: '',
    description: '',
    value: 0
  }
}

// 添加技能
async function handleAddSkill(data) {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.AddWeaponSkill(
      activeWeapon.value.id,
      data.name,
      data.description
    )
    // 重新加载技能列表
    await loadWeaponDetail(activeWeapon.value.id)
  } catch (error) {
    console.error('添加技能失败:', error)
  }
}

// 删除技能
async function handleDeleteSkill(id) {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.DeleteWeaponSkill(id)
    // 重新加载技能列表
    await loadWeaponDetail(activeWeapon.value.id)
  } catch (error) {
    console.error('删除技能失败:', error)
  }
}

// 开始编辑技能
function startEditSkill(skill) {
  editingSkill.value = {
    id: skill.id,
    name: skill.name,
    description: skill.description || ''
  }
}

// 保存技能
async function saveSkill() {
  if (!activeWeapon.value) return
  try {
    await window.go.main.app.UpdateWeaponSkill(
      editingSkill.value.id,
      editingSkill.value.name,
      editingSkill.value.description
    )
    // 重新加载技能列表
    await loadWeaponDetail(activeWeapon.value.id)
    // 清除编辑状态
    editingSkill.value = {
      id: null,
      name: '',
      description: ''
    }
  } catch (error) {
    console.error('更新技能失败:', error)
  }
}

// 取消编辑技能
function cancelEditSkill() {
  editingSkill.value = {
    id: null,
    name: '',
    description: ''
  }
}

// 开始编辑基本信息
function startEditBasicInfo(field) {
  if (!activeWeapon.value) return
  editingBasicInfo.value.field = field
  if (field === 'holder') {
    editingBasicInfo.value.holder = activeWeapon.value.holder || ''
  } else if (field === 'level') {
    editingBasicInfo.value.level = activeWeapon.value.level
  }
}

// 保存基本信息
async function saveBasicInfo(field) {
  if (!activeWeapon.value) return
  
  try {
    if (field === 'holder') {
      await window.go.main.app.UpdateWeaponBasicInfo(
        activeWeapon.value.id,
        activeWeapon.value.name,
        editingBasicInfo.value.holder,
        activeWeapon.value.level
      )
    } else if (field === 'level') {
      await window.go.main.app.UpdateWeaponBasicInfo(
        activeWeapon.value.id,
        activeWeapon.value.name,
        activeWeapon.value.holder,
        editingBasicInfo.value.level
      )
    }
    
    // 重新加载武器详情
    await loadWeaponDetail(activeWeapon.value.id)
    
    // 清除编辑状态
    editingBasicInfo.value = {
      field: '',
      holder: '',
      level: 0
    }
  } catch (error) {
    console.error('更新基本信息失败:', error)
  }
}
</script>

<style scoped>
.weapon-page {
  display: flex;
  height: 100%;
  background-color: var(--app-bg);
}

.main-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.weapon-detail {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.info-section {
  grid-column: 1;
}

.attributes-section,
.skills-section {
  grid-column: 1;
  display: flex;
  flex-direction: column;
  min-height: 400px;
}

/* 横纵均匀布局：属性和技能并排显示 */
@media (min-width: 1024px) {
  .weapon-detail {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
  }

  .info-section {
    grid-column: 1 / -1;
    grid-row: 1;
  }

  .attributes-section {
    grid-column: 1;
    grid-row: 2;
  }

  .skills-section {
    grid-column: 2;
    grid-row: 2;
  }
}

.info-section,
.attributes-section,
.skills-section {
  background-color: var(--app-surface);
  border-radius: 8px;
  padding: 24px;
  box-shadow: var(--app-shadow-sm);
  height: fit-content;
}

.section-header {
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--app-border);
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
}

.info-item {
  padding: 16px;
  border-radius: 8px;
  border: 1px solid var(--app-border);
  transition: all 0.3s ease;
  background: linear-gradient(135deg, var(--app-surface) 0%, var(--app-item-bg) 100%);
}

.info-item:hover {
  border-color: var(--app-accent);
  transform: translateY(-2px);
  box-shadow: var(--app-shadow-accent-lg);
}

.info-item.editing {
  border-color: var(--app-accent);
  background-color: var(--app-hover-bg);
  box-shadow: var(--app-shadow-accent-md);
}

.info-item label {
  display: block;
  font-size: 12px;
  color: var(--app-text-muted);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 16px;
  color: var(--app-text-primary);
  font-weight: 500;
}

.edit-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--app-border);
  border-radius: 4px;
  font-size: 14px;
  color: var(--app-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.edit-input:focus {
  border-color: var(--app-accent);
}

.attributes-list,
.skills-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
  overflow-y: auto;
  max-height: 500px;
  padding-right: 8px;
}

/* 自定义滚动条样式 */
.attributes-list::-webkit-scrollbar,
.skills-list::-webkit-scrollbar {
  width: 6px;
}

.attributes-list::-webkit-scrollbar-track,
.skills-list::-webkit-scrollbar-track {
  background: var(--app-divider-soft);
  border-radius: 3px;
}

.attributes-list::-webkit-scrollbar-thumb,
.skills-list::-webkit-scrollbar-thumb {
  background: var(--app-border);
  border-radius: 3px;
}

.attributes-list::-webkit-scrollbar-thumb:hover,
.skills-list::-webkit-scrollbar-thumb:hover {
  background: var(--app-text-muted);
}

.attribute-item,
.skill-item {
  display: flex;
  align-items: center;
  padding: 18px;
  background-color: var(--app-item-bg);
  border-radius: 8px;
  border: 1px solid var(--app-border);
  transition: all 0.3s ease;
  position: relative;
}

.attribute-item:hover,
.skill-item:hover {
  border-color: var(--app-accent);
  background-color: var(--app-hover-bg);
  transform: translateX(4px);
  box-shadow: var(--app-shadow-accent-sm);
}

.attribute-item.editing,
.skill-item.editing {
  border-color: var(--app-accent);
  background-color: var(--app-hover-bg);
  box-shadow: var(--app-shadow-accent-md);
}

.attribute-content,
.skill-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.attribute-content.editing,
.skill-content.editing {
  flex-direction: column;
  align-items: stretch;
}

.attribute-name,
.skill-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--app-text-primary);
  min-width: 120px;
}

.attribute-description,
.skill-description {
  font-size: 14px;
  color: var(--app-text-secondary);
  flex: 1;
}

.attribute-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--app-accent);
  min-width: 60px;
  text-align: right;
}

.edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.save-btn,
.cancel-btn {
  padding: 6px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.save-btn {
  background-color: var(--app-accent);
  color: var(--app-text-inverse);
}

.save-btn:hover {
  background-color: var(--app-accent-hover);
}

.cancel-btn {
  background-color: var(--app-muted-bg);
  color: var(--app-text-secondary);
}

.cancel-btn:hover {
  background-color: var(--app-divider-soft);
}

.delete-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 24px;
  height: 24px;
  border: none;
  background-color: var(--app-danger);
  color: var(--app-text-inverse);
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  opacity: 0;
  transition: all 0.2s;
}

.attribute-item:hover .delete-btn,
.skill-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background-color: var(--app-danger-hover);
}

.add-btn {
  width: 100%;
  padding: 14px;
  border: 2px dashed var(--app-border);
  border-radius: 8px;
  background-color: var(--app-surface);
  color: var(--app-text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 12px;
  flex-shrink: 0;
}

.add-btn:hover {
  border-color: var(--app-accent);
  color: var(--app-accent);
  background-color: var(--app-hover-bg);
  transform: translateY(-2px);
  box-shadow: var(--app-shadow-accent-md);
}

.add-btn:active {
  transform: translateY(0);
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: var(--app-text-secondary);
  font-size: 16px;
}

.empty-skills {
  text-align: center;
  padding: 40px 20px;
  color: var(--app-text-secondary);
}

.empty-skills p {
  margin: 0 0 8px 0;
}

.empty-skills .tip {
  font-size: 12px;
  color: var(--app-border);
}
</style>
