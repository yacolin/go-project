package middlewares

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()

	// 注册自定义正则校验器
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		pattern := fl.Param() // 获取正则表达式
		value := fl.Field().String()
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	})
}

func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化自定义校验器
		InitValidator()

		// 继续处理请求
		c.Next()
	}
}
