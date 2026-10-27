<template>
  <div class="user-management">
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名、昵称或邮箱"
            clearable
            style="width: 220px"
            @keyup.enter="searchUsers"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchUsers">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="user-table-scroll">
          <el-table
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
        </div>
        <div class="pagination">
          <el-pagination
            v-model:current-page="pagination.currentPage"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import axios from 'axios'

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const userList = ref([])

const getUserList = async () => {
  try {
    const response = await axios.get('/admin/users', {
      params: {
        page: pagination.currentPage,
        page_size: pagination.pageSize,
        keyword: searchForm.keyword.trim()
      }
    })
    userList.value = response.data.data.list || []
    pagination.total = response.data.data.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

const searchUsers = () => {
  pagination.currentPage = 1
  getUserList()
}

const resetSearch = () => {
  searchForm.keyword = ''
  pagination.currentPage = 1
  getUserList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  getUserList()
}

const handleCurrentChange = (current) => {
  pagination.currentPage = current
  getUserList()
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
