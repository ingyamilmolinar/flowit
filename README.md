# flowit [WIP]
A flexible git workflow manager

## Background
Enforcing git workflows within a team, project or across a company is hard. In the current software development landscape, there's not an easy way to enforce specific workflow rules other than with diligent training or by introducing bloated tools that make your processes more complex, instead of simplifying them. Ensuring consistency, reducing mental overhead and avoiding bad practices while still keeping the flexibility to adapt your git workflow as the project evolves with minimum work is the main goal of this project.

## Overview
`flowit` uses a declarative approach to define a git workflow. Writting a single `yaml` file is all that's needed to start enjoying the benefits of a managed workflow. This file tells `flowit` everything it needs to enforce and execute the rules of your workflow whatever those might be. The CLI integrates nicely with `git` so users don't need to "leave" git to make use of it.

## Usage
> `git flowit <workflow-id> [workflow-instance-id] <stage-id> [args...]`

## User concepts
From the user perspective, there are only two concepts that need to be understood to make use of the tool.

### Workflow
Each worfklow refers to one full development cycle, it is usually tied to a branch type but it can be defined as any well defined process that has a beginning and an end. Each workflow step is called a `stage` and the rules on how stages are linked together are called `transitions`.

### Stages
Each workflow is comprised of at least two `stages` representing the beginning and the termination of the workflow cycle. Each stage represents a set of commands that are executed once the workflow arrives at that particular stage. Each stage can contain a previously defined number of required or optional command line arguments that will be referenced on the commands themselves.

### Example
In the workflow `life` there are four stages: `birth`, `growth`, `reproduction`, `death`. Where `reproduction` stage is optional but all the rest of them are mandatory. We can model `life` workflow as follows:
```yaml
  workflow: life
  stages: [ birth, growth, reprodution, death ]
  transitions: [ birth -> growth, growth -> reproduction, growth -> death, reproduction -> death ]
```

## Designer concepts
From the workflow designer perspective there are several concepts that need to be understood to be able to design a robust workflow.

### Flowit workflow examples
Before digging into the nitty gritty details, there are some already designed workflow definitions under the `samples` directory that show how some of the most popular git workflows could be defined in `flowit`. Feel free to modify them as needed to suit your specific project needs.

### Flowit workflow definition
```yaml
flowit:
# Each section purpose and contents will be explained below
```

#### Version (Required)
Number describing to which specification version this particular workflow definition is complying to. The current version is `0.1`.
```yaml
  version: "0.1"
```

#### Config (Optional)
The workflow designer can tweek `flowit` behavior to address their specific needs.
- `abort-on-failed-action`: Wether or not to abort a workflow stage if an action command returns a non zero status code. The default is `true`.
- `strict`: Wether or not to disallow regular git refs modification. Setting this to `true` will disallow git to alter state for the project. The default is `false`.
- `shell`: Location of the executable shell in which the stage `conditions` and `actions` commands will run. It defaults to the default shell. This value is OS dependent.
```yaml
  config:
    abort-on-failed-action: true
    strict: true
    shell: /usr/bin/env bash
```

#### Variables (Optional)
Convenient centralized definition of workflow variables. These can be harcoded or read from the environment and can be used anywhere in the workflow.
```yaml
  variables:
    circleci-username: ${CIRCLECI_USERNAME}
    circleci-project-name: ${CIRCLECI_PROJECT_NAME}
    circleci-token: ${CIRCLECI_TOKEN}
```

#### Branches (Required)
Definition of all the git branches that will be involved throughout the workflow lifecycle. Each branch can contain the following properties:
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the workflow to refer to this specific branch.
- `name` (Required): This property is the actual branch name, for long living branches, this is a fixed string like `master` or `develop`. For ephemeral branches this is usually a variable or group of variables that represent the branch name and that will be set by the user whenever the branch creation stage is run.
- `eternal` (Required): Wether or not this branch is long lived or not.
- `protected` (Required): Wether or not this branch is allowed to be committed on directly.
- `transitions` (Required if `eternal` is `false`): A list of transitions specifying the allowed base branches that his particular branch can branch from and the target branches where this branch is allowed to merge back to.
```yaml
  branches:
  - id: master
    name: master
    eternal: true
    protected: true

  - id: feature
    name: feature/$<jira-issue-id>/$<feature-branch-suffix>
    eternal: false
    protected: false
    transitions:
    - from: master
      to: [ master:local ]
```
The `feature` branch above name will be defined at runtime from the variables `jira-issue-id` and `feature-branch-suffix` and can be only created from longliving and protected `master` branch and then it must be merged back to `master` branch. This transition should happen locally as opposed to going through a PR process.

#### Tags (Optional)
These are the various types of git tags that will be created at some point in your workflow lifecycle. 
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the workflow to refer to this specific tag.
- `format` (Required): A valid regular expression that defines the valid tag name format.
- `stages` (Optional): A map where each key is a valid workflow ID and the value is a list of stage IDs in which the tagging action is permitted
- `branches` (Required): A list of branch IDs in which the tagging is permitted
```yaml
  tags:
  - id: release
    format: "[0-9]+\.[0-9]+\.[0-9]+"
    stages:
      release: [ finish ]
      hotfix: [ finish ]
    branches: [ master ]
```
The `release` tag above which naming format complies with a semantic versioning can only be set during the `release finish` and `hotfix finish` stages and can only be set in `master` branch.

#### State Machines (Required)
State machines codify the stages and transitions that are going to be allowed as part of a specific workflow. 
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the workflow to refer to this specific state machine.
- `stages` (Required): List of all possible stage IDs
- `initial-stage` (Required): Stage ID of the first stage in the workflow.
- `final-stages` (Required): List of stage IDs which are the end stages of the workflow.
- `transitions` (Required): List which represent the allowed transitions between stages.
```yaml
  state-machines:
  - id: simple-machine
    stages: [ start, publish, finish ]
    initial-stage: start
    final-stages: [ finish ]
    transitions:
    - from: [ "!finish" ]
      to: [ "!start" ]
```
The `simple-machine` state machine above contains `start`, `publish` and `finish` stages where `start` is the initial stage and `finish` is the final stage. The special `!` prefixed to a valid stage indicates that we are excluding that specific stage but considering all the rest of them. That transition rule can be translated as follows: "A transition is allowed from `ALL EXCEPT finish` stages to `ALL EXCEPT start` stages". This means we can transition from `start` to `finish` or from `publish` to itself but not from `ANY` stage to `start` or from `finish` to `ANY` stage. Below is the equivalent state machine definition without using the special `!` prefix syntax.
```yaml
  state-machines:
  - id: simple-machine
    stages: [ start, publish, finish ]
    initial-stage: start
    final-stages: [ finish ]
    transitions:
    - from: [ start ]
      to: [ publish, finish ]
    - from: [ publish ]
      to: [ publish, finish ]
```

#### Workflows (Required)
Workflows are the main section of the specification. They define the amount of workflows supported, which state machine rules they comform to and exactly how the workflow stages are defined.
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the workflow definition to refer to this specific workflow.
- `state-machine` (Required): ID of the state machine which will be used to validate the allowed stages and transitions for this specific workflow instace.
- `stages` (Required): List of stages that make up the workflow. The stage IDs should match the referenced state machine stage list.
```yaml
  workflows:
  - id: feature
    state-machine: simple-machine
    stages:
    - ... # This will be explained in detail in the following subsection
```

##### Stages (Required)
Stages define the conditions and actions that will take place in the workflow lifecycle when a command is issued.
- `args` (Optional): This section defines the number of arguments a specific command will accept and which workflow variables they will populate.
- `conditions` (Optional): This section defines a list of commands that will be executed in order before the main stage actions. If any condition fails, the stage actions execution will be aborted.
- `actions` (Required): This section defines a list of commands that will be executed in order once the conditions ran succesfully.
```yaml
  ... # workflow definition
  stages:
  - id: start
    args:
    - < feature-branch-suffix | Branch name without prefix >
    - < jira-issue-id | Related Jira Issue ID >
    conditions:
    - "[[ $(jira list --status $<jira-issue-id>) == *'Open'* ]]"
    actions:
    - git checkout develop
    - git pull origin develop
    - git checkout -b $<branches[feature].name> develop
    - jira transition $<jira-issue-id> 'In progress'

  - id: publish
    conditions:
    - ./run-tests.sh
    - "[[ $(jira list --status $<jira-issue-id>) == *'In Progress'* ]]"
    actions:
    - git checkout $<branches[feature].name>
    - git push origin $<branches[feature].name>
    - jira transition $<jira-issue-id> 'In code review'

  - id: finish
    conditions:
    - "[[ $(curl https://circleci.com/api/v1.1/project/github/$<circleci-username>/$<circleci-project-name>?circle-token=$<circleci-token>) == *'Passed'* ]]"
    - "[[ $(hub pr list --base $<branches[feature].name>) == *'Merged'* ]]"
    actions:
    - jira transition $<jira-issue-id> 'Done'
    - git checkout develop
    - git pull origin develop
    - git branch -D $<branches[feature].name>
    - git push --delete origin $<branches[feature].name>
```
These stages are part of the `feature` workflow. This means that each stage will be run in the command line as `git flowit feature <stage-id>`. We can see in the section above that `feature` workflow referenced `simple-machine` as its state machine and we can see in the state machine definition that `simple-machine` has `start` as the initial stage.

 On the `start` stage definition we can see that there are two arguments defined. This means that in order to start a new `feature` workflow we will need to run `git flowit feature start <arg-1> <arg-2>`. `feature-branch-suffix` and `jira-issue-id` workflow variables will be set to whatever value of `arg-1` and `arg-2` we specify in the command line. This feature will allow the workflow designer to refer to instances of values specified in previous stages without having the need to specify them as arguments in each stage they are needed.
 
 Each of the conditions will be sequentially run and in case of all succeeding, the actions will be performed in the same manner. In case of any action failing, the value of `abort-on-failed-action` will be taken into account in wether or not to abort or continue the stage execution. 
 
 One last important thing to note is that for every initial stage command that is run, a new unique workflow instance identifier will be generated so we can reference a specific workflow in case multiple workflows are run in parallel (which is normally the case). In order to run a following allowed stage such as `publish` or `finish`, we should specify the workflow instance ID like: `git flowit feature <workflow-instance-id> <stage-id> [args...]`.

## Inspiration
This project was inspired heavily on Vincent Driessen's [gitflow](https://github.com/nvie/gitflow) project and it's most active [fork](https://github.com/petervanderdoes/gitflow-avh) maintained by Peter van der Does.

