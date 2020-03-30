package configs

// RawFlowitConfig is the typed data structure for populating the input configuration
type RawFlowitConfig struct {
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

// FlowitConfig is the user friendly data structure for consuming the input configuration
type FlowitConfig struct {
	Version   string
	Config    config
	Variables variables
	Workflow  workflow
}

type config struct {
	AbortOnFailedAction bool
	Strict              bool
	Shell               string
}
type variables map[string]interface{}

type workflow struct {
	Branches []branch
	Tags     []tag
	Stages   []branchType
}

type branch struct {
	ID          string
	Name        string
	Prefix      string
	Suffix      string
	Eternal     bool
	Protected   bool
	Transitions []transition
}

type tag struct {
	ID       string
	Format   string
	Stages   map[string][]string
	Branches []string
}

type branchType map[string][]stage

type stage struct {
	//Name       *map[string]interface{} `valid:"required"`
	Start      string
	Fetch      string
	Sync       string
	Publish    string
	Finish     string
	Conditions []string
	Actions    []string
}

type transition struct {
	From string
	To   []string
}
