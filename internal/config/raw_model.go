package config

// rawWorkflowDefinition is the typed data structure for populating the input configuration
// pointers are used to be able to signal between unset values and zero values
type rawWorkflowDefinition struct {
	Flowit *rawMainDefinition
}

type rawMainDefinition struct {
	Version   *string
	Config    *rawConfig
	Variables *rawVariables
	Branches  []*rawBranch
	Tags      []*rawTag
	Workflows []*rawWorkflow
}

type rawConfig struct {
	AbortOnFailedAction *bool `mapstructure:"abort-on-failed-action"`
	Strict              *bool
	Shell               *string
}

type rawVariables map[string]interface{}

type rawBranch struct {
	ID          *string
	Name        *string
	Prefix      *string
	Suffix      *string
	Eternal     *bool
	Protected   *bool
	Transitions []*rawTransition
}

type rawTag struct {
	ID       *string
	Format   *string
	Stages   *map[string][]*string
	Branches []*string
}

type rawWorkflow map[string][]*rawStage

// TODO: Can this be map[string]string | map[string][]string?
type rawStage map[string]interface{}

type rawTransition struct {
	From *string
	To   []*string
}
