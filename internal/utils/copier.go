package utils

import (
	"bytes"
	"encoding/gob"

	"github.com/pkg/errors"
)

// DeepCopy copies recursively an object into another
func DeepCopy(from interface{}, to interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(from); err != nil {
		return errors.Wrap(err, "Error encoding config")
	}
	if err := dec.Decode(to); err != nil {
		return errors.Wrap(err, "Error decoding config")
	}
	return nil
}
