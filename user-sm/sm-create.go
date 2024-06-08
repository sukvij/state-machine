package usersm

import "vijju/sm"

type UserSMImpl struct {
	innerSM *sm.SM
}

func NewSM(name string) *UserSMImpl {
	return &UserSMImpl{
		innerSM: sm.New(name),
	}
}

// add transition define the one transition for the sm

func (s *UserSMImpl) AddTransition(name string, currentState sm.State, ev sm.Event, transFunc UserMethods, nextState ...sm.State) {
	s.innerSM.AddTransition(name, currentState, ev, transFunc, nextState...)
}
