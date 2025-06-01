package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @description: 获取所有歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /songs [get]
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

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_songs"},
			fmt.Errorf("查询总计失败：%w", err),
		))
		return
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_songs"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  songs,
		Total: count,
	})
}

/**
 * @description: 创建歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /songs [post]
func CreateSong(c *gin.Context) {
	// 绑定请求数据
	var createReq models.SongForm

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig("song"))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, createReq.AlbumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorNotFound,
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
			utils.ErrorDatabaseCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_song"},
			fmt.Errorf("song创建失败：%w", err),
		))
		return
	}

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newSong,
	)
}

/**
 * @description: 获取单个歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /songs/:id [get]
func GetSongByID(c *gin.Context) {
	id := c.Param("id")

	var song models.Song
	if err := configs.DB.First(&song, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorNotFound,
			http.StatusNotFound,
			gin.H{"resource": "song"},
			fmt.Errorf("歌曲不存在：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, song)
}

/**
 * @description: 更新歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /songs/:id [put]
func UpdateSong(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 绑定请求数据
	var updateReq models.SongForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig("song"))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 3. 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, updateReq.AlbumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorNotFound,
			http.StatusNotFound,
			gin.H{"resource": "album"},
			fmt.Errorf("专辑不存在：%w", err),
		))
		return
	}

	// 4. 更新数据库
	if err := configs.DB.Model(&models.Song{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_song"},
			fmt.Errorf("song更新失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Song updated successfully"})
}

/**
 * @description: 删除歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /songs/:id [delete]
func DeleteSong(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 删除数据库记录
	if err := configs.DB.Delete(&models.Song{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_song"},
			fmt.Errorf("song删除失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Song deleted successfully"})
}

/**
 * @description: 获取专辑下的所有歌曲
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums/:id/songs [get]
func GetSongsByAlbumID(c *gin.Context) {
	// 1. 获取专辑 ID
	albumID := c.Param("id")

	// 2. 验证专辑是否存在
	var album models.Album
	if err := configs.DB.First(&album, albumID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorNotFound,
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
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_album_songs"},
			fmt.Errorf("查询总计失败：%w", err),
		))
		return
	}

	// 获取分页数据
	if err := baseQuery.Order("track_number").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
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
