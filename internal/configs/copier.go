package configs

import (
	"bytes"
	"encoding/gob"

	"github.com/pkg/errors"
)

func deepCopy(from interface{}, to interface{}) error {
	// TODO: Check that both inputs are pointers
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	err := enc.Encode(from)
	if err != nil {
		return errors.Wrap(err, "Error encoding config")
	}
	err = dec.Decode(to)
	if err != nil {
		return errors.Wrap(err, "Error decoding config")
	}
	return nil
}
