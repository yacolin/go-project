package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				errors[field] = "该字段为必填字段"
			case "min":
				errors[field] = fmt.Sprintf("最小长度/值为 %s", fe.Param())
			case "max":
				errors[field] = fmt.Sprintf("最大长度/值为 %s", fe.Param())
			case "gte":
				errors[field] = fmt.Sprintf("值必须大于或等于 %s", fe.Param())
			case "lte":
				errors[field] = fmt.Sprintf("值必须小于或等于 %s", fe.Param())
			case "email":
				errors[field] = "必须是有效的电子邮件地址"
			case "url":
				errors[field] = "必须是有效的 URL"
			case "len":
				errors[field] = fmt.Sprintf("长度必须为 %s", fe.Param())
			case "numeric":
				errors[field] = "必须是数字"
			case "alpha":
				errors[field] = "只能包含字母"
			case "alphanum":
				errors[field] = "只能包含字母和数字"
			case "uuid":
				errors[field] = "必须是有效的 UUID"
			case "regexp":
				errors[field] = "格式不符合要求" // 自定义正则表达式的错误信息
			default:
				errors[field] = fmt.Sprintf("格式无效（%s）", fe.Tag())
			}
		}
	}
	return errors
}
