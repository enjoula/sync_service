<template>
  <div class="video-management">
    <!-- 搜索和筛选区域 -->
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small">
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="请输入视频标题" style="width: 200px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select
            v-model="searchForm.type"
            multiple
            collapse-tags
            collapse-tags-tooltip
            clearable
            placeholder="请选择视频类型"
            style="width: 180px"
          >
            <el-option label="电影" value="movie" />
            <el-option label="电视剧" value="tv" />
            <el-option label="动漫" value="anime" />
            <el-option label="综艺" value="tvshow" />
            <el-option label="纪录片" value="doc" />
          </el-select>
        </el-form-item>
        <el-form-item label="上映日期">
          <el-date-picker
            v-model="searchForm.release_date"
            type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            multiple
            collapse-tags
            collapse-tags-tooltip
            clearable
            placeholder="请选择状态"
            style="width: 160px"
          >
            <el-option label="上架" value="1" />
            <el-option label="下架" value="0" />
            <el-option label="未上架" value="__empty__" />
          </el-select>
        </el-form-item>
        <el-form-item label="完结">
          <el-select v-model="searchForm.is_completed" placeholder="请选择" style="width: 90px" clearable>
            <el-option label="已" value="1" />
            <el-option label="未" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="更新">
          <el-select v-model="searchForm.is_update" placeholder="请选择" style="width: 90px" clearable>
            <el-option label="有" value="1" />
            <el-option label="无" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" @click="searchVideos">搜索</el-button>
          <el-button size="small" @click="resetSearch">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="success" size="small" @click="openVideoForm">添加视频</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 视频列表（列多时在容器内横向滚动） -->
    <el-card class="list-card" shadow="hover">
      <div class="list-card-stretch">
        <div class="video-table-scroll">
          <el-table
            :data="videoList"
            class="video-list-table"
            border
            allow-drag-last-column
            height="100%"
            style="width: 100%"
            @row-click="handleVideoClick"
          >
        <el-table-column prop="id" label="ID" width="120" fixed />
        <el-table-column prop="source_id" label="来源ID" width="100">
          <template #default="scope">
            {{ scope.row.source_id != null ? scope.row.source_id : '—' }}
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="100" show-overflow-tooltip />
        <el-table-column label="封面" width="52" align="center">
          <template #default="scope">
            <div
              v-if="coverSrc(scope.row.cover_url)"
              class="cover-thumb"
              @mouseenter="showCoverPreview($event, coverSrc(scope.row.cover_url))"
              @mouseleave="hideCoverPreview"
            >
              <img
                :src="coverSrc(scope.row.cover_url)"
                alt=""
                class="cover-thumb-img"
                referrerpolicy="no-referrer"
              />
            </div>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="140" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            <el-tag size="small" :type="getTypeTagType(scope.row.type)">
              {{ getTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="简介" min-width="160" show-overflow-tooltip />
        <el-table-column prop="release_date" label="上映日期" width="120">
          <template #default="scope">
            {{ scope.row.release_date || '—' }}
          </template>
        </el-table-column>
        <el-table-column prop="score" label="评分" width="80">
          <template #default="scope">
            {{ scope.row.score != null ? scope.row.score : '—' }}
          </template>
        </el-table-column>
        <el-table-column label="国家/地区" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            {{ formatJsonDisplay(scope.row.country_json) }}
          </template>
        </el-table-column>
        <el-table-column label="标签" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            {{ formatJsonDisplay(scope.row.tags_json) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="scope">
            <el-tag size="small" :type="statusTagType(scope.row.status)">
              {{ statusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="episode_count" label="剧集数" width="90" />
        <el-table-column prop="is_completed" label="是否完结" width="90">
          <template #default="scope">
            <el-tag size="small" :type="scope.row.is_completed ? 'success' : 'warning'">
              {{ scope.row.is_completed ? '已完结' : '未完结' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_update" label="是否有更新" width="100">
          <template #default="scope">
            <el-tag size="small" :type="scope.row.is_update ? 'danger' : 'info'">
              {{ scope.row.is_update ? '有更新' : '无更新' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="165" />
        <el-table-column label="操作" width="96" fixed="right" align="center">
          <template #default="scope">
            <div class="table-ops-cell">
              <el-button size="small" type="primary" link @click.stop="openVideoForm(scope.row)">
                编辑
              </el-button>
              <el-button size="small" type="info" link @click.stop="openEpisodeDrawer(scope.row)">
                剧集
              </el-button>
            </div>
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

    <!-- 剧集管理抽屉 -->
    <el-drawer
      v-model="episodeDrawerVisible"
      title="剧集管理"
      direction="rtl"
      size="80%"
    >
      <div v-if="selectedVideo">
        <h3>{{ selectedVideo.title }} - 剧集列表</h3>
        <el-button type="primary" style="margin-bottom: 20px" @click="openEpisodeForm">添加剧集</el-button>
        <el-table :data="episodeList" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="episode_number" label="集数" width="80" />
          <el-table-column prop="name" label="名称" min-width="200" />
          <el-table-column prop="play_urls" label="播放地址" min-width="300" />
          <el-table-column prop="duration_seconds" label="时长(秒)" width="100" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button size="small" type="primary" @click="openEpisodeForm(scope.row)">
                编辑
              </el-button>
              <el-button size="small" type="danger" @click="deleteEpisode(scope.row.id)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else>
        <el-empty description="请选择一个视频" />
      </div>
    </el-drawer>

    <!-- 视频表单弹窗（与接口可写字段对齐） -->
    <el-dialog
      v-model="videoFormVisible"
      class="video-form-dialog"
      :title="videoForm.id ? '编辑视频' : '添加视频'"
      :width="videoForm.id ? '1180px' : '480px'"
      top="5vh"
      destroy-on-close
    >
      <el-form
        :model="videoForm"
        :rules="videoFormRules"
        ref="videoFormRef"
        label-width="96px"
        size="small"
      >
        <template v-if="!videoForm.id">
          <el-form-item label="来源ID" prop="source_id">
            <el-input-number
              v-model="videoForm.source_id"
              :min="1"
              :precision="0"
              :controls="false"
              placeholder="请填写来源ID"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="来源" prop="source">
            <el-input v-model="videoForm.source" placeholder="如 douban、xiaoya" />
          </el-form-item>
          <el-form-item label="标题" prop="title">
            <el-input v-model="videoForm.title" placeholder="请输入视频标题" />
          </el-form-item>
          <el-form-item label="类型" prop="type">
            <el-select v-model="videoForm.type" placeholder="请选择视频类型" style="width: 100%">
              <el-option label="电影" value="movie" />
              <el-option label="电视剧" value="tv" />
              <el-option label="动漫" value="anime" />
              <el-option label="综艺" value="variety" />
              <el-option label="纪录片" value="doc" />
            </el-select>
          </el-form-item>
          <el-form-item label="视频封面" prop="cover_url">
            <el-input v-model="videoForm.cover_url" placeholder="请填写封面图片地址" />
          </el-form-item>
        </template>
        <template v-else>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="视频ID">
              <el-input :model-value="videoForm.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="来源ID">
              <el-input-number
                v-model="videoForm.source_id"
                :min="0"
                :precision="0"
                :controls="false"
                placeholder="可选，可清空"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="来源" prop="source">
              <el-input v-model="videoForm.source" placeholder="如 douban、xiaoya" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="类型" prop="type">
              <el-select v-model="videoForm.type" placeholder="请选择视频类型" style="width: 100%">
                <el-option label="电影" value="movie" />
                <el-option label="电视剧" value="tv" />
                <el-option label="动漫" value="anime" />
                <el-option label="综艺" value="variety" />
                <el-option label="纪录片" value="doc" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="标题" prop="title">
              <el-input v-model="videoForm.title" placeholder="请输入视频标题" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-select v-model="videoForm.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="上架" value="1" />
                <el-option label="下架" value="0" />
                <el-option label="未上架" value="__empty__" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="评分" prop="score">
              <el-input-number v-model="videoForm.score" :min="0" :max="10" :step="0.1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="上映日期" prop="release_date">
              <el-date-picker
                v-model="videoForm.release_date"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="请选择上映日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="IMDb ID" prop="imdb_id">
              <el-input
                v-model="videoForm.imdb_id"
                placeholder="IMDb 编号"
                maxlength="20"
                show-word-limit
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="6">
            <el-form-item label="片长(分钟)" prop="runtime" label-width="92px">
              <el-input-number
                v-model="videoForm.runtime"
                :min="0"
                :precision="0"
                :controls="false"
                placeholder="可选"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="剧集数" prop="episode_count" label-width="72px">
              <el-input-number
                v-model="videoForm.episode_count"
                :min="0"
                :precision="0"
                :controls="false"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="是否完结" label-width="80px">
              <el-switch v-model="videoForm.is_completed" active-text="已完结" inactive-text="未完结" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="有无更新" label-width="80px">
              <el-switch v-model="videoForm.is_update" active-text="有更新" inactive-text="无更新" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 封面/简介占用较高，保持单列 -->
        <el-form-item label="封面URL" prop="cover_url">
          <el-input v-model="videoForm.cover_url" placeholder="封面图片地址" />
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input
            v-model="videoForm.description"
            type="textarea"
            rows="3"
            placeholder="请输入视频简介"
          />
        </el-form-item>

        <!-- JSON字段编辑 -->
        <el-form-item label="国家/地区">
          <div class="tag-input-container">
            <el-tag
              v-for="(tag, index) in videoForm.country_json"
              :key="index"
              closable
              @close="removeTag('country_json', index)"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-model="tagInput"
              placeholder="输入后按回车添加"
              @keyup.enter="addTag('country_json')"
              style="width: 150px; margin-left: 10px"
            />
          </div>
        </el-form-item>
        <el-form-item label="导演">
          <div class="tag-input-container">
            <el-tag
              v-for="(tag, index) in videoForm.director_json"
              :key="index"
              closable
              @close="removeTag('director_json', index)"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-model="tagInput"
              placeholder="输入后按回车添加"
              @keyup.enter="addTag('director_json')"
              style="width: 150px; margin-left: 10px"
            />
          </div>
        </el-form-item>
        <el-form-item label="演员">
          <div class="tag-input-container">
            <el-tag
              v-for="(tag, index) in videoForm.actors_json"
              :key="index"
              closable
              @close="removeTag('actors_json', index)"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-model="tagInput"
              placeholder="输入后按回车添加"
              @keyup.enter="addTag('actors_json')"
              style="width: 150px; margin-left: 10px"
            />
          </div>
        </el-form-item>
        <el-form-item label="标签">
          <div class="tag-input-container">
            <el-tag
              v-for="(tag, index) in videoForm.tags_json"
              :key="index"
              closable
              @close="removeTag('tags_json', index)"
            >
              {{ tag }}
            </el-tag>
            <el-input
              v-model="tagInput"
              placeholder="输入后按回车添加"
              @keyup.enter="addTag('tags_json')"
              style="width: 150px; margin-left: 10px"
            />
          </div>
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="创建时间">
              <el-input :model-value="videoForm.created_at || '—'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="更新时间">
              <el-input :model-value="videoForm.updated_at || '—'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="videoFormVisible = false">取消</el-button>
          <el-button type="primary" @click="saveVideo">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 剧集表单弹窗 -->
    <el-dialog
      v-model="episodeFormVisible"
      :title="episodeForm.id ? '编辑剧集' : '添加剧集'"
      width="600px"
    >
      <el-form :model="episodeForm" :rules="episodeFormRules" ref="episodeFormRef" label-width="100px">
        <el-form-item label="集数" prop="episode_number">
          <el-input-number v-model="episodeForm.episode_number" :min="1" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="episodeForm.name" placeholder="请输入剧集名称" />
        </el-form-item>
        <el-form-item label="播放地址" prop="play_urls">
          <el-input v-model="episodeForm.play_urls" placeholder="请输入播放地址" />
        </el-form-item>
        <el-form-item label="时长(秒)" prop="duration_seconds">
          <el-input-number v-model="episodeForm.duration_seconds" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="episodeFormVisible = false">取消</el-button>
          <el-button type="primary" @click="saveEpisode">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
  <teleport to="body">
    <div v-if="coverPreview.visible" class="cover-float" :style="coverPreview.style">
      <img :src="coverPreview.url" alt="封面预览" referrerpolicy="no-referrer" />
    </div>
  </teleport>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import axios from 'axios'

// 搜索表单
const searchForm = reactive({
  title: '',
  type: [],
  release_date: null,
  status: [],
  is_completed: '',
  is_update: ''
})

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 20, // 默认每页20条
  total: 0
})

// 视频列表
const videoList = ref([])
// 选中的视频
const selectedVideo = ref(null)
// 剧集列表
const episodeList = ref([])
// 剧集管理抽屉
const episodeDrawerVisible = ref(false)

// 视频表单
const videoFormVisible = ref(false)
const videoForm = reactive({
  id: '',
  source_id: null,
  source: '',
  title: '',
  type: '',
  cover_url: '',
  description: '',
  release_date: '',
  score: null,
  status: '1',
  imdb_id: '',
  runtime: null,
  episode_count: 0,
  is_completed: false,
  is_update: false,
  created_at: '',
  updated_at: '',
  country_json: [],
  director_json: [],
  actors_json: [],
  tags_json: []
})
const videoFormRef = ref(null)
const videoFormRules = computed(() => {
  if (videoForm.id) {
    return {
      title: [{ required: true, message: '请输入视频标题', trigger: 'blur' }],
      type: [{ required: true, message: '请选择视频类型', trigger: 'change' }]
    }
  }
  return {
    source_id: [{
      validator: (_rule, value, callback) => {
        if (value == null || value === '' || Number(value) <= 0) {
          callback(new Error('请填写来源ID'))
          return
        }
        callback()
      },
      trigger: 'blur'
    }],
    source: [{ required: true, whitespace: true, message: '请填写来源', trigger: 'blur' }],
    title: [{ required: true, whitespace: true, message: '请填写标题', trigger: 'blur' }],
    type: [{ required: true, message: '请选择类型', trigger: 'change' }],
    cover_url: [{ required: true, whitespace: true, message: '请填写视频封面', trigger: 'blur' }]
  }
})

// 剧集表单
const episodeFormVisible = ref(false)
const episodeForm = reactive({
  id: '',
  episode_number: 1,
  name: '',
  play_urls: '',
  duration_seconds: 0
})
const episodeFormRef = ref(null)
const episodeFormRules = {
  episode_number: [{ required: true, message: '请输入集数', trigger: 'blur' }],
  play_urls: [{ required: true, message: '请输入播放地址', trigger: 'blur' }]
}

// Tag输入
const tagInput = ref('')
const coverPreview = reactive({
  visible: false,
  url: '',
  style: {}
})

const coverSrc = (url) => {
  if (!url || typeof url !== 'string') return ''
  const trimmed = url.trim()
  if (!trimmed) return ''
  const normalized = trimmed.startsWith('//') ? `https:${trimmed}` : trimmed
  try {
    const parsed = new URL(normalized)
    if (/(\.|^)doubanio\.com$/i.test(parsed.hostname)) {
      return `/api/admin/cover?url=${encodeURIComponent(parsed.href)}`
    }
  } catch {
    return ''
  }
  return normalized
}

const showCoverPreview = (event, url) => {
  if (!url) return
  const rect = event.currentTarget.getBoundingClientRect()
  const width = 220
  const height = 320
  let left = rect.right + 12
  if (left + width > window.innerWidth - 8) {
    left = Math.max(8, rect.left - width - 12)
  }
  let top = rect.top - 8
  if (top + height > window.innerHeight - 8) {
    top = Math.max(8, window.innerHeight - height - 8)
  }
  coverPreview.url = url
  coverPreview.style = { left: `${left}px`, top: `${top}px` }
  coverPreview.visible = true
}

const hideCoverPreview = () => {
  coverPreview.visible = false
}

// 获取视频列表
const getVideoList = async () => {
  try {
    const response = await axios.get('/admin/videos', {
      params: {
        page: pagination.currentPage,
        page_size: pagination.pageSize,
        title: searchForm.title,
        type: searchForm.type.join(','),
        status: searchForm.status.join(','),
        release_date_start: searchForm.release_date?.[0] || '',
        release_date_end: searchForm.release_date?.[1] || '',
        is_completed: searchForm.is_completed,
        is_update: searchForm.is_update
      }
    })
    videoList.value = response.data.data.list
    pagination.total = response.data.data.total
  } catch (error) {
    console.error('获取视频列表失败:', error)
  }
}

// 搜索视频
const searchVideos = () => {
  pagination.currentPage = 1
  getVideoList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.title = ''
  searchForm.type = []
  searchForm.release_date = null
  searchForm.status = []
  searchForm.is_completed = ''
  searchForm.is_update = ''
  pagination.currentPage = 1
  getVideoList()
}

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size
  getVideoList()
}

const handleCurrentChange = (current) => {
  pagination.currentPage = current
  getVideoList()
}

// 视频行点击：仅打开剧集管理（编辑按钮已 @click.stop，不会进这里）
const handleVideoClick = (row) => {
  openEpisodeDrawer(row)
}

const openEpisodeDrawer = (row) => {
  selectedVideo.value = row
  getEpisodeList(row.id)
  episodeDrawerVisible.value = true
}

// 获取剧集列表
const getEpisodeList = async (videoId) => {
  try {
    const response = await axios.get('/admin/episodes', {
      params: { video_id: videoId }
    })
    episodeList.value = response.data.data
  } catch (error) {
    console.error('获取剧集列表失败:', error)
  }
}

const jsonFieldToTags = (val) => {
  if (val == null || val === '') return []
  if (Array.isArray(val)) return val.map((x) => String(x))
  if (typeof val === 'string') {
    try {
      const p = JSON.parse(val)
      return Array.isArray(p) ? p.map((x) => String(x)) : []
    } catch {
      return []
    }
  }
  return []
}

const resetVideoForm = () => {
  videoForm.id = ''
  videoForm.source_id = null
  videoForm.source = ''
  videoForm.title = ''
  videoForm.type = ''
  videoForm.cover_url = ''
  videoForm.description = ''
  videoForm.release_date = ''
  videoForm.score = null
  videoForm.status = '1'
  videoForm.imdb_id = ''
  videoForm.runtime = null
  videoForm.episode_count = 0
  videoForm.is_completed = false
  videoForm.is_update = false
  videoForm.created_at = ''
  videoForm.updated_at = ''
  videoForm.country_json = []
  videoForm.director_json = []
  videoForm.actors_json = []
  videoForm.tags_json = []
}

// 打开视频表单（@click.stop 避免触发行点击打开剧集抽屉）
const openVideoForm = (video = null) => {
  if (video) {
    videoForm.id = video.id
    videoForm.source_id = video.source_id != null ? Number(video.source_id) : null
    videoForm.source = video.source ?? ''
    videoForm.title = video.title ?? ''
    videoForm.type = video.type ?? ''
    videoForm.cover_url = video.cover_url ?? ''
    videoForm.description = video.description ?? ''
    videoForm.release_date = video.release_date || ''
    videoForm.score = video.score != null ? Number(video.score) : null
    videoForm.status = video.status === '0' || video.status === '1' ? String(video.status) : '__empty__'
    videoForm.imdb_id = video.imdb_id ?? ''
    videoForm.runtime = video.runtime != null ? Number(video.runtime) : null
    videoForm.episode_count = video.episode_count != null ? Number(video.episode_count) : 0
    videoForm.is_completed = !!video.is_completed
    videoForm.is_update = !!video.is_update
    videoForm.created_at = video.created_at ?? ''
    videoForm.updated_at = video.updated_at ?? ''
    videoForm.country_json = jsonFieldToTags(video.country_json)
    videoForm.director_json = jsonFieldToTags(video.director_json)
    videoForm.actors_json = jsonFieldToTags(video.actors_json)
    videoForm.tags_json = jsonFieldToTags(video.tags_json)
  } else {
    resetVideoForm()
  }
  videoFormVisible.value = true
}

// 保存视频
const saveVideo = async () => {
  if (!videoFormRef.value) return
  
  try {
    await videoFormRef.value.validate()

    if (videoForm.id) {
      await axios.put(`/admin/videos/${videoForm.id}`, {
        source_id: videoForm.source_id != null ? Number(videoForm.source_id) : null,
        source: videoForm.source,
        title: videoForm.title,
        type: videoForm.type,
        cover_url: videoForm.cover_url,
        description: videoForm.description,
        release_date: videoForm.release_date ? videoForm.release_date : null,
        score: videoForm.score,
        country_json: JSON.stringify(videoForm.country_json),
        director_json: JSON.stringify(videoForm.director_json),
        actors_json: JSON.stringify(videoForm.actors_json),
        tags_json: JSON.stringify(videoForm.tags_json),
        status: videoForm.status === '__empty__' ? '' : videoForm.status,
        imdb_id: videoForm.imdb_id,
        runtime: videoForm.runtime != null ? Number(videoForm.runtime) : null,
        episode_count: Number(videoForm.episode_count) || 0,
        is_completed: !!videoForm.is_completed,
        is_update: !!videoForm.is_update
      })
    } else {
      await axios.post('/admin/videos', {
        source_id: Number(videoForm.source_id),
        source: videoForm.source.trim(),
        title: videoForm.title.trim(),
        type: videoForm.type,
        cover_url: videoForm.cover_url.trim()
      })
    }
    
    videoFormVisible.value = false
    getVideoList()
  } catch (error) {
    console.error('保存视频失败:', error)
  }
}



// 打开剧集表单
const openEpisodeForm = (episode = null) => {
  if (episode) {
    // 编辑模式
    episodeForm.id = episode.id
    episodeForm.episode_number = episode.episode_number
    episodeForm.name = episode.name
    episodeForm.play_urls = episode.play_urls
    episodeForm.duration_seconds = episode.duration_seconds
  } else {
    // 新增模式
    episodeForm.id = ''
    episodeForm.episode_number = 1
    episodeForm.name = ''
    episodeForm.play_urls = ''
    episodeForm.duration_seconds = 0
  }
  episodeFormVisible.value = true
}

// 保存剧集
const saveEpisode = async () => {
  if (!episodeFormRef.value) return
  
  try {
    await episodeFormRef.value.validate()
    
    if (episodeForm.id) {
      // 更新
      await axios.put(`/admin/episodes/${episodeForm.id}`, episodeForm)
    } else {
      // 新增
      await axios.post('/admin/episodes', {
        ...episodeForm,
        video_id: selectedVideo.value.id
      })
    }
    
    episodeFormVisible.value = false
    getEpisodeList(selectedVideo.value.id)
  } catch (error) {
    console.error('保存剧集失败:', error)
  }
}

// 删除剧集
const deleteEpisode = async (id) => {
  try {
    await axios.delete(`/admin/episodes/${id}`)
    getEpisodeList(selectedVideo.value.id)
  } catch (error) {
    console.error('删除剧集失败:', error)
  }
}

// 添加Tag
const addTag = (field) => {
  if (tagInput.value && !videoForm[field].includes(tagInput.value)) {
    videoForm[field].push(tagInput.value)
    tagInput.value = ''
  }
}

// 移除Tag
const removeTag = (field, index) => {
  videoForm[field].splice(index, 1)
}

// 与 /api/admin/videos 返回的 JSON 数组字段一致展示（数组 / JSON 字符串 / 其他）
const formatJsonDisplay = (val) => {
  if (val == null || val === '') return '—'
  if (Array.isArray(val)) return val.length ? val.join('、') : '—'
  if (typeof val === 'string') {
    try {
      const parsed = JSON.parse(val)
      if (Array.isArray(parsed)) return parsed.length ? parsed.join('、') : '—'
      if (parsed && typeof parsed === 'object') return JSON.stringify(parsed)
      return String(parsed)
    } catch {
      return val
    }
  }
  if (typeof val === 'object') return JSON.stringify(val)
  return String(val)
}

// 获取类型标签类型
const statusText = (status) => {
  if (status === '1') return '上架'
  if (status === '0') return '下架'
  return '未上架'
}

const statusTagType = (status) => {
  if (status === '1') return 'success'
  if (status === '0') return 'info'
  return 'warning'
}

const getTypeTagType = (type) => {
  const typeMap = {
    movie: 'primary',
    tv: 'success',
    anime: 'warning',
    tvshow: 'info',
    variety: 'info',
    doc: 'danger'
  }
  return typeMap[type] || 'default'
}

// 获取类型文本
const getTypeText = (type) => {
  const typeMap = {
    movie: '电影',
    tv: '电视剧',
    anime: '动漫',
    tvshow: '综艺',
    variety: '综艺',
    doc: '纪录片'
  }
  return typeMap[type] || type
}

// 初始化
onMounted(() => {
  getVideoList()
})
</script>

<style scoped>
.video-management {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 16px;
  width: 100%;
  max-width: 100%;
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
  min-width: 0;
  width: 100%;
  max-width: 100%;
  min-height: 0;
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  overflow: hidden;
}

:deep(.search-card .el-card__body) {
  padding: 16px;
}

/* 列表卡片铺满主区剩余高度，el-card__body 贴近视口底部方向伸展 */
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
  min-width: 0;
  min-height: 0;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.search-form {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

/* 让搜索区域更紧凑，同时字号不至于太小 */
:deep(.search-form .el-form-item) {
  margin-bottom: 6px;
}

:deep(.search-form .el-form-item__label) {
  font-size: 13px;
}

:deep(.search-form .el-input__wrapper),
:deep(.search-form .el-select__wrapper),
:deep(.search-form .el-date-editor .el-input__wrapper) {
  padding-top: 0;
  padding-bottom: 0;
}

.pagination {
  flex-shrink: 0;
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
  padding-top: 4px;
}

.episode-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tag-input-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border: 1px solid #1E2A3F;
  border-radius: 8px;
  min-height: 32px;
  width: 100%;
  background: #0E1524;
}

.tag-input-container .el-tag {
  margin-bottom: 8px;
}

.cover-thumb {
  display: inline-flex;
  line-height: 0;
}

.cover-thumb-img {
  width: 28px;
  height: 40px;
  border-radius: 3px;
  object-fit: cover;
  display: block;
  cursor: zoom-in;
  background: #061022;
}

.cover-float {
  position: fixed;
  z-index: 4000;
  padding: 6px;
  background: #0E1524;
  border: 1px solid #78F0FF;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.45);
  pointer-events: none;
}

.cover-float img {
  display: block;
  width: 220px;
  height: 320px;
  object-fit: cover;
  border-radius: 4px;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}

/* 表格区吃掉卡片内剩余高度，横向/纵向滚动由 el-table 内部处理 */
.video-table-scroll {
  flex: 1;
  min-width: 0;
  min-height: 0;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
}

.video-table-scroll :deep(.video-list-table.el-table) {
  width: 100%;
  max-width: 100%;
  --el-table-border-color: #1E2A3F;
  font-size: 12px;
}

.video-table-scroll :deep(.video-list-table .el-table__header th),
.video-table-scroll :deep(.video-list-table .cell),
.video-table-scroll :deep(.video-list-table .el-tag),
.video-table-scroll :deep(.video-list-table .el-button) {
  font-size: 12px;
}

.video-table-scroll :deep(.video-list-table .el-table__cell) {
  padding: 2px 0;
}

.video-table-scroll :deep(.video-list-table .cell) {
  line-height: 18px;
  padding: 2px 6px;
}

.video-table-scroll :deep(.video-list-table .el-table__column-resize-proxy) {
  border-left-color: #78F0FF;
}

/* 确保卡片宽度适应容器 */
:deep(.el-card) {
  width: 100%;
}

.table-ops-cell {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 0;
}

.table-ops-cell .el-button {
  margin: 0;
  padding: 0 4px;
  height: auto;
  line-height: 1.35;
}

:deep(.video-form-dialog) {
  max-width: calc(100vw - 48px);
}

:deep(.video-form-dialog.el-dialog) {
  padding: 0;
}

:deep(.video-form-dialog .el-dialog__body) {
  max-height: calc(100vh - 140px);
  overflow-y: auto;
  padding: 16px 20px 8px;
}

:deep(.video-form-dialog .el-form-item__label) {
  white-space: nowrap;
}

:deep(.video-form-dialog .el-form-item__content) {
  min-width: 0;
}

/* 表格更现代：表头底色与 hover */
:deep(.el-table__header th) {
  background-color: #0E1524;
  color: #8E9AAF;
  font-weight: 600;
}

:deep(.el-table__body tr:hover > td) {
  background-color: #1E3048 !important;
}

/* 弹窗头部铺满，标题和关闭按钮同一行 */
:deep(.video-form-dialog .el-dialog__header) {
  background: #0E1524;
  border-bottom: 1px solid #1E2A3F;
  margin: 0;
  padding: 16px 48px 16px 20px;
}

:deep(.video-form-dialog .el-dialog__title) {
  color: #EAF0F8;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

:deep(.video-form-dialog .el-dialog__headerbtn) {
  top: 0;
  right: 0;
  width: 48px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.video-form-dialog .el-dialog__footer) {
  padding: 8px 20px 16px;
}
</style>