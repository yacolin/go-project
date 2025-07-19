package controllers

import (
	"encoding/json"
	"fmt"
	"go-project/configs"
	"go-project/constants"
	"go-project/models"
	"go-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// SongCache 歌曲相关的缓存键
	SongCache = utils.NewCacheKeys(constants.SONG)
)

// @Summary 获取所有歌曲
// @Description 获取所有歌曲，支持分页
// @Tags songs
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ListResponse
// @Failure 500 {object} utils.BusinessError
// @Router /songs [get]
func GetAllSongs(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return
	}

	// 2. 数据库操作
	var (
		songs []models.Song
		count int64
	)

	baseQuery := configs.DB.Model(&models.Song{})

	// 3.1 获取数据总数缓存
	totalStr, err := configs.GetCache(SongCache.TotalKey)
	if err == nil {
		// 缓存命中，解析总数
		count, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {

		// 获取数据总数
		if err := baseQuery.Count(&count).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.DBQuery,
				http.StatusInternalServerError,
				gin.H{"operation": "query_songs"},
				fmt.Errorf("查询总计失败：%w", err),
			))
			return
		}

		// 设置总数缓存，过期时间5分钟
		_ = configs.SetCache(SongCache.TotalKey, fmt.Sprintf("%d", count), utils.DefaultCacheTime)
	}

	// 3.2 尝试获取分页数据缓存
	listCacheKey := utils.GenListCacheKey(SongCache.ListPrefix, limit, offset)
	listCache, err := configs.GetCache(listCacheKey)
	if err == nil {
		// 缓存命中，解析数据
		if err := json.Unmarshal([]byte(listCache), &songs); err == nil {
			utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
				List:  songs,
				Total: count,
			})
			return
		}
	}

	// 3.3 缓存未命中或解析失败，从数据库查询
	if err := baseQuery.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_songs"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3.4 设置分页数据缓存，过期时间5分钟
	if listData, err := json.Marshal(songs); err == nil {
		_ = configs.SetCache(listCacheKey, string(listData), utils.DefaultCacheTime)
	}

	// 4. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  songs,
		Total: count,
	})
}

// @Summary 创建歌曲
// @Description 创建一首新歌
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongForm true "歌曲信息"
// @Success 201 {object} models.Song
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /songs [post]
func CreateSong(c *gin.Context) {
	// 绑定请求数据
	var createReq models.SongForm

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.SONG))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, createReq.AlbumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": "album"},
			fmt.Errorf("专辑不存在：%w", err),
		))
		return
	}

	// 创建新的 Song 实例
	newSong := models.Song{
		Title:       createReq.Title,
		Duration:    createReq.Duration,
		TrackNumber: createReq.TrackNumber,
		AlbumID:     createReq.AlbumID,
	}

	// 写入数据库
	if err := configs.DB.Create(&newSong).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_song"},
			fmt.Errorf("song创建失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(SongCache)

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newSong,
	)
}

// @Summary 获取单个歌曲
// @Description 根据ID获取歌曲详情
// @Tags songs
// @Produce json
// @Param id path int true "歌曲ID"
// @Success 200 {object} models.Song
// @Failure 404 {object} utils.BusinessError
// @Router /songs/{id} [get]
func GetSongByID(c *gin.Context) {
	id := c.Param("id")

	var song models.Song
	if err := configs.DB.First(&song, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": constants.SONG},
			fmt.Errorf("歌曲不存在：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, song)
}

// @Summary 更新歌曲
// @Description 根据ID更新歌曲信息
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "歌曲ID"
// @Param song body models.SongForm true "歌曲信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 绑定请求数据
	var updateReq models.SongForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.SONG))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 3. 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, updateReq.AlbumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": "album"},
			fmt.Errorf("专辑不存在：%w", err),
		))
		return
	}

	// 4. 更新数据库
	if err := configs.DB.Model(&models.Song{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_song"},
			fmt.Errorf("song更新失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(SongCache)

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Song updated successfully"})
}

// @Summary 删除歌曲
// @Description 根据ID删除歌曲
// @Tags songs
// @Produce json
// @Param id path int true "歌曲ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} utils.BusinessError
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 删除数据库记录
	if err := configs.DB.Delete(&models.Song{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_song"},
			fmt.Errorf("song删除失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(SongCache)

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Song deleted successfully"})
}

// @Summary 获取专辑下的所有歌曲
// @Description 根据专辑ID获取该专辑下的所有歌曲，支持分页
// @Tags songs
// @Produce json
// @Param id path int true "专辑ID"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ListResponse
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /albums/{id}/songs [get]
func GetSongsByAlbumID(c *gin.Context) {
	// 1. 获取专辑 ID
	albumID := c.Param("id")

	// 2. 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, albumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": "album"},
			fmt.Errorf("专辑不存在：%w", err),
		))
		return
	}

	// 3. 获取分页参数
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return
	}

	// 4. 查询歌曲
	var (
		songs []models.Song
		count int64
	)

	baseQuery := configs.DB.Model(&models.Song{}).Where("album_id = ?", albumID)

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_album_songs"},
			fmt.Errorf("查询总计失败：%w", err),
		))
		return
	}

	// 获取分页数据
	if err := baseQuery.Order("track_number").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_album_songs"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 5. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  songs,
		Total: count,
	})
}
