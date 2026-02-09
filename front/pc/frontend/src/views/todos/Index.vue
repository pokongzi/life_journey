<template>
  <div class="todos-page">
    <div class="toolbar">
      <el-input
        v-model="newTitle"
        placeholder="输入待办，回车添加"
        style="width: 320px"
        @keyup.enter="addTodo"
      >
        <template #append>
          <el-button type="primary" @click="addTodo">添加</el-button>
        </template>
      </el-input>
      <el-radio-group v-model="filter" @change="loadTodos">
        <el-radio-button label="all">全部</el-radio-button>
        <el-radio-button label="pending">未完成</el-radio-button>
        <el-radio-button label="done">已完成</el-radio-button>
      </el-radio-group>
    </div>
    <div class="list" v-loading="loading">
      <el-empty v-if="!loading && todos.length === 0" description="暂无待办" />
      <div v-else>
        <div v-for="t in todos" :key="t.id" class="todo-item">
          <el-checkbox :model-value="t.status === 'done'" @change="(v) => toggleTodo(t.id, v)">
            {{ t.title }}
          </el-checkbox>
          <el-button type="danger" link size="small" @click.stop="removeTodo(t.id)">删除</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listTodos, createTodo, updateTodo, deleteTodo } from '@/api/todos'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const todos = ref([])
const newTitle = ref('')
const filter = ref('all')

async function loadTodos() {
  loading.value = true
  try {
    const res = await listTodos({ status: filter.value === 'all' ? undefined : filter.value })
    todos.value = res?.list || res?.data || []
  } catch {
    todos.value = []
    ElMessage.warning('Go 后端待办接口未就绪，请先实现 /api/todos')
  } finally {
    loading.value = false
  }
}

async function addTodo() {
  if (!newTitle.value.trim()) return
  try {
    await createTodo({ title: newTitle.value.trim() })
    newTitle.value = ''
    loadTodos()
  } catch (e) {
    ElMessage.error(e?.message || '添加失败')
  }
}

async function toggleTodo(id, checked) {
  const status = checked ? 'done' : 'pending'
  try {
    await updateTodo(id, { status })
    loadTodos()
  } catch (e) {
    ElMessage.error(e?.message || '更新失败')
  }
}

async function removeTodo(id) {
  try {
    await deleteTodo(id)
    loadTodos()
  } catch (e) {
    ElMessage.error(e?.message || '删除失败')
  }
}

onMounted(loadTodos)
</script>

<style lang="less" scoped>
.todos-page {
  padding: 20px;
}
.toolbar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
}
.todo-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
</style>
