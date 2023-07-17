package models

type User struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(16)" json:"username"`
	Password []byte `gorm:"type:bytea" json:"-"`
}
