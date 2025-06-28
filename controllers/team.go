package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// 2. 数据库操作
	var (
		teams []models.Team
		count int64
	)

	baseQuery := configs.DB.Model(&models.Team{})

	if err := baseQuery.Count(&count).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_teams"},
			fmt.Errorf("查询总计失败：%w", err),
		))
	}

	if err := baseQuery.Limit(limit).Offset(offset).Find(&teams).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "query_teams"},
			fmt.Errorf("查询失败：%w", err),
		))
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
		utils.QueryOptions{ResourceName: "team"},
	); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, team)
}
