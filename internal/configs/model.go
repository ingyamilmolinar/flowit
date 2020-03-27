package configs

// Flowit is the typed data structure for populating the input configuration
type Flowit struct {
	Version   *string    `valid:"flowitversion~Unsupported flowit version,required"`
	Config    *config    `valid:"optional"`
	Variables *variables `valid:"-"`
	Workflow  *workflow  `valid:"required"`
}

type config struct {
	AbortOnFailedAction *bool   `valid:"optional" mapstructure:"abort-on-failed-action"`
	Strict              *bool   `valid:"optional"`
	Shell               *string `valid:"flowitconfigshell~Invalid config shell,optional"`
}
type variables map[string]interface{}

type workflow struct {
	Branches []*branch `valid:"flowitbranches~Invalid branches,required"`
	Tags     []*tag    `valid:"-"` //`valid:"optional"`
	Stages   []*stage  `valid:"-"` //`valid:"required"`
}

type branch struct {
	ID          *string       `valid:"flowitbranchid~Invalid branch ID,required"`
	Name        *string       `valid:"flowitbranchname~Invalid branch name,required"`
	Prefix      *string       `valid:"flowitbranchpreffix~Invalid branch preffix,optional"`
	Suffix      *string       `valid:"flowitbranchsuffix~Invalid branch suffix,optional"`
	Eternal     *bool         `valid:"required"`
	Protected   *bool         `valid:"required"`
	Transitions []*transition `valid:"flowitbranchtransitions~Invalid branch transitions,optional"`
}

type tag map[string]interface{}

type stage map[string]interface{}

type transition struct {
	From *string   `valid:"required"`
	To   []*string `valid:"required"`
}
