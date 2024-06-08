package usersm

type User struct {
	Name string
}

type UserMethods interface {
	GetAllUsers() []*User
	CreaateUser() *User
}

func GetAllUsers() {

}

func CreaateUser() {

}
