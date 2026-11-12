<template>
  <div v-if="checking" class="boot-screen" />
  <Login v-else-if="!authed" @success="onLogin" />
  <div v-else-if="authed" class="app-container">
    <el-container class="app-layout" direction="vertical">
      <el-header height="40px" class="header">
        <h1>WTV管理系统</h1>
        <div class="header-user">
          <div v-if="isSuper" class="header-scheduler">
            <span>定时任务</span>
            <el-switch
              v-model="schedulerEnabled"
              :loading="schedulerSaving"
              inline-prompt
              active-text="开"
              inactive-text="关"
              @change="toggleScheduler"
            />
          </div>
          <span class="header-name">{{ profile.nickname || profile.username }}</span>
          <el-tag size="small" effect="plain">{{ profile.role_name || '管理员' }}</el-tag>
          <el-button link @click="passwordVisible = true">修改密码</el-button>
          <el-button link @click="logout">退出</el-button>
        </div>
      </el-header>
      <el-container class="app-body">
        <el-aside width="130px" class="aside">
          <el-menu
            :default-active="activeMenu"
            :key="activeMenu"
            class="el-menu-vertical-demo"
            @select="handleMenuSelect"
          >
            <el-menu-item v-for="item in visibleMenus" :key="item.index" :index="item.index">
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.label }}</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="app-main">
          <VideoManagement v-if="activeMenu === 'video'" :can-write="can('video:write')" />
          <UserManagement v-else-if="activeMenu === 'user'" />
          <FeedbackManagement
            v-else-if="activeMenu === 'feedback'"
            :can-write="can('feedback:write')"
          />
          <AdminManagement v-else-if="activeMenu === 'admin'" />
          <div v-else-if="activeMenu === 'member'" class="placeholder">
            <el-empty description="会员管理开发中" />
          </div>
        </el-main>
      </el-container>
    </el-container>

    <el-dialog v-model="passwordVisible" title="修改密码" width="400px">
      <el-form :model="passwordForm" label-width="72px">
        <el-form-item label="原密码">
          <el-input v-model="passwordForm.old_password" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.new_password" type="password" show-password placeholder="6-64位" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordVisible = false">取消</el-button>
        <el-button type="primary" :loading="passwordSaving" @click="changePassword">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import VideoManagement from './views/VideoManagement.vue'
import UserManagement from './views/UserManagement.vue'
import FeedbackManagement from './views/FeedbackManagement.vue'
import AdminManagement from './views/AdminManagement.vue'
import Login from './views/Login.vue'
import { VideoCamera, User, ChatDotRound, Star, Setting } from '@element-plus/icons-vue'
import { clearSession, getProfile, getToken, hasPerm, setSession } from './auth'

const menus = [
  { index: 'video', label: '视频管理', perm: 'video:read', icon: VideoCamera },
  { index: 'user', label: '用户管理', perm: 'user:read', icon: User },
  { index: 'feedback', label: '用户反馈', perm: 'feedback:read', icon: ChatDotRound },
  { index: 'member', label: '会员管理', perm: 'member:read', icon: Star },
  { index: 'admin', label: '管理员', perm: 'admin:manage', icon: Setting }
]

const checking = ref(true)
const authed = ref(false)
const profile = ref(getProfile() || {})
const activeMenu = ref('video')
const passwordVisible = ref(false)
const passwordSaving = ref(false)
const passwordForm = reactive({
  old_password: '',
  new_password: ''
})
const schedulerEnabled = ref(true)
const schedulerSaving = ref(false)
const schedulerReady = ref(false)

const visibleMenus = computed(() => menus.filter((item) => hasPerm(profile.value, item.perm)))
const isSuper = computed(() => profile.value?.role === 'super')

const can = (perm) => hasPerm(profile.value, perm)

const loadScheduler = async () => {
  if (!isSuper.value) return
  schedulerReady.value = false
  try {
    const response = await axios.get('/admin/scheduler')
    if (response.data.code === 0) {
      schedulerEnabled.value = !!response.data.data.enabled
    }
  } catch {
    schedulerEnabled.value = true
  } finally {
    await nextTick()
    schedulerReady.value = true
  }
}

const toggleScheduler = async (value) => {
  if (!schedulerReady.value || schedulerSaving.value) return
  schedulerSaving.value = true
  try {
    const response = await axios.put('/admin/scheduler', { enabled: value })
    if (response.data.code !== 0) {
      schedulerReady.value = false
      schedulerEnabled.value = !value
      await nextTick()
      schedulerReady.value = true
      ElMessage.error(response.data.message || '操作失败')
      return
    }
    ElMessage.success(value ? '定时任务已打开' : '定时任务已关闭')
  } catch {
    schedulerReady.value = false
    schedulerEnabled.value = !value
    await nextTick()
    schedulerReady.value = true
    ElMessage.error('操作失败')
  } finally {
    schedulerSaving.value = false
  }
}

const ensureMenu = () => {
  if (!visibleMenus.value.some((item) => item.index === activeMenu.value)) {
    activeMenu.value = visibleMenus.value[0]?.index || ''
  }
}

const handleMenuSelect = (key) => {
  activeMenu.value = key
}

const onLogin = (admin) => {
  profile.value = admin
  authed.value = true
  ensureMenu()
  loadScheduler()
}

const logout = () => {
  clearSession()
  profile.value = {}
  authed.value = false
  passwordVisible.value = false
}

const onExpired = () => {
  profile.value = {}
  authed.value = false
}

const restore = async () => {
  if (!getToken()) {
    checking.value = false
    return
  }
  try {
    const response = await axios.get('/admin/me')
    if (response.data.code === 0) {
      profile.value = response.data.data
      setSession(getToken(), response.data.data)
      authed.value = true
      ensureMenu()
      loadScheduler()
    } else {
      clearSession()
    }
  } catch {
    clearSession()
  } finally {
    checking.value = false
  }
}

const changePassword = async () => {
  if (!passwordForm.old_password || !passwordForm.new_password) {
    ElMessage.warning('请填写原密码和新密码')
    return
  }
  passwordSaving.value = true
  try {
    const response = await axios.put('/admin/password', {
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '修改失败')
      return
    }
    ElMessage.success('密码已更新')
    passwordVisible.value = false
    passwordForm.old_password = ''
    passwordForm.new_password = ''
  } finally {
    passwordSaving.value = false
  }
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  document.title = 'WTV管理系统'
  window.addEventListener('admin-auth-expired', onExpired)
  restore()
})

onUnmounted(() => {
  window.removeEventListener('admin-auth-expired', onExpired)
})
</script>

<style scoped>
.boot-screen {
  width: 100%;
  height: 100vh;
  background: #070B12;
}

.app-container {
  width: 100%;
  height: 100vh;
  background: #070B12;
}

.app-layout {
  height: 100vh;
}

.app-body {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
}

.app-main {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0 !important;
}

.placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.aside {
  background-color: #0A101E;
  color: #EAF0F8;
  padding-top: 10px;
  border-right: 1px solid #1E2A3F;
}

.el-menu-vertical-demo {
  width: 100%;
  background-color: transparent;
  border-right: none;
}

.el-menu-item {
  color: #8E9AAF;
}

.el-menu-item.is-active {
  background-color: #1E3048;
  color: #78F0FF;
}

:deep(.el-menu-item:hover) {
  background-color: #151C2E !important;
  color: #EAF0F8;
}

:deep(.el-menu-item) {
  height: 42px;
  line-height: 42px;
  border-radius: 10px;
  margin: 4px 6px;
  padding: 0 10px 0 6px;
}

:deep(.el-menu-item__content) {
  padding-left: 0;
  justify-content: flex-start;
  font-size: 14px;
}

:deep(.el-menu-item__content span) {
  font-size: 14px;
}

.header {
  --el-header-height: 40px;
  height: 40px !important;
  padding: 0 16px;
  background-color: #0A101E;
  color: #EAF0F8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #1E2A3F;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.5);
}

.header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
  color: #EAF0F8;
}

.header-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-scheduler {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 8px;
  font-size: 13px;
  color: #8E9AAF;
}

.header-name {
  font-size: 13px;
  color: #EAF0F8;
}

.header-user :deep(.el-button.is-link) {
  color: #8E9AAF;
  font-size: 13px;
}

.header-user :deep(.el-button.is-link:hover) {
  color: #78F0FF;
}
</style>
