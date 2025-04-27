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
