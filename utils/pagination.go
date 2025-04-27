package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// 校验分页
func ValidatePagination(pageStr, sizeStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, errors.New("页码必须大于1")
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		return 0, 0, errors.New("分页搭笑必需在1-100之间")
	}

	return size, (page - 1) * size, nil
}

// 获取分页
func GetPaginationQuery(c *gin.Context) (limit, offset int, isAbort bool) {
	pageStr := c.DefaultQuery("current", "1")
	sizeStr := c.DefaultQuery("pageSize", "10")

	limit, offset, err := ValidatePagination(pageStr, sizeStr)
	if err != nil {
		c.Error(NewBusinessError(
			ErrorParamInvalidPagination,
			http.StatusBadRequest,
			gin.H{"fields": []string{"page", "size"}, "reason": err.Error()},
			err,
		))
		return 0, 0, true // 返回终止标志
	}

	return limit, offset, false
}
