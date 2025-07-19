package controllers

import (
	"encoding/json"
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// PetCache 宠物相关的缓存键
	PetCache = utils.NewCacheKeys("pet")
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

	// 3.1 获取数据总数
	totalStr, err := configs.GetCache(PetCache.TotalKey)
	if err == nil {
		// 缓存命中，解析总数
		count, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {
		// 如果缓存不存在，则查询数据库
		if err := baseQuery.Count(&count).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.DBQuery,
				http.StatusInternalServerError,
				gin.H{"operation": "query_pets"},
				fmt.Errorf("查询总计失败：%w", err),
			))
			return
		}

		// 设置总数缓存
		_ = configs.SetCache(PetCache.TotalKey, fmt.Sprintf("%d", count), utils.DefaultCacheTime)
	}

	// 3.2 尝试获取分页数据缓存
	listCacheKey := utils.GenListCacheKey(PetCache.ListPrefix, limit, offset)
	listCache, err := configs.GetCache(listCacheKey)
	if err == nil {
		// 缓存命中，解析数据
		if err := json.Unmarshal([]byte(listCache), &pets); err == nil {
			utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
				List:  pets,
				Total: count,
			})
			return
		}
	}

	// 3.3 缓存未命中或解析失败，从数据库查询
	if err := baseQuery.Limit(limit).Offset(offset).Find(&pets).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_pets"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}
	// 3.4 设置分页数据缓存
	if listData, err := json.Marshal(pets); err == nil {
		_ = configs.SetCache(listCacheKey, string(listData), utils.DefaultCacheTime)
	}

	// 4. 返回结果
	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  pets,
		Total: count,
	})
}
