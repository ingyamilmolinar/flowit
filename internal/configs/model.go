package configs

type Flowit struct {
	Version   *string    `valid:"flowitversion~Unsupported flowit version,required"`
	Config    *Config    `valid:"optional"`
	Variables *Variables `valid:"-"`
	Workflow  *Workflow  `valid:"required"`
}

type Config struct {
	AbortOnFailedAction *bool   `valid:"optional" mapstructure:"abort-on-failed-action"`
	Strict              *bool   `valid:"optional"`
	Shell               *string `valid:"optional"`
}
type Variables map[string]interface{}

type Workflow struct {
	Branches []*Branch `valid:"-"` //`valid:"flowitbranches,required"`
	Tags     []*Tag    `valid:"-"` //`valid:"optional"`
	Stages   []*Stage  `valid:"-"` //`valid:"required"`
}

type Tag map[string]interface{}

type Stage map[string]interface{}

type Branch struct {
	Id          *string       `valid:"required"`
	Name        *string       `valid:"required"`
	Prefix      *string       `valid:"optional"`
	Suffix      *string       `valid:"optional"`
	Eternal     *bool         `valid:"required"`
	Protected   *bool         `valid:"required"`
	Transitions []*Transition `valid:"optional"`
}

type Transition struct {
	From *string   `valid:"required"`
	To   []*string `valid:"required"`
}
