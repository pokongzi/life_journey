<template>
  <div class="tools-page">
    <h3>实用工具</h3>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="tool-card" shadow="hover" @click="compressVisible = true">
          <template #header>
            <span><el-icon><Picture /></el-icon> 图片压缩</span>
          </template>
          <p>上传图片并压缩为 JPEG，可选质量等级</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="tool-card disabled" shadow="hover">
          <template #header>
            <span><el-icon><Document /></el-icon> PDF 处理</span>
          </template>
          <p>合并、拆分、转换 PDF（v1.1）</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="tool-card disabled" shadow="hover">
          <template #header>
            <span><el-icon><Folder /></el-icon> 文件压缩</span>
          </template>
          <p>压缩、解压 ZIP/7Z（v1.1）</p>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图片压缩弹窗 -->
    <el-dialog v-model="compressVisible" title="图片压缩" width="600px" destroy-on-close>
      <div class="compress-body">
        <el-form label-width="80px" style="margin-bottom: 16px">
          <el-form-item label="压缩质量">
            <el-slider v-model="quality" :min="10" :max="100" :step="5" show-input />
          </el-form-item>
        </el-form>

        <el-upload
          ref="uploadRef"
          drag
          :auto-upload="false"
          :limit="1"
          :on-change="onFileChange"
          :on-remove="onFileRemove"
          accept="image/*"
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">拖拽图片到此处，或 <em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip">支持 JPEG / PNG / GIF，压缩后输出为 JPEG</div>
          </template>
        </el-upload>

        <!-- 压缩结果 -->
        <div v-if="compressResult" class="result-info">
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="原始大小">{{ formatSize(compressResult.originalSize) }}</el-descriptions-item>
            <el-descriptions-item label="压缩后大小">{{ formatSize(compressResult.compressedSize) }}</el-descriptions-item>
            <el-descriptions-item label="压缩率">{{ compressResult.ratio }}</el-descriptions-item>
            <el-descriptions-item label="文件名">{{ compressResult.filename }}</el-descriptions-item>
          </el-descriptions>
          <el-button type="primary" style="margin-top: 12px" @click="downloadResult">
            下载压缩后图片
          </el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="compressVisible = false">关闭</el-button>
        <el-button type="primary" :loading="compressing" :disabled="!selectedFile" @click="doCompress">
          开始压缩
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Picture, Document, Folder, UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const compressVisible = ref(false)
const quality = ref(75)
const selectedFile = ref(null)
const compressing = ref(false)
const compressResult = ref(null)
const compressedBlob = ref(null)
const uploadRef = ref()

function onFileChange(uploadFile) {
  selectedFile.value = uploadFile.raw
  compressResult.value = null
  compressedBlob.value = null
}

function onFileRemove() {
  selectedFile.value = null
  compressResult.value = null
  compressedBlob.value = null
}

async function doCompress() {
  if (!selectedFile.value) return
  compressing.value = true
  compressResult.value = null
  compressedBlob.value = null

  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('quality', String(quality.value))

    const token = localStorage.getItem('lj_token')
    const res = await axios.post('/api/tools/image/compress', formData, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      responseType: 'blob',
    })

    // 从响应头读取压缩信息
    const originalSize = parseInt(res.headers['x-original-size'] || '0', 10)
    const compressedSize = parseInt(res.headers['x-compressed-size'] || '0', 10)
    const ratio = res.headers['x-compression-ratio'] || '-'

    compressedBlob.value = res.data
    compressResult.value = {
      originalSize,
      compressedSize: compressedSize || res.data.size,
      ratio,
      filename: `compressed_${selectedFile.value.name}`,
    }

    ElMessage.success('压缩完成')
  } catch (e) {
    ElMessage.error('压缩失败：' + (e?.message || '请确认 Go 后端已启动'))
  } finally {
    compressing.value = false
  }
}

function downloadResult() {
  if (!compressedBlob.value) return
  const url = URL.createObjectURL(compressedBlob.value)
  const a = document.createElement('a')
  a.href = url
  a.download = compressResult.value?.filename || 'compressed.jpg'
  a.click()
  URL.revokeObjectURL(url)
}

function formatSize(bytes) {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}
</script>

<style lang="less" scoped>
.tools-page {
  padding: 20px;
  h3 {
    margin-bottom: 20px;
  }
}
.tool-card {
  cursor: pointer;
  p {
    margin: 0;
    font-size: 14px;
    color: #909399;
  }
  &.disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}
.compress-body {
  padding: 0 4px;
}
.result-info {
  margin-top: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}
</style>
