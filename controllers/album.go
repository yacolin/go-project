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
 * @description: 获取所有专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums [get]
func GetAllAlbums(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return // 直接终止
	}

	// 2. 数据库操作
	var (
		albums []models.Album
		count  int64
	)

	baseQuery := configs.DB.Model(&models.Album{})

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_albums"},
			fmt.Errorf("查询总计失败：%w", err),
		))
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&albums).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_albums"},
			fmt.Errorf("查询失败：%w", err),
		))
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  albums,
		Total: count,
	})
}

/**
 * @description: 创建专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums [post]
func CreateAlbum(c *gin.Context) {
	// 绑定请求数据
	var createReq models.AlbumForm

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 创建新的 Album 实例
	newAlbum := models.Album{
		Name:        createReq.Name,
		Author:      createReq.Author,
		Description: createReq.Description,
		Liked:       createReq.Liked,
	}

	// 写入数据库
	if err := configs.DB.Create(&newAlbum).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_album"},
			fmt.Errorf("album创建失败：%w", err),
		))
		return
	}

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newAlbum,
	)
}

/**
 * @description: 获取单个专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var album models.Album
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&album,
		utils.QueryOptions{ResourceName: "album"},
	); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, album)
}

/**
 * @description: 更新专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums/:id [put]
func UpdateAlbum(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 绑定请求数据
	var updateReq models.AlbumForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 3. 更新数据库
	if err := configs.DB.Model(&models.Album{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_album"},
			fmt.Errorf("album更新失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Album updated successfully"})
}

/**
 * @description: 删除专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums/:id [delete]
func DeleteAlbum(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 删除数据库记录
	if err := configs.DB.Delete(&models.Album{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_album"},
			fmt.Errorf("album删除失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Album deleted successfully"})
}

/**
 * @description: 搜索专辑信息
 * @param {*gin.Context} c
 * @return {*}
 */
// @router /albums/search [get]
func SearchAlbums(c *gin.Context) {
	// 1. 获取查询参数
	query := c.Query("author")
	if query == "" {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": "query parameter is required"},
			fmt.Errorf("查询参数不能为空"),
		))
		return
	}

	// 2. 查询数据库
	var albums []models.Album
	if err := configs.DB.Where("author LIKE ?", "%"+query+"%").Find(&albums).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "search_albums"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, albums)
}
