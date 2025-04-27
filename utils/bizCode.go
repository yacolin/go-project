package utils

// ----------------------------
// 成功状态码 (2000-2999)
// ----------------------------
const (
	OK = 2000 + iota
	Created
	NoContent
	Accepted
	Deleted
	Updated
)

// ----------------------------
// 客户端错误 (4000-4999)
// ----------------------------
const (
	// 通用客户端错误 (4000-4009)
	ErrorBadRequest   = 4000 + iota // 无效的请求参数
	ErrorUnauthorized               // 身份验证失败
	_                               // 预留占位符
	ErrorForbidden                  // 没有访问权限
	ErrorNotFound                   // 资源不存在

	// 参数校验错误 (4100-4199)
	ErrorParamInvalidPwd        = 4100 + iota // 密码错误
	ErrorParamInvalidPagination               // 无效的分页参数
	ErrorParamInvalidPrice                    // 无效的价格参数
	ErrorParamMissingID                       // 缺少ID参数
	ErrorParamMissingName                     // 缺少Name参数

	// 数据错误 (4200-4299)
	ErrorDataNotFound = 4200 + iota // 数据不存在

	// 令牌错误 (4300-4399)
	ErrorTokenGenFailed             = 4300 + iota // 生成token失败
	ErrorTokenInvalid                             // 无效的token
	ErrorTokenExpired                             // token已过期
	ErrorTokenNotFound                            // token不存在
	ErrorTokenInvalidFormat                       // token格式错误
	ErrorTokenInvalidSignature                    // token签名错误
	ErrorTokenInvalidClaims                       // token解析错误
	ErrorTokenInvalidClaimsUserID                 // token中缺少user_id
	ErrorTokenInvalidClaimsUserName               // token中缺少user_name
	ErrorTokenInvalidAudience                     // token受众错误

)

// ----------------------------
// 服务器错误 (5000-5999)
// ----------------------------
const (
	// 通用服务器错误 (5000-5009)
	ErrorInternal = 5000 + iota // 服务器内部错误

	// 数据库错误 (5100-5199)
	ErrorDatabaseQuery          = 5100 + iota // 数据查询失败
	ErrorDatabaseCount                        // 数据统计失败
	ErrorDatabaseDelete                       // 数据库删除失败
	ErrorDatabaseUpdate                       // 数据库更新失败
	ErrorDatabaseCreate                       // 数据库创建失败
	ErrorDatabaseDuplicateEntry               // 数据重复冲突
	ErrorUserNotFound                         // 用户不存在
)

// ----------------------------
// 状态码映射表
// ----------------------------
var CodeMessages = map[int]string{
	// 成功状态码
	OK:        "请求成功",
	Created:   "创建成功",
	NoContent: "操作执行成功",
	Accepted:  "请求已被接受",
	Deleted:   "删除成功",
	Updated:   "更新成功",

	// 客户端错误
	ErrorBadRequest:             "无效的请求参数",
	ErrorUnauthorized:           "身份验证失败",
	ErrorForbidden:              "没有访问权限",
	ErrorNotFound:               "资源不存在",
	ErrorParamInvalidPwd:        "密码错误",
	ErrorParamInvalidPrice:      "无效的价格参数",
	ErrorParamInvalidPagination: "无效的分页参数",
	ErrorParamMissingID:         "缺少ID参数",
	ErrorParamMissingName:       "缺少Name参数",
	ErrorDataNotFound:           "数据不存在",

	ErrorTokenGenFailed:             "生成token失败",
	ErrorTokenInvalid:               "无效的token",
	ErrorTokenExpired:               "token已过期",
	ErrorTokenNotFound:              "token不存在",
	ErrorTokenInvalidFormat:         "token格式错误",
	ErrorTokenInvalidSignature:      "token签名错误",
	ErrorTokenInvalidClaims:         "token解析错误",
	ErrorTokenInvalidClaimsUserID:   "token中缺少user_id",
	ErrorTokenInvalidClaimsUserName: "token中缺少user_name",
	ErrorTokenInvalidAudience:       "token受众错误",

	// 服务器错误
	ErrorInternal:               "服务器内部错误",
	ErrorDatabaseQuery:          "数据查询失败",
	ErrorDatabaseCount:          "数据统计失败",
	ErrorDatabaseDelete:         "数据库删除失败",
	ErrorDatabaseUpdate:         "数据库更新失败",
	ErrorDatabaseCreate:         "数据库创建失败",
	ErrorDatabaseDuplicateEntry: "数据重复冲突",
	ErrorUserNotFound:           "用户不存在",
}
