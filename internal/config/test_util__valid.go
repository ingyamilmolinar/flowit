package config

func validConfigJustMandatoryFields() rawWorkflowDefinition {

	var config rawWorkflowDefinition
	var flowit rawMainDefinition

	config.Flowit = &flowit

	version := "0.1"

	var branch rawBranch
	branchID := "master"
	branchName := "master"
	branchEternal := true
	branchProtected := true
	branch.ID = &branchID
	branch.Name = &branchName
	branch.Eternal = &branchEternal
	branch.Protected = &branchProtected

	var stateMachine rawStateMachine
	stateMachineID := "simple-machine"
	startStageID := "start"
	finishStageID := "finish"
	stateMachineStages := []*string{
		&startStageID, &finishStageID,
	}
	stateMachineInitialStage := startStageID
	stateMachineFinalStages := []*string{
		&finishStageID,
	}
	stateMachineTransitions := []*rawStateMachineTransition{
		{
			From: []*string{
				&startStageID,
			},
			To: []*string{
				&finishStageID,
			},
		},
	}

	stateMachine.ID = &stateMachineID
	stateMachine.Stages = stateMachineStages
	stateMachine.InitialStage = &stateMachineInitialStage
	stateMachine.FinalStages = stateMachineFinalStages
	stateMachine.Transitions = stateMachineTransitions

	startStageAction1 := "start action1"
	startStageAction2 := "start action2"
	startStage := rawStage{
		ID:      &startStageID,
		Actions: []*string{&startStageAction1, &startStageAction2},
	}
	finishStageAction1 := "finish action1"
	finishStageAction2 := "finish action2"
	finishStage := rawStage{
		ID:      &finishStageID,
		Actions: []*string{&finishStageAction1, &finishStageAction2},
	}
	workflowID := "feature"
	workflowType := rawWorkflow{
		ID:           &workflowID,
		StateMachine: &stateMachineID,
		Stages: []*rawStage{
			&startStage,
			&finishStage,
		},
	}

	mainConfig := rawMainDefinition{
		Version: &version,
		Branches: []*rawBranch{
			&branch,
		},
		StateMachines: []*rawStateMachine{
			&stateMachine,
		},
		Workflows: []*rawWorkflow{
			&workflowType,
		},
	}

	config.Flowit = &mainConfig
	return config
}

func validConfigWithOptionalFields() WorkflowDefinition {

	var flowit mainDefinition

	flowit.Version = "0.1"
	flowit.Config = Config{
		AbortOnFailedAction: true,
		Strict:              false,
		Shell:               "/usr/bin/env bash",
	}
	flowit.Variables = map[string]interface{}{
		"var1": "value",
		"var2": 12345,
		"var3": "${env-variable}",
	}
	flowit.Branches = []Branch{
		{
			ID:        "master",
			Name:      "master",
			Eternal:   true,
			Protected: true,
		},
		{
			ID:        "feature",
			Name:      "$<prefix>$<suffix>",
			Prefix:    "feature/$<jira-issue-id>",
			Suffix:    "$<feature-branch-suffix>",
			Eternal:   false,
			Protected: false,
			Transitions: []Transition{
				{
					From: "feature",
					To: []string{
						"master:local",
					},
				},
			},
		},
	}
	flowit.Tags = []Tag{
		{
			ID:     "release-tag",
			Format: "release-[0-9]+",
			Stages: map[string][]string{
				"feature": {
					"start",
				},
			},
			Branches: []string{
				"master",
			},
		},
	}
	flowit.StateMachines = []StateMachine{
		{
			ID: "simple-machine",
			Stages: []string{
				"start", "finish",
			},
			InitialStage: "start",
			FinalStages:  []string{"finish"},
			Transitions: []StateMachineTransition{
				{
					From: []string{"start"},
					To:   []string{"finish"},
				},
			},
		},
	}
	flowit.Workflows = []Workflow{
		{
			ID:           "feature",
			StateMachine: flowit.StateMachines[0].ID,
			Stages: []Stage{
				{
					ID:   "start",
					Args: []string{"< my-var-1 | My-desc-1 >", "< my-var-2 | My-desc-2 >"},
					Conditions: []string{
						"start condition1",
					},
					Actions: []string{"start action1", "start action2"},
				},
				{
					ID:   "finish",
					Args: []string{"< my-var-1 | My-desc-1 >", "< my-var-2 | My-desc-2 >"},
					Conditions: []string{
						"finish condition1",
					},
					Actions: []string{"finish action1", "finish action2"},
				},
			},
		},
	}

	var config WorkflowDefinition
	config.Flowit = flowit
	return config
}
