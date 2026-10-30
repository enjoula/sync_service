// repository 包提供数据访问层，封装数据库操作
package repository

import (
	"strconv"
	"time"
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"gorm.io/gorm"
)

// VideoRef 状态收尾日志用的视频标识。ID 用字符串，避免雪花 ID 在 JSON 里丢精度。
type VideoRef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type idTitleRow struct {
	ID    int64  `gorm:"column:id"`
	Title string `gorm:"column:title"`
}

func toVideoRefs(rows []idTitleRow) []VideoRef {
	refs := make([]VideoRef, 0, len(rows))
	for _, row := range rows {
		refs = append(refs, VideoRef{
			ID:    strconv.FormatInt(row.ID, 10),
			Title: row.Title,
		})
	}
	return refs
}

func queryIDTitles(query string, args ...interface{}) ([]idTitleRow, error) {
	rows := make([]idTitleRow, 0)
	err := database.DB.Raw(query, args...).Scan(&rows).Error
	return rows, err
}

// VideoRepository 视频仓库接口
type VideoRepository interface {
	// FindBySourceID 根据来源ID查找视频
	FindBySourceID(sourceID int64) (*model.Video, error)

	// Create 创建视频记录
	Create(video *model.Video) error

	// CreateBasic 只写入来源 ID、来源、标题、类型、封面，以及雪花主键和创建/更新时间
	CreateBasic(video *model.Video) error

	// Update 更新视频记录
	Update(video *model.Video) error

	// UpdateDetails 更新视频详情字段（只更新指定的详情字段，不会覆盖其他字段）
	UpdateDetails(video *model.Video) error

	// FindNeedDetailVideos 查找需要补充详情的视频（source_id不为空且release_date和country_json都为空）
	FindNeedDetailVideos(limit int) ([]*model.Video, error)

	// FindNeedDetailVideosByType 根据类型查找需要补充详情的视频（source_id不为空且release_date和country_json都为空）
	FindNeedDetailVideosByType(videoType string, limit int) ([]*model.Video, error)

	// FindAllVideos 查找所有视频（仅返回 id 和 title）
	FindAllVideos() ([]*model.Video, error)

	// FindVideosWithPagination 查找视频列表（带分页、排序和筛选）
	FindVideosWithPagination(limit, offset int, title string, videoTypes, statuses []string, releaseStart, releaseEnd *time.Time, isCompleted *bool, isUpdate *bool) ([]*model.Video, int64, error)

	// FindVideosByStatusNotEqual 查找 status 不等于指定值的视频（返回 id、type、title）
	FindVideosByStatusNotEqual(status string) ([]*model.Video, error)

	// FindVideosNeedUpdateEpisodes 查找需要搜索播放地址的视频（跳过 status=0 和 is_completed=1，title有内容）
	FindVideosNeedUpdateEpisodes() ([]*model.Video, error)

	// FindVideosNeedRefreshEpisodeCount 查找需要从豆瓣刷新集数的视频（status!=0、is_completed!=1、有 source_id，不含 movie、tvshow）
	FindVideosNeedRefreshEpisodeCount() ([]*model.Video, error)

	// UpdateEpisodeCountWithoutUpdatedAt 只更新 episode_count，不修改 updated_at
	UpdateEpisodeCountWithoutUpdatedAt(videoID int64, episodeCount int64) error

	// UpdateVideosStatusByEpisodes 更新存在 episodes 记录的 videos 的 status。status 为 0、episode_count 为 0 或空的不处理。
	UpdateVideosStatusByEpisodes(status string) error

	// ShelfMoviesWithEpisodes 有剧集的电影上架并标为完结。status 为 0 或已完结的不处理。
	ShelfMoviesWithEpisodes() (int64, error)

	// UpdateVideoStatus 更新指定视频的status
	// 返回 (是否真的执行了更新, 错误)
	UpdateVideoStatus(videoID int64, status string) (bool, error)

	// UpdateVideoIsUpdate 更新指定视频的is_update字段
	// 返回 (是否真的执行了更新, 错误)
	UpdateVideoIsUpdate(videoID int64, isUpdate bool) (bool, error)

	// UpdateVideoIsCompleted 更新指定视频的is_completed字段
	// 返回 (是否真的执行了更新, 错误)
	UpdateVideoIsCompleted(videoID int64, isCompleted bool) (bool, error)

	// UpdateVideoUpdatedAt 更新指定视频的updated_at字段
	UpdateVideoUpdatedAt(videoID int64, updatedAt time.Time) error

	// UpdateVideoStateAfterEpisodeInsert 插入剧集后写上架、有更新，并把 updated_at 设为最后一集的 created_at。markCompleted 为真时同时标完结。
	UpdateVideoStateAfterEpisodeInsert(videoID int64, updatedAt time.Time, markCompleted bool) error

	// ClearIsUpdateOlderThan 将 updated_at 早于 since、且 status 不为 0 的视频 is_update 设为 0。不改 updated_at。返回实际改到的视频。
	ClearIsUpdateOlderThan(since time.Time) ([]VideoRef, error)

	// ApplyInsertedEpisodeVideoState 对本次新插入剧集的视频统一写上架、有更新和 updated_at。电影标完结；其它类型集数已齐才标完结。status 为 0 的不处理。返回实际改到的视频。
	ApplyInsertedEpisodeVideoState(videoIDs []int64) ([]VideoRef, error)

	// MarkSeriesCompletedWhenCountMatches 非电影现有集数等于 episode_count 时只把 is_completed 设为 1。status 为 0、集数为 0 的不处理。返回实际改到的视频。
	MarkSeriesCompletedWhenCountMatches() ([]VideoRef, error)

	// ShelfPendingMoviesWithEpisodes 已有剧集且未完结的电影上架并标完结，updated_at 取最后一集，is_update 先记为 1。status 为 0 或已完结的不处理。返回实际改到的视频。
	ShelfPendingMoviesWithEpisodes() ([]VideoRef, error)

	// AlignUpdatedAtWithLastEpisode 将有剧集且 status 不为 0 的视频 updated_at 对齐为最后一条剧集的 created_at，已一致的不更新
	AlignUpdatedAtWithLastEpisode() (int64, error)

	// FindByID 根据ID查找视频（返回完整信息）
	FindByID(videoID int64) (*model.Video, error)

	// FindVideosByIsUpdate 查找 is_update = 1 的视频（返回 id）
	FindVideosByIsUpdate() ([]*model.Video, error)

	// FindVideosWithEpisodes 查找所有有episode的视频（返回 id）
	FindVideosWithEpisodes() ([]*model.Video, error)

	// FindVideosWithEpisodesAndStatus1 查找 status=1 且有episode的视频（返回 id）
	FindVideosWithEpisodesAndStatus1() ([]*model.Video, error)

	// FindVideosWithEpisodesStatus1AndRecentRelease 查找 status=1、release_date在最近两个月内且有episode的视频（返回 id）
	FindVideosWithEpisodesStatus1AndRecentRelease() ([]*model.Video, error)

	// FindVideosStatus1AndRecentRelease 查找 status=1、release_date在最近两个月内的视频（返回 id）
	FindVideosStatus1AndRecentRelease() ([]*model.Video, error)

	// FindVideosStatus1AndRecentCreated 查找 status=1、created_at在最近两个月内的视频（返回 id）
	FindVideosStatus1AndRecentCreated() ([]*model.Video, error)

	// FindVideosStatus1OrIsUpdate1AndRecentCreated 查找需要维护 is_update 的视频（返回 id）。status 为 0 的不处理。
	FindVideosStatus1OrIsUpdate1AndRecentCreated() ([]*model.Video, error)

	// DeleteByID 根据ID删除视频
	DeleteByID(videoID int64) error
}

// videoRepository 视频仓库实现
type videoRepository struct{}

// NewVideoRepository 创建视频仓库实例
func NewVideoRepository() VideoRepository {
	return &videoRepository{}
}

// FindBySourceID 根据来源ID查找视频
func (r *videoRepository) FindBySourceID(sourceID int64) (*model.Video, error) {
	var video model.Video
	err := database.DB.Where("source_id = ?", sourceID).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// Create 创建视频记录
func (r *videoRepository) Create(video *model.Video) error {
	return database.DB.Create(video).Error
}

// CreateBasic 只写入来源 ID、来源、标题、类型、封面，以及雪花主键和创建/更新时间。
// 其余列留给数据库默认值，避免把未填写字段写成空串。
func (r *videoRepository) CreateBasic(video *model.Video) error {
	return database.DB.Select("ID", "SourceID", "Source", "Title", "Type", "CoverURL", "CreatedAt", "UpdatedAt").Create(video).Error
}

// Update 更新视频记录（更新所有字段）
func (r *videoRepository) Update(video *model.Video) error {
	return database.DB.Save(video).Error
}

// UpdateDetails 更新视频详情字段（只更新指定的详情字段，不会覆盖其他字段）
func (r *videoRepository) UpdateDetails(video *model.Video) error {
	// 使用 Select 配合 UpdateColumns 来更新指定字段
	// UpdateColumns 会更新所有字段（包括零值），Select 确保只更新指定的字段
	updates := map[string]interface{}{
		"description":   video.Description,
		"release_date":  video.ReleaseDate,
		"country_json":  video.CountryJSON,
		"director_json": video.DirectorJSON,
		"actors_json":   video.ActorsJSON,
		"tags_json":     video.TagsJSON,
		"imdb_id":       video.IMDbID,
		"runtime":       video.Runtime,
		"score":         video.Score,
		"episode_count": video.EpisodeCount,
		"updated_at":    video.UpdatedAt,
	}

	return database.DB.Model(video).
		Where("id = ?", video.ID).
		Select(
			"description",
			"release_date",
			"country_json",
			"director_json",
			"actors_json",
			"tags_json",
			"imdb_id",
			"runtime",
			"score",
			"episode_count",
			"updated_at",
		).
		UpdateColumns(updates).Error
}

// FindNeedDetailVideos 查找需要补充详情的视频
// 条件：source_id不为空且country_json为空（country_json为空说明详情未获取）。status 为 0 的不处理。
func (r *videoRepository) FindNeedDetailVideos(limit int) ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Where("source_id IS NOT NULL AND source_id != 0 AND (status IS NULL OR status != ?) AND (country_json IS NULL OR country_json = '[]')", "0").
		Limit(limit).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindNeedDetailVideosByType 根据类型查找需要补充详情的视频
// 条件：source_id不为空且release_date和country_json都为空。status 为 0 的不处理。
func (r *videoRepository) FindNeedDetailVideosByType(videoType string, limit int) ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Where("source_id IS NOT NULL AND source_id != 0 AND type = ? AND (status IS NULL OR status != ?) AND (release_date IS NULL) AND (country_json IS NULL OR country_json = '[]')", videoType, "0").
		Limit(limit).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindAllVideos 查找所有视频（返回完整信息）
func (r *videoRepository) FindAllVideos() ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Order("created_at DESC"). // 按创建时间倒序
							Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func applyStatusFilter(query *gorm.DB, statuses []string) *gorm.DB {
	concrete := make([]string, 0, len(statuses))
	includeEmpty := false
	for _, status := range statuses {
		if status == "__empty__" {
			includeEmpty = true
			continue
		}
		concrete = append(concrete, status)
	}
	switch {
	case includeEmpty && len(concrete) > 0:
		return query.Where("(status IN ? OR status IS NULL OR status = '')", concrete)
	case includeEmpty:
		return query.Where("status IS NULL OR status = ''")
	case len(concrete) > 0:
		return query.Where("status IN ?", concrete)
	default:
		return query
	}
}

// FindVideosWithPagination 查找视频列表（带分页、排序和筛选）
func (r *videoRepository) FindVideosWithPagination(limit, offset int, title string, videoTypes, statuses []string, releaseStart, releaseEnd *time.Time, isCompleted *bool, isUpdate *bool) ([]*model.Video, int64, error) {
	var videos []*model.Video
	var total int64

	// 构建查询
	query := database.DB.Model(&model.Video{})

	// 应用筛选条件
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if len(videoTypes) > 0 {
		query = query.Where("type IN ?", videoTypes)
	}

	query = applyStatusFilter(query, statuses)

	if releaseStart != nil {
		query = query.Where("release_date >= ?", releaseStart.Format("2006-01-02"))
	}
	if releaseEnd != nil {
		query = query.Where("release_date <= ?", releaseEnd.Format("2006-01-02"))
	}

	if isCompleted != nil {
		query = query.Where("is_completed = ?", *isCompleted)
	}

	if isUpdate != nil {
		query = query.Where("is_update = ?", *isUpdate)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 按创建时间倒序查询，带分页
	err := query.Order("created_at DESC"). // 按创建时间倒序
						Limit(limit).Offset(offset).
						Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

// FindVideosByStatusNotEqual 查找 status 不等于指定值的视频（返回 id、type、title）
func (r *videoRepository) FindVideosByStatusNotEqual(status string) ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Select("id", "type", "title").
		Where("status != ? OR status IS NULL", status).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosNeedUpdateEpisodes 查找需要搜索播放地址的视频
// 跳过 status=0 和 is_completed=1。未上架（status 为空）仍会搜索。title 必须有内容。
func (r *videoRepository) FindVideosNeedUpdateEpisodes() ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Select("id", "type", "title", "source_id", "episode_count").
		Where("(status IS NULL OR status != ?) AND (is_completed IS NULL OR is_completed != ?) AND title IS NOT NULL AND title != ?", "0", true, "").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosNeedRefreshEpisodeCount 查找需要从豆瓣刷新集数的视频
// 条件：status != 0、is_completed != 1、source_id 有效，且类型不是 movie、tvshow。
func (r *videoRepository) FindVideosNeedRefreshEpisodeCount() ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Select("id", "type", "title", "source_id", "episode_count").
		Where("status != ? AND is_completed != ? AND source_id IS NOT NULL AND source_id != 0 AND type NOT IN ?", "0", true, []string{"movie", "tvshow"}).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// UpdateEpisodeCountWithoutUpdatedAt 只更新 episode_count，不修改 updated_at
func (r *videoRepository) UpdateEpisodeCountWithoutUpdatedAt(videoID int64, episodeCount int64) error {
	return database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumns(map[string]interface{}{
			"episode_count": episodeCount,
		}).Error
}

// UpdateVideosStatusByEpisodes 更新存在 episodes 记录的 videos 的 status。
// status 为 0、episode_count 为 0 或空的不处理。
func (r *videoRepository) UpdateVideosStatusByEpisodes(status string) error {
	// 有剧集、status 不是 0、剧集数不是 0 的视频上架。剧集数为 0 或空时不改 status。
	err := database.DB.Exec(`
		UPDATE videos v
		JOIN (
			SELECT DISTINCT video_id
			FROM episodes
		) e ON v.id = e.video_id
		SET v.status = ?
		WHERE (v.episode_count IS NOT NULL AND v.episode_count != 0)
		  AND (v.status IS NULL OR (v.status != ? AND v.status != '0'))
	`, status, status).Error
	return err
}

// ShelfMoviesWithEpisodes 有剧集的电影改为上架且已完结。status 为 0 或 is_completed 为 1 的不处理。
// 只改 status 和 is_completed，不改 updated_at。
func (r *videoRepository) ShelfMoviesWithEpisodes() (int64, error) {
	result := database.DB.Exec(`
		UPDATE videos v
		INNER JOIN (
			SELECT DISTINCT video_id
			FROM episodes
		) e ON v.id = e.video_id
		SET v.status = '1', v.is_completed = 1
		WHERE v.type = 'movie'
		  AND (v.status IS NULL OR v.status != '0')
		  AND (v.is_completed IS NULL OR v.is_completed != 1)
	`)
	return result.RowsAffected, result.Error
}

// UpdateVideoStatus 更新指定视频的status
// 如果当前值已经是目标值，则不执行数据库更新
// 返回 (是否真的执行了更新, 错误)
func (r *videoRepository) UpdateVideoStatus(videoID int64, status string) (bool, error) {
	// 先查询当前值，如果已经是目标值，则跳过更新
	var currentVideo model.Video
	err := database.DB.Select("status").Where("id = ?", videoID).First(&currentVideo).Error
	if err != nil {
		return false, err
	}

	// 如果当前值已经是目标值，直接返回，不执行更新
	if currentVideo.Status == status {
		return false, nil
	}

	// 只更新 status，避免自动把 updated_at 改成当前时间
	err = database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("status", status).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateVideoIsUpdate 更新指定视频的is_update字段
// 如果当前值已经是目标值，则不执行数据库更新
// 返回 (是否真的执行了更新, 错误)
func (r *videoRepository) UpdateVideoIsUpdate(videoID int64, isUpdate bool) (bool, error) {
	// 先查询当前值，如果已经是目标值，则跳过更新
	var currentVideo model.Video
	err := database.DB.Select("is_update").Where("id = ?", videoID).First(&currentVideo).Error
	if err != nil {
		return false, err
	}

	// 如果当前值已经是目标值，直接返回，不执行更新
	if currentVideo.IsUpdate == isUpdate {
		return false, nil
	}

	// 只更新 is_update，避免自动把 updated_at 改成当前时间
	err = database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("is_update", isUpdate).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateVideoIsCompleted 更新指定视频的is_completed字段
// 如果当前值已经是目标值，则不执行数据库更新
// 返回 (是否真的执行了更新, 错误)
func (r *videoRepository) UpdateVideoIsCompleted(videoID int64, isCompleted bool) (bool, error) {
	// 先查询当前值，如果已经是目标值，则跳过更新
	var currentVideo model.Video
	err := database.DB.Select("is_completed").Where("id = ?", videoID).First(&currentVideo).Error
	if err != nil {
		return false, err
	}

	// 如果当前值已经是目标值，直接返回，不执行更新
	if currentVideo.IsCompleted == isCompleted {
		return false, nil
	}

	// 只更新 is_completed，避免自动把 updated_at 改成当前时间
	err = database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("is_completed", isCompleted).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateVideoUpdatedAt 更新指定视频的updated_at字段
func (r *videoRepository) UpdateVideoUpdatedAt(videoID int64, updatedAt time.Time) error {
	return database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("updated_at", updatedAt).Error
}

// UpdateVideoStateAfterEpisodeInsert 插入剧集后：status=1，is_update=1，updated_at 为最后一集 created_at。
// markCompleted 为真时同时把 is_completed 设为 1。不改其它字段。
func (r *videoRepository) UpdateVideoStateAfterEpisodeInsert(videoID int64, updatedAt time.Time, markCompleted bool) error {
	updates := map[string]interface{}{
		"status":     "1",
		"is_update":  true,
		"updated_at": updatedAt,
	}
	if markCompleted {
		updates["is_completed"] = true
	}
	return database.DB.Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumns(updates).Error
}

// ClearIsUpdateOlderThan 将 updated_at 早于 since、且 status 不为 0 的视频 is_update 设为 0。
func (r *videoRepository) ClearIsUpdateOlderThan(since time.Time) ([]VideoRef, error) {
	rows, err := queryIDTitles(`
		SELECT id, title
		FROM videos
		WHERE (status IS NULL OR status != '0')
		  AND updated_at IS NOT NULL
		  AND updated_at < ?
		  AND (is_update IS NULL OR is_update != 0)
	`, since)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []VideoRef{}, nil
	}
	result := database.DB.Exec(`
		UPDATE videos
		SET is_update = 0
		WHERE (status IS NULL OR status != '0')
		  AND updated_at IS NOT NULL
		  AND updated_at < ?
		  AND (is_update IS NULL OR is_update != 0)
	`, since)
	if result.Error != nil {
		return nil, result.Error
	}
	return toVideoRefs(rows), nil
}

// ApplyInsertedEpisodeVideoState 对本次新插入剧集的视频统一写状态。status 为 0 的不处理。
func (r *videoRepository) ApplyInsertedEpisodeVideoState(videoIDs []int64) ([]VideoRef, error) {
	if len(videoIDs) == 0 {
		return []VideoRef{}, nil
	}
	rows := make([]idTitleRow, 0)
	for start := 0; start < len(videoIDs); start += 500 {
		end := start + 500
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		chunk, err := queryIDTitles(`
			SELECT v.id, v.title
			FROM videos v
			INNER JOIN (
				SELECT video_id
				FROM episodes
				WHERE video_id IN ?
				GROUP BY video_id
			) e ON v.id = e.video_id
			WHERE v.status IS NULL OR v.status != '0'
		`, videoIDs[start:end])
		if err != nil {
			return nil, err
		}
		rows = append(rows, chunk...)
	}
	if len(rows) == 0 {
		return []VideoRef{}, nil
	}
	for start := 0; start < len(videoIDs); start += 500 {
		end := start + 500
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		result := database.DB.Exec(`
			UPDATE videos v
			INNER JOIN (
				SELECT video_id, MAX(created_at) AS last_created_at, COUNT(*) AS episode_total
				FROM episodes
				WHERE video_id IN ?
				GROUP BY video_id
			) e ON v.id = e.video_id
			SET v.status = '1',
				v.is_update = 1,
				v.updated_at = e.last_created_at,
				v.is_completed = CASE
					WHEN v.type = 'movie' THEN 1
					WHEN v.episode_count > 0 AND e.episode_total = v.episode_count THEN 1
					ELSE v.is_completed
				END
			WHERE v.status IS NULL OR v.status != '0'
		`, videoIDs[start:end])
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return toVideoRefs(rows), nil
}

// MarkSeriesCompletedWhenCountMatches 非电影现有集数等于 episode_count 时只更新 is_completed。
func (r *videoRepository) MarkSeriesCompletedWhenCountMatches() ([]VideoRef, error) {
	rows, err := queryIDTitles(`
		SELECT v.id, v.title
		FROM videos v
		INNER JOIN (
			SELECT video_id, COUNT(*) AS episode_total
			FROM episodes
			GROUP BY video_id
		) e ON v.id = e.video_id
		WHERE (v.status IS NULL OR v.status != '0')
		  AND v.type <> 'movie'
		  AND v.episode_count > 0
		  AND e.episode_total = v.episode_count
		  AND (v.is_completed IS NULL OR v.is_completed <> 1)
	`)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []VideoRef{}, nil
	}
	result := database.DB.Exec(`
		UPDATE videos v
		INNER JOIN (
			SELECT video_id, COUNT(*) AS episode_total
			FROM episodes
			GROUP BY video_id
		) e ON v.id = e.video_id
		SET v.is_completed = 1
		WHERE (v.status IS NULL OR v.status != '0')
		  AND v.type <> 'movie'
		  AND v.episode_count > 0
		  AND e.episode_total = v.episode_count
		  AND (v.is_completed IS NULL OR v.is_completed <> 1)
	`)
	if result.Error != nil {
		return nil, result.Error
	}
	return toVideoRefs(rows), nil
}

// ShelfPendingMoviesWithEpisodes 已有剧集、未下架且未完结的电影上架并标完结。
func (r *videoRepository) ShelfPendingMoviesWithEpisodes() ([]VideoRef, error) {
	rows, err := queryIDTitles(`
		SELECT v.id, v.title
		FROM videos v
		INNER JOIN (
			SELECT video_id
			FROM episodes
			GROUP BY video_id
		) e ON v.id = e.video_id
		WHERE v.type = 'movie'
		  AND (v.status IS NULL OR v.status != '0')
		  AND (v.is_completed IS NULL OR v.is_completed <> 1)
	`)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []VideoRef{}, nil
	}
	result := database.DB.Exec(`
		UPDATE videos v
		INNER JOIN (
			SELECT video_id, MAX(created_at) AS last_created_at
			FROM episodes
			GROUP BY video_id
		) e ON v.id = e.video_id
		SET v.status = '1',
			v.is_completed = 1,
			v.is_update = 1,
			v.updated_at = e.last_created_at
		WHERE v.type = 'movie'
		  AND (v.status IS NULL OR v.status != '0')
		  AND (v.is_completed IS NULL OR v.is_completed <> 1)
	`)
	if result.Error != nil {
		return nil, result.Error
	}
	return toVideoRefs(rows), nil
}

// AlignUpdatedAtWithLastEpisode 将有剧集的视频 updated_at 对齐为该剧最后一条剧集的 created_at。
// 已经等于最后一条 created_at 的记录不更新。status 为 0 的不处理。
func (r *videoRepository) AlignUpdatedAtWithLastEpisode() (int64, error) {
	result := database.DB.Exec(`
		UPDATE videos v
		INNER JOIN (
			SELECT video_id, MAX(created_at) AS last_created_at
			FROM episodes
			GROUP BY video_id
		) last_ep ON v.id = last_ep.video_id
		SET v.updated_at = last_ep.last_created_at
		WHERE (v.status IS NULL OR v.status != '0')
		  AND (v.updated_at IS NULL OR v.updated_at <> last_ep.last_created_at)
	`)
	return result.RowsAffected, result.Error
}

// FindByID 根据ID查找视频（返回完整信息）
func (r *videoRepository) FindByID(videoID int64) (*model.Video, error) {
	var video model.Video
	err := database.DB.Where("id = ?", videoID).First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// FindVideosByIsUpdate 查找 is_update = 1 的视频（返回 id）
func (r *videoRepository) FindVideosByIsUpdate() ([]*model.Video, error) {
	var videos []*model.Video
	err := database.DB.Select("id").
		Where("is_update = ?", true).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosWithEpisodes 查找所有有episode的视频（返回 id）
func (r *videoRepository) FindVideosWithEpisodes() ([]*model.Video, error) {
	var videos []*model.Video
	// 使用 DISTINCT 查询所有有episode的视频ID
	err := database.DB.Table("videos").
		Select("DISTINCT videos.id").
		Joins("INNER JOIN episodes ON videos.id = episodes.video_id").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosWithEpisodesAndStatus1 查找 status=1 且有episode的视频（返回 id）
func (r *videoRepository) FindVideosWithEpisodesAndStatus1() ([]*model.Video, error) {
	var videos []*model.Video
	// 使用 DISTINCT 查询 status=1 且有episode的视频ID
	err := database.DB.Table("videos").
		Select("DISTINCT videos.id").
		Joins("INNER JOIN episodes ON videos.id = episodes.video_id").
		Where("videos.status = ?", "1").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosWithEpisodesStatus1AndRecentRelease 查找 status=1、release_date在最近两个月内且有episode的视频（返回 id）
func (r *videoRepository) FindVideosWithEpisodesStatus1AndRecentRelease() ([]*model.Video, error) {
	var videos []*model.Video
	// 计算两个月前的日期（只取日期部分，忽略时间）
	now := time.Now()
	twoMonthsAgo := time.Date(now.Year(), now.Month()-2, now.Day(), 0, 0, 0, 0, now.Location())

	// 使用 DISTINCT 查询 status=1、release_date在最近两个月内且有episode的视频ID
	// release_date IS NOT NULL 确保只处理有上映日期的视频
	err := database.DB.Table("videos").
		Select("DISTINCT videos.id").
		Joins("INNER JOIN episodes ON videos.id = episodes.video_id").
		Where("videos.status = ?", "1").
		Where("videos.release_date IS NOT NULL").
		Where("videos.release_date >= ?", twoMonthsAgo).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosStatus1AndRecentRelease 查找 status=1、release_date在最近两个月内的视频（返回 id）
func (r *videoRepository) FindVideosStatus1AndRecentRelease() ([]*model.Video, error) {
	var videos []*model.Video
	// 计算两个月前的日期（只取日期部分，忽略时间）
	now := time.Now()
	twoMonthsAgo := time.Date(now.Year(), now.Month()-2, now.Day(), 0, 0, 0, 0, now.Location())

	// 查询 status=1、release_date在最近两个月内的视频ID
	// release_date IS NOT NULL 确保只处理有上映日期的视频
	err := database.DB.Select("id").
		Where("status = ?", "1").
		Where("release_date IS NOT NULL").
		Where("release_date >= ?", twoMonthsAgo).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosStatus1AndRecentCreated 查找 status=1、created_at在最近两个月内的视频（返回 id）
func (r *videoRepository) FindVideosStatus1AndRecentCreated() ([]*model.Video, error) {
	var videos []*model.Video
	// 计算两个月前的时间
	twoMonthsAgo := time.Now().AddDate(0, -2, 0)

	// 查询 status=1、created_at在最近两个月内的视频ID
	err := database.DB.Select("id").
		Where("status = ?", "1").
		Where("created_at >= ?", twoMonthsAgo).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// FindVideosStatus1OrIsUpdate1AndRecentCreated 查找需要维护 is_update 的视频（返回 id）
// status 为 0 的不处理。其余满足：已上架或已是有更新，未完结或已是有更新，并且已是有更新或两个月内有更新。
func (r *videoRepository) FindVideosStatus1OrIsUpdate1AndRecentCreated() ([]*model.Video, error) {
	var videos []*model.Video
	twoMonthsAgo := time.Now().AddDate(0, -2, 0)

	err := database.DB.Select("id").
		Where("(status IS NULL OR status != ?)", "0").
		Where("(status = ? OR is_update = ?)", "1", true).
		Where("(is_completed != ? OR is_update = ?)", true, true).
		Where("(is_update = ? OR updated_at >= ?)", true, twoMonthsAgo).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// DeleteByID 根据ID删除视频
func (r *videoRepository) DeleteByID(videoID int64) error {
	return database.DB.Delete(&model.Video{}, videoID).Error
}
