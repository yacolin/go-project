package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @router /comments [post]
func CreateComment(c *gin.Context) {
	var createReq models.CommentForm
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误: %w", err),
		))
		return
	}

	// 检查photo是否存在
	var photo models.Photo
	if err := configs.DB.First(&photo, createReq.PhotoID).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorNotFound,
			http.StatusNotFound,
			gin.H{"resource": "photo"},
			fmt.Errorf("photo不存在: %w", err),
		))
		return
	}

	newComment := models.Comment{
		PhotoID: createReq.PhotoID,
		Content: createReq.Content,
		Author:  createReq.Author,
	}
	if err := configs.DB.Create(&newComment).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_comment"},
			fmt.Errorf("comment创建失败: %w", err),
		))
		return
	}
	utils.Success(c, http.StatusCreated, utils.Created, newComment)
}

// @router /comments/:id [put]
func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var updateReq models.CommentForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorBadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err)},
			fmt.Errorf("参数错误: %w", err),
		))
		return
	}
	if err := configs.DB.Model(&models.Comment{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_comment"},
			fmt.Errorf("comment更新失败: %w", err),
		))
		return
	}
	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Comment updated successfully"})
}

// @router /comments/:id [delete]
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := configs.DB.Delete(&models.Comment{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_comment"},
			fmt.Errorf("comment删除失败: %w", err),
		))
		return
	}
	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Comment deleted successfully"})
}

// @router /photos/:id/comments [get]
func GetCommentsByPhotoID(c *gin.Context) {
	photoID := c.Param("id")
	var comments []models.Comment
	if err := configs.DB.Where("photo_id = ?", photoID).Order("created_at desc").Find(&comments).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_comments"},
			fmt.Errorf("查询评论失败: %w", err),
		))
		return
	}
	utils.Success(c, http.StatusOK, utils.OK, comments)
}
