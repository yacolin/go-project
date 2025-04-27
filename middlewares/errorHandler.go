package middlewares

import (
	"errors"
	"go-project/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			var bizErr *utils.BusinessError
			if errors.As(err, &bizErr) {
				// 记录原始错误日志
				log.Printf("Business Error: %v\n", bizErr.Err)

				// 返回结构化得错误
				utils.ErrorResponse(c, bizErr.HttpCode, bizErr.BizCode, bizErr.Details)
			} else {
				// 处理未知错误
				log.Printf("System Error: %v\n", err)
				utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorInternal, "系统繁忙，请稍后再试")
			}
		}
	}
}
