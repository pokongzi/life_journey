package entity

type User struct {
	ID           int64  `gorm:"id"`
	UID          string `gorm:"uid"`
	Name         string `gorm:"name"`
	Birthday     string `gorm:"birthday"`
	Phone        string `gorm:"phone"`
	BirthdayType int    `gorm:"birthday_type"`
	Icon         string `gorm:"icon"`
	Sex          int    `gorm:"sex"`
	Status       int    `gorm:"status"`
	CreateTime   int64  `gorm:"create_time"`
	UpdateTime   int64  `gorm:"update_time"`
}

func (u *User) TableName() string {
	return "user"
}
