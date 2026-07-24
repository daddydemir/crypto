package domain

type User struct {
	ID       string `gorm:"primary_key"`
	Username string `gorm:"size:25;unique" json:"username"`
	Password string `gorm:"size:255" json:"password"`
	IsActive bool   `json:"-" gorm:"default:true"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserAuthRepository interface {
	CreateUser(user User) error
	Login(username, password string) (*LoginResponse, error)
}
