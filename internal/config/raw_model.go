package config

// rawWorkflowDefinition is the typed data structure used for populating and validating the workflow configuration
// Pointers are used extensibly to be able to differentiate between unset values and default zero values
type rawWorkflowDefinition struct {
	Flowit *rawMainDefinition
}

type rawMainDefinition struct {
	Version       *string
	Config        *rawConfig
	Variables     *rawVariables
	StateMachines []*rawStateMachine `mapstructure:"state-machines"`
	Workflows     []*rawWorkflow
}

type rawConfig struct {
	AbortOnFailedAction *bool `mapstructure:"abort-on-failed-action"`
	Shell               *string
}

type rawVariables map[string]interface{}

type rawStateMachine struct {
	ID           *string
	Stages       []*string
	InitialStage *string   `mapstructure:"initial-stage"`
	FinalStages  []*string `mapstructure:"final-stages"`
	Transitions  []*rawStateMachineTransition
}

type rawStateMachineTransition struct {
	From []*string
	To   []*string
}

type rawStages map[string][]*string

type rawWorkflow struct {
	ID           *string
	StateMachine *string `mapstructure:"state-machine"`
	Stages       []*rawStage
}

type rawStage struct {
	ID         *string
	Args       []*string
	Conditions []*string
	Actions    []*string
}

type rawTransition struct {
	From *string
	To   []*string
}
