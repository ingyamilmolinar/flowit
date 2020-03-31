package config

// FlowitConfig is the consumer friendly data structure for reading the input configuration
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
