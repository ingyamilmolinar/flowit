package fsm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yamil-rivera/flowit/internal/fsm"
)

var _ = Describe("FSM", func() {

	stateMachines := []fsm.StateMachine{
		{
			ID: "state-machine-1",
			States: []string{
				"stage-1",
				"stage-2",
				"stage-3",
				"stage-4",
			},
			InitialState: "stage-1",
			FinalStates: []string{
				"stage-4",
			},
			Transitions: []fsm.StateMachineTransition{
				{
					From: []string{
						"stage-1",
					},
					To: []string{
						"stage-2",
						"stage-3",
					},
				},
				{
					From: []string{
						"stage-2",
						"stage-3",
					},
					To: []string{
						"stage-4",
					},
				},
			},
		},
	}

	service := fsm.NewService(stateMachines)

	Context("Retrieving available states", func() {

		It("should successfully return the initial state", func() {

			states := service.InitialState("state-machine-1")
			Expect(states).To(BeEquivalentTo("stage-1"))

		})

		It("should successfully return the available states", func() {

			transitions := stateMachines[0].Transitions
			for _, transition := range transitions {
				for _, from := range transition.From {

					states := service.AvailableStates("state-machine-1", from)
					Expect(states).To(ConsistOf(transition.To))

				}
			}

			states := service.AvailableStates("state-machine-1", "stage-4")
			Expect(len(states)).To(Equal(0))

		})

	})

	Context("Verifying valid transitions", func() {

		It("should successfully verify if a transition is valid", func() {

			valid := service.IsTransitionValid("state-machine-1", "stage-1")
			Expect(valid).To(BeTrue())

			transitions := stateMachines[0].Transitions
			for _, transition := range transitions {

				for _, from := range transition.From {
					for _, to := range transition.To {

						valid := service.IsTransitionValid("state-machine-1", from, to)
						Expect(valid).To(BeTrue())

					}
				}

				for _, from := range transition.From {
					for _, to := range transition.To {

						valid := service.IsTransitionValid("state-machine-1", to, from)
						Expect(valid).To(BeFalse())

					}
				}
			}
		})

		It("should successfully verify invalid parameters", func() {

			valid := service.IsTransitionValid("state-machine-1")
			Expect(valid).To(BeFalse())

			valid = service.IsTransitionValid("state-machine-1", "stage-1", "stage-2", "stage-3")
			Expect(valid).To(BeFalse())

		})

		Context("Verifying active state", func() {

			It("should successfully verify active state", func() {

				for _, state := range stateMachines[0].States {
					isFinal := false
					for _, finalState := range stateMachines[0].FinalStates {
						if state == finalState {
							isFinal = true
						}
					}

					active := service.IsActiveState("state-machine-1", state)
					// final state is not considered active
					if isFinal {
						Expect(active).To(BeFalse())
					} else {
						Expect(active).To(BeTrue())
					}

				}

			})

		})

		Context("Verifying final state", func() {

			It("should successfully verify final state", func() {

				for _, state := range stateMachines[0].States {
					isFinal := false
					for _, finalState := range stateMachines[0].FinalStates {
						if state == finalState {
							isFinal = true
						}
					}

					final := service.IsFinalState("state-machine-1", state)
					if isFinal {
						Expect(final).To(BeTrue())
					} else {
						Expect(final).To(BeFalse())
					}

				}

			})

		})

	})

})
