package service

import (
	"goframe/dao"
	"goframe/lib/myredis"
	"goframe/model"
)

type TestService struct{
	testDao dao.TestDao
}



func (self *TestService) GetUsers() []model.TestUser {
	return self.testDao.GetUsers()
}
func (self *TestService) GetCacheInfo() map[string]string {
	//hset user:1 aa hello
	data, _ := myredis.StringMap(myredis.HGetAll("user:1"))
	return data
}