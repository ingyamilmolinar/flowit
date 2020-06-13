package config

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"

	"hash/fnv"

	validator "github.com/go-ozzo/ozzo-validation/v4"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

const transitionExceptionPrefix = "!"

func stateMachineValidator(stateMachine interface{}) error {
	switch stateMachine := stateMachine.(type) {
	case rawStateMachine:
		if err := validator.ValidateStruct(&stateMachine,
			validator.Field(&stateMachine.ID,
				validator.Required,
				validator.By(stateMachineIDValidator)),
			validator.Field(&stateMachine.Stages,
				validator.Required,
				validator.Each(validator.By(stateMachineStageValidator))),
			validator.Field(&stateMachine.InitialStage,
				validator.Required,
				validator.NewStringRule(isStateMachineStageValid(stateMachine.Stages), "State Machine Initial Stage is invalid")),
			validator.Field(&stateMachine.FinalStages,
				validator.Required,
				validator.Each(
					validator.NewStringRule(isStateMachineStageValid(stateMachine.Stages), "State Machine Final Stage is invalid"))),
			validator.Field(&stateMachine.Transitions,
				validator.Required,
				validator.Each(
					validator.By(stateMachineTransitionValidator(stateMachine.Stages)))),
		); err != nil {
			return errors.WithStack(err)
		}
		if err := validateStateMachineGraph(stateMachine); err != nil {
			return errors.WithStack(err)
		}
		return nil
	default:
		return errors.New("Invalid state machine type. Got " + reflect.TypeOf(stateMachine).Name())
	}
}

func stateMachineIDValidator(stateMachineID interface{}) error {
	return validIdentifier(stateMachineID)
}

func stateMachineStageValidator(stateMachineStage interface{}) error {
	return validIdentifier(stateMachineStage)
}

func isStateMachineStageValid(stateMachineStages []*string) func(string) bool {
	return func(stateMachineInitialStage string) bool {
		return utils.FindStringInPtrArray(stateMachineInitialStage, stateMachineStages)
	}
}

func stateMachineTransitionValidator(stateMachineStages []*string) func(interface{}) error {
	return func(stateMachineTransition interface{}) error {
		switch transition := stateMachineTransition.(type) {
		case rawStateMachineTransition:
			parsedTransition, err := parseStateMachineTransition(transition, stateMachineStages)
			if err != nil {
				return errors.WithStack(err)
			}
			return validator.ValidateStruct(&parsedTransition,
				validator.Field(&parsedTransition.From,
					validator.Required,
					validator.Each(
						validator.NewStringRule(
							isStateMachineStageValid(stateMachineStages), "State Machine Transition 'From' is invalid"))),
				validator.Field(&parsedTransition.To,
					validator.Required,
					validator.Each(
						validator.NewStringRule(
							isStateMachineStageValid(stateMachineStages), "State Machine Transition 'To' is invalid"))),
			)
		default:
			return errors.New("Invalid state machine transition type. Got " + reflect.TypeOf(transition).Name())
		}
	}
}

func parseStateMachineTransition(transition rawStateMachineTransition, stages []*string) (rawStateMachineTransition, error) {
	var result rawStateMachineTransition

	from, err := parseTransitionStages(transition.From, stages)
	if err != nil {
		return rawStateMachineTransition{}, errors.WithStack(err)
	}
	result.From = from

	to, err := parseTransitionStages(transition.To, stages)
	if err != nil {
		return rawStateMachineTransition{}, errors.WithStack(err)
	}
	result.To = to

	return result, nil
}

func parseTransitionStages(transitionStages []*string, stages []*string) ([]*string, error) {
	var parsedTransitionStages []*string
	for _, transitionStage := range transitionStages {

		if strings.HasPrefix(*transitionStage, transitionExceptionPrefix) {
			transitionStages, err := getAllOtherStages((*transitionStage)[1:], stages)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			parsedTransitionStages = append(parsedTransitionStages, transitionStages...)
		} else {
			parsedTransitionStages = append(parsedTransitionStages, transitionStage)
		}

	}
	return parsedTransitionStages, nil
}

// TODO: Where to locate this?
func getAllOtherStages(transitionStage string, stages []*string) ([]*string, error) {
	var transitionStages []*string
	found := false
	for _, stage := range stages {
		if *stage == transitionStage {
			found = true
		} else {
			transitionStages = append(transitionStages, stage)
		}
	}
	if !found {
		return nil, errors.New("Stage " + transitionStage + " is not defined")
	}
	return transitionStages, nil
}

func validateStateMachineGraph(sm rawStateMachine) error {
	dg := buildDirectedGraph(sm)
	if err := validateInitialStage(dg, *sm.InitialStage); err != nil {
		return errors.WithStack(err)
	}
	if err := validatePaths(dg, sm); err != nil {
		return errors.WithStack(err)
	}
	if err := validateFinalStages(dg, sm.FinalStages); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func validateInitialStage(dg graph.Directed, initialStage string) error {
	if dg.To(generateNodeID(initialStage)).Len() > 0 {
		return errors.New("Initial Stage '" + initialStage + "' cannot be the destination in a transition")
	}
	return nil
}

func validatePaths(dg graph.Directed, sm rawStateMachine) error {
	for _, stage := range sm.Stages {
		if found := utils.FindStringInPtrArray(*stage, sm.FinalStages); found {
			continue
		}

		reachableFinalStages := 0
		for _, reachableNode := range getReachableNodes(dg, dg.Node(generateNodeID(*stage))) {
			for _, finalStage := range sm.FinalStages {
				if reachableNode.ID() == generateNodeID(*finalStage) {
					reachableFinalStages++
				}
			}
		}
		if reachableFinalStages == 0 {
			return errors.New("Cannot reach a final node from '" + *stage + "' stage")
		}

	}
	return nil
}

func validateFinalStages(dg graph.Directed, finalStages []*string) error {
	for _, finalStage := range finalStages {
		if dg.From(generateNodeID(*finalStage)).Len() > 0 {
			return errors.New("Final stage: '" + *finalStage + "' cannot be the source in a transition")
		}
	}
	return nil
}

func getReachableNodes(dg graph.Directed, node graph.Node) []graph.Node {
	visited := make(map[int64]bool)
	var stack, reachable []graph.Node
	stack = append(stack, node)
	for len(stack) > 0 {
		popped := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		reachable = append(reachable, popped)
		nodes := dg.From(popped.ID())
		for nodes.Next() {
			next := nodes.Node()
			if !visited[next.ID()] {
				stack = append(stack, next)
				visited[next.ID()] = true
			}
		}
	}
	return reachable
}

func buildDirectedGraph(sm rawStateMachine) graph.Directed {
	parsedTransitions := make([]rawStateMachineTransition, len(sm.Transitions))
	for _, transition := range sm.Transitions {
		parsedTransition, _ := parseStateMachineTransition(*transition, sm.Stages)
		parsedTransitions = append(parsedTransitions, parsedTransition)
	}

	digraph := simple.NewDirectedGraph()
	for _, stage := range sm.Stages {
		digraph.AddNode(newNode(*stage))
	}
	for _, transition := range parsedTransitions {
		for _, from := range transition.From {
			for _, to := range transition.To {
				fromNode := digraph.Node(generateNodeID(*from))
				toNode := digraph.Node(generateNodeID(*to))
				if fromNode.ID() == toNode.ID() {
					// TODO: Workaround panic for self-referencing nodes!!
					intermediateNode := newNode(*from + "_" + *to)
					digraph.AddNode(intermediateNode)
					digraph.SetEdge(newEdge(fromNode, intermediateNode))
					digraph.SetEdge(newEdge(intermediateNode, toNode))
				} else {
					digraph.SetEdge(newEdge(fromNode, toNode))
				}
			}
		}
	}
	return digraph
}

type node struct {
	state string
	id    int64
}

func newNode(state string) node {
	return node{state, generateNodeID(state)}
}

func (n node) ID() int64 {
	return n.id
}

func generateNodeID(state string) int64 {
	h := fnv.New64a()
	h.Write([]byte(state))
	return int64(h.Sum64())
}

type edge struct {
	from graph.Node
	to   graph.Node
}

func newEdge(from graph.Node, to graph.Node) edge {
	return edge{from, to}
}

func (e edge) From() graph.Node {
	return e.from
}

func (e edge) To() graph.Node {
	return e.to
}

func (e edge) ReversedEdge() graph.Edge {
	return newEdge(e.to, e.from)
}
