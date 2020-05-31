package repository

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/pkg/errors"
	"github.com/xujiajun/nutsdb"
	"github.com/yamil-rivera/flowit/internal/io"
	"github.com/yamil-rivera/flowit/internal/models"
)

var db *nutsdb.DB

// TODO: This file has to be created on the project root (along with .git)
func openDB() error {
	opt := nutsdb.DefaultOptions
	opt.Dir = ".flowitDS"
	openedDB, err := nutsdb.Open(opt)
	if err != nil {
		return errors.WithStack(err)
	}
	db = openedDB
	return nil
}

func closeDB() {
	if err := db.Close(); err != nil {
		io.Logger.Errorf("%+v", err)
	}
}

// DeleteDB wipes the DB clean
func DeleteDB() error {
	return io.RemoveDirectory(".flowitDS")
}

// PutWorkflow takes a models.Workflow struct and saves it into the DB
func PutWorkflow(workflow models.Workflow) error {
	bytes, err := encodeWorkflow(workflow)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := openDB(); err != nil {
		return errors.WithStack(err)
	}
	defer closeDB()

	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put("workflows_"+workflow.DefinitionID, []byte(workflow.ID), bytes, 0); err != nil {
				return errors.Wrap(err, "Error trying to save workflow")
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "Error trying to open update transaction")
	}
	return nil
}

// UpdateWorkflow takes a models.Workflow struct and saves it into the DB
// It overrides it if an existing struct with same ID is found
func UpdateWorkflow(workflow models.Workflow) error {
	return PutWorkflow(workflow)
}

// DeleteWorkflow takes a definitionID and workflowID and removes the workflow from the DB
// If the workflow or bucket does not exist, an error is returned
func DeleteWorkflow(definitionID, workflowID string) error {
	if err := openDB(); err != nil {
		return errors.WithStack(err)
	}
	defer closeDB()

	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			if _, err := tx.Get("workflows_"+definitionID, []byte(workflowID)); err != nil {
				return errors.Wrap(err, "Error trying to delete workflow")
			}
			if err := tx.Delete("workflows_"+definitionID, []byte(workflowID)); err != nil {
				return errors.Wrap(err, "Error trying to delete workflow")
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "Error trying to open update transaction")
	}
	return nil
}

// GetWorkflowFromPreffix takes a definitionID and workflowPreffix and returns a workflow
// which ID begins with the preffix wrapped in an optional.
// If no workflow is found or the bucket does not exist, an empty optional is returned
func GetWorkflowFromPreffix(definitionID, workflowPreffix string) (models.OptionalWorkflow, error) {
	if err := openDB(); err != nil {
		return models.OptionalWorkflow{}, errors.WithStack(err)
	}
	defer closeDB()

	var workflow models.Workflow
	workflowSet := false
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.PrefixScan("workflows_"+definitionID, []byte(workflowPreffix), 1)
			if err != nil {
				return nil
			}
			w, err := decodeWorkflow(entries[0].Value)
			if err != nil {
				return errors.WithStack(err)
			}
			workflow = *w
			workflowSet = true
			return nil
		}); err != nil {
		return models.OptionalWorkflow{}, errors.Wrap(err, "Error within happened within transaction")
	}
	if workflowSet {
		return models.NewWorkflow(workflow), nil
	}
	return models.OptionalWorkflow{}, nil
}

// GetWorkflow takes a definitionID and workflowID and returns the workflow which ID exactly matches the workflowID
// wrapped in an optional.
// If no workflow is found or the bucket does not exist, an empty optional is returned
func GetWorkflow(definitionID, workflowID string) (models.OptionalWorkflow, error) {
	if err := openDB(); err != nil {
		return models.OptionalWorkflow{}, errors.WithStack(err)
	}
	defer closeDB()

	var workflow models.Workflow
	workflowSet := false
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entry, err := tx.Get("workflows_"+definitionID, []byte(workflowID))
			if err != nil {
				return nil
			}
			w, err := decodeWorkflow(entry.Value)
			if err != nil {
				return errors.WithStack(err)
			}
			workflow = *w
			workflowSet = true
			return nil
		}); err != nil {
		return models.OptionalWorkflow{}, errors.Wrap(err, "Error within happened within transaction")
	}
	if workflowSet {
		return models.NewWorkflow(workflow), nil
	}
	return models.OptionalWorkflow{}, nil
}

// GetWorkflows takes a definitionID, an integer 'n' and whether or not inactive workflows are excluded
// and returns a list of 'n' workflows that match the criteria.
// If n is 0, all existing workflows that match the criteria are returned
// TODO: Make the filtering flexible
func GetWorkflows(definitionID string, n int, excludeInactive bool) ([]models.Workflow, error) {
	if n < 0 {
		return nil, nil
	}
	if err := openDB(); err != nil {
		return nil, errors.WithStack(err)
	}
	defer closeDB()

	var workflows []models.Workflow
	var copied bool
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll("workflows_" + definitionID)
			if err != nil {
				return nil
			}
			count := n
			numberOfEntries := len(entries)
			if numberOfEntries > 0 {
				for i := 0; i < numberOfEntries; i++ {
					w, err := decodeWorkflow(entries[i].Value)
					if err != nil {
						return errors.WithStack(err)
					}
					if (!excludeInactive || w.IsActive) && (n == 0 || count > 0) {
						workflows = append(workflows, *w)
						count--
						copied = true
					}
				}
			}
			if n > 0 && (numberOfEntries-count) < n {
				return errors.New("Number of requested workflows: " + strconv.Itoa(n) +
					" is greater than the number of items: " + strconv.Itoa(numberOfEntries-count))
			}
			return nil
		}); err != nil {
		return workflows, errors.Wrap(err, "Error within happened within transaction")
	}
	if copied {
		return workflows, nil
	}
	return nil, nil
}

func decodeWorkflow(buf []byte) (*models.Workflow, error) {
	var target models.Workflow
	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&target); err != nil {
		return nil, errors.Wrap(err, "Error trying to decode workflow")
	}
	return &target, nil
}

func encodeWorkflow(source models.Workflow) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(source); err != nil {
		return nil, errors.Wrap(err, "Error trying to encode workflow")
	}
	return buf.Bytes(), nil
}
