<template>
  <div class="notes-page">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-select
        v-model="currentNotebookId"
        placeholder="选择笔记本"
        clearable
        style="width: 200px"
        @change="loadNotes"
      >
        <el-option
          v-for="nb in notebooks"
          :key="nb.id"
          :label="nb.name"
          :value="nb.id"
        />
      </el-select>
      <el-button @click="openNotebookDialog()">新建笔记本</el-button>
      <el-button type="primary" :disabled="!currentNotebookId" @click="openNoteEditor()">
        新建笔记
      </el-button>
      <el-input
        v-model="keyword"
        placeholder="搜索笔记标题"
        clearable
        style="width: 200px; margin-left: auto"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 笔记列表 -->
    <div class="list" v-loading="loading">
      <el-empty v-if="!loading && filteredNotes.length === 0" :description="emptyTip" />
      <div v-else class="note-cards">
        <div
          v-for="n in filteredNotes"
          :key="n.id"
          class="note-card"
          @click="openNoteEditor(n)"
        >
          <div class="note-card-header">
            <h4>{{ n.title || '无标题' }}</h4>
            <el-button
              class="delete-btn"
              type="danger"
              link
              size="small"
              @click.stop="removeNote(n)"
            >
              删除
            </el-button>
          </div>
          <p class="file-path">{{ n.file_path }}</p>
          <span class="time">{{ formatTime(n.updated_at) }}</span>
        </div>
      </div>
    </div>

    <!-- 笔记编辑器弹窗 -->
    <el-dialog
      v-model="editorVisible"
      :title="editingNote.id ? '编辑笔记' : '新建笔记'"
      width="80%"
      top="5vh"
      destroy-on-close
      @open="onEditorOpen"
    >
      <el-input
        v-model="editingNote.title"
        placeholder="笔记标题"
        size="large"
        style="margin-bottom: 12px"
      />
      <el-input
        v-model="editingNote.content"
        type="textarea"
        :autosize="{ minRows: 12, maxRows: 24 }"
        placeholder="输入 Markdown 内容…"
        class="editor-textarea"
      />
      <div class="editor-footer-info">
        <span v-if="editingNote.file_path">文件：{{ editingNote.file_path }}</span>
      </div>
      <template #footer>
        <el-button @click="editorVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveNote">保存</el-button>
      </template>
    </el-dialog>

    <!-- 新建/编辑笔记本弹窗 -->
    <el-dialog
      v-model="notebookDialogVisible"
      :title="editingNotebook.id ? '编辑笔记本' : '新建笔记本'"
      width="420px"
    >
      <el-form label-width="60px">
        <el-form-item label="名称">
          <el-input v-model="editingNotebook.name" placeholder="笔记本名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editingNotebook.description" placeholder="可选描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="notebookDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveNotebook">
          {{ editingNotebook.id ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { listNotebooks, createNotebook, updateNotebook, deleteNotebook } from '@/api/notebooks'
import { listNotes, getNote, createNote, updateNote, deleteNote } from '@/api/notes'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

// ---------- 笔记本 ----------
const notebooks = ref([])
const currentNotebookId = ref(null)
const notebookDialogVisible = ref(false)
const editingNotebook = reactive({ id: null, name: '', description: '' })

async function loadNotebooks() {
  try {
    const res = await listNotebooks()
    notebooks.value = res?.data || []
    // 默认选中第一个笔记本
    if (notebooks.value.length > 0 && !currentNotebookId.value) {
      currentNotebookId.value = notebooks.value[0].id
    }
  } catch {
    notebooks.value = []
  }
}

function openNotebookDialog(nb) {
  if (nb) {
    editingNotebook.id = nb.id
    editingNotebook.name = nb.name
    editingNotebook.description = nb.description || ''
  } else {
    editingNotebook.id = null
    editingNotebook.name = ''
    editingNotebook.description = ''
  }
  notebookDialogVisible.value = true
}

async function saveNotebook() {
  if (!editingNotebook.name.trim()) {
    ElMessage.warning('请输入笔记本名称')
    return
  }
  try {
    if (editingNotebook.id) {
      await updateNotebook(editingNotebook.id, {
        name: editingNotebook.name,
        description: editingNotebook.description,
      })
      ElMessage.success('笔记本已更新')
    } else {
      const res = await createNotebook({
        name: editingNotebook.name,
        description: editingNotebook.description,
      })
      const created = res?.data
      if (created?.id) currentNotebookId.value = created.id
      ElMessage.success('笔记本已创建')
    }
    notebookDialogVisible.value = false
    await loadNotebooks()
    loadNotes()
  } catch (e) {
    ElMessage.error(e?.message || '操作失败')
  }
}

// ---------- 笔记 ----------
const loading = ref(false)
const saving = ref(false)
const notes = ref([])
const keyword = ref('')
const editorVisible = ref(false)
const editingNote = reactive({ id: null, title: '', content: '', file_path: '' })

const filteredNotes = computed(() => {
  if (!keyword.value) return notes.value
  const kw = keyword.value.toLowerCase()
  return notes.value.filter((n) => (n.title || '').toLowerCase().includes(kw))
})

const emptyTip = computed(() => {
  if (!currentNotebookId.value) return '请先选择或创建一个笔记本'
  return '该笔记本暂无笔记，点击"新建笔记"开始写作'
})

async function loadNotes() {
  if (!currentNotebookId.value) {
    notes.value = []
    return
  }
  loading.value = true
  try {
    const res = await listNotes({ notebook_id: currentNotebookId.value })
    notes.value = res?.data || []
  } catch {
    notes.value = []
    ElMessage.warning('加载笔记失败，请确认 Go 后端已启动')
  } finally {
    loading.value = false
  }
}

function openNoteEditor(note) {
  if (note) {
    // 编辑已有笔记 — 先填元数据，content 在 onEditorOpen 时加载
    editingNote.id = note.id
    editingNote.title = note.title
    editingNote.content = '' // 待加载
    editingNote.file_path = note.file_path || ''
  } else {
    // 新建笔记
    editingNote.id = null
    editingNote.title = ''
    editingNote.content = ''
    editingNote.file_path = ''
  }
  editorVisible.value = true
}

async function onEditorOpen() {
  // 如果是编辑已有笔记，从后端加载完整内容
  if (editingNote.id) {
    try {
      const res = await getNote(editingNote.id)
      const data = res?.data || res
      editingNote.content = data?.content || ''
      editingNote.file_path = data?.file_path || editingNote.file_path
    } catch {
      ElMessage.warning('加载笔记内容失败')
    }
  }
}

async function saveNote() {
  if (!editingNote.title.trim()) {
    ElMessage.warning('请输入笔记标题')
    return
  }
  saving.value = true
  try {
    if (editingNote.id) {
      // 更新
      await updateNote(editingNote.id, {
        title: editingNote.title,
        content: editingNote.content,
      })
      ElMessage.success('笔记已保存')
    } else {
      // 创建
      await createNote({
        notebook_id: currentNotebookId.value,
        title: editingNote.title,
        content: editingNote.content,
      })
      ElMessage.success('笔记已创建')
    }
    editorVisible.value = false
    loadNotes()
  } catch (e) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function removeNote(note) {
  try {
    await ElMessageBox.confirm(`确定删除「${note.title}」吗？对应的 .md 文件也会被删除。`, '确认删除', {
      type: 'warning',
    })
    await deleteNote(note.id)
    ElMessage.success('已删除')
    loadNotes()
  } catch {
    // 用户取消
  }
}

// ---------- 工具函数 ----------
function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString()
}

// ---------- 初始化 ----------
onMounted(async () => {
  await loadNotebooks()
  loadNotes()
})
</script>

<style lang="less" scoped>
.notes-page {
  padding: 20px;
}
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
}
.note-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}
.note-card {
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  &:hover {
    border-color: #409eff;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
    .delete-btn {
      opacity: 1;
    }
  }
}
.note-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  h4 {
    margin: 0;
    font-size: 16px;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
.delete-btn {
  opacity: 0;
  transition: opacity 0.2s;
}
.file-path {
  margin: 6px 0;
  font-size: 12px;
  color: #a8abb2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.time {
  font-size: 12px;
  color: #c0c4cc;
}
.editor-textarea :deep(textarea) {
  font-family: 'Consolas', 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  line-height: 1.6;
}
.editor-footer-info {
  margin-top: 8px;
  font-size: 12px;
  color: #a8abb2;
}
</style>
