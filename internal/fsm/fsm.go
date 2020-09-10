package fsm

import (
	"github.com/looplab/fsm"
)

// Service exposes the methods to interact with the FSM service
type Service struct {
	stateMachines map[string]*fsm.FSM
}

// StateMachine is the data structure representing the state machine properties
// that will initialize the FSM service
type StateMachine struct {
	ID           string
	States       []string
	InitialState string
	FinalStates  []string
	Transitions  []StateMachineTransition
}

// StateMachineTransition encodes the allowed transitions between state machine states
type StateMachineTransition struct {
	From []string
	To   []string
}

// NewService initializes and returns a new instance of the FSM service
func NewService(stateMachines []StateMachine) *Service {
	var smMap = make(map[string]*fsm.FSM, len(stateMachines))
	for _, stateMachine := range stateMachines {
		stateMachineID := stateMachine.ID
		states := stateMachine.States

		fsmEvents := make([]fsm.EventDesc, len(states))
		allStates := make([]string, len(states)+1)
		allStates[0] = originState()
		for j, state := range states {
			allStates[j+1] = state
		}

		for j, state := range states {

			if state == stateMachine.InitialState {
				fsmEvents[j] = fsm.EventDesc{
					Name: state,
					Src:  []string{originState()},
					Dst:  state,
				}
			} else {
				src, dst := generateStates(state, stateMachine.Transitions)
				fsmEvents[j] = fsm.EventDesc{
					Name: state,
					Src:  src,
					Dst:  dst,
				}
			}

		}
		smMap[stateMachineID] = fsm.NewFSM(originState(), fsmEvents, map[string]fsm.Callback{})
	}
	return &Service{stateMachines: smMap}
}

// IsTransitionValid verifies whether or not a state machine can transition between two given states
func (s Service) IsTransitionValid(stateMachineID string, states ...string) bool {
	if len(states) == 0 || len(states) > 2 {
		return false
	}

	stateMachine := s.stateMachines[stateMachineID]
	originalState := stateMachine.Current()

	var fromState, toState string
	if len(states) == 1 {
		fromState = originalState
		toState = states[0]
	} else {
		fromState = states[0]
		toState = states[1]
	}

	stateMachine.SetState(fromState)
	canTransition := stateMachine.Can(toState)
	stateMachine.SetState(originalState)

	return canTransition
}

// AvailableStates returns the states that are immediately available to transition to
// for a given state machine
func (s Service) AvailableStates(stateMachineID string, currentState string) []string {
	stateMachine := s.stateMachines[stateMachineID]
	originalState := stateMachine.Current()
	stateMachine.SetState(currentState)
	availableTransitions := stateMachine.AvailableTransitions()
	stateMachine.SetState(originalState)

	return availableTransitions
}

// InitialState returns the initial state of a state machine
func (s Service) InitialState(stateMachineID string) string {
	stateMachine := s.stateMachines[stateMachineID]
	return stateMachine.AvailableTransitions()[0]
}

// IsActiveState validates whether or not a particular state is active
// for a given state machine. Active states are all state machine states
// except the origin state and the final state
func (s Service) IsActiveState(stateMachineID, state string) bool {
	stateMachine := s.stateMachines[stateMachineID]
	originState := stateMachine.Current()
	stateMachine.SetState(state)
	availableTransitions := len(stateMachine.AvailableTransitions())
	stateMachine.SetState(originState)
	return originState != state && availableTransitions > 0
}

// IsFinalState validates whether or not a particular state is final
// for a given state machine
func (s Service) IsFinalState(stateMachineID, state string) bool {
	originState := s.stateMachines[stateMachineID].Current()
	return !s.IsActiveState(stateMachineID, state) && originState != state
}

func generateStates(stage string, transitions []StateMachineTransition) ([]string, string) {
	var srcStages []string
	for _, transition := range transitions {
		for _, to := range transition.To {
			if to == stage {
				srcStages = append(srcStages, transition.From...)
			}
		}
	}
	return srcStages, stage
}

func originState() string {
	return "origin"
}
