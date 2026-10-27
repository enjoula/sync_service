<template>
  <div class="app-container">
    <!-- 顶栏全宽，压在侧栏之上 -->
    <el-container class="app-layout" direction="vertical">
      <el-header height="40px" class="header">
        <h1>WTV管理系统</h1>
      </el-header>
      <el-container class="app-body">
        <el-aside width="130px" class="aside">
          <el-menu
            :default-active="activeMenu"
            class="el-menu-vertical-demo"
            router
            @select="handleMenuSelect"
          >
            <el-menu-item index="video">
              <el-icon><VideoCamera /></el-icon>
              <span>视频管理</span>
            </el-menu-item>
            <el-menu-item index="user">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
            <el-menu-item index="member">
              <el-icon><Star /></el-icon>
              <span>会员管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="app-main">
          <VideoManagement v-if="activeMenu === 'video'" />
          <UserManagement v-else-if="activeMenu === 'user'" />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import VideoManagement from './views/VideoManagement.vue'
import UserManagement from './views/UserManagement.vue'
import { VideoCamera, User, Star } from '@element-plus/icons-vue'

const activeMenu = ref('video')

const handleMenuSelect = (key) => {
  activeMenu.value = key
  // 这里可以添加路由跳转逻辑
  console.log('Selected menu:', key)
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  document.title = 'WTV管理系统'
})
</script>

<style scoped>
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
</style>