package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取所有专辑信息
// @Description 获取所有专辑信息，支持分页
// @Tags albums
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ApiRes
// @Failure 500 {object} utils.ApiRes
// @Router /albums [get]
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
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_albums"},
			fmt.Errorf("查询总计失败：%w", err),
		))
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&albums).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
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

// @Summary 创建专辑信息
// @Description 创建一个新的专辑
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.AlbumForm true "专辑信息"
// @Success 201 {object} models.Album
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /albums [post]
func CreateAlbum(c *gin.Context) {
	// 绑定请求数据
	var createReq models.AlbumForm

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig("album"))},
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
			utils.DBCreate,
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

// @Summary 获取单个专辑信息
// @Description 根据ID获取专辑详情
// @Tags albums
// @Produce json
// @Param id path int true "专辑ID"
// @Success 200 {object} models.Album
// @Failure 404 {object} utils.BusinessError
// @Router /albums/{id} [get]
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

// @Summary 更新专辑信息
// @Description 根据ID更新专辑信息
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "专辑ID"
// @Param album body models.AlbumForm true "专辑信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /albums/{id} [put]
func UpdateAlbum(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 绑定请求数据
	var updateReq models.AlbumForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig("album"))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 3. 更新数据库
	if err := configs.DB.Model(&models.Album{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_album"},
			fmt.Errorf("album更新失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Album updated successfully"})
}

// @Summary 删除专辑信息
// @Description 根据ID删除专辑
// @Tags albums
// @Produce json
// @Param id path int true "专辑ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} utils.BusinessError
// @Router /albums/{id} [delete]
func DeleteAlbum(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 删除数据库记录
	if err := configs.DB.Delete(&models.Album{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_album"},
			fmt.Errorf("album删除失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Album deleted successfully"})
}

// @Summary 搜索专辑信息
// @Description 根据作者名称模糊搜索专辑
// @Tags albums
// @Produce json
// @Param author query string true "作者名称"
// @Success 200 {array} models.Album
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /albums/search [get]
func SearchAlbums(c *gin.Context) {
	// 1. 获取查询参数
	query := c.Query("author")
	if query == "" {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
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
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "search_albums"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, albums)
}
