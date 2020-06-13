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

func setDefaults(workflowDefinition *rawWorkflowDefinition) {
	setDefaultValues(workflowDefinition, generateDefaultValues(workflowDefinition))
}

func generateDefaultValues(workflowDefinition *rawWorkflowDefinition) *defaults {

	var defaultValues defaults

	defaultValues.AbortOnFailedAction = true
	defaultValues.Strict = false
	defaultValues.Shell = generateDefaultShell()
	defaultValues.Stages = generateDefaultStages(workflowDefinition.Flowit.Workflows)
	defaultValues.Branches = generateDefaultBranches(workflowDefinition.Flowit.Branches)

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

func generateDefaultStages(workflows []*rawWorkflow) rawStages {
	allStages := make(rawStages)
	for _, workflow := range workflows {
		workflowID := workflow.ID
		stagesIDs := make([]*string, len(workflow.Stages))
		for i, stage := range workflow.Stages {
			stagesIDs[i] = stage.ID
		}
		allStages[*workflowID] = stagesIDs
	}
	return allStages
}

func generateDefaultBranches(branches []*rawBranch) []*string {
	allBranches := make([]*string, len(branches))
	for i, branch := range branches {
		allBranches[i] = branch.ID
	}
	return allBranches
}

func setDefaultValues(workflowDefinition *rawWorkflowDefinition, defaultValues *defaults) {
	// In case 'config' section is missing all together
	if workflowDefinition.Flowit.Config == nil {
		workflowDefinition.Flowit.Config = &rawConfig{}
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
