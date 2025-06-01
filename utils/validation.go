package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationConfig 用于配置验证错误信息
type ValidationConfig struct {
	FieldMap map[string]string // 字段映射，key为字段名，value为显示名
}

// NewValidationConfig 创建一个新的验证配置
func NewValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		FieldMap: make(map[string]string),
	}
}

// SetFieldMap 设置字段映射
func (c *ValidationConfig) SetFieldMap(fieldMap map[string]string) *ValidationConfig {
	c.FieldMap = fieldMap
	return c
}

// FormatValidationErrors 格式化验证错误信息
// 如果config为nil，将使用原始字段名
func FormatValidationErrors(err error, config *ValidationConfig) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := strings.ToLower(fe.Field())

			// 获取字段显示名
			var fieldName string
			if config != nil {
				fieldName = config.FieldMap[fe.Field()]
			}
			if fieldName == "" {
				fieldName = field
			}

			switch fe.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s为必填字段", fieldName)
			case "min":
				if fe.Type().Kind().String() == "string" {
					errors[field] = fmt.Sprintf("%s的长度不能少于%s个字符", fieldName, fe.Param())
				} else {
					errors[field] = fmt.Sprintf("%s不能小于%s", fieldName, fe.Param())
				}
			case "max":
				if fe.Type().Kind().String() == "string" {
					errors[field] = fmt.Sprintf("%s的长度不能超过%s个字符", fieldName, fe.Param())
				} else {
					errors[field] = fmt.Sprintf("%s不能大于%s", fieldName, fe.Param())
				}
			case "gte":
				errors[field] = fmt.Sprintf("%s必须大于或等于%s", fieldName, fe.Param())
			case "lte":
				errors[field] = fmt.Sprintf("%s必须小于或等于%s", fieldName, fe.Param())
			case "email":
				errors[field] = fmt.Sprintf("%s必须是有效的电子邮件地址", fieldName)
			case "url":
				errors[field] = fmt.Sprintf("%s必须是有效的URL", fieldName)
			case "len":
				errors[field] = fmt.Sprintf("%s的长度必须为%s", fieldName, fe.Param())
			case "numeric":
				errors[field] = fmt.Sprintf("%s必须是数字", fieldName)
			case "alpha":
				errors[field] = fmt.Sprintf("%s只能包含字母", fieldName)
			case "alphanum":
				errors[field] = fmt.Sprintf("%s只能包含字母和数字", fieldName)
			case "uuid":
				errors[field] = fmt.Sprintf("%s必须是有效的UUID", fieldName)
			case "regexp":
				errors[field] = fmt.Sprintf("%s格式不符合要求", fieldName)
			default:
				errors[field] = fmt.Sprintf("%s的格式无效（%s）", fieldName, fe.Tag())
			}
		}
	}
	return errors
}
