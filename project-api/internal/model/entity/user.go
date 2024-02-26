package entity

type User struct {
	Id   int64  `gorm:"column:id" db:"column:id" json:"id" form:"id"`
	Name string `gorm:"column:name" db:"column:name" json:"name" form:"name"`
}

func (User) TableName() string {
	return "user"
}
