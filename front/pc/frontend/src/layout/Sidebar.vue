<template>
  <div class="sidebar">
    <div class="logo">Life Journey</div>
    <el-menu
      :default-active="activeMenu"
      class="menu"
      background-color="#f5f7fa"
      text-color="#303133"
      active-text-color="#409eff"
      router
    >
      <el-menu-item index="/notes">
        <el-icon><Document /></el-icon>
        <span>笔记</span>
      </el-menu-item>
      <el-menu-item index="/todos">
        <el-icon><List /></el-icon>
        <span>待办</span>
      </el-menu-item>
      <el-menu-item index="/tools">
        <el-icon><Tools /></el-icon>
        <span>工具</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Document, List, Tools } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style lang="less" scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.logo {
  height: 48px;
  line-height: 48px;
  padding-left: 20px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  border-bottom: 1px solid #e4e7ed;
}
.menu {
  flex: 1;
  border-right: none;
}
</style>
