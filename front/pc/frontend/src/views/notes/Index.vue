<template>
  <div class="notes-page">
    <div class="toolbar">
      <el-button type="primary" @click="createNote">新建笔记</el-button>
      <el-input
        v-model="keyword"
        placeholder="搜索笔记"
        clearable
        style="width: 200px"
        @keyup.enter="loadNotes"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>
    <div class="list" v-loading="loading">
      <el-empty v-if="!loading && notes.length === 0" description="暂无笔记，请新建" />
      <div v-else class="note-cards">
        <div
          v-for="n in notes"
          :key="n.id"
          class="note-card"
          @click="openNote(n)"
        >
          <h4>{{ n.title || '无标题' }}</h4>
          <p class="summary">{{ n.summary || '暂无摘要' }}</p>
          <span class="time">{{ formatTime(n.updatedAt) }}</span>
        </div>
      </div>
    </div>
    <el-dialog v-model="dialogVisible" title="笔记编辑" width="80%" destroy-on-close>
      <div class="editor-placeholder">
        Markdown 编辑器占位，待集成 Milkdown/Vditor
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveNote">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listNotes } from '@/api/notes'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const notes = ref([])
const keyword = ref('')
const dialogVisible = ref(false)

async function loadNotes() {
  loading.value = true
  try {
    const res = await listNotes({ keyword: keyword.value })
    notes.value = res?.list || res?.data || []
  } catch {
    notes.value = []
    ElMessage.warning('Go 后端笔记接口未就绪，请先实现 /api/notes')
  } finally {
    loading.value = false
  }
}

function formatTime(t) {
  if (!t) return ''
  const d = new Date(t)
  return d.toLocaleString()
}

function createNote() {
  dialogVisible.value = true
}

function openNote(n) {
  dialogVisible.value = true
}

function saveNote() {
  ElMessage.info('保存接口待实现')
}

onMounted(loadNotes)
</script>

<style lang="less" scoped>
.notes-page {
  padding: 20px;
}
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}
.note-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
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
  }
}
.note-card h4 {
  margin: 0 0 8px;
  font-size: 16px;
}
.note-card .summary {
  margin: 0 0 8px;
  font-size: 13px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.note-card .time {
  font-size: 12px;
  color: #c0c4cc;
}
.editor-placeholder {
  min-height: 200px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 4px;
  color: #909399;
}
</style>
