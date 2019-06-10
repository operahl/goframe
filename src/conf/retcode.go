package conf

//所有的参数错误以400开头，资源位找到404开头，内部错误以500开头
const (
	CodeOk                  = 0  // seccess
	CodeAesErr              = 1  // aes error1
	CodeParaErr             = 3  // param error
	SystemError             = 5  // 服务器错误
)
