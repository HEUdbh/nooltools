<script setup>
import { ref, onMounted, watch } from 'vue'
import { GetMarkdownFiles, ReadMarkdownFile, SaveMarkdownFile, DeleteMarkdownFile, RenameMarkdownFile } from '../../wailsjs/go/main/App.js'
import { marked } from 'marked'

// 当前编辑的文件
const currentFile = ref({
  name: '',
  title: '',
  content: ''
})

// 文件列表
const fileList = ref([])

// 编辑器内容
const editorContent = ref('')

// 预览内容
const previewContent = ref('')

// 编辑器视图模式（editor/preview）
const viewMode = ref('editor')

// 是否显示文件列表侧边栏
const showFileSidebar = ref(true)

// 是否显示目录侧边栏
const showTocSidebar = ref(true)

// 标题目录列表
const tocList = ref([])

// 是否正在加载
const loading = ref(false)

// 错误信息
const errorMessage = ref('')

// 加载文件列表
async function loadFileList() {
  try {
    loading.value = true
    const files = await GetMarkdownFiles()
    fileList.value = files || []
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = '加载文件列表失败: ' + error
    console.error('加载文件列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载文件内容
async function loadFile(filename) {
  try {
    loading.value = true
    const content = await ReadMarkdownFile(filename)
    currentFile.value = {
      name: filename,
      title: filename.replace('.md', ''),
      content: content
    }
    editorContent.value = content
    updatePreview()
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = '加载文件失败: ' + error
    console.error('加载文件失败:', error)
  } finally {
    loading.value = false
  }
}

// 创建新文件
function createNewFile() {
  const title = prompt('请输入新文件标题:', '未命名文档')
  if (title && title.trim()) {
    const trimmedTitle = title.trim()
    currentFile.value = {
      name: trimmedTitle + '.md',
      title: trimmedTitle,
      content: ''
    }
    editorContent.value = ''
    updatePreview()
  }
}

// 保存文件
async function saveFile() {
  let fileName = (currentFile.value.name || '').trim()
  if (!fileName) {
    const fallbackTitle = (currentFile.value.title || '未命名文档').trim() || '未命名文档'
    fileName = `${fallbackTitle}.md`
  }
  const normalizedFileName = /\.md$/i.test(fileName) ? fileName : `${fileName}.md`
  const title = normalizedFileName.replace(/\.md$/i, '')

  try {
    loading.value = true
    await SaveMarkdownFile(title, editorContent.value)
    
    // 更新当前文件信息
    currentFile.value = {
      name: normalizedFileName,
      title: title,
      content: editorContent.value
    }
    
    errorMessage.value = ''
    await loadFileList()
    alert('保存成功!')
  } catch (error) {
    errorMessage.value = '保存失败: ' + error
    console.error('保存失败:', error)
    alert('保存失败: ' + error)
  } finally {
    loading.value = false
  }
}

// 删除文件
async function deleteFile(filename) {
  if (confirm(`确定要删除文件 "${filename}" 吗?`)) {
    try {
      loading.value = true
      await DeleteMarkdownFile(filename)
      errorMessage.value = ''
      
      // 如果删除的是当前文件，清空编辑器
      if (currentFile.value.name === filename) {
        currentFile.value = { name: '', title: '', content: '' }
        editorContent.value = ''
        previewContent.value = ''
      }
      
      await loadFileList()
      alert('删除成功!')
    } catch (error) {
      errorMessage.value = '删除失败: ' + error
      console.error('删除失败:', error)
      alert('删除失败: ' + error)
    } finally {
      loading.value = false
    }
  }
}

// 重命名文件
async function renameFile(oldName) {
  const newName = prompt('请输入新文件名:', oldName.replace('.md', ''))
  if (newName && newName.trim() && newName.trim() !== oldName.replace('.md', '')) {
    try {
      loading.value = true
      await RenameMarkdownFile(oldName, newName.trim())
      errorMessage.value = ''
      
      // 如果重命名的是当前文件，更新当前文件信息
      if (currentFile.value.name === oldName) {
        currentFile.value.name = newName.trim() + '.md'
        currentFile.value.title = newName.trim()
      }
      
      await loadFileList()
      alert('重命名成功!')
    } catch (error) {
      errorMessage.value = '重命名失败: ' + error
      console.error('重命名失败:', error)
      alert('重命名失败: ' + error)
    } finally {
      loading.value = false
    }
  }
}

// 更新预览
function updatePreview() {
  if (!editorContent.value) {
    previewContent.value = ''
    tocList.value = []
    return
  }
  
  // 配置marked的自定义渲染器
  const renderer = new marked.Renderer()
  let headingCounter = 0
  
  // 自定义标题渲染，添加data-heading-id属性
  renderer.heading = function(text, level, raw) {
    const headingId = `heading-${headingCounter++}`
    return `<h${level} data-heading-id="${headingId}">${text}</h${level}>`
  }
  
  // 使用自定义渲染器
  marked.setOptions({
    renderer: renderer
  })
  
  previewContent.value = marked(editorContent.value)
  
  // 更新目录
  updateToc()
}

// 解析markdown标题，生成目录
function updateToc() {
  if (!editorContent.value) {
    tocList.value = []
    return
  }
  
  const lines = editorContent.value.split('\n')
  const toc = []
  let idCounter = 0
  
  for (const line of lines) {
    // 匹配markdown标题（# 到 ######）
    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (match) {
      const level = match[1].length // 标题级别 1-6
      const text = match[2].trim() // 标题文本
      const id = `heading-${idCounter++}`
      
      toc.push({
        id,
        level,
        text,
        line: lines.indexOf(line)
      })
    }
  }
  
  tocList.value = toc
}

// 监听编辑器内容变化
watch(editorContent, () => {
  updatePreview()
})

// 切换编辑器/预览模式
function toggleViewMode() {
  viewMode.value = viewMode.value === 'editor' ? 'preview' : 'editor'
}

// 切换文件列表侧边栏显示
function toggleFileSidebar() {
  showFileSidebar.value = !showFileSidebar.value
}

// 切换目录侧边栏显示
function toggleTocSidebar() {
  showTocSidebar.value = !showTocSidebar.value
}

// 滚动到指定标题
function scrollToHeading(headingId) {
  // 找到对应的标题元素
  const headingElement = document.querySelector(`[data-heading-id="${headingId}"]`)
  if (headingElement) {
    headingElement.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

// 组件挂载时加载文件列表
onMounted(() => {
  loadFileList()
})
</script>

<template>
  <div class="markdown-editor-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <button class="btn btn-secondary" @click="toggleFileSidebar" title="切换文件列表">
          {{ showFileSidebar ? '📁 隐藏列表' : '📁 显示列表' }}
        </button>
        <button class="btn btn-primary" @click="saveFile" :disabled="loading">
          💾 保存
        </button>
        <button class="btn btn-secondary" @click="createNewFile" :disabled="loading">
          ➕ 新建
        </button>
        <span v-if="currentFile.title" class="current-file">
          当前文件: {{ currentFile.title }}
        </span>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-secondary" @click="toggleTocSidebar" title="切换目录">
          {{ showTocSidebar ? '📑 隐藏目录' : '📑 显示目录' }}
        </button>
        <button class="btn btn-secondary" @click="toggleViewMode">
          {{ viewMode === 'editor' ? '👁️ 切换到预览' : '📝 切换到编辑' }}
        </button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="editor-content">
      <!-- 文件列表侧边栏 -->
      <transition name="sidebar-slide">
        <div v-if="showFileSidebar" class="file-sidebar">
          <div class="file-sidebar-header">
            <h3>📄 文件列表</h3>
          </div>
          <div class="file-list">
            <div
              v-for="file in fileList"
              :key="file.name"
              :class="['file-item', { active: currentFile.name === file.name }]"
              @click="loadFile(file.name)"
            >
              <span class="file-name">{{ file.title }}</span>
              <div class="file-actions">
                <button class="btn-icon" @click.stop="renameFile(file.name)" title="重命名">✏️</button>
                <button class="btn-icon" @click.stop="deleteFile(file.name)" title="删除">🗑️</button>
              </div>
            </div>
            <div v-if="fileList.length === 0" class="empty-state">
              暂无文件，点击"新建"创建
            </div>
          </div>
        </div>
      </transition>

      <!-- 编辑器和预览区域 -->
      <div class="editor-preview-container">
        <!-- 编辑器 -->
        <div v-if="viewMode === 'editor'" class="editor-section full-width">
          <div class="section-header">
            <h3>📝 编辑器</h3>
          </div>
          <textarea
            v-model="editorContent"
            class="markdown-editor"
            placeholder="开始编写你的Markdown文档..."
            spellcheck="false"
          ></textarea>
        </div>

        <!-- 预览 -->
        <div v-else class="preview-section full-width">
          <div class="section-header">
            <h3>👁️ 预览</h3>
          </div>
          <div class="markdown-preview" v-html="previewContent"></div>
        </div>

        <!-- 目录侧边栏 -->
        <transition name="sidebar-slide">
          <div v-if="showTocSidebar" class="toc-sidebar">
            <div class="toc-sidebar-header">
              <h3>📑 目录</h3>
            </div>
            <div class="toc-list">
              <div
                v-for="item in tocList"
                :key="item.id"
                :class="['toc-item', `toc-level-${item.level}`]"
                @click="scrollToHeading(item.id)"
              >
                {{ item.text }}
              </div>
              <div v-if="tocList.length === 0" class="empty-state">
                暂无标题
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<style scoped>
.markdown-editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--app-bg);
}

/* 工具栏样式 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background-color: var(--app-surface);
  border-bottom: 1px solid var(--app-border);
  box-shadow: var(--app-shadow-sm);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-file {
  color: var(--app-text-secondary);
  font-size: 14px;
  padding: 6px 12px;
  background-color: var(--app-hover-bg);
  border-radius: 4px;
}

/* 按钮样式 */
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background-color: var(--app-accent);
  color: var(--app-text-inverse);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--app-accent-hover);
}

.btn-secondary {
  background-color: var(--app-muted-bg);
  color: var(--app-text-secondary);
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--app-hover-bg);
  color: var(--app-accent);
}

.btn-icon {
  padding: 4px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.btn-icon:hover {
  background-color: var(--app-hover-bg);
}

/* 主内容区域 */
.editor-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 文件列表侧边栏 */
.file-sidebar {
  width: 280px;
  background-color: var(--app-surface);
  border-right: 1px solid var(--app-border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

/* 侧边栏过渡动画 */
.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: all 0.3s ease;
}

.sidebar-slide-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}

.sidebar-slide-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.sidebar-slide-enter-to,
.sidebar-slide-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.file-sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.file-sidebar-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--app-text-primary);
}

.file-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  margin-bottom: 8px;
  background-color: var(--app-item-bg);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-item:hover {
  background-color: var(--app-hover-bg);
}

.file-item.active {
  background-color: var(--app-active-bg);
  border: 1px solid var(--app-accent);
}

.file-name {
  flex: 1;
  font-size: 14px;
  color: var(--app-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.file-item:hover .file-actions {
  opacity: 1;
}

.empty-state {
  text-align: center;
  color: var(--app-text-muted);
  padding: 40px 20px;
  font-size: 14px;
}

/* 目录侧边栏样式 */
.toc-sidebar {
  width: 280px;
  display: flex;
  flex-direction: column;
  background-color: var(--app-surface);
  border-left: 1px solid var(--app-border);
  flex-shrink: 0;
}

.toc-sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--app-divider-soft);
}

.toc-sidebar-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--app-text-primary);
}

.toc-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.toc-item {
  padding: 8px 12px;
  margin-bottom: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  color: var(--app-text-primary);
  border-radius: 4px;
  line-height: 1.5;
}

.toc-item:hover {
  background-color: var(--app-hover-bg);
  color: var(--app-accent);
}

/* 标题层级缩进 */
.toc-level-1 {
  padding-left: 12px;
  font-weight: 600;
}

.toc-level-2 {
  padding-left: 24px;
  font-weight: 500;
}

.toc-level-3 {
  padding-left: 36px;
  font-weight: 400;
}

.toc-level-4 {
  padding-left: 48px;
  font-weight: 400;
}

.toc-level-5 {
  padding-left: 60px;
  font-weight: 400;
}

.toc-level-6 {
  padding-left: 72px;
  font-weight: 400;
}

/* 编辑器和预览容器 */
.editor-preview-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-section,
.preview-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-section.full-width {
  flex: 1;
}

.section-header {
  padding: 12px 20px;
  background-color: var(--app-surface);
  border-bottom: 1px solid var(--app-border);
}

.section-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--app-text-secondary);
}

/* 编辑器样式 */
.markdown-editor {
  flex: 1;
  padding: 20px;
  border: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  background-color: var(--app-surface);
  color: var(--app-text-primary);
  outline: none;
}

.markdown-editor::placeholder {
  color: var(--app-text-muted);
}

/* 预览样式 */
.markdown-preview {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: var(--app-surface);
  color: var(--app-text-primary);
}

/* Markdown预览样式 */
.markdown-preview :deep(h1) {
  font-size: 2.5em;
  font-weight: 700;
  margin: 0.5em 0;
  border-bottom: 2px solid var(--md-divider);
  padding-bottom: 0.3em;
  color: var(--md-heading-color);
}

.markdown-preview :deep(h2) {
  font-size: 1.5em;
  font-weight: 600;
  margin: 0.83em 0;
  border-bottom: 1px solid var(--md-divider);
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h3) {
  font-size: 1.25em;
  font-weight: 600;
  margin: 1em 0;
}

.markdown-preview :deep(p) {
  margin: 1em 0;
  line-height: 1.6;
}

.markdown-preview :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: var(--md-inline-code-bg);
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.markdown-preview :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: var(--md-code-block-bg);
  border-radius: 6px;
}

.markdown-preview :deep(pre code) {
  padding: 0;
  margin: 0;
  font-size: 100%;
  background-color: transparent;
}

.markdown-preview :deep(blockquote) {
  padding: 0 1em;
  color: var(--md-blockquote-text);
  border-left: 0.25em solid var(--md-blockquote-border);
  margin: 1em 0;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  padding-left: 2em;
  margin: 1em 0;
}

.markdown-preview :deep(li) {
  margin: 0.5em 0;
}

.markdown-preview :deep(a) {
  color: var(--md-link-color);
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.markdown-preview :deep(table) {
  border-spacing: 0;
  border-collapse: collapse;
  margin: 1em 0;
  width: 100%;
}

.markdown-preview :deep(table th),
.markdown-preview :deep(table td) {
  padding: 6px 13px;
  border: 1px solid var(--md-table-border);
}

.markdown-preview :deep(table th) {
  font-weight: 600;
  background-color: var(--md-code-block-bg);
}

.markdown-preview :deep(img) {
  max-width: 100%;
  height: auto;
}

/* 错误提示 */
.error-message {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 12px 20px;
  background-color: var(--app-danger);
  color: var(--app-text-inverse);
  border-radius: 6px;
  box-shadow: var(--app-shadow-lg);
  z-index: 1000;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    transform: translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
