package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 统一响应结构体
type ApiRes struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func Success(c *gin.Context, httpStatus, bizCode int, data interface{}) {
	c.JSON(httpStatus, ApiRes{
		Code:      0,
		Message:   CodeMessages[bizCode],
		Data:      data,
		Errors:    nil,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func ErrorResponse(c *gin.Context, httpStatus, bizCode int, Errors interface{}) {
	c.JSON(httpStatus, ApiRes{
		Code:      bizCode,
		Message:   CodeMessages[bizCode],
		Data:      nil,
		Errors:    Errors,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// BusinessError 自定义业务错误
type BusinessError struct {
	BizCode  int
	HttpCode int
	Err      error
	Details  interface{}
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("BizCode%d, HttpCode:%d, Details%v", e.BizCode, e.HttpCode, e.Details)
}

func NewBusinessError(bizCode, httpCode int, details interface{}, err error) *BusinessError {
	return &BusinessError{
		BizCode:  bizCode,
		HttpCode: httpCode,
		Details:  details,
		Err:      err,
	}
}
