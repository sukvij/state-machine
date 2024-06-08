package usersm

const (
	A = "a"
	B = "b"
	C = "c"
)

var UserSM = NewSM("user-sm")

func init() {
	// UserSM.AddTransition(A, sm.StateInit, sm.EventAuth, GetAllUsers, sm.StateAuth, sm.StateAuth, sm.StateCapt)
	// UserSM.AddTransition(A, sm.StateInit, sm.EventAuth, CreaateUser, sm.StateAuth, sm.StateAuth, sm.StateCapt)

}
