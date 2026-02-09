<template>
  <div class="tools-page">
    <h3>实用工具</h3>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="tool-card" shadow="hover" @click="goCompress">
          <template #header>
            <span><el-icon><Picture /></el-icon> 图片压缩</span>
          </template>
          <p>支持批量压缩图片，可选质量等级</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="tool-card" shadow="hover" disabled>
          <template #header>
            <span><el-icon><Document /></el-icon> PDF 处理</span>
          </template>
          <p>合并、拆分、转换 PDF（v1.1）</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="tool-card" shadow="hover" disabled>
          <template #header>
            <span><el-icon><Folder /></el-icon> 文件压缩</span>
          </template>
          <p>压缩、解压 ZIP/7Z（v1.1）</p>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="compressVisible" title="图片压缩" width="600px" destroy-on-close>
      <div class="compress-placeholder">
        <el-upload
          drag
          :auto-upload="false"
          :limit="10"
          accept="image/*"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">拖拽图片到此处，或点击上传</div>
          <template #tip>
            <div class="el-upload__tip">Go 后端需实现 POST /api/tools/image/compress</div>
          </template>
        </el-upload>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Picture, Document, Folder, UploadFilled } from '@element-plus/icons-vue'

const compressVisible = ref(false)

function goCompress() {
  compressVisible.value = true
}
</script>

<style lang="less" scoped>
.tools-page {
  padding: 20px;
}
.tools-page h3 {
  margin-bottom: 20px;
}
.tool-card {
  cursor: pointer;
  p {
    margin: 0;
    font-size: 14px;
    color: #909399;
  }
}
.tool-card[disabled] {
  opacity: 0.6;
  cursor: not-allowed;
}
.compress-placeholder {
  padding: 20px;
}
</style>
