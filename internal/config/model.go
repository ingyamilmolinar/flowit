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
	Stages   []workflowType
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

type workflowType map[string][]stage

type stage map[string]interface{}

type transition struct {
	From string
	To   []string
}
