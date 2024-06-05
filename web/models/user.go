package models

import "web/dao"

type User struct {
	Openid int
	Score  int
}

func (User) TableName() string {
	return "user"
}

func GetUserTest(openid int) (User, error) {
	var user User
	err := dao.DB.Where("openid = ?", openid).First(&user).Error
	return user, err
}
