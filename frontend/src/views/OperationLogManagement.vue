<template>
  <div class="oplog-management">
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small" @submit.prevent>
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.action" clearable placeholder="全部" style="width: 110px">
            <el-option label="新增" value="create" />
            <el-option label="编辑" value="update" />
            <el-option label="删除" value="delete" />
          </el-select>
        </el-form-item>
        <el-form-item label="对象">
          <el-select v-model="searchForm.object_type" clearable placeholder="全部" style="width: 130px">
            <el-option label="视频" value="video" />
            <el-option label="剧集" value="episode" />
            <el-option label="管理员账号" value="account" />
            <el-option label="用户反馈" value="feedback" />
            <el-option label="定时任务" value="scheduler" />
            <el-option label="登录密码" value="password" />
          </el-select>
        </el-form-item>
        <el-form-item label="结果">
          <el-select v-model="searchForm.result" clearable placeholder="全部" style="width: 100px">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failure" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="操作人 / 对象名 / ID"
            clearable
            style="width: 200px"
            @keyup.enter="searchLogs"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchLogs">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="oplog-table-scroll">
          <el-table :data="logList" class="oplog-list-table" border height="100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80" fixed />
            <el-table-column prop="created_at" label="操作时间" width="165" />
            <el-table-column prop="admin_username" label="操作人" width="120" show-overflow-tooltip />
            <el-table-column label="类型" width="80">
              <template #default="scope">
                {{ scope.row.action_label || scope.row.action }}
              </template>
            </el-table-column>
            <el-table-column label="对象" width="110">
              <template #default="scope">
                {{ scope.row.object_label || scope.row.object_type }}
              </template>
            </el-table-column>
            <el-table-column prop="object_name" label="对象名称" min-width="140" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.object_name || '—' }}</template>
            </el-table-column>
            <el-table-column prop="object_id" label="对象ID" min-width="140" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.object_id || '—' }}</template>
            </el-table-column>
            <el-table-column label="结果" width="80">
              <template #default="scope">
                <el-tag size="small" :type="scope.row.result === 'success' ? 'success' : 'danger'">
                  {{ scope.row.result_label || scope.row.result }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="client_ip" label="IP" width="130" show-overflow-tooltip>
              <template #default="scope">{{ scope.row.client_ip || '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right" align="center">
              <template #default="scope">
                <el-button size="small" type="primary" link @click="openDetail(scope.row)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <div class="pagination">
          <el-pagination
            v-model:current-page="pagination.currentPage"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-drawer v-model="detailVisible" title="操作日志详情" size="520px" destroy-on-close>
      <div v-if="current" class="detail-body">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="ID">{{ current.id }}</el-descriptions-item>
          <el-descriptions-item label="操作时间">{{ current.created_at }}</el-descriptions-item>
          <el-descriptions-item label="操作人">
            {{ current.admin_username }}（ID: {{ current.admin_id }}）
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            {{ current.action_label || current.action }}
          </el-descriptions-item>
          <el-descriptions-item label="操作对象">
            {{ current.object_label || current.object_type }}
          </el-descriptions-item>
          <el-descriptions-item label="对象名称">{{ current.object_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="对象ID">{{ current.object_id || '—' }}</el-descriptions-item>
          <el-descriptions-item label="结果">
            <el-tag size="small" :type="current.result === 'success' ? 'success' : 'danger'">
              {{ current.result_label || current.result }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="失败原因">{{ current.error_msg || '—' }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ current.client_ip || '—' }}</el-descriptions-item>
        </el-descriptions>

        <div class="json-block">
          <div class="json-title">操作前</div>
          <pre>{{ formatJSON(current.before_data) }}</pre>
        </div>
        <div class="json-block">
          <div class="json-title">操作后</div>
          <pre>{{ formatJSON(current.after_data) }}</pre>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const searchForm = reactive({
  action: '',
  object_type: '',
  result: '',
  keyword: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const loading = ref(false)
const logList = ref([])
const detailVisible = ref(false)
const current = ref(null)

const formatJSON = (raw) => {
  if (!raw) return '—'
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

const getLogList = async () => {
  loading.value = true
  try {
    const response = await axios.get('/admin/operation-logs', {
      params: {
        page: pagination.currentPage,
        page_size: pagination.pageSize,
        action: searchForm.action || undefined,
        object_type: searchForm.object_type || undefined,
        result: searchForm.result || undefined,
        keyword: searchForm.keyword.trim() || undefined
      }
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '查询失败')
      return
    }
    logList.value = response.data.data.list || []
    pagination.total = response.data.data.total || 0
  } catch (error) {
    console.error('获取操作日志失败:', error)
    ElMessage.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

const searchLogs = () => {
  pagination.currentPage = 1
  getLogList()
}

const resetSearch = () => {
  searchForm.action = ''
  searchForm.object_type = ''
  searchForm.result = ''
  searchForm.keyword = ''
  searchLogs()
}

const handleSizeChange = () => {
  pagination.currentPage = 1
  getLogList()
}

const handleCurrentChange = () => {
  getLogList()
}

const openDetail = async (row) => {
  try {
    const response = await axios.get(`/admin/operation-logs/${row.id}`)
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '获取详情失败')
      return
    }
    current.value = response.data.data
    detailVisible.value = true
  } catch (error) {
    console.error('获取操作日志详情失败:', error)
    ElMessage.error('获取操作日志详情失败')
  }
}

onMounted(() => {
  getLogList()
})
</script>

<style scoped>
.oplog-management {
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

.oplog-table-scroll {
  flex: 1;
  min-height: 0;
  width: 100%;
  overflow: hidden;
}

.detail-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 24px;
}

.json-block {
  border: 1px solid #1E2A3F;
  border-radius: 8px;
  overflow: hidden;
}

.json-title {
  padding: 8px 12px;
  background: #0A101E;
  color: #8E9AAF;
  font-size: 13px;
}

.json-block pre {
  margin: 0;
  padding: 12px;
  max-height: 240px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  color: #EAF0F8;
  background: #070B12;
}
</style>
