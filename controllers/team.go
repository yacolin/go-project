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
	// TeamCache 专辑相关的缓存键
	TeamCache = utils.NewCacheKeys(constants.TEAM)
)

// @Summary 获取所有团队
// @Description 获取所有团队，支持分页
// @Tags teams
// @Produce json
// @Param limit query int false "每页数量"
// @Param offset query int false "偏移量"
// @Success 200 {object} utils.ApiRes
// @Failure 500 {object} utils.ApiRes
// @Router /teams [get]
func GetAllTeams(c *gin.Context) {
	// 1. 参数解析与校验
	limit, offset, isAbort := utils.GetPaginationQuery(c)
	if isAbort {
		return // 直接终止
	}

	var (
		teams []models.Team
		count int64
	)

	baseQuery := configs.DB.Model(&models.Team{})

	// 获取数据总数缓存
	totalStr, err := configs.GetCache(TeamCache.TotalKey)
	if err == nil {
		count, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {
		if err := baseQuery.Count(&count).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.DBQuery,
				http.StatusInternalServerError,
				gin.H{"operation": "query_teams"},
				fmt.Errorf("查询总计失败：%w", err),
			))
			return
		}

		// 设置总数缓存，过期时间5分钟
		_ = configs.SetCache(TeamCache.TotalKey, fmt.Sprintf("%d", count), utils.DefaultCacheTime)
	}

	// 尝试获取分页数据缓存
	listCacheKey := utils.GenListCacheKey(TeamCache.ListPrefix, limit, offset)
	listCache, err := configs.GetCache(listCacheKey)
	if err == nil {
		if err := json.Unmarshal([]byte(listCache), &teams); err == nil {
			utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
				List:  teams,
				Total: count,
			})
			return
		}
	}

	if err := baseQuery.Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_teams"},
			fmt.Errorf("查询失败：%w", err),
		))
		return
	}

	if listData, err := json.Marshal(teams); err == nil {
		_ = configs.SetCache(listCacheKey, string(listData), utils.DefaultCacheTime)
	}

	utils.Success(c, http.StatusOK, utils.OK, utils.ListResponse{
		List:  teams,
		Total: count,
	})
}

// @Summary 获取单个团队
// @Description 根据ID获取团队详情
// @Tags teams
// @Produce json
// @Param id path int true "团队ID"
// @Success 200 {object} models.Team
// @Failure 404 {object} utils.BusinessError
// @Router /teams/{id} [get]
func GetTeamByID(c *gin.Context) {
	id := c.Param("id")

	var team models.Team
	if err := utils.FindByID(
		c,
		configs.DB,
		id,
		&team,
		utils.QueryOptions{ResourceName: constants.TEAM},
	); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, team)
}

// @Summary 创建团队
// @Description 创建一个新的团队
// @Tags teams
// @Accept json
// @Produce json
// @Param team body models.TeamForm true "团队信息"
// @Success 201 {object} models.Team
// @Failure 400,500 {object} utils.BusinessError
// @Router /teams [post]
func CreateTeam(c *gin.Context) {
	var createReq models.TeamForm
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.TEAM))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	newTeam := models.Team{
		Name: createReq.Name,
		City: createReq.City,
	}

	if err := configs.DB.Create(&newTeam).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"operation": "create_team"},
			fmt.Errorf("team创建失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(TeamCache)

	utils.Success(c, http.StatusCreated, utils.Created, newTeam)
}

// @Summary 更新团队
// @Description 更新团队信息
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "团队ID"
// @Param team body models.TeamForm true "团队信息"
// @Success 200 {object} map[string]string
// @Failure 400,500 {object} utils.BusinessError
// @Router /teams/{id} [put]
func UpdateTeam(c *gin.Context) {
	id := c.Param("id")

	var updateReq models.TeamForm
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig(constants.TEAM))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	if err := configs.DB.Model(&models.Team{}).Where("id = ?", id).Updates(updateReq.ToMap()).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBUpdate,
			http.StatusInternalServerError,
			gin.H{"operation": "update_team"},
			fmt.Errorf("team更新失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(TeamCache)

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"message": "Team updated successfully"})
}

// @Summary 删除团队
// @Description 删除团队
// @Tags teams
// @Produce json
// @Param id path int true "团队ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} utils.BusinessError
// @Router /teams/{id} [delete]
func DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	if err := configs.DB.Delete(&models.Team{}, id).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBDelete,
			http.StatusInternalServerError,
			gin.H{"operation": "delete_team"},
			fmt.Errorf("team删除失败：%w", err),
		))
		return
	}

	// 清除缓存
	utils.ClearListCache(TeamCache)

	utils.Success(c, http.StatusOK, utils.Deleted, gin.H{"message": "Team deleted successfully"})
}
