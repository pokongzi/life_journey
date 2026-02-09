<template>
  <el-container class="main-layout">
    <el-aside width="220px" class="aside">
      <Sidebar />
    </el-aside>
    <el-container direction="vertical">
      <el-header height="48" class="header">
        <TopBar />
      </el-header>
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import Sidebar from './Sidebar.vue'
import TopBar from './TopBar.vue'
</script>

<style lang="less" scoped>
.main-layout {
  height: 100vh;
}
.aside {
  background: #f5f7fa;
  border-right: 1px solid #e4e7ed;
}
.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 16px;
  display: flex;
  align-items: center;
}
.main {
  background: #fff;
  overflow: auto;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
