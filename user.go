package model

type User struct {
	Iduser int
	Username string
	Password string
	Name string
	IsCompany int
}

func (user *User) String() string {
	return "user : {" + user.Username + ", " + user.Password + ", " + user.Name + "}"
}
