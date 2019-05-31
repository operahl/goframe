package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"goframe/conf"
	"goframe/lib/myredis"
)

type BaseDao struct {
}

var db *xorm.Engine


// MysqlPool mysql pool
func InitDB() {
	host := conf.Config.Mysql.Host
	username := conf.Config.Mysql.UserName
	password := conf.Config.Mysql.Password
	port := conf.Config.Mysql.Port
	database := conf.Config.Mysql.Database
	maxOpenConns := conf.Config.Mysql.MaxOpenConns
	maxIdleConns := conf.Config.Mysql.MaxIdleConns

	dbsource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4", username, password, host, port, database)
	db, _ = xorm.NewEngine("mysql", dbsource)
	db.ShowSQL(true)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
}

func InitRedis()  {
	host := conf.Config.Redis.Host + ":" + conf.Config.Redis.Port
	auth := conf.Config.Redis.Auth
	db := conf.Config.Redis.Db
	maxIdle := conf.Config.Redis.MaxIdle
	MaxActive := conf.Config.Redis.MaxActive
	idleTimeout := conf.Config.Redis.IdleTimeout

	myredis.InitRedis(host,auth,db ,maxIdle,MaxActive,idleTimeout)


}