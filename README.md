# flowit
A flexible git workflow manager

## Overview
Managing git workflows within a team, project or across a company is hard. In the current software development landscape, there's not an easy way to enforce specific workflow rules other than with intensive engineering training or by introducing tools that make your processes more complex, instead of simplifying them. Ensuring consistency, reducing mental overhead and avoiding bad practices while still keeping the flexibility to adapt your workflow as the project evolves with minimum work is the main goal of this project.

`flowit` uses a declarative appproach to define a workflow. Writting a single file is all that's needed to start enjoying the benefits of a managed workflow. There are some concepts that need to be understood before designing any workflow though.

### Branches
These are all the git branches that will be involved throughout the workflow lifecycle. You can define longliving branches as well as ephemeral ones. In this section, you will need to also define whether or not the branch is protected and in case of ephemeral braches, the different transitions that are allowed. Ex: The `feature` branch can be created only from longliving `master` branch and then it must be merged back to `master` branch. This should happen locally as opposed to via a remote PR process.

### Tags
These are the various types of git tags that will be created at some point in your workflow lifecycle. The branches where the tags are going to be placed, the naming pattern, and when those tags are going to be created can be specified in this section. This section is optional since some workflows do not need tagging support.

### Stages
Stages are the main part of the workflow definition and defines the conditions and actions that will take place in the workflow lifecycle when a command is issued. The stages are divided by short lived branch type. Each stage can accept a variable amount of command line arguments that are completely configurable. A stage is composed by two subsections:

1. conditions: A list of commands that will be executed and when all of them terminate with zero status, the workflow actions will occur.
2. actions: A list of commands that will be executed in order once the conditions were succesfully ran. These are the main commands linked to the particular stage.

There are 5 possible stages in a lifecycle:

1. start: Kicks off the workflow for the specified branch. It will normally create the particular branch based on another but it can perform supporting actions such as pulling the base branch latest changes and transitioning the linked project ticket to the 'In Progress' state.
2. fetch: An alternative way of kicking off the workflow. It will normally fetch and pull the specified branch from a remote repository to the local one. This stage can be useful if a workflow supports collaborative development in a common branch. This is an optional stage.
3. sync: Allows syncing a branch to get the latest remote changes. This stage is useful to update your local branch to contain the most recent merged changes in the base branch. This is an optional stage.
4. publish: Pushes the branch to the remote repository. This can also be used to transition the linked project ticket to the 'Ready for Review' state. This is an optional stage.
5. finish: This stage ends the lifecycle for a particular branch. It usually means that the branch will be merged back to the base branch, the project ticket transitioned to 'Done' or 'Ready for QA', a possible tagging can take place and the local and remote repositories could be cleaned up.

### Config
The workflow designer can tweek Flowit behavoir to address their specific needs.
List of configuration flags:
- abort-on-failed-action [boolean]: Wether or not to abort a workflow stage if an action command returns a non zero status code.
- strict [boolean]: Wether or not to disallow regular git refs modifications. Setting this to `true` will disallow regular git to alter state for the project.
- shell [string]: Location of the executable shell in which the condition and action commands will run.

### Variables
Convenient centralized definition of workflow variables to avoid changing values in multiple places. Variables can be harcoded or read from the environment and can be used anywhere in the workflow.

## Stage lifecycle
 start  --> [sync] --> [publish] --> finish
[fetch] ------^

## Usage
> git flowit <branch-type> <stage> <args...>

## Workflow specification
Check out the SPEC.md specification markdown for all the technical specification details.

## Workflow samples
There are some sample workflow definitions under `samples` directory for some popular workflows such as trunk based development, feature branching and Git Flow.
