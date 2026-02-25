<template>
  <div class="character-page">
    <CharacterSidebar
      :characters="characters"
      :active-character-id="activeCharacterId"
      @select="handleSelectCharacter"
      @add="handleAddCharacter"
      @delete="handleDeleteCharacter"
    />
    <div class="main-content">
      <div v-if="activeCharacter" class="character-detail">
        <!-- 人物基本信息 -->
        <section class="info-section">
          <div class="section-header">
            <h2>基本信息</h2>
          </div>
          <div class="info-grid">
            <div class="info-item">
              <label>名称</label>
              <div class="info-value">{{ activeCharacter.name }}</div>
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'faction' }" @dblclick="startEditBasicInfo('faction')">
              <label>势力</label>
              <div v-if="editingBasicInfo.field !== 'faction'" class="info-value">{{ activeCharacter.faction || '无' }}</div>
              <input
                v-else
                v-model="editingBasicInfo.faction"
                class="edit-input"
                @blur="saveBasicInfo('faction')"
                @keyup.enter="saveBasicInfo('faction')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'wealth' }" @dblclick="startEditBasicInfo('wealth')">
              <label>财产</label>
              <div v-if="editingBasicInfo.field !== 'wealth'" class="info-value">{{ activeCharacter.wealth || 0 }}</div>
              <input
                v-else
                v-model.number="editingBasicInfo.wealth"
                type="number"
                class="edit-input"
                @blur="saveBasicInfo('wealth')"
                @keyup.enter="saveBasicInfo('wealth')"
              />
            </div>
            <div class="info-item" :class="{ editing: editingBasicInfo.field === 'level' }" @dblclick="startEditBasicInfo('level')">
              <label>等级</label>
              <div v-if="editingBasicInfo.field !== 'level'" class="info-value">Lv.{{ activeCharacter.level }}</div>
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

        <!-- 个人属性 -->
        <section class="attributes-section">
          <div class="section-header">
            <h2>个人属性</h2>
          </div>
          <div class="attributes-list">
            <div
              v-for="attr in activeCharacter.attributes"
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
              v-for="skill in activeCharacter.skills"
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
            <div v-if="activeCharacter.skills.length === 0" class="empty-skills">
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
        <p>请选择或添加一个人物</p>
      </div>
    </div>

    <!-- 弹窗组件 -->
    <AddCharacterModal
      v-model:visible="showCharacterModal"
      @confirm="handleAddCharacterConfirm"
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
import { ref, computed, onMounted, nextTick } from 'vue'
import CharacterSidebar from '../components/CharacterSidebar.vue'
import AddAttributeModal from '../components/AddAttributeModal.vue'
import AddSkillModal from '../components/AddSkillModal.vue'
import AddCharacterModal from '../components/AddCharacterModal.vue'

// 人物数据
const characters = ref([])
const activeCharacterId = ref(null)
const showAttributeModal = ref(false)
const showSkillModal = ref(false)
const showCharacterModal = ref(false)
const loading = ref(false)

// 基本信息编辑状态
const editingBasicInfo = ref({
  field: '', // 当前正在编辑的字段
  name: '',
  faction: '',
  wealth: 0,
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

const activeCharacter = computed(() => {
  return characters.value.find(c => c.id === activeCharacterId.value)
})

// 组件挂载时加载人物列表
onMounted(async () => {
  await loadCharacters()
})

// 从数据库加载人物列表
async function loadCharacters() {
  loading.value = true
  try {
    const result = await window.go.main.app.GetAllCharacters()
    characters.value = result.map(char => ({
      id: parseInt(char.id),
      name: char.name,
      faction: char.shili,
      wealth: parseInt(char.property),
      level: parseInt(char.level),
      attributes: [],
      skills: []
    }))
    
    // 如果有数据，默认选中第一个
    if (characters.value.length > 0 && !activeCharacterId.value) {
      activeCharacterId.value = characters.value[0].id
      await loadCharacterDetail(activeCharacterId.value)
    }
  } catch (error) {
    console.error('加载人物列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 从数据库加载人物详细信息
async function loadCharacterDetail(characterId) {
  try {
    const result = await window.go.main.app.GetCharacterInfo(parseInt(characterId))
    const character = characters.value.find(c => c.id === parseInt(characterId))
    if (character) {
      // 更新基本信息
      character.name = result.name
      character.faction = result.shili
      character.wealth = parseInt(result.property)
      character.level = parseInt(result.level)
      // 更新属性和技能
      character.attributes = result.attributes || []
      character.skills = result.skills || []
    }
  } catch (error) {
    console.error('加载人物详情失败:', error)
  }
}

// 选择人物
async function handleSelectCharacter(id) {
  activeCharacterId.value = id
  await loadCharacterDetail(id)
}

// 添加人物 - 显示弹窗
function handleAddCharacter() {
  showCharacterModal.value = true
}

// 确认添加人物
async function handleAddCharacterConfirm(data) {
  console.log('开始添加人物，数据:', data)
  try {
    console.log('调用 CreateCharacter 接口，参数:', data.name, data.faction, data.wealth, data.level)
    const newId = await window.go.main.app.CreateCharacter(
      data.name,
      data.faction,
      parseInt(data.wealth),
      parseInt(data.level)
    )
    console.log('创建人物成功，ID:', newId)
    
    // 重新加载列表
    await loadCharacters()
    activeCharacterId.value = parseInt(newId)
    await loadCharacterDetail(newId)
  } catch (error) {
    console.error('添加人物失败:', error)
    alert('添加人物失败: ' + error.message)
  }
}

// 删除人物
async function handleDeleteCharacter(characterId) {
  if (!confirm('确定要删除这个人物吗？删除后无法恢复。')) {
    return
  }
  
  try {
    await window.go.main.app.DeleteCharacter(characterId)
    
    // 如果删除的是当前选中的人物，清空选中状态
    if (activeCharacterId.value === characterId) {
      activeCharacterId.value = null
    }
    
    // 重新加载列表
    await loadCharacters()
  } catch (error) {
    console.error('删除人物失败:', error)
    alert('删除人物失败: ' + error.message)
  }
}

// 添加属性
async function handleAddAttribute(data) {
  if (!activeCharacter.value) return
  try {
    await window.go.main.app.AddCharacterAttribute(
      activeCharacter.value.id,
      data.name,
      data.description,
      parseInt(data.value)
    )
    // 重新加载属性列表
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('添加属性失败:', error)
  }
}

// 删除属性
async function handleDeleteAttribute(id) {
  if (!activeCharacter.value) return
  try {
    await window.go.main.app.DeleteCharacterAttribute(id)
    // 重新加载属性列表
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('删除属性失败:', error)
  }
}

// 添加技能
async function handleAddSkill(data) {
  if (!activeCharacter.value) return
  try {
    await window.go.main.app.AddCharacterSkill(
      activeCharacter.value.id,
      data.name,
      data.description
    )
    // 重新加载技能列表
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('添加技能失败:', error)
  }
}

// 删除技能
async function handleDeleteSkill(id) {
  if (!activeCharacter.value) return
  try {
    await window.go.main.app.DeleteCharacterSkill(id)
    // 重新加载技能列表
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('删除技能失败:', error)
  }
}

// 开始编辑基本信息
function startEditBasicInfo(field) {
  if (!activeCharacter.value) return
  
  // 设置当前编辑的字段
  editingBasicInfo.value.field = field
  
  // 复制当前值到编辑状态
  editingBasicInfo.value[field] = activeCharacter.value[field]
  
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
  if (!activeCharacter.value) return
  
  const value = editingBasicInfo.value[field]
  
  try {
    // 调用后端接口更新基本信息
    await window.go.main.app.UpdateCharacterBasicInfo(
      parseInt(activeCharacter.value.id),
      activeCharacter.value.name,
      field === 'faction' ? value : activeCharacter.value.faction,
      field === 'wealth' ? parseInt(value) : parseInt(activeCharacter.value.wealth),
      field === 'level' ? parseInt(value) : parseInt(activeCharacter.value.level)
    )
    
    // 重新加载人物信息
    await loadCharacterDetail(activeCharacter.value.id)
    
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
    await window.go.main.app.UpdateCharacterAttribute(
      editingAttribute.value.id,
      editingAttribute.value.name,
      editingAttribute.value.description,
      editingAttribute.value.value
    )
    editingAttribute.value = { id: null, name: '', description: '', value: 0 }
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('保存属性失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑属性
function cancelEditAttribute() {
  editingAttribute.value = { id: null, name: '', description: '', value: 0 }
}

// 开始编辑技能
function startEditSkill(skill) {
  editingSkill.value = {
    id: skill.id,
    name: skill.name,
    description: skill.description || ''
  }
  nextTick(() => {
    const inputs = document.querySelectorAll('.skill-item.editing .edit-input')
    if (inputs.length > 0) {
      inputs[0].focus()
    }
  })
}

// 保存技能
async function saveSkill() {
  try {
    await window.go.main.app.UpdateCharacterSkill(
      editingSkill.value.id,
      editingSkill.value.name,
      editingSkill.value.description
    )
    editingSkill.value = { id: null, name: '', description: '' }
    await loadCharacterDetail(activeCharacter.value.id)
  } catch (error) {
    console.error('保存技能失败:', error)
    alert('保存失败: ' + error.message)
  }
}

// 取消编辑技能
function cancelEditSkill() {
  editingSkill.value = { id: null, name: '', description: '' }
}
</script>

<style scoped>
.character-page {
  display: flex;
  height: 100%;
  background-color: var(--app-bg);
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 40px;
}

.character-detail {
  width: 100%;
}

.info-section,
.attributes-section,
.skills-section {
  background-color: var(--app-surface);
  border-radius: 12px;
  padding: 32px;
  margin-bottom: 28px;
  box-shadow: var(--app-shadow-sm);
}

.section-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin: 0;
}

/* 基本信息 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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
  background-color: var(--app-item-bg);
  border-radius: 8px;
  border: 1px solid var(--app-border);
  transition: all 0.2s ease;
}

.attribute-item:hover {
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-accent-sm);
}

.attribute-content {
  flex: 1;
  min-width: 0;
}

.attribute-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin-bottom: 8px;
}

.attribute-description {
  font-size: 13px;
  color: var(--app-text-muted);
  margin-bottom: 10px;
  line-height: 1.4;
}

.attribute-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--app-accent);
}

/* 技能列表 */
.skills-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.skill-item {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  background-color: var(--app-item-bg);
  border-radius: 8px;
  border: 1px solid var(--app-border);
  transition: all 0.2s ease;
}

.skill-item:hover {
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-accent-sm);
}

.skill-content {
  flex: 1;
  min-width: 0;
}

.skill-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--app-text-primary);
  margin-bottom: 8px;
}

.skill-description {
  font-size: 13px;
  color: var(--app-text-muted);
  line-height: 1.6;
}

.empty-skills {
  grid-column: 1 / -1;
  text-align: center;
  padding: 48px 20px;
  color: var(--app-text-muted);
}

.empty-skills p {
  margin: 10px 0;
}

.tip {
  font-size: 13px;
  color: var(--app-text-muted);
}

/* 删除按钮 */
.delete-btn {
  width: 30px;
  height: 30px;
  border: none;
  background: none;
  color: var(--app-text-muted);
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
  background-color: var(--app-danger-soft-bg);
  color: var(--app-danger);
}

/* 添加按钮 */
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

/* 空状态 */
.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 500px;
  color: var(--app-text-muted);
  font-size: 16px;
}

/* 编辑输入框 */
.edit-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--app-border);
  border-radius: 4px;
  font-size: 14px;
  margin-bottom: 8px;
  transition: all 0.2s ease;
}

.edit-input:focus {
  outline: none;
  border-color: var(--app-accent);
  box-shadow: var(--app-shadow-focus);
}

.edit-input::placeholder {
  color: var(--app-text-muted);
}

/* 编辑状态 */
.attribute-item.editing,
.skill-item.editing {
  border-color: var(--app-accent);
  background-color: var(--app-surface);
  box-shadow: var(--app-shadow-accent-sm);
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
  .character-detail {
    max-width: 900px;
  }
  
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
  }
  
  .attributes-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 24px;
  }
  
  .info-grid,
  .attributes-list,
  .skills-list {
    grid-template-columns: 1fr;
  }
  
  .info-section,
  .attributes-section,
  .skills-section {
    padding: 24px;
  }
}
</style>