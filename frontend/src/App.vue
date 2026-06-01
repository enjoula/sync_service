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
          <VideoManagement />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import VideoManagement from './views/VideoManagement.vue'
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
  background: linear-gradient(180deg, #121212 0%, #0a0a0a 100%);
}

.app-layout {
  height: 100vh;
}

.app-body {
  flex: 1;
  min-height: 0;
  display: flex;
}

.app-main {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0 !important;
}

.aside {
  background-color: #0d0d0d;
  color: #e5e5e5;
  padding-top: 10px;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
}

.el-menu-vertical-demo {
  width: 100%;
  background-color: transparent;
  border-right: none;
}

.el-menu-item {
  color: #a3a3a3;
}

.el-menu-item.is-active {
  background-color: #2a2a2a;
  color: #f5f5f5;
}

:deep(.el-menu-item:hover) {
  background-color: #1a1a1a !important;
  color: #e5e5e5;
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
  background-color: #050505;
  color: #e5e5e5;
  display: flex;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.5);
}

.header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
  color: #fafafa;
}
</style>