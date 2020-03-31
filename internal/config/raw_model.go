package config

// rawFlowitConfig is the typed data structure for populating the input configuration
// pointers are used to be able to signal between unset values and zero values
type rawFlowitConfig struct {
	Version   *string       `valid:"flowitversion~Unsupported flowit version,required"`
	Config    *rawConfig    `valid:"optional"`
	Variables *rawVariables `valid:"-"`
	Workflow  *rawWorkflow  `valid:"required"`
}

type rawConfig struct {
	AbortOnFailedAction *bool   `valid:"optional" mapstructure:"abort-on-failed-action"`
	Strict              *bool   `valid:"optional"`
	Shell               *string `valid:"flowitconfigshell~Invalid config shell,optional"`
}

type rawVariables map[string]interface{}

type rawWorkflow struct {
	Branches []*rawBranch     `valid:"flowitbranches~Invalid branches,required"`
	Tags     []*rawTag        `valid:"optional"`
	Stages   []*rawBranchType `valid:"required"`
}

type rawBranch struct {
	ID          *string          `valid:"flowitbranchid~Invalid branch ID,required"`
	Name        *string          `valid:"flowitbranchname~Invalid branch name,required"`
	Prefix      *string          `valid:"flowitbranchpreffix~Invalid branch preffix,optional"`
	Suffix      *string          `valid:"flowitbranchsuffix~Invalid branch suffix,optional"`
	Eternal     *bool            `valid:"required"`
	Protected   *bool            `valid:"required"`
	Transitions []*rawTransition `valid:"flowitbranchtransitions~Invalid branch transitions,optional"`
}

type rawTag struct {
	ID       *string               `valid:"required"`
	Format   *string               `valid:"required"`
	Stages   *map[string][]*string `valid:"required"`
	Branches []*string             `valid:"required"`
}

type rawBranchType map[string][]*rawStage

type rawStage struct {
	//Name       *map[string]interface{} `valid:"required"`
	// TODO: Can we avoid harcoding stages?
	Start      *string   `valid:"required"`
	Fetch      *string   `valid:"optional"`
	Sync       *string   `valid:"optional"`
	Publish    *string   `valid:"optional"`
	Finish     *string   `valid:"required"`
	Conditions []*string `valid:"optional"`
	Actions    []*string `valid:"required"`
}

type rawTransition struct {
	From *string   `valid:"required"`
	To   []*string `valid:"required"`
}
