package config

func validConfigJustMandatoryFields() rawFlowitConfig {

	var flowit rawFlowitConfig

	version := "0.1"
	flowit.Version = &version

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
	workflowType := rawWorkflowType{
		"feature": []*rawStage{
			&stage,
		},
	}

	workflow := rawWorkflow{
		Branches: []*rawBranch{
			&branch,
		},
		Stages: []*rawWorkflowType{
			&workflowType,
		},
	}
	flowit.Workflow = &workflow

	return flowit
}

func validConfigWithOptionalFields() FlowitConfig {

	var flowit FlowitConfig

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
	flowit.Workflow.Branches = []branch{
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
	flowit.Workflow.Tags = []tag{
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
	flowit.Workflow.Stages = []workflowType{
		{
			"feature": []stage{
				{
					"id":   "start",
					"args": "arg1 arg2",
					"conditions": []string{
						"condition1",
					},
					"actions": []string{"action1", "action2"},
				},
			},
		},
	}

	return flowit
}
