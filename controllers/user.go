package controllers

import (
	"fmt"
	"go-project/configs"
	"go-project/models"
	"go-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary 用户登录
// @Description 用户登录，获取access_token和refresh_token
// @Tags users
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "登录信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.BusinessError
// @Failure 404 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /login [post]
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"validation": utils.FormatValidationErrors(err, utils.GetValidationConfig("login"))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	var user models.User
	if err := configs.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusInternalServerError,
			gin.H{"error": utils.CodeMessages[utils.UserNotFound]},
			fmt.Errorf("用户不存在：%w", err),
		))
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"error": utils.CodeMessages[utils.InvalidPwd]},
			fmt.Errorf("密码错误：%w", nil),
		))
		return
	}

	// 生成Token
	// token, err := utils.GenerateToken(user.ID)
	// if err != nil {
	// 	c.Error(utils.NewBusinessError(
	// 		utils.BadRequest,
	// 		http.StatusBadRequest,
	// 		gin.H{"error": utils.CodeMessages[utils.TkGen]},
	// 		fmt.Errorf("生成token失败：%w", err),
	// 	))
	// 	return
	// }
	// utils.Success(c, http.StatusOK, utils.OK, gin.H{"token": token})

	// 生成accessToken
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"error": utils.CodeMessages[utils.AccessTkGen]},
			fmt.Errorf("生成accessToken失败：%w", err),
		))
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"error": utils.CodeMessages[utils.RefreshTkGen]},
			fmt.Errorf("生成refreshToken失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

// @Summary 用户注册
// @Description 用户注册
// @Tags users
// @Accept json
// @Produce json
// @Param register body models.LoginInput true "注册信息"
// @Success 200 {object} nil
// @Failure 400 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /register [post]
func Register(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"register": utils.FormatValidationErrors(err, utils.GetValidationConfig("register"))},
			fmt.Errorf("参数错误：%w", err),
		))
		return
	}

	// 检查用户是否已存在
	var existUser models.User
	if err := configs.DB.Where("username = ?", input.Username).First(&existUser).Error; err == nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"register": "用户名已存在"},
			fmt.Errorf("用户名已存在：%w", err),
		))
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"register": "密码加密失败"},
			fmt.Errorf("密码加密失败：%w", err),
		))
		return
	}

	// 创建用户
	user := models.User{
		Username: input.Username,
		Password: hashedPassword,
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		c.Error(utils.NewBusinessError(
			utils.DBCreate,
			http.StatusInternalServerError,
			gin.H{"register": "注册失败"},
			fmt.Errorf("注册失败：%w", err),
		))
		return
	}

	utils.Success(c, http.StatusOK, utils.OK, nil)
}

// @Summary 刷新Token
// @Description 使用refresh_token刷新access_token
// @Tags users
// @Accept json
// @Produce json
// @Param refresh body models.RefreshInput true "刷新信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.BusinessError
// @Failure 401 {object} utils.BusinessError
// @Failure 500 {object} utils.BusinessError
// @Router /refresh [post]
func Refresh(c *gin.Context) {
	var input models.RefreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.NewBusinessError(
			utils.BadRequest,
			http.StatusBadRequest,
			gin.H{"refresh": utils.FormatValidationErrors(err, utils.GetValidationConfig("refresh"))},
			fmt.Errorf("参数错误：%w", err),
		))
	}

	// 解析 refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return configs.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.Error(utils.NewBusinessError(
			utils.Unauthorized,
			http.StatusUnauthorized,
			gin.H{"error": utils.CodeMessages[utils.TkInvalid]},
			fmt.Errorf("无效的refresh token：%w", nil),
		))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Error(utils.NewBusinessError(
			utils.Unauthorized,
			http.StatusUnauthorized,
			gin.H{"error": utils.CodeMessages[utils.TkClaims]},
			fmt.Errorf("token解析失败：%w", nil),
		))
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.Error(utils.NewBusinessError(
			utils.Unauthorized,
			http.StatusUnauthorized,
			gin.H{"error": utils.CodeMessages[utils.TkUserID]},
			fmt.Errorf("refresh token缺少user_id：%w", nil),
		))
		return
	}
	userID := uint(userIDFloat)

	// 生成新的 access token
	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.Error(utils.NewBusinessError(
			utils.ErrInternal,
			http.StatusInternalServerError,
			gin.H{"error": utils.CodeMessages[utils.TkGen]},
			fmt.Errorf("生成新的access token失败：%w", err),
		))
		return
	}
	utils.Success(c, http.StatusOK, utils.OK, gin.H{"access_token": newAccessToken})
}
