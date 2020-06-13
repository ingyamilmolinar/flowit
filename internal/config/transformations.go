package config

import "strings"

func applyTransformations(workflowDefinition *rawWorkflowDefinition) {
	transformStateMachines(workflowDefinition.Flowit.StateMachines)
}

func transformStateMachines(stateMachines []*rawStateMachine) {
	for _, sm := range stateMachines {
		transformTransitions(sm.Transitions, sm.Stages)
	}
}

func transformTransitions(transitions []*rawStateMachineTransition, stages []*string) {
	for _, transition := range transitions {
		var transformedFrom []*string
		for _, from := range transition.From {
			if strings.HasPrefix(*from, transitionExceptionPrefix()) {
				// We can safely ignore error since the transitions were validated already
				otherStages, _ := getAllOtherStages((*from)[1:], stages)
				transformedFrom = append(transformedFrom, otherStages...)
			} else {
				transformedFrom = append(transformedFrom, from)
			}
		}
		transition.From = transformedFrom
		var transformedTo []*string
		for _, to := range transition.To {
			if strings.HasPrefix(*to, transitionExceptionPrefix()) {
				// We can safely ignore error since the transitions were validated already
				otherStages, _ := getAllOtherStages((*to)[1:], stages)
				transformedTo = append(transformedTo, otherStages...)
			} else {
				transformedTo = append(transformedTo, to)
			}
		}
		transition.To = transformedTo
	}
}
