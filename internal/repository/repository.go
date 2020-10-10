package repository

import (
	"bytes"
	"encoding/gob"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/io"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

// Service is the data structure from which to use the persistance methods
type Service struct{}

// NewService creates and returns a Service instance
// TODO: Make DB location configurable
func NewService() *Service {
	return &Service{}
}

// Drop wipes the DB clean
func (rs Service) Drop() error {
	db, err := openDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer closeDB(db)

	if err := db.Update(
		func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				return tx.DeleteBucket(name)
			})
		}); err != nil {
		return errors.Wrap(err, "Error trying to open update transaction")
	}
	return os.RemoveAll(dbLocation())
}

// PutWorkflow takes a workflow.Workflow struct and saves it into the DB
func (rs Service) PutWorkflow(workflow w.Workflow) error {
	db, err := openDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer closeDB(db)

	bytes, err := encodeWorkflow(workflow)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := db.Update(
		func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("workflows_" + workflow.Name))
			if err != nil {
				return errors.WithStack(err)
			}
			if err := bucket.Put([]byte(workflow.ID), bytes); err != nil {
				return errors.Wrap(err, "Error trying to save workflow")
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "Error trying to open update transaction")
	}
	return nil
}

// DeleteWorkflow takes a workflowName and workflowID and removes the workflow from the DB
// If the workflow or bucket does not exist, an error is returned
func (rs Service) DeleteWorkflow(workflowName, workflowID string) error {
	db, err := openDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer closeDB(db)

	if err := db.Update(
		func(tx *bolt.Tx) error {
			bucketName := "workflows_" + workflowName
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return errors.New("Bucket " + bucketName + " does not exist")
			}
			if err := b.Delete([]byte(workflowID)); err != nil {
				return errors.Wrap(err, "Error trying to delete workflow")
			}
			return nil
		}); err != nil {
		return errors.Wrap(err, "Error trying to open update transaction")
	}
	return nil
}

// GetWorkflowFromPreffix takes a workflowName and workflowPreffix and returns a workflow
// which ID begins with the preffix wrapped in an optional.
// If no workflow is found or the bucket does not exist, an empty optional is returned
// TODO: Type alias
func (rs Service) GetWorkflowFromPreffix(workflowName, workflowPreffix string) (w.OptionalWorkflow, error) {
	db, err := openDB()
	if err != nil {
		return w.OptionalWorkflow{}, errors.WithStack(err)
	}
	defer closeDB(db)

	var workflow w.Workflow
	workflowSet := false
	if err := db.View(
		func(tx *bolt.Tx) error {
			bucketName := "workflows_" + workflowName
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return errors.New("Bucket " + bucketName + " does not exist")
			}
			c := b.Cursor()
			for k, v := c.Seek([]byte(workflowPreffix)); k != nil && bytes.HasPrefix(k, []byte(workflowPreffix)); k, v = c.Next() {
				w, err := decodeWorkflow(v)
				if err != nil {
					return errors.WithStack(err)
				}
				workflow = *w
				workflowSet = true
				break
			}
			return nil
		}); err != nil {
		return w.OptionalWorkflow{}, errors.Wrap(err, "Error within happened within transaction")
	}
	if workflowSet {
		return w.NewWorkflowOptional(workflow), nil
	}
	return w.OptionalWorkflow{}, nil
}

// GetWorkflow takes a workflowName and workflowID and returns the workflow which ID exactly matches the workflowID
// wrapped in an optional.
// If no workflow is found or the bucket does not exist, an empty optional is returned
func (rs Service) GetWorkflow(workflowName, workflowID string) (w.OptionalWorkflow, error) {
	db, err := openDB()
	if err != nil {
		return w.OptionalWorkflow{}, errors.WithStack(err)
	}
	defer closeDB(db)

	var workflow w.Workflow
	workflowSet := false
	if err := db.View(
		func(tx *bolt.Tx) error {
			bucketName := "workflows_" + workflowName
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return nil
			}
			entry := b.Get([]byte(workflowID))
			if entry == nil {
				return nil
			}
			w, err := decodeWorkflow(entry)
			if err != nil {
				return errors.WithStack(err)
			}
			workflow = *w
			workflowSet = true
			return nil
		}); err != nil {
		return w.OptionalWorkflow{}, errors.Wrap(err, "Error within happened within transaction")
	}
	if workflowSet {
		return w.NewWorkflowOptional(workflow), nil
	}
	return w.OptionalWorkflow{}, nil
}

// GetWorkflows takes a workflowName, an integer 'n' and whether or not inactive workflows are excluded
// and returns a list of 'n' workflows that match the criteria.
// If n is 0, all existing workflows that match the criteria are returned
func (rs Service) GetWorkflows(workflowName string, n int, excludeInactive bool) ([]w.Workflow, error) {
	db, err := openDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer closeDB(db)

	if n < 0 {
		return nil, nil
	}

	var workflows []w.Workflow
	var copied bool
	if err := db.View(
		func(tx *bolt.Tx) error {
			bucketName := "workflows_" + workflowName
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return nil
			}
			count := n
			if err := b.ForEach(func(k, v []byte) error {
				w, err := decodeWorkflow(v)
				if err != nil {
					return errors.WithStack(err)
				}
				if (!excludeInactive || w.IsActive) && (n == 0 || count > 0) {
					workflows = append(workflows, *w)
					count--
					copied = true
				}
				return nil
			}); err != nil {
				return errors.WithStack(err)
			}
			if n > 0 && count > 0 {
				return errors.New("Number of requested workflows: " + strconv.Itoa(n) +
					" is greater than the number of items: " + strconv.Itoa(n-count))
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

func dbLocation() string {
	return ".flowitDS"
}

func openDB() (*bolt.DB, error) {
	dbPath := dbLocation()
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 0})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

func closeDB(db *bolt.DB) {
	if err := db.Close(); err != nil {
		io.Logger.Errorf("%+v", err)
	}
}

func decodeWorkflow(buf []byte) (*w.Workflow, error) {
	var target w.Workflow
	dec := gob.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&target); err != nil {
		return nil, errors.Wrap(err, "Error trying to decode workflow")
	}
	return &target, nil
}

func encodeWorkflow(source w.Workflow) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(source); err != nil {
		return nil, errors.Wrap(err, "Error trying to encode workflow")
	}
	return buf.Bytes(), nil
}
