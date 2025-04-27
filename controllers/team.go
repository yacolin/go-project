package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
