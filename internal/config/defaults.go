package config

import (
	"os"
)

type defaults struct {
	AbortOnFailedAction bool
	Strict              bool
	Shell               string
	Stages              rawStages
	Branches            []*string
}

// We need a package global variable to be the target of our rawWorkflowDefinition pointers
var defaultValues defaults
var emptyConfig rawConfig

// setDefaults receives a validated raw workflow definition and includes default values where nil pointers are seen
// TODO: Can we make this function stateless??
func setDefaults(workflowDefinition *rawWorkflowDefinition) {

	defaultValues.AbortOnFailedAction = true

	defaultValues.Strict = false

	defaultValues.Shell = func() string {
		envShell := os.Getenv("SHELL")
		if envShell != "" {
			return envShell
		}
		// TODO: Set shell according to the OS
		return "/usr/bin/env bash"
	}()

	defaultValues.Stages = func() rawStages {
		// Defaults to all stages
		allStages := make(rawStages)
		for _, workflow := range workflowDefinition.Flowit.Workflows {
			for workflowID, stages := range *workflow {
				stagesIDs := make([]*string, len(stages))
				for i, stage := range stages {
					stagesIDs[i] = stage.ID
				}
				allStages[workflowID] = stagesIDs
			}
		}
		return allStages
	}()

	defaultValues.Branches = func() []*string {
		// Defaults to all branches
		allBranches := make([]*string, len(workflowDefinition.Flowit.Branches))
		for i, branch := range workflowDefinition.Flowit.Branches {
			allBranches[i] = branch.ID
		}
		return allBranches
	}()

	setDefaultValues(workflowDefinition, &defaultValues)
}

func setDefaultValues(workflowDefinition *rawWorkflowDefinition, defaultValues *defaults) {
	if workflowDefinition.Flowit.Config == nil {
		// In case 'config' section is missing all together
		// TODO: This is extremely ugly
		workflowDefinition.Flowit.Config = &emptyConfig
	}
	if workflowDefinition.Flowit.Config.AbortOnFailedAction == nil {
		workflowDefinition.Flowit.Config.AbortOnFailedAction = &defaultValues.AbortOnFailedAction
	}
	if workflowDefinition.Flowit.Config.Strict == nil {
		workflowDefinition.Flowit.Config.Strict = &defaultValues.Strict
	}
	if workflowDefinition.Flowit.Config.Shell == nil {
		workflowDefinition.Flowit.Config.Shell = &defaultValues.Shell
	}
	for _, tag := range workflowDefinition.Flowit.Tags {
		if tag.Stages == nil {
			tag.Stages = &defaultValues.Stages
		}
		if tag.Branches == nil {
			tag.Branches = defaultValues.Branches
		}
	}

}
