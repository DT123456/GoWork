package models

type User struct {
	Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
	Roles    []Role `gorm:"many2many:user_roles" json:"roles"`
}
