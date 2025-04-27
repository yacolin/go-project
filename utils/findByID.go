package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QueryOptions struct {
	ResourceName string
	Preloads     []string
}

func FindByID[T any](
	c *gin.Context,
	db *gorm.DB,
	id string,
	record *T,
	opts QueryOptions,
) error {
	// 参数校验
	if id == "" {
		return NewBusinessError(
			ErrorParamMissingID,
			http.StatusBadRequest,
			gin.H{"field": "id"},
			errors.New(CodeMessages[ErrorParamMissingID]),
		)
	}

	// 构建查询
	query := db.Model(record)
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// 执行查询
	if err := query.First(record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewBusinessError(
				ErrorNotFound,
				http.StatusNotFound,
				gin.H{"resource": opts.ResourceName, "id": id},
				fmt.Errorf("%s不存在", opts.ResourceName),
			)
		}

		return NewBusinessError(
			ErrorDatabaseQuery,
			http.StatusInternalServerError,
			gin.H{"operation": "get_" + opts.ResourceName},
			fmt.Errorf("数据库查询失败：%w", err),
		)
	}
	return nil
}
