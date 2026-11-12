<template>
  <div class="admin-management">
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名或昵称"
            clearable
            style="width: 220px"
            @keyup.enter="searchAccounts"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchAccounts">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="success" size="small" @click="openCreate">新增管理员</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="admin-table-scroll">
          <el-table :data="accountList" class="admin-list-table" border height="100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" min-width="120" show-overflow-tooltip />
            <el-table-column prop="nickname" label="昵称" min-width="120" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.nickname || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="role_name" label="角色" width="120" />
            <el-table-column label="状态" width="90">
              <template #default="scope">
                <el-tag size="small" :type="scope.row.status === 1 ? 'success' : 'info'">
                  {{ scope.row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_login_at" label="最后登录" width="165">
              <template #default="scope">
                {{ scope.row.last_login_at || '—' }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="165" />
            <el-table-column label="操作" width="80" fixed="right" align="center">
              <template #default="scope">
                <el-button size="small" type="primary" link @click="openEdit(scope.row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="pagination">
          <el-pagination
            v-model:current-page="pagination.currentPage"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            class="admin-pagination"
            layout="sizes, prev, pager, next"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="formVisible" :title="form.id ? '编辑管理员' : '新增管理员'" width="420px">
      <el-form :model="form" label-width="72px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" :disabled="!!form.id" placeholder="4-32位字母、数字或下划线" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="可选" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role" :disabled="isSelf" style="width: 100%">
            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.id" label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" :disabled="isSelf" />
        </el-form-item>
        <el-form-item :label="form.id ? '重置密码' : '密码'">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="form.id ? '留空则不修改' : '6-64位'"
          />
        </el-form-item>
        <p class="role-hint">超级管理员可管理全部功能。内容编辑负责视频并可处理反馈。运营可查看视频、用户、会员并处理反馈。只读可查看视频与反馈。</p>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAccount">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getProfile } from '../auth'

const roleOptions = [
  { value: 'super', label: '超级管理员' },
  { value: 'editor', label: '内容编辑' },
  { value: 'operator', label: '运营' },
  { value: 'viewer', label: '只读' }
]

const searchForm = reactive({ keyword: '' })
const pagination = reactive({ currentPage: 1, pageSize: 20, total: 0 })
const accountList = ref([])
const formVisible = ref(false)
const saving = ref(false)
const form = reactive({
  id: 0,
  username: '',
  nickname: '',
  role: 'viewer',
  status: 1,
  password: ''
})

const isSelf = computed(() => {
  const profile = getProfile()
  return !!form.id && profile && profile.id === form.id
})

const getAccountList = async () => {
  const response = await axios.get('/admin/accounts', {
    params: {
      page: pagination.currentPage,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword.trim()
    }
  })
  if (response.data.code !== 0) {
    return
  }
  accountList.value = response.data.data.list || []
  pagination.total = response.data.data.total || 0
}

const searchAccounts = () => {
  pagination.currentPage = 1
  getAccountList()
}

const resetSearch = () => {
  searchForm.keyword = ''
  pagination.currentPage = 1
  getAccountList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  getAccountList()
}

const handleCurrentChange = (current) => {
  pagination.currentPage = current
  getAccountList()
}

const openCreate = () => {
  form.id = 0
  form.username = ''
  form.nickname = ''
  form.role = 'viewer'
  form.status = 1
  form.password = ''
  formVisible.value = true
}

const openEdit = (row) => {
  form.id = row.id
  form.username = row.username
  form.nickname = row.nickname || ''
  form.role = row.role
  form.status = row.status
  form.password = ''
  formVisible.value = true
}

const saveAccount = async () => {
  saving.value = true
  try {
    let response
    if (!form.id) {
      response = await axios.post('/admin/accounts', {
        username: form.username.trim(),
        password: form.password,
        nickname: form.nickname.trim(),
        role: form.role
      })
    } else {
      const payload = {
        nickname: form.nickname.trim(),
        role: form.role,
        status: form.status
      }
      if (form.password) {
        payload.password = form.password
      }
      response = await axios.put(`/admin/accounts/${form.id}`, payload)
    }
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '保存失败')
      return
    }
    ElMessage.success('已保存')
    formVisible.value = false
    getAccountList()
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  getAccountList()
})
</script>

<style scoped>
.admin-management {
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

.pagination {
  flex-shrink: 0;
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.admin-table-scroll {
  flex: 1;
  min-height: 0;
  width: 100%;
  overflow: hidden;
}

.role-hint {
  margin: 0;
  color: #8E9AAF;
  font-size: 12px;
  line-height: 1.6;
}
</style>
