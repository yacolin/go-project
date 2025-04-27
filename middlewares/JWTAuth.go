package middlewares

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头拿 Authorization: Bearer tokenxxx
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(utils.NewBusinessError(
				utils.ErrorUnauthorized,
				http.StatusUnauthorized,
				gin.H{"error": utils.CodeMessages[utils.ErrorTokenNotFound]},
				fmt.Errorf("参数错误：%w", nil),
			))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Error(utils.NewBusinessError(
				utils.ErrorUnauthorized,
				http.StatusUnauthorized,
				gin.H{"error": utils.CodeMessages[utils.ErrorTokenInvalidFormat]},
				fmt.Errorf("参数错误：%w", nil),
			))
			c.Abort()
			return
		}

		// 解析token
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return configs.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.Error(utils.NewBusinessError(
				utils.ErrorUnauthorized,
				http.StatusUnauthorized,
				gin.H{"error": utils.CodeMessages[utils.ErrorTokenInvalid]},
				fmt.Errorf("参数错误：%w", nil),
			))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(utils.NewBusinessError(
				utils.ErrorUnauthorized,
				http.StatusUnauthorized,
				gin.H{"error": utils.CodeMessages[utils.ErrorTokenInvalidClaims]},
				fmt.Errorf("参数错误：%w", nil),
			))
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.Error(utils.NewBusinessError(
				utils.ErrorUnauthorized,
				http.StatusUnauthorized,
				gin.H{"error": utils.CodeMessages[utils.ErrorTokenInvalidClaimsUserID]},
				fmt.Errorf("参数错误：%w", nil),
			))
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// 查找用户是否还存在（可选）
		var user models.User
		if err := configs.DB.First(&user, userID).Error; err != nil {
			c.Error(utils.NewBusinessError(
				utils.ErrorNotFound,
				http.StatusUnauthorized,
				// gin.H{"error": "用户不存在"},
				nil,
				fmt.Errorf("参数错误：%w", err),
			))
			c.Abort()
			return
		}

		// 保存当前登录用户到上下文
		c.Set("user", user)

		c.Next()
	}
}
