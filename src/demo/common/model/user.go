package model

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "User"
}
