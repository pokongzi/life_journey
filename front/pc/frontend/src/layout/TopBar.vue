<template>
  <div class="top-bar">
    <div class="title">{{ currentTitle }}</div>
    <div class="actions">
      <el-dropdown trigger="click" @command="handleCommand">
        <span class="user-trigger">
          <el-icon><User /></el-icon>
          <span>{{ authStore.user?.nickname || authStore.user?.email || '用户' }}</span>
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, ArrowDown, SwitchButton } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const currentTitle = computed(() => route.meta?.title || 'Life Journey')

async function handleCommand(cmd) {
  if (cmd === 'logout') {
    await authStore.logout()
    router.push('/login')
  }
}
</script>

<style lang="less" scoped>
.top-bar {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.title {
  font-size: 16px;
  font-weight: 500;
}
.user-trigger {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  &:hover {
    background: #f5f7fa;
  }
}
</style>
