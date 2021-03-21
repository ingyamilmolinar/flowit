package fsm

import (
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/utils"
)

// FsmServiceFactory defines the methods that an FSM Service Factory should implement
type FsmServiceFactory interface {
	NewFsmService(config config.Flowit) (Service, error)
}

// ServiceFactory exposes the FSM Service Factory methods
type ServiceFactory struct{}

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

// NewServiceFactory returns the default implementation of the FSM Service Factory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
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

// OriginState returns the very first state that ALL state machines start with`
// This is different than the InitialState and is the same for ALL state machines
func (s Service) OriginState() string {
	return originState()
}

// InitialState returns the initial state of a state machine given a state machine ID
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

// IsFinalState validates whether or not a particular state is the last state
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

// NewFsmService receives a Flowit configuration and returns a new FSM Service
// TODO: Unit test
func (s ServiceFactory) NewFsmService(definition config.Flowit) (service Service, err error) {
	fsms, err := buildFSMs(definition.StateMachines, definition.Workflows)
	if err != nil {
		return service, errors.WithStack(err)
	}
	return *NewService(fsms), nil
}

func buildFSMs(stateMachines []config.StateMachine, workflows []config.Workflow) ([]StateMachine, error) {

	fsmWorkflows := make([]StateMachine, len(workflows))
	for i, workflow := range workflows {
		fsmWorkflows[i].ID = workflow.StateMachine

		var stateMachine StateMachine
		for _, sm := range stateMachines {

			if workflow.StateMachine == sm.ID {
				stateMachine.States = sm.Stages
				stateMachine.InitialState = sm.InitialStage
				stateMachine.FinalStates = sm.FinalStages
				fsmTransitions, err := buildTransitions(sm.Transitions)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				stateMachine.Transitions = fsmTransitions
			}

		}

		fsmWorkflows[i].States = stateMachine.States
		fsmWorkflows[i].InitialState = stateMachine.InitialState
		fsmWorkflows[i].FinalStates = stateMachine.FinalStates
		fsmWorkflows[i].Transitions = stateMachine.Transitions

	}
	return fsmWorkflows, nil
}

func buildTransitions(configTransitions []config.StateMachineTransition) ([]StateMachineTransition, error) {

	var fsmTransitions []StateMachineTransition
	if err := utils.DeepCopy(configTransitions, &fsmTransitions); err != nil {
		return nil, errors.WithStack(err)
	}
	return fsmTransitions, nil

}
