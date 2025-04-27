package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	// 单文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "File upload failed",
		})
		return
	}

	// 验证文件类型
	if !isImage(file) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Only image files are allowed",
		})
		return
	}

	// 验证文件大小（限制5MB）
	if file.Size > 5<<20 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 5MB limit",
		})
		return
	}

	// 生成唯一文件名
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))

	// 文件保存路径
	uploadPath := "./static/uploads/"
	filePath := filepath.Join(uploadPath, newFileName)

	// 确保目录存在
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	// 保存到数据库
	fileRecord := models.File{
		FileName:  newFileName,
		FilePath:  filePath,
		MimeType:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		CreatedAt: time.Now(),
	}

	if result := configs.DB.Create(&fileRecord); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file record",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fmt.Sprintf("/static/uploads/%s", newFileName),
		"file_id":  fileRecord.ID,
	})
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
