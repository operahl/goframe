package conf

//所有的参数错误以400开头，资源位找到404开头，内部错误以500开头
const (
	CodeOk                  = 0  // seccess
	CodeAesErr              = 1  // aes error1
	CodeParaErr             = 3  // param error
	SystemError             = 5  // 服务器错误
	CodeEmailIllegal        = 6  // 邮箱格式不合法
	UsernameOrPasswordError = 7  // 用户名或者密码错误
	CodeSendMailError       = 8  // 邮件发送失败
	CodeVerifyCodeError     = 10 // verifycode error
	VideoNumExceed          = 11 // 视频个数超出限制
	CodeEmailRegistered     = 13 // 邮箱已注册
	CodeEmailNotRegister    = 14 // 邮箱未注册
	UsernameIrregular       = 24 // 用户名包含特殊字符
	StatusError             = 27 //
	FbLoginCheckFailure     = 28 //fb 登录，去fb验证错误，服务器发送http请求错误，客户端应该重试
	FbuidMatchError         = 29 //fb 登录，去fb验证了，服务器获取到的fbid和客户端传的不一样
	MomentUpNoDate			= 36
	MomentDownNoDate		= 37
	MylikeNoData            = 37 //我喜欢的列表为空
	Nodata                  = 40 //没有数据
	BlockUser               = 44 //黑名单用户
	NoNewBanner             = 45 // 编辑精选封面没有新数据

	AccessDbMissMatch    = 3003 // accesstoken refreshtoken 与 数据库不匹配
	AccessCacheMissMatch = 3004 // accesstoken refreshtoken 与 缓存不匹配
	UserNotExists        = 3005 // 用户不存在
	FileUploadError      = 8001 // 文件上传失败
	FileExtError         = 8002 // 上传文件类型不允许
)
