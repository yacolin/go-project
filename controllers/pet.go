package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取所有宠物
// @Description 获取所有宠物，支持分页
// @Tags pets
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ApiRes
// @Failure 500 {object} utils.ApiRes
// @Router /pets [get]
func GetAllPets(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return // 直接终止
	}

	// 2. 数据库操作
	var (
		pets  []models.Pet
		count int64
	)

	baseQuery := configs.DB.Model(&models.Pet{})

	// 获取数据总数
	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_pets"},
			fmt.Errorf("查询总计失败：%w", err),
		))
	}

	// 获取分页数据
	if err := baseQuery.Limit(limit).Offset(offset).Find(&pets).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_pets"},
			fmt.Errorf("查询失败：%w", err),
		))
	}

	// 3. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  pets,
		Total: count,
	})
}
