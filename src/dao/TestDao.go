package dao

import (
	"goframe/model"
)

type TestDao struct {

}

func (t *TestDao) GetUsers() []model.TestUser {
	var data []model.TestUser
	db.Table("user").Limit(30).Find(&data)
	return data
}
