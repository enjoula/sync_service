<template>
  <div class="feedback-management">
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small" @submit.prevent>
        <el-form-item label="来源">
          <el-select v-model="searchForm.source" clearable placeholder="全部" style="width: 130px">
            <el-option v-for="item in sourceOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" clearable placeholder="全部" style="width: 120px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.problem_type" clearable placeholder="全部" style="width: 160px">
            <el-option v-for="item in problemTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="标题 / 描述 / 搜索词"
            clearable
            style="width: 200px"
            @keyup.enter="searchFeedbacks"
          />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input
            v-model="searchForm.user_id"
            placeholder="用户ID"
            clearable
            style="width: 160px"
            @keyup.enter="searchFeedbacks"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchFeedbacks">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="feedback-table-scroll">
          <el-table :data="feedbackList" class="feedback-list-table" border height="100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80" fixed />
            <el-table-column label="来源" width="110">
              <template #default="scope">
                <el-tag size="small" :type="scope.row.source === 'search' ? 'warning' : 'primary'" effect="plain">
                  {{ scope.row.source_label || scope.row.source }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="scope">
                <el-tag size="small" :type="statusTagType(scope.row.status)">
                  {{ scope.row.status_label || scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="问题类型" min-width="140" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.problem_label || scope.row.problem_type }}
              </template>
            </el-table-column>
            <el-table-column label="视频/关键词" min-width="180" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.video_title || scope.row.search_keyword || '—' }}
              </template>
            </el-table-column>
            <el-table-column label="剧集" width="80">
              <template #default="scope">
                {{ scope.row.episode_number != null ? `第${scope.row.episode_number}集` : '—' }}
              </template>
            </el-table-column>
            <el-table-column label="反馈内容" min-width="200" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.content || '—' }}
              </template>
            </el-table-column>
            <el-table-column label="用户" min-width="140" show-overflow-tooltip>
              <template #default="scope">
                {{ userDisplay(scope.row) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="反馈/搜索时间" width="165" />
            <el-table-column label="处理人" width="110" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.handled_by_name || '—' }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="90" fixed="right" align="center">
              <template #default="scope">
                <el-button size="small" type="primary" link @click="openDetail(scope.row)">处理</el-button>
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
            layout="total, sizes, prev, pager, next"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-drawer v-model="detailVisible" title="反馈详情" size="480px" destroy-on-close>
      <div v-if="current" class="detail-body">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="ID">{{ current.id }}</el-descriptions-item>
          <el-descriptions-item label="来源">
            <el-tag size="small" :type="current.source === 'search' ? 'warning' : 'primary'" effect="plain">
              {{ current.source_label || current.source }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag size="small" :type="statusTagType(current.status)">
              {{ current.status_label || current.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="问题类型">
            {{ current.problem_label || current.problem_type }}
          </el-descriptions-item>
          <el-descriptions-item label="用户">{{ userDisplay(current) }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ current.user_id }}</el-descriptions-item>
          <el-descriptions-item label="视频标题">{{ current.video_title || '—' }}</el-descriptions-item>
          <el-descriptions-item v-if="current.source === 'search'" label="搜索关键词">
            {{ current.search_keyword || '—' }}
          </el-descriptions-item>
          <el-descriptions-item label="视频ID">{{ current.video_id || '—' }}</el-descriptions-item>
          <el-descriptions-item label="剧集">
            {{ current.episode_number != null ? `第${current.episode_number}集` : '—' }}
          </el-descriptions-item>
          <el-descriptions-item label="反馈内容">{{ current.content || '—' }}</el-descriptions-item>
          <el-descriptions-item label="联系方式">{{ current.contact || '—' }}</el-descriptions-item>
          <el-descriptions-item label="设备">{{ current.client_device || '—' }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ current.client_ip || '—' }}</el-descriptions-item>
          <el-descriptions-item label="反馈/搜索时间">{{ current.created_at }}</el-descriptions-item>
          <el-descriptions-item label="处理人">{{ current.handled_by_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="处理时间">{{ current.handled_at || '—' }}</el-descriptions-item>
        </el-descriptions>

        <el-form v-if="canWrite" class="handle-form" label-width="72px" size="small">
          <el-form-item label="处理状态">
            <el-select v-model="handleForm.status" style="width: 100%">
              <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="处理备注">
            <el-input
              v-model="handleForm.admin_remark"
              type="textarea"
              :rows="4"
              maxlength="1000"
              show-word-limit
              placeholder="记录处理说明，便于后续追踪"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="saving" @click="saveHandle">保存处理</el-button>
          </el-form-item>
        </el-form>
        <el-alert
          v-else
          type="info"
          :closable="false"
          show-icon
          title="当前账号仅可查看反馈，处理需具备反馈写权限"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const props = defineProps({
  canWrite: {
    type: Boolean,
    default: false
  }
})

const sourceOptions = [
  { value: 'feedback', label: '主动反馈' },
  { value: 'search', label: '搜索无结果' }
]

const statusOptions = [
  { value: 'pending', label: '待处理' },
  { value: 'processing', label: '处理中' },
  { value: 'resolved', label: '已解决' },
  { value: 'rejected', label: '已驳回' },
  { value: 'closed', label: '已关闭' }
]

const problemTypeOptions = [
  { value: 'play_error', label: '无法播放/卡顿/花屏' },
  { value: 'wrong_content', label: '内容错误/不对版' },
  { value: 'missing_episode', label: '缺少剧集' },
  { value: 'missing_video', label: '缺少视频' },
  { value: 'subtitle_error', label: '字幕问题' },
  { value: 'audio_error', label: '音画不同步/无声' },
  { value: 'other', label: '其他' }
]

const searchForm = reactive({
  source: '',
  status: '',
  problem_type: '',
  keyword: '',
  user_id: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const loading = ref(false)
const feedbackList = ref([])
const detailVisible = ref(false)
const current = ref(null)
const saving = ref(false)
const handleForm = reactive({
  status: 'pending',
  admin_remark: ''
})

const statusTagType = (status) => {
  switch (status) {
    case 'pending':
      return 'warning'
    case 'processing':
      return ''
    case 'resolved':
      return 'success'
    case 'rejected':
      return 'danger'
    case 'closed':
      return 'info'
    default:
      return 'info'
  }
}

const userDisplay = (row) => {
  if (!row) return '—'
  if (row.nickname && row.username) return `${row.nickname} (${row.username})`
  return row.nickname || row.username || row.user_id || '—'
}

const getFeedbackList = async () => {
  loading.value = true
  try {
    const response = await axios.get('/admin/feedbacks', {
      params: {
        page: pagination.currentPage,
        page_size: pagination.pageSize,
        source: searchForm.source || undefined,
        status: searchForm.status || undefined,
        problem_type: searchForm.problem_type || undefined,
        keyword: searchForm.keyword.trim() || undefined,
        user_id: searchForm.user_id.trim() || undefined
      }
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '查询失败')
      return
    }
    feedbackList.value = response.data.data.list || []
    pagination.total = response.data.data.total || 0
  } catch (error) {
    console.error('获取反馈列表失败:', error)
    ElMessage.error('获取反馈列表失败')
  } finally {
    loading.value = false
  }
}

const searchFeedbacks = () => {
  pagination.currentPage = 1
  getFeedbackList()
}

const resetSearch = () => {
  searchForm.source = ''
  searchForm.status = ''
  searchForm.problem_type = ''
  searchForm.keyword = ''
  searchForm.user_id = ''
  searchFeedbacks()
}

const handleSizeChange = () => {
  pagination.currentPage = 1
  getFeedbackList()
}

const handleCurrentChange = () => {
  getFeedbackList()
}

const openDetail = async (row) => {
  try {
    const response = await axios.get(`/admin/feedbacks/${row.id}`)
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '获取详情失败')
      return
    }
    current.value = response.data.data
    handleForm.status = response.data.data.status || 'pending'
    handleForm.admin_remark = response.data.data.admin_remark || ''
    detailVisible.value = true
  } catch (error) {
    console.error('获取反馈详情失败:', error)
    ElMessage.error('获取反馈详情失败')
  }
}

const saveHandle = async () => {
  if (!current.value || !props.canWrite) return
  if (!handleForm.status) {
    ElMessage.warning('请选择处理状态')
    return
  }
  saving.value = true
  try {
    const response = await axios.put(`/admin/feedbacks/${current.value.id}`, {
      status: handleForm.status,
      admin_remark: handleForm.admin_remark
    })
    if (response.data.code !== 0) {
      ElMessage.error(response.data.message || '保存失败')
      return
    }
    current.value = response.data.data
    ElMessage.success('已更新处理状态')
    getFeedbackList()
  } catch (error) {
    console.error('更新反馈失败:', error)
    ElMessage.error('更新反馈失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  getFeedbackList()
})
</script>

<style scoped>
.feedback-management {
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

.feedback-table-scroll {
  flex: 1;
  min-height: 0;
  width: 100%;
  overflow: hidden;
}

.detail-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-bottom: 24px;
}

.handle-form {
  margin-top: 4px;
}
</style>
