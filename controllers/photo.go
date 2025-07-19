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
	// PhotoCache 照片相关的缓存键
	PhotoCache = utils.NewCacheKeys(constants.PHOTO)
)

// @Summary 获取所有照片
// @Description 获取所有照片，支持分页
// @Tags photos
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ApiRes
// @Failure 500 {object} utils.ApiRes
// @Router /photos [get]
func GetAllPhotos(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return
	}

	// 2. 数据库操作
	var (
		photos []models.Photo
		count  int64
	)

	baseQuery := configs.DB.Model(&models.Photo{})

	// 3.1 获取数据总数缓存
	totalStr, err := configs.GetCache(PhotoCache.TotalKey)
	if err == nil {
		// 缓存命中，解析总数
		count, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {
		// 如果缓存不存在，则查询数据库
		if err := baseQuery.Count(&count).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.DBQuery,
				http.StatusInternalServerError,
				gin.H{"operation": "query_photos"},
				fmt.Errorf("查询总计失败：%w", err),
			))
			return
		}
		// 设置总数缓存，过期时间5分钟
		_ = configs.SetCache(PhotoCache.TotalKey, fmt.Sprintf("%d", count), utils.DefaultCacheTime)

	}

	// 3.2 尝试获取分页数据缓存
	listCacheKey := utils.GenListCacheKey(AlbumCache.ListPrefix, limit, offset)
	listCache, err := configs.GetCache(listCacheKey)
	if err == nil {
		// 缓存命中，解析数据
		if err := json.Unmarshal([]byte(listCache), &photos); err == nil {
			utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
				List:  photos,
				Total: count,
			})
			return
		}
	}

	// 3.3 缓存未命中或解析失败，从数据库查询
	if err := baseQuery.Limit(limit).Offset(offset).Find(&photos).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_photos"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3.4 设置分页数据缓存，过期时间5分钟
	if listData, err := json.Marshal(photos); err == nil {
		_ = configs.SetCache(listCacheKey, string(listData), utils.DefaultCacheTime)
	}

	// 4. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  photos,
		Total: count,
	})
}

// @Summary 创建照片
// @Description 创建一条新的照片记录
// @Tags photos
// @Accept json
// @Produce json
// @Param photo body models.PhotoForm true "照片信息"
// @Success 201 {object} models.Photo
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /photos [post]
func CreatePhoto(c *gin.Context) {
	// 绑定请求数据
	var createReq models.PhotoForm

	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.PHOTO))},
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

	// 创建新的 Photo 实例
	newPhoto := models.Photo{
		Title:       createReq.Title,
		URL:         createReq.URL,
		Description: createReq.Description,
		AlbumID:     createReq.AlbumID,
	}

	// 写入数据库
	if err := configs.DB.Create(&newPhoto).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_photo"},
			fmt.Errorf("photo创建失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(PhotoCache)

	// 返回创建结果
	utils.Success(
		c,
		http.StatusCreated,
		utils.Created,
		newPhoto,
	)
}

// @Summary 获取单个照片
// @Description 根据ID获取照片详情
// @Tags photos
// @Produce json
// @Param id path int true "照片ID"
// @Success 200 {object} models.Photo
// @Failure 404 {object} utils.BusinessError
// @Router /photos/{id} [get]
func GetPhotoByID(c *gin.Context) {
	id := c.Param("id")

	var photo models.Photo
	if err := configs.DB.First(&photo, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": constants.PHOTO},
			fmt.Errorf("照片不存在：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, photo)
}

// @Summary 更新照片
// @Description 根据ID更新照片信息
// @Tags photos
// @Accept json
// @Produce json
// @Param id path int true "照片ID"
// @Param photo body models.PhotoForm true "照片信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /photos/{id} [put]
func UpdatePhoto(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 绑定请求数据
	var updateReq models.PhotoForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.PHOTO))},
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
	if err := configs.DB.Model(&models.Photo{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_photo"},
			fmt.Errorf("photo更新失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(PhotoCache)

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Photo updated successfully"})
}

// @Summary 删除照片
// @Description 根据ID删除照片
// @Tags photos
// @Produce json
// @Param id path int true "照片ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} utils.BusinessError
// @Router /photos/{id} [delete]
func DeletePhoto(c *gin.Context) {
	// 1. 获取 ID 参数
	id := c.Param("id")

	// 2. 删除数据库记录
	if err := configs.DB.Delete(&models.Photo{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_photo"},
			fmt.Errorf("photo删除失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(PhotoCache)

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Photo deleted successfully"})
}

// @Summary 获取专辑下的所有照片
// @Description 根据专辑ID获取该专辑下的所有照片，支持分页
// @Tags photos
// @Produce json
// @Param id path int true "专辑ID"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ListResponse
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /albums/{id}/photos [get]
func GetPhotosByAlbumID(c *gin.Context) {
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

	// 4. 查询照片
	var (
		photos []models.Photo
		count  int64
	)

	baseQuery := configs.DB.Model(&models.Photo{}).Where("album_id = ?", albumID)

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_album_photos"},
			fmt.Errorf("查询总计失败：%w", err),
		))
		return
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&photos).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_album_photos"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 5. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  photos,
		Total: count,
	})
}

// @Summary 获取照片的所有评论
// @Description 根据照片ID获取所有评论，支持分页
// @Tags comments
// @Produce json
// @Param id path int true "照片ID"
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ListResponse
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /photos/{id}/comments [get]
func GetPhotoComments(c *gin.Context) {
	photoID := c.Param("id")

	// 验证照片是否存在
	var photo models.Photo
	if err := configs.DB.First(&photo, photoID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": constants.PHOTO},
			fmt.Errorf("照片不存在：%w", err),
		))
		return
	}

	// 分页参数
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return
	}

	// 查询评论
	var (
		comments []models.Comment
		count    int64
	)
	baseQuery := configs.DB.Model(&models.Comment{}).Where("photo_id = ?", photoID)

	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_photo_comments"},
			fmt.Errorf("查询评论总数失败：%w", err),
		))
		return
	}

	if err := baseQuery.Limit(limit).Offset(offset).Find(&comments).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_photo_comments"},
			fmt.Errorf("查询评论失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  comments,
		Total: count,
	})
}
