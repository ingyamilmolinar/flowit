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

	stage := rawStage{
		"id":      "start",
		"actions": []interface{}{"action1", "action2"},
	}
	workflowType := rawWorkflow{
		"feature": []*rawStage{
			&stage,
		},
	}

	mainConfig := rawMainDefinition{
		Version: &version,
		Branches: []*rawBranch{
			&branch,
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
	flowit.Config = config{
		AbortOnFailedAction: true,
		Strict:              false,
		Shell:               "/usr/bin/env bash",
	}
	flowit.Variables = map[string]interface{}{
		"var1": "value",
		"var2": 12345,
		"var3": "${env-variable}",
	}
	flowit.Branches = []branch{
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
			Transitions: []transition{
				{
					From: "feature",
					To: []string{
						"master:local",
					},
				},
			},
		},
	}
	flowit.Tags = []tag{
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
	flowit.Workflows = []workflow{
		{
			"feature": []stage{
				{
					"id":   "start",
					"args": []string{"arg1", "arg2"},
					"conditions": []string{
						"condition1",
					},
					"actions": []string{"action1", "action2"},
				},
			},
		},
	}

	var config WorkflowDefinition
	config.Flowit = flowit

	return config
}
