package model

type Address struct {
	UserIdType
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UserIdType struct {
	UserId uint `json:"user_id" form:"user_id" gorm:"not null;index:UserIdIdx;comment:'用户ID'"`
}
