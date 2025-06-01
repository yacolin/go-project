package utils

// ValidationConfigs 存储所有表单的验证配置
var ValidationConfigs = struct {
	Album    *ValidationConfig
	User     *ValidationConfig
	Login    *ValidationConfig
	Register *ValidationConfig
	Refresh  *ValidationConfig
	Book     *ValidationConfig
	Photo    *ValidationConfig
	Comment  *ValidationConfig
	Song     *ValidationConfig
}{
	Album: NewValidationConfig().SetFieldMap(map[string]string{
		"Name":        "名称",
		"Author":      "作者",
		"Description": "描述",
		"Liked":       "点赞数",
	}),
	User: NewValidationConfig().SetFieldMap(map[string]string{
		"Username": "用户名",
		"Password": "密码",
		"Email":    "邮箱",
		"Role":     "角色",
	}),
	Login: NewValidationConfig().SetFieldMap(map[string]string{
		"Username": "用户名",
		"Password": "密码",
	}),
	Register: NewValidationConfig().SetFieldMap(map[string]string{
		"Username": "用户名",
		"Password": "密码",
		"Email":    "邮箱",
	}),
	Refresh: NewValidationConfig().SetFieldMap(map[string]string{
		"RefreshToken": "刷新令牌",
	}),
	Book: NewValidationConfig().SetFieldMap(map[string]string{
		"ISBN":        "ISBN",
		"Title":       "书名",
		"Author":      "作者",
		"Stock":       "库存",
		"Publisher":   "出版社",
		"PublishDate": "出版日期",
	}),
	Photo: NewValidationConfig().SetFieldMap(map[string]string{
		"Title":       "标题",
		"URL":         "图片地址",
		"Description": "描述",
		"AlbumID":     "专辑ID",
	}),
	Comment: NewValidationConfig().SetFieldMap(map[string]string{
		"Content": "评论内容",
		"Author":  "评论作者",
		"PhotoID": "照片ID",
	}),
	Song: NewValidationConfig().SetFieldMap(map[string]string{
		"Title":       "歌曲名称",
		"Duration":    "时长",
		"TrackNumber": "曲目编号",
		"AlbumID":     "专辑ID",
	}),
}

// GetValidationConfig 根据表单类型获取对应的验证配置
func GetValidationConfig(formType string) *ValidationConfig {
	switch formType {
	case "album":
		return ValidationConfigs.Album
	case "user":
		return ValidationConfigs.User
	case "login":
		return ValidationConfigs.Login
	case "register":
		return ValidationConfigs.Register
	case "refresh":
		return ValidationConfigs.Refresh
	case "book":
		return ValidationConfigs.Book
	case "photo":
		return ValidationConfigs.Photo
	case "comment":
		return ValidationConfigs.Comment
	case "song":
		return ValidationConfigs.Song
	default:
		return nil
	}
}
