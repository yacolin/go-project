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
	BadRequest   = 4000 + iota // 无效的请求参数
	Unauthorized               // 身份验证失败
	_                          // 预留占位符
	Forbidden                  // 没有访问权限
	NotFound                   // 资源不存在

	// 参数校验错误 (4100-4199)
	InvalidPwd   = 4100 + iota // 密码错误
	InvalidPage                // 无效的分页参数
	InvalidPrice               // 无效的价格参数
	MissingID                  // 缺少ID参数
	MissingName                // 缺少Name参数

	// 数据错误 (4200-4299)
	NoData = 4200 + iota // 数据不存在

	// 令牌错误 (4300-4399)
	TkGen        = 4300 + iota // 生成token失败
	AccessTkGen                // 生成AccessToken失败
	RefreshTkGen               // 生成RefreshToken失败
	TkInvalid                  // 无效的token
	TkExpired                  // token已过期
	TkNotFound                 // token不存在
	TkFormat                   // token格式错误
	TkSign                     // token签名错误
	TkClaims                   // token解析错误
	TkUserID                   // token中缺少user_id
	TkUserName                 // token中缺少user_name
	TkAudience                 // token受众错误

)

// ----------------------------
// 服务器错误 (5000-5999)
// ----------------------------
const (
	// 通用服务器错误 (5000-5009)
	ErrInternal = 5000 + iota // 服务器内部错误

	// 数据库错误 (5100-5199)
	DBQuery  = 5100 + iota // 数据查询失败
	DBCount                // 数据统计失败
	DBDelete               // 数据库删除失败
	DBUpdate               // 数据库更新失败
	DBCreate               // 数据库创建失败
	DBDup                  // 数据重复冲突

	UserNotFound // 用户不存在
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
	BadRequest:   "无效的请求参数",
	Unauthorized: "身份验证失败",
	Forbidden:    "没有访问权限",
	NotFound:     "资源不存在",
	InvalidPwd:   "密码错误",
	InvalidPrice: "无效的价格参数",
	InvalidPage:  "无效的分页参数",
	MissingID:    "缺少ID参数",
	MissingName:  "缺少Name参数",
	NoData:       "数据不存在",

	TkGen:        "生成token失败",
	AccessTkGen:  "生成AccessToken失败",
	RefreshTkGen: "生成RefreshToken失败",
	TkInvalid:    "无效的token",
	TkExpired:    "token已过期",
	TkNotFound:   "token不存在",
	TkFormat:     "token格式错误",
	TkSign:       "token签名错误",
	TkClaims:     "token解析错误",
	TkUserID:     "token中缺少user_id",
	TkUserName:   "token中缺少user_name",
	TkAudience:   "token受众错误",

	// 服务器错误
	ErrInternal:  "服务器内部错误",
	DBQuery:      "数据查询失败",
	DBCount:      "数据统计失败",
	DBDelete:     "数据库删除失败",
	DBUpdate:     "数据库更新失败",
	DBCreate:     "数据库创建失败",
	DBDup:        "数据重复冲突",
	UserNotFound: "用户不存在",
}
