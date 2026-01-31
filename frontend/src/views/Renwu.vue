<template>
  <div class="character-page">
    <CharacterSidebar
      :characters="characters"
      :active-character-id="activeCharacterId"
      @select="handleSelectCharacter"
      @add="handleAddCharacter"
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
            <div class="info-item">
              <label>势力</label>
              <div class="info-value">{{ activeCharacter.faction || '无' }}</div>
            </div>
            <div class="info-item">
              <label>财产</label>
              <div class="info-value">{{ activeCharacter.wealth || 0 }}</div>
            </div>
            <div class="info-item">
              <label>等级</label>
              <div class="info-value">Lv.{{ activeCharacter.level }}</div>
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
              <div class="attribute-content">
                <div class="attribute-name">{{ attr.name }}</div>
                <div v-if="attr.description" class="attribute-description">{{ attr.description }}</div>
                <div class="attribute-value">{{ attr.value }}</div>
              </div>
              <button class="delete-btn" @click="handleDeleteAttribute(attr.id)" title="删除属性">
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
              <div class="skill-content">
                <div class="skill-name">{{ skill.name }}</div>
                <div v-if="skill.description" class="skill-description">{{ skill.description }}</div>
              </div>
              <button class="delete-btn" @click="handleDeleteSkill(skill.id)" title="删除技能">
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
import { ref, computed, onMounted } from 'vue'
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
      id: char.id,
      name: char.name,
      faction: char.shili,
      wealth: char.property,
      level: char.level,
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
    const result = await window.go.main.app.GetCharacterInfo(characterId)
    const character = characters.value.find(c => c.id === characterId)
    if (character) {
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
      data.wealth,
      data.level
    )
    console.log('创建人物成功，ID:', newId)
    
    // 重新加载列表
    await loadCharacters()
    activeCharacterId.value = newId
    await loadCharacterDetail(newId)
  } catch (error) {
    console.error('添加人物失败:', error)
    alert('添加人物失败: ' + error.message)
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
</script>

<style scoped>
.character-page {
  display: flex;
  height: 100%;
  background-color: #f5f7fa;
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
  grid-template-columns: repeat(4, 1fr);
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
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e8eaed;
  transition: all 0.2s ease;
}

.skill-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.skill-content {
  flex: 1;
  min-width: 0;
}

.skill-name {
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.skill-description {
  font-size: 13px;
  color: #8c92a0;
  line-height: 1.6;
}

.empty-skills {
  grid-column: 1 / -1;
  text-align: center;
  padding: 48px 20px;
  color: #8c92a0;
}

.empty-skills p {
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