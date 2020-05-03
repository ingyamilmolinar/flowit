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
	Stages   map[string][]string
	Branches []string
}

type workflow map[string][]stage

// TODO: Can this be map[string]string | map[string][]string?
type stage map[string]interface{}

type transition struct {
	From string
	To   []string
}
