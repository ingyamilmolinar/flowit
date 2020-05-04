package config

// WorkflowDefinition is the consumer friendly data structure for reading the input configuration
type WorkflowDefinition struct {
	Flowit mainDefinition
}

type mainDefinition struct {
	Version   string
	Config    config
	Variables variables
	Branches  []branch
	Tags      []tag
	Workflows []workflow
}

type config struct {
	AbortOnFailedAction bool
	Strict              bool
	Shell               string
}
type variables map[string]interface{}

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
	Stages   stages
	Branches []string
}

type stages map[string][]string

type workflow map[string][]stage

type stage struct {
	ID         string
	Args       []string
	Conditions []string
	Actions    []string
}

type transition struct {
	From string
	To   []string
}
