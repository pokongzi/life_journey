package model

import (
	"life_journey/common/mysql"
	"life_journey/module/entity"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser() *User {
	return &User{
		db: mysql.DB,
	}
}

func (rep *User) Select(id string, args ...interface{}) entity.User {
	user := entity.User{}
	rep.db.Where("id = ?", id).Find(&user)
	return user
}

func (rep *User) Create(query interface{}, args ...interface{}) (id string, err error) {
	user := entity.User{}
	result := rep.db.Create(&user)
	if result.Error == nil {
		return "", result.Error
	}
	return user.UID, nil
}
