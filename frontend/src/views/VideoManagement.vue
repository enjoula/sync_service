<template>
  <div class="video-management">
    <!-- 搜索和筛选区域 -->
    <el-card class="search-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" class="search-form" size="small">
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="请输入视频标题" style="width: 200px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="请选择视频类型" style="width: 120px">
            <el-option label="电影" value="movie" />
            <el-option label="电视剧" value="tv" />
            <el-option label="动漫" value="anime" />
            <el-option label="综艺" value="variety" />
            <el-option label="纪录片" value="doc" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" style="width: 120px">
            <el-option label="返回" value="1" />
            <el-option label="不返回" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否完结">
          <el-select v-model="searchForm.is_completed" placeholder="请选择" style="width: 120px">
            <el-option label="已完结" value="1" />
            <el-option label="未完结" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="是否有更新">
          <el-select v-model="searchForm.is_update" placeholder="请选择" style="width: 120px">
            <el-option label="有更新" value="1" />
            <el-option label="无更新" value="0" />
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
            height="100%"
            @row-click="handleVideoClick"
          >
        <el-table-column prop="id" label="ID" width="120" fixed />
        <el-table-column prop="source_id" label="来源ID" width="100">
          <template #default="scope">
            {{ scope.row.source_id != null ? scope.row.source_id : '—' }}
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="100" show-overflow-tooltip />
        <el-table-column label="封面" width="80">
          <template #default="scope">
            <div class="cover-container">
              <el-image
                :src="scope.row.cover_url || 'https://via.placeholder.com/100x150'"
                fit="cover"
                style="width: 60px; height: 90px; cursor: pointer; transition: transform 0.3s"
              />
              <div class="cover-preview">
                <el-image
                  :src="scope.row.cover_url || 'https://via.placeholder.com/100x150'"
                  fit="cover"
                  style="width: 300px; height: 400px"
                />
              </div>
            </div>
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
        <el-table-column label="导演" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            {{ formatJsonDisplay(scope.row.director_json) }}
          </template>
        </el-table-column>
        <el-table-column label="演员" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            {{ formatJsonDisplay(scope.row.actors_json) }}
          </template>
        </el-table-column>
        <el-table-column label="标签" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            {{ formatJsonDisplay(scope.row.tags_json) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag size="small" :type="scope.row.status === '1' ? 'success' : 'info'">
              {{ scope.row.status === '1' ? '返回' : '不返回' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="imdb_id" label="IMDb ID" width="110" show-overflow-tooltip />
        <el-table-column prop="runtime" label="片长(分)" width="100">
          <template #default="scope">
            {{ scope.row.runtime != null ? scope.row.runtime : '—' }}
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
        <el-table-column prop="updated_at" label="更新时间" width="165" />
        <el-table-column label="操作" width="76" fixed="right" align="center">
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
          layout="total, sizes, prev, pager, next, jumper"
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
      width="720px"
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
        <el-form-item v-if="videoForm.id" label="视频ID">
          <el-input :model-value="videoForm.id" disabled />
        </el-form-item>

        <!-- 两列网格：把字段高度减半 -->
        <el-row :gutter="12">
          <el-col :span="12">
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
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="videoForm.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="返回" value="1" />
                <el-option label="不返回" value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="来源" prop="source">
              <el-input v-model="videoForm.source" placeholder="如 douban、xiaoya" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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
        </el-row>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="标题" prop="title">
              <el-input v-model="videoForm.title" placeholder="请输入视频标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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
        </el-row>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="评分" prop="score">
              <el-input-number v-model="videoForm.score" :min="0" :max="10" :step="0.1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="片长(分钟)" prop="runtime">
              <el-input-number
                v-model="videoForm.runtime"
                :min="0"
                :precision="0"
                :controls="true"
                placeholder="可选"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="剧集数" prop="episode_count">
              <el-input-number
                v-model="videoForm.episode_count"
                :min="0"
                :precision="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="是否完结">
              <el-switch v-model="videoForm.is_completed" active-text="已完结" inactive-text="未完结" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="是否有更新">
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

        <template v-if="videoForm.id">
          <el-row :gutter="12">
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
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'

// 搜索表单
const searchForm = reactive({
  title: '',
  type: '',
  status: '',
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
const videoFormRules = {
  title: [{ required: true, message: '请输入视频标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择视频类型', trigger: 'change' }]
}

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

// 获取视频列表
const getVideoList = async () => {
  try {
    const response = await axios.get('/admin/videos', {
      params: {
        page: pagination.currentPage,
        page_size: pagination.pageSize,
        title: searchForm.title,
        type: searchForm.type,
        status: searchForm.status,
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
  searchForm.type = ''
  searchForm.status = ''
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
    videoForm.status = video.status != null ? String(video.status) : '1'
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

    const videoData = {
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
      status: videoForm.status,
      imdb_id: videoForm.imdb_id,
      runtime: videoForm.runtime != null ? Number(videoForm.runtime) : null,
      episode_count: Number(videoForm.episode_count) || 0,
      is_completed: !!videoForm.is_completed,
      is_update: !!videoForm.is_update
    }

    if (videoForm.id) {
      // 更新
      await axios.put(`/admin/videos/${videoForm.id}`, videoData)
    } else {
      // 新增
      await axios.post('/admin/videos', videoData)
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
const getTypeTagType = (type) => {
  const typeMap = {
    movie: 'primary',
    tv: 'success',
    anime: 'warning',
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
  border: 1px solid #3a3a3a;
  border-radius: 8px;
  min-height: 32px;
  width: 100%;
  background: #1a1a1a;
}

.tag-input-container .el-tag {
  margin-bottom: 8px;
}

/* 封面预览样式 */
.cover-container {
  position: relative;
  display: inline-block;
}

.cover-preview {
  position: absolute;
  top: -10px;
  left: 70px;
  z-index: 1000;
  padding: 5px;
  background: rgba(0, 0, 0, 0.8);
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  pointer-events: none;
}

.cover-container:hover .cover-preview {
  display: block;
}

.cover-container .cover-preview {
  display: none;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}

/* 表格区吃掉卡片内剩余高度，横向/纵向滚动由 el-table 内部处理 */
.video-table-scroll {
  flex: 1;
  min-height: 0;
  width: 100%;
  max-width: 100%;
  overflow: hidden;
  -webkit-overflow-scrolling: touch;
}

.video-table-scroll :deep(.video-list-table.el-table) {
  width: max-content;
  min-width: 100%;
}

/* 确保卡片宽度适应容器 */
:deep(.el-card) {
  width: 100%;
}

.table-ops-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 4px 0;
}

.table-ops-cell .el-button {
  margin: 0;
  padding: 0 4px;
  height: auto;
  line-height: 1.35;
}

:deep(.video-form-dialog .el-dialog__body) {
  max-height: calc(100vh - 140px);
  overflow-y: auto;
  padding-right: 8px;
}

/* 表格更现代：表头底色与 hover */
:deep(.el-table__header th) {
  background-color: #242424;
  color: #c9c9c9;
  font-weight: 600;
}

:deep(.el-table__body tr:hover > td) {
  background-color: #2a2a2a !important;
}

/* 弹窗头部：黑灰层次 */
:deep(.video-form-dialog .el-dialog__header) {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.06), rgba(255, 255, 255, 0));
  border-bottom: 1px solid #333;
  padding: 14px 20px;
}
</style>