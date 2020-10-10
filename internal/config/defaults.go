package config

import (
	"os"
)

type defaults struct {
	CheckpointExecution bool
	Shell               string
	Stages              rawStages
	Branches            []*string
}

func setDefaults(workflowDefinition *rawWorkflowDefinition) {
	setDefaultValues(workflowDefinition, generateDefaultValues(workflowDefinition))
}

func generateDefaultValues(workflowDefinition *rawWorkflowDefinition) *defaults {

	var defaultValues defaults

	defaultValues.CheckpointExecution = true
	defaultValues.Shell = generateDefaultShell()

	return &defaultValues
}

// TODO: Set shell according to the OS
func generateDefaultShell() string {
	envShell := os.Getenv("SHELL")
	if envShell != "" {
		return envShell
	}
	return "/usr/bin/env bash"
}

func setDefaultValues(workflowDefinition *rawWorkflowDefinition, defaultValues *defaults) {
	// In case 'config' section is missing all together
	if workflowDefinition.Flowit.Config == nil {
		workflowDefinition.Flowit.Config = &rawConfig{}
	}
	if workflowDefinition.Flowit.Config.CheckpointExecution == nil {
		workflowDefinition.Flowit.Config.CheckpointExecution = &defaultValues.CheckpointExecution
	}
	if workflowDefinition.Flowit.Config.Shell == nil {
		workflowDefinition.Flowit.Config.Shell = &defaultValues.Shell
	}
}
