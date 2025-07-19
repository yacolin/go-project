package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize    = 5 << 20 // 5MB
	UploadBasePath = "./static/uploads/"
)

// @Summary 获取所有文件
// @Description 获取所有文件，支持分页
// @Tags files
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ListResponse
// @Failure 500 {object} utils.BusinessError
// @Router /files [get]
func GetAllFiles(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return
	}

	// 2. 数据库操作
	var (
		files []models.File
		count int64
	)

	baseQuery := configs.DB.Model(&models.File{})

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_files"},
			fmt.Errorf("查询总计失败：%w", err),
		))
		return
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_files"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  files,
		Total: count,
	})
}

// @Summary 上传文件
// @Description 上传图片文件到OSS
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /files [post]
func UploadFile(c *gin.Context) {
	// 1. 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"field": "file"},
			fmt.Errorf("文件上传失败：%w", err),
		))
		return
	}

	// 2. 验证文件类型
	if !isImage(file) {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": "file_type"},
			fmt.Errorf("仅支持图片文件"),
		))
		return
	}

	// 3. 验证文件大小
	if file.Size > MaxFileSize {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": "file_size"},
			fmt.Errorf("文件大小超过5MB限制"),
		))
		return
	}

	// 4. 生成唯一文件名
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

	// 5. 打开文件
	src, err := file.Open()
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"operation": "open_file"},
			fmt.Errorf("打开文件失败：%w", err),
		))
		return
	}
	defer src.Close()

	// 6. 上传到OSS
	bucket := configs.GetOSSBucket()
	if bucket == nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"operation": "get_oss_bucket"},
			fmt.Errorf("获取OSS Bucket失败"),
		))
		return
	}

	ossPath := fmt.Sprintf("uploads/%s", newFileName)
	err = bucket.PutObject(ossPath, src)
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"operation": "upload_to_oss"},
			fmt.Errorf("上传到OSS失败：%w", err),
		))
		return
	}

	// 7. 获取文件的OSS URL
	ossURL := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, bucket.Client.Config.Endpoint, ossPath)

	// 8. 保存到数据库
	fileRecord := models.File{
		FileName: newFileName,
		FilePath: ossPath,
		OssURL:   ossURL,
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}

	if err := configs.DB.Create(&fileRecord).Error; err != nil {
		// 删除已上传的OSS文件
		bucket.DeleteObject(ossPath)
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_file_record"},
			fmt.Errorf("保存文件记录失败：%w", err),
		))
		return
	}

	// 9. 返回成功响应
	utils.Success(c, http.StatusCreated, utils.Created, gin.H{
		"id":      fileRecord.ID,
		"oss_url": ossURL,
	})
}

// @Summary 获取单个文件信息
// @Description 根据ID获取文件详情
// @Tags files
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} models.File
// @Failure 404 {object} utils.BusinessError
// @Router /files/{id} [get]
func GetFileByID(c *gin.Context) {
	id := c.Param("id")

	var file models.File
	if err := configs.DB.First(&file, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": "file"},
			fmt.Errorf("文件不存在：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, file)
}

// @Summary 删除文件
// @Description 根据ID删除文件
// @Tags files
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /files/{id} [delete]
func DeleteFile(c *gin.Context) {
	id := c.Param("id")

	// 1. 查找文件记录
	var file models.File
	if err := configs.DB.First(&file, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.NotFound,
			http.StatusNotFound,
			gin.H{"resource": "file"},
			fmt.Errorf("文件不存在：%w", err),
		))
		return
	}

	// 2. 从OSS删除文件
	bucket := configs.GetOSSBucket()
	if bucket != nil {
		err := bucket.DeleteObject(file.FilePath)
		if err != nil {
			c.Error(utils.NewBusinessError(
				utils.ErrInternal,
				http.StatusInternalServerError,
				gin.H{"operation": "delete_oss_file"},
				fmt.Errorf("从OSS删除文件失败：%w", err),
			))
			return
		}
	}

	// 3. 删除数据库记录
	if err := configs.DB.Delete(&file).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_file_record"},
			fmt.Errorf("删除文件记录失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "File deleted successfully"})
}

// 验证是否为图片文件
func isImage(file *multipart.FileHeader) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	buff := make([]byte, 512)
	if _, err = io.ReadFull(src, buff); err != nil {
		return false
	}

	mimeType := http.DetectContentType(buff)
	switch mimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}
