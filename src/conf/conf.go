package conf

import (
	"encoding/hex"

	"github.com/BurntSushi/toml"
)

var Config Conf

type Conf struct {
	Server struct {
		Port               string //http服务端口
		Aeskey             string //aes加密秘钥
		AeskeyBytes        []byte //同上，二进制状态，初始访问会对其赋值，进行Cache
		WhiteIp            string //input接口IP白名单
		Mode               string
		Ginmode            string
		PassSalt           string
		RootPath           string
		DeviceId           int    // 0 imei, other
		Domain             string //域名
	}
	Log struct{
		LogPath            string //错误日志落地目录
		LogLevel           string // debug  info  error
		MaxSize			   int	 //单个日志最大容量
		MaxAge			   int	//日志最大保留天数
		MaxBackups		   int //备份最大天数

	}
	Redis struct {
		Host 			string //地址
		Port 			string //端口
		Auth 			string //口令
		Db   			int    //select db
		MaxIdle 		int		//空闲连接数
		MaxActive 		int		//最大连接数
		IdleTimeout		int		//连接超时时间
	}
	Resource struct {
		VideoPath              string //视频源路径
		ImagePath              string //图片源路径
		AvatarPath             string //头像源路径
		VideoCompressPath      string //视频压缩路径
		ImageCompressPath      string //图片压缩路径
		AvatarPathCompressPath string //头像压缩路径
	}
	Mysql struct {
		Host              string //地址
		Port              string //端口
		UserName          string //用户名
		Password          string //密码
		Database          string //数据库
		MaxOpenConns      int    //最大连接数
		MaxIdleConns      int    //最少启动连接数

	}

	IpWhiteList struct {
		Ips []string
	}
	Common struct {
		TokenTtl                 int
		RedisExpireTime          int
	}

}


func ReadCfg(path string) Conf {
	if Config.Mysql.Host != "" {
		return Config
	}
	if _, err := toml.DecodeFile(path, &Config); err != nil {
		panic("Parse File Error:" + err.Error())
	}
	if aeskey, err := hex.DecodeString(Config.Server.Aeskey); err != nil {
		panic("aes key error:" + err.Error())
	} else {
		Config.Server.AeskeyBytes = aeskey
	}
	return Config
}
