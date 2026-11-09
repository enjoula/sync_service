<template>
  <div class="user-management">
    <el-tabs v-model="activeTab" class="user-tabs" @tab-change="onTabChange">
      <el-tab-pane label="用户列表" name="users" />
      <el-tab-pane label="用户观看记录" name="watches" />
      <el-tab-pane label="用户登录信息" name="logins" />
    </el-tabs>

    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" class="search-form" size="small" @submit.prevent>
        <el-form-item label="关键词">
          <el-input
            v-model="keywords[activeTab]"
            :placeholder="keywordPlaceholder"
            clearable
            style="width: 260px"
            @keyup.enter="searchCurrent"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchCurrent">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="user-table-scroll">
          <el-table
            v-if="activeTab === 'users'"
            :data="userList"
            class="user-list-table"
            border
            allow-drag-last-column
            height="100%"
          >
            <el-table-column prop="id" label="ID" width="180" fixed show-overflow-tooltip />
            <el-table-column label="头像" width="72" align="center">
              <template #default="scope">
                <el-image
                  v-if="scope.row.avatar"
                  :src="scope.row.avatar"
                  fit="cover"
                  class="avatar"
                />
                <span v-else>—</span>
              </template>
            </el-table-column>
            <el-table-column prop="username" label="用户名" min-width="120" show-overflow-tooltip />
            <el-table-column prop="nickname" label="昵称" min-width="120" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.nickname || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.email || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="watch_count" label="观看次数" width="90" />
            <el-table-column prop="last_watch_at" label="最后观看" width="165">
              <template #default="scope">
                {{ scope.row.last_watch_at || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="last_open_at" label="最近打开" width="165">
              <template #default="scope">
                {{ scope.row.last_open_at || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="165" />
            <el-table-column prop="updated_at" label="更新时间" width="165" />
          </el-table>

          <el-table
            v-else-if="activeTab === 'watches'"
            :data="watchList"
            class="user-list-table"
            border
            allow-drag-last-column
            height="100%"
          >
            <el-table-column prop="user_id" label="用户ID" width="180" fixed show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" min-width="120" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.username || '—' }}</template>
            </el-table-column>
            <el-table-column prop="nickname" label="昵称" min-width="120" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.nickname || '—' }}</template>
            </el-table-column>
            <el-table-column prop="video_title" label="视频" min-width="180" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.video_title || '—' }}</template>
            </el-table-column>
            <el-table-column label="剧集" min-width="160" show-overflow-tooltip>
              <template #default="scope">{{ episodeText(scope.row) }}</template>
            </el-table-column>
            <el-table-column prop="play_count" label="观看次数" width="90" />
            <el-table-column label="进度" width="90">
              <template #default="scope">{{ formatProgress(scope.row.last_position_ms) }}</template>
            </el-table-column>
            <el-table-column prop="last_played_at" label="最后观看" width="165">
              <template #default="scope">{{ scope.row.last_played_at || '—' }}</template>
            </el-table-column>
            <el-table-column prop="created_at" label="首次观看" width="165">
              <template #default="scope">{{ scope.row.created_at || '—' }}</template>
            </el-table-column>
          </el-table>

          <el-table
            v-else
            :data="loginList"
            class="user-list-table"
            border
            allow-drag-last-column
            height="100%"
          >
            <el-table-column prop="user_id" label="用户ID" width="180" fixed show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" min-width="120" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.username || '—' }}</template>
            </el-table-column>
            <el-table-column prop="nickname" label="昵称" min-width="120" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.nickname || '—' }}</template>
            </el-table-column>
            <el-table-column prop="device" label="设备" min-width="140" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.device || '—' }}</template>
            </el-table-column>
            <el-table-column prop="ip_address" label="IP" min-width="140" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.ip_address || '—' }}</template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="scope">
                <el-tag size="small" :type="loginStatusType(scope.row.status)">
                  {{ loginStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="expires_at" label="过期时间" width="165">
              <template #default="scope">{{ scope.row.expires_at || '—' }}</template>
            </el-table-column>
            <el-table-column prop="created_at" label="登录时间" width="165">
              <template #default="scope">{{ scope.row.created_at || '—' }}</template>
            </el-table-column>
          </el-table>
        </div>
        <div class="pagination">
          <el-pagination
            :current-page="currentPagination.currentPage"
            :page-size="currentPagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            class="admin-pagination"
            layout="sizes, prev, pager, next"
            :total="currentPagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { computed, reactive, ref, onMounted } from 'vue'
import axios from 'axios'

const activeTab = ref('users')
const keywords = reactive({
  users: '',
  watches: '',
  logins: ''
})
const paginations = reactive({
  users: { currentPage: 1, pageSize: 20, total: 0 },
  watches: { currentPage: 1, pageSize: 20, total: 0 },
  logins: { currentPage: 1, pageSize: 20, total: 0 }
})
const loaded = reactive({
  users: false,
  watches: false,
  logins: false
})

const userList = ref([])
const watchList = ref([])
const loginList = ref([])

const currentPagination = computed(() => paginations[activeTab.value])
const keywordPlaceholder = computed(() => {
  if (activeTab.value === 'watches') return '用户名、昵称或视频标题'
  if (activeTab.value === 'logins') return '用户名、昵称、设备或 IP'
  return '用户名、昵称或邮箱'
})

const listParams = (tab) => ({
  page: paginations[tab].currentPage,
  page_size: paginations[tab].pageSize,
  keyword: keywords[tab].trim()
})

const getUserList = async () => {
  try {
    const response = await axios.get('/admin/users', { params: listParams('users') })
    userList.value = response.data.data.list || []
    paginations.users.total = response.data.data.total || 0
    loaded.users = true
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

const getWatchList = async () => {
  try {
    const response = await axios.get('/admin/users/watch-records', { params: listParams('watches') })
    watchList.value = response.data.data.list || []
    paginations.watches.total = response.data.data.total || 0
    loaded.watches = true
  } catch (error) {
    console.error('获取观看记录失败:', error)
  }
}

const getLoginList = async () => {
  try {
    const response = await axios.get('/admin/users/logins', { params: listParams('logins') })
    loginList.value = response.data.data.list || []
    paginations.logins.total = response.data.data.total || 0
    loaded.logins = true
  } catch (error) {
    console.error('获取登录信息失败:', error)
  }
}

const loaders = {
  users: getUserList,
  watches: getWatchList,
  logins: getLoginList
}

const onTabChange = (name) => {
  if (!loaded[name]) {
    loaders[name]()
  }
}

const searchCurrent = () => {
  paginations[activeTab.value].currentPage = 1
  loaders[activeTab.value]()
}

const resetSearch = () => {
  keywords[activeTab.value] = ''
  paginations[activeTab.value].currentPage = 1
  loaders[activeTab.value]()
}

const handleSizeChange = (size) => {
  paginations[activeTab.value].pageSize = size
  paginations[activeTab.value].currentPage = 1
  loaders[activeTab.value]()
}

const handleCurrentChange = (current) => {
  paginations[activeTab.value].currentPage = current
  loaders[activeTab.value]()
}

const episodeText = (row) => {
  const number = row.episode_number ? `第${row.episode_number}集` : ''
  const name = row.episode_name || ''
  if (number && name) return `${number} ${name}`
  return number || name || '—'
}

const formatProgress = (ms) => {
  const total = Math.max(0, Math.floor(Number(ms || 0) / 1000))
  const hours = Math.floor(total / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const seconds = total % 60
  const pad = (value) => String(value).padStart(2, '0')
  if (hours > 0) return `${hours}:${pad(minutes)}:${pad(seconds)}`
  return `${pad(minutes)}:${pad(seconds)}`
}

const loginStatusText = (status) => {
  if (status === 'active') return '有效'
  if (status === 'expired') return '已过期'
  if (status === 'inactive') return '已失效'
  return '—'
}

const loginStatusType = (status) => {
  if (status === 'active') return 'success'
  if (status === 'expired') return 'warning'
  return 'info'
}

onMounted(() => {
  getUserList()
})
</script>

<style scoped>
.user-management {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 16px;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

.user-tabs {
  flex-shrink: 0;
}

.user-tabs :deep(.el-tabs__header) {
  margin: 0 0 12px;
}

.user-tabs :deep(.el-tabs__item) {
  color: #8E9AAF;
  font-size: 14px;
}

.user-tabs :deep(.el-tabs__item.is-active) {
  color: #78F0FF;
}

.user-tabs :deep(.el-tabs__active-bar) {
  background-color: #78F0FF;
}

.user-tabs :deep(.el-tabs__nav-wrap::after) {
  background-color: #1E2A3F;
}

.user-tabs :deep(.el-tabs__content) {
  display: none;
}

.search-card {
  flex-shrink: 0;
  margin-bottom: 16px;
  width: 100%;
  border-radius: 12px;
}

.list-card {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
}

:deep(.search-card .el-card__body) {
  padding: 16px;
}

:deep(.list-card .el-card__body) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 16px;
}

.list-card-stretch {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.search-form {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.search-form .el-form-item) {
  margin-bottom: 6px;
}

:deep(.search-form .el-form-item__label) {
  font-size: 13px;
}

.pagination {
  flex-shrink: 0;
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
  padding-top: 4px;
}

.user-table-scroll {
  flex: 1;
  min-height: 0;
  width: 100%;
  overflow: hidden;
}

.user-table-scroll :deep(.user-list-table.el-table) {
  width: max-content;
  min-width: 100%;
  --el-table-border-color: #1E2A3F;
  font-size: 12px;
}

.user-table-scroll :deep(.user-list-table .el-table__header th),
.user-table-scroll :deep(.user-list-table .cell) {
  font-size: 12px;
}

.user-table-scroll :deep(.user-list-table .el-table__column-resize-proxy) {
  border-left-color: #78F0FF;
}

:deep(.el-table__header th) {
  background-color: #0E1524;
  color: #8E9AAF;
  font-weight: 600;
}

:deep(.el-table__body tr:hover > td) {
  background-color: #1E3048 !important;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 6px;
}
</style>
