# flowit [WIP]
A flexible workflow manager

## Overview
`flowit` is a CLI utility that manages user-defined workflows and ensures execution consistency.

## Usage
> `flowit <workflow-id> [workflow-instance-id] <stage-id> [args...]`

## User concepts
From the user perspective, there are only two concepts that need to be understood to make use of `flowit`.

### Workflow
Each worfklow can be defined as any well defined process that has a beginning and an end. Each workflow step is called a `stage` and the rules that explain how stages relate to each other are called `transitions`.

### Stages
Each workflow is comprised of at least two `stages` representing the beginning and the termination of the workflow cycle. Each stage represents a set of commands that are executed once the workflow arrives at that particular `stage`. Each `stage` execution can contain any number of command line arguments.

### Example
In the workflow `life` there are four stages: `birth`, `growth`, `reproduction`, `death`. Where the `reproduction` stage is optional but all the rest of them are mandatory. We can model our `life` workflow as follows:
```yaml
  workflow: life
  stages: [ birth, growth, reprodution, death ]
  transitions: [ birth -> growth, growth -> reproduction, growth -> death, reproduction -> death ]
```

## Workflow designer concepts
From the workflow designer perspective there are some more concepts that need to be understood to be able to design a workflow.

### Workflow definition
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
- `shell`: Location of the executable shell in which the stage `conditions` and `actions` commands will run. It defaults to the default shell. This value is OS dependent.
```yaml
  config:
    abort-on-failed-action: true
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

#### State Machines (Required)
State machines codify the stages and transitions that are going to be allowed as part of a specific workflow. 
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the workflow to refer to this specific state machine.
- `stages` (Required): List of all possible stages
- `initial-stage` (Required): Stage ID of the first stage in the workflow.
- `final-stages` (Required): List of the final stages in the workflow.
- `transitions` (Required): List which represent the relationships between stages.
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
The `simple-machine` state machine above contains `start`, `publish` and `finish` stages where `start` is the initial stage and `finish` is the final stage. The special `!` prefixed to a valid stage indicates that we are excluding that specific stage but considering all the rest of them. That transition rule can be translated as follows: `A transition is allowed from all stages EXCEPT 'finish' to all stages EXCEPT 'start'`.
This means we can transition from `start` to `finish` or from `publish` to `publish` but not from `ANY` stage to `start` or from `finish` to `ANY` stage. Below is the equivalent state machine definition without using the special `!` prefix syntax.
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
Workflows are usually the largest section of the specification. They define the workflows supported, which state machine rules they comform to and exactly how the workflow stages are composed by conditions and actions.
- `id` (Required): This property can be arbitrarily defined by the workflow designer. It is the main handler allowing the CLI to refer to this specific workflow.
- `state-machine` (Required): ID of the state machine which will be used to validate the allowed stages and transitions for this specific workflow instance.
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
- `conditions` (Optional): This section defines a list of commands that will be executed in order before the main stage actions. If any condition fails, the stage actions execution will be aborted. Conditions should avoid altering any state and they should be idempotent operations.
- `actions` (Required): This section defines a list of commands that will be executed in order once the conditions ran succesfully. Actions can alter state and are not required to be idempotent.
```yaml
  ... # workflow definition
  stages:
  - id: start
    args:
    - < feature-branch-suffix | Branch name without prefix >
    actions:
    - git checkout master
    - git pull origin master
    - git checkout -b $<branches[feature].name> master

  - id: publish
    conditions:
    - ./run-tests.sh
    actions:
    - git checkout $<branches[feature].name>
    - git push origin $<branches[feature].name>

  - id: finish
    actions:
    - git checkout master
    - git pull origin master
    - git checkout $<branches[feature].name>
    - git rebase master
    - git checkout master
    - git merge $<branches[feature].name>
```
These stages are part of the `feature` workflow. This means that each stage will be run in the command line as `flowit feature <stage-id>`. We can see in the section above that `feature` workflow referenced `simple-machine` as its state machine and we can see in the state machine definition that `simple-machine` has `start` as the initial stage.

 On the `start` stage definition we can see that there are two arguments defined. This means that in order to start a new `feature` workflow we will need to run `flowit feature start <arg-1>`. `feature-branch-suffix` workflow variable will be set to whatever value of `arg-1` we specify in the command line. This feature will allow the workflow designer to refer to instances of values specified in previous stages without having the need to specify them as arguments in each stage they are needed.
 
 Each of the conditions will be sequentially run and in case of all succeeding, the actions will be performed in the same manner. In case of any action failing, the value of `abort-on-failed-action` will be taken into account in wether or not to abort or continue the stage actions execution. 
 
 One last important thing to note is that for every initial stage command that is run, a new unique workflow instance identifier will be generated so we can reference a specific workflow in case multiple workflows are run in parallel (which is normally the case). In order to run a following allowed stage such as `publish` or `finish`, we should specify the workflow instance ID (short version): `flowit feature <workflow-instance-id> <stage-id> [args...]`.

## Inspiration
This project was inspired on Vincent Driessen's [gitflow](https://github.com/nvie/gitflow) project and it's most active [fork](https://github.com/petervanderdoes/gitflow-avh).
