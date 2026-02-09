<template>
  <div class="login-page">
    <div class="login-box">
      <h1 class="title">Life Journey</h1>
      <p class="subtitle">高效工作助手 · 笔记 · 待办 · 工具</p>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        class="form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            size="large"
            prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item v-if="loginType === 'password'" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item v-else prop="code">
          <el-input
            v-model="form.code"
            placeholder="验证码"
            size="large"
            prefix-icon="Key"
          >
            <template #append>
              <el-button :disabled="countdown > 0" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s 后重发` : '发送验证码' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleLogin">
            登录
          </el-button>
        </el-form-item>
        <div class="links">
          <el-link type="primary" @click="toggleLoginType">
            {{ loginType === 'password' ? '验证码登录' : '密码登录' }}
          </el-link>
          <el-link type="primary" @click="goRegister">注册</el-link>
        </div>
      </el-form>
      <p class="tip">前端仅调用本地 Go 接口，Go 转发到远程服务端</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formRef = ref()
const loading = ref(false)
const loginType = ref('password') // password | code
const countdown = ref(0)

const form = reactive({
  email: '',
  password: '',
  code: '',
})

const rules = computed(() =>
  loginType.value === 'password'
    ? { email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }], password: [{ required: true, message: '请输入密码', trigger: 'blur' }] }
    : { email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }], code: [{ required: true, message: '请输入验证码', trigger: 'blur' }] }
)

function toggleLoginType() {
  loginType.value = loginType.value === 'password' ? 'code' : 'password'
}

function sendCode() {
  if (!form.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  ElMessage.info('Go 后端尚未实现验证码接口，请使用密码登录')
  countdown.value = 60
  const t = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) clearInterval(t)
  }, 1000)
}

async function handleLogin() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    if (loginType.value === 'password') {
      await authStore.loginByPassword(form.email, form.password)
    } else {
      await authStore.loginByCode(form.email, form.code)
    }
    ElMessage.success('登录成功')
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (e) {
    ElMessage.error(e?.message || '登录失败，请确保 Go 后端已启动（端口 13245）')
  } finally {
    loading.value = false
  }
}

function goRegister() {
  ElMessage.info('注册页待实现，Go 后端需提供 /api/auth/register')
}
</script>

<style lang="less" scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
}
.login-box {
  width: 380px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}
.title {
  font-size: 28px;
  font-weight: 600;
  text-align: center;
  margin-bottom: 8px;
}
.subtitle {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-bottom: 32px;
}
.form :deep(.el-form-item) {
  margin-bottom: 20px;
}
.form :deep(.el-button) {
  width: 100%;
}
.links {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 16px;
}
.tip {
  margin-top: 24px;
  font-size: 12px;
  color: #c0c4cc;
  text-align: center;
}
</style>
