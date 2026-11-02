<template>
  <div class="login-page">
    <div class="login-card">
      <h1>WTV管理系统</h1>
      <p class="login-tip">请使用管理员账号登录</p>
      <el-form :model="form" @submit.prevent="submit">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" @keyup.enter="submit" />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            show-password
            size="large"
            @keyup.enter="submit"
          />
        </el-form-item>
        <el-button type="primary" size="large" class="login-button" :loading="loading" @click="submit">
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { setSession } from '../auth'

const emit = defineEmits(['success'])

const form = reactive({
  username: '',
  password: ''
})
const loading = ref(false)

const submit = async () => {
  const username = form.username.trim()
  if (!username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const response = await axios.post('/admin/login', {
      username,
      password: form.password
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '登录失败')
      return
    }
    setSession(response.data.data.token, response.data.data.admin)
    emit('success', response.data.data.admin)
  } catch (error) {
    ElMessage.error('登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #070B12;
}

.login-card {
  width: 360px;
  padding: 32px 28px 28px;
  background: #0A101E;
  border: 1px solid #1E2A3F;
  border-radius: 16px;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.35);
}

.login-card h1 {
  margin: 0;
  font-size: 22px;
  line-height: 1.3;
  color: #EAF0F8;
  letter-spacing: 0;
}

.login-tip {
  margin: 8px 0 24px;
  color: #8E9AAF;
  font-size: 13px;
}

.login-button {
  width: 100%;
}
</style>
