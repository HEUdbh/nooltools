<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
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

// 是否显示预览
const showPreview = ref(true)

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
    currentFile.value = {
      name: title.trim() + '.md',
      title: title.trim(),
      content: ''
    }
    // 第一行作为标题
    editorContent.value = `# ${title.trim()}\n\n`
    updatePreview()
  }
}

// 保存文件
async function saveFile() {
  // 从第一行读取标题
  const lines = editorContent.value.split('\n')
  let title = '未命名文档'
  
  // 查找第一个非空行作为标题
  for (const line of lines) {
    const trimmedLine = line.trim()
    if (trimmedLine) {
      // 移除markdown标题符号（#）
      title = trimmedLine.replace(/^#+\s*/, '').trim()
      break
    }
  }
  
  // 如果标题为空，使用默认标题
  if (!title) {
    title = '未命名文档'
  }

  try {
    loading.value = true
    await SaveMarkdownFile(title, editorContent.value)
    
    // 更新当前文件信息
    currentFile.value = {
      name: title + '.md',
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
    return
  }
  
  const lines = editorContent.value.split('\n')
  let processedContent = editorContent.value
  
  // 检查第一行是否已经是标题格式
  if (lines.length > 0) {
    const firstLine = lines[0].trim()
    // 如果第一行不是标题格式，且不为空，则将其转换为标题
    if (firstLine && !firstLine.startsWith('#')) {
      // 将第一行转换为一级标题
      lines[0] = `# ${firstLine}`
      processedContent = lines.join('\n')
    }
  }
  
  previewContent.value = marked(processedContent)
}

// 监听编辑器内容变化
watch(editorContent, () => {
  updatePreview()
})

// 切换预览显示
function togglePreview() {
  showPreview.value = !showPreview.value
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
        <button class="btn btn-secondary" @click="togglePreview">
          {{ showPreview ? '👁️ 隐藏预览' : '👁️ 显示预览' }}
        </button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="editor-content">
      <!-- 文件列表侧边栏 -->
      <div class="file-sidebar">
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

      <!-- 编辑器和预览区域 -->
      <div class="editor-preview-container">
        <!-- 编辑器 -->
        <div class="editor-section" :class="{ 'full-width': !showPreview }">
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
        <div v-if="showPreview" class="preview-section">
          <div class="section-header">
            <h3>👁️ 预览</h3>
          </div>
          <div class="markdown-preview" v-html="previewContent"></div>
        </div>
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
  background-color: #f5f7fa;
}

/* 工具栏样式 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background-color: #ffffff;
  border-bottom: 1px solid #e8eaed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.current-file {
  color: #5a6c7d;
  font-size: 14px;
  padding: 6px 12px;
  background-color: #f0f7ff;
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
  background-color: #1890ff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #40a9ff;
}

.btn-secondary {
  background-color: #f0f2f5;
  color: #5a6c7d;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #e6f7ff;
  color: #1890ff;
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
  background-color: #f0f7ff;
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
  background-color: #ffffff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
}

.file-sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.file-sidebar-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
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
  background-color: #f8f9fa;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-item:hover {
  background-color: #f0f7ff;
}

.file-item.active {
  background-color: #e6f7ff;
  border: 1px solid #1890ff;
}

.file-name {
  flex: 1;
  font-size: 14px;
  color: #2c3e50;
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
  color: #999;
  padding: 40px 20px;
  font-size: 14px;
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
  background-color: #ffffff;
  border-bottom: 1px solid #e8eaed;
}

.section-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #5a6c7d;
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
  background-color: #ffffff;
  color: #2c3e50;
  outline: none;
}

.markdown-editor::placeholder {
  color: #bdc3c7;
}

/* 预览样式 */
.markdown-preview {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #ffffff;
  border-left: 1px solid #e8eaed;
}

/* Markdown预览样式 */
.markdown-preview :deep(h1) {
  font-size: 2.5em;
  font-weight: 700;
  margin: 0.5em 0;
  border-bottom: 2px solid #eaecef;
  padding-bottom: 0.3em;
  color: #1a1a1a;
}

.markdown-preview :deep(h2) {
  font-size: 1.5em;
  font-weight: 600;
  margin: 0.83em 0;
  border-bottom: 1px solid #eaecef;
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
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.markdown-preview :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
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
  color: #6a737d;
  border-left: 0.25em solid #dfe2e5;
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
  color: #0366d6;
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
  border: 1px solid #dfe2e5;
}

.markdown-preview :deep(table th) {
  font-weight: 600;
  background-color: #f6f8fa;
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
  background-color: #f5222d;
  color: white;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
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