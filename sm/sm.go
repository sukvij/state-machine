package sm

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

const (
	logTag = "SM"
	// this is initial state
	StateInit = State(0)
	StateAuth = State(1)
	StateCapt = State(2)
	// EventNoNeed denotes the default event for trigger FSM run without outside caller
	EventNoNeed  = Event(0)
	EventAuth    = Event(1)
	EventCapture = Event(2)
)

type State int16

type Event int16

type Function interface{}

type SM struct {
	trans map[stateEvent]*cmdStates // The command and destination states for each (source state, event)
	name  string
}

// Command is an abstraction of action that's done on a state
type command struct {
	name string
	do   Function
}

// stateEvent is a tuple (state, event)
// this struct is only used as the key of the transition map in state machine
type stateEvent struct {
	src   State
	event Event
}

// cmdState is a tuple (command, states), where states is a set of states
// this struct is only used as the value of the transition map in state machine
type cmdStates struct {
	cmd  *command
	dest map[State]bool
}

// StateContext is the persistent form of the state machine, usually it's a row in db.
type StateContext interface {
	GetMsgID() string
	GetState() State
	SetState(state State)
	Update(ctx context.Context, stateCtx StateContext) (err error)
}

// New creates a state machine. Context is the real entity of fsm which can be save to db.
// Deprecated: should not use this func directly, use smgenerator instead, which is similar to C++ template.
func New(name string) *SM {
	return &SM{
		trans: map[stateEvent]*cmdStates{},
		name:  name,
	}
}

func (sm *SM) AddTransition(name string, currState State, ev Event, goFunc Function, nextStates ...State) {
	if len(nextStates) == 0 {
		log.Fatal(logTag, "no next event found, event = %v, currentState = %v", ev, currState)
		panic("nextState cannot be empty")
	}

	if errValidate := validateTransitionFunction(goFunc); errValidate != nil {
		panic("invalid transition function")
	}

	se := stateEvent{currState, ev}

	if _, ok := sm.trans[se]; ok {
		panic("currState and ev already registered.") // to avoid infinite loop
	}

	command := &command{name: name, do: goFunc}
	sm.trans[se] = &cmdStates{command, map[State]bool{}}

	for _, nextStates := range nextStates {
		sm.trans[se].dest[nextStates] = true
	}
}

func (sm *SM) validateNextState(currState State, ev Event, nextState State) bool {
	se := stateEvent{currState, ev}
	if _, ok := sm.trans[se]; !ok {
		return false
	}
	return sm.trans[se].dest[nextState]
}

// get function to make the transition
func (sm *SM) getCommand(stateCtx StateContext, ev Event) *command {
	se := stateEvent{stateCtx.GetState(), ev}
	if _, ok := sm.trans[se]; !ok {
		return nil
	}
	return sm.trans[se].cmd
}

func (sm *SM) Run(ctx context.Context, stateCtx StateContext, ev Event, arguments ...interface{}) (nextCtx StateContext, err error) {
	for { // move SM as far as possible

		stepCount := 0
		cmd := sm.getCommand(stateCtx, ev)
		if cmd == nil {
			if stepCount > 0 { // SM has finished success
				return nextCtx, nil
			}
			smErr := fmt.Errorf("no transition found, msgID = %v, currState = %v, event=%v", stateCtx.GetMsgID(), stateCtx.GetState(), ev)
			return nil, smErr
		}

		// nextCtx, err = sm.executeTransitionFunc(cmd.do, ctx, stateCtx, arguments) // we will create new nextCtx for next state
		// nextCtx.Update(ctx, stateCtx)

		// stateCtx = nextCtx
		ev = EventNoNeed
		stepCount++
	}
}

func validateTransitionFunction(goFunc Function) error {
	funcType := reflect.TypeOf(goFunc)

	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("%v is not a function", goFunc)
	}

	// funcType.NumIn() count inbound parameteres
	// funcType.NumOut() count outbound parameters

	return nil
}
