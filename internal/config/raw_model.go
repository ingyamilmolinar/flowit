package config

// rawFlowitConfig is the typed data structure for populating the input configuration
// pointers are used to be able to signal between unset values and zero values
type rawFlowitConfig struct {
	Version   *string
	Config    *rawConfig
	Variables *rawVariables
	Workflow  *rawWorkflow
}

type rawConfig struct {
	AbortOnFailedAction *bool `mapstructure:"abort-on-failed-action"`
	Strict              *bool
	Shell               *string
}

type rawVariables map[string]interface{}

type rawWorkflow struct {
	Branches []*rawBranch
	Tags     []*rawTag
	Stages   []*rawWorkflowType
}

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

type rawWorkflowType map[string][]*rawStage

// TODO: Can this be map[string]string | map[string][]string?
type rawStage map[string]interface{}

type rawTransition struct {
	From *string
	To   []*string
}
