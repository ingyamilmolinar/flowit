package utils

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// DeepCopy copies recursively an object into another using encoding/json
// encoding/gob was not used even though is faster because it does not preserve zero values
// See: https://github.com/golang/go/issues/4609
func DeepCopy(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		fmt.Println(err.Error())
		return errors.Wrap(err, "Deep copy error while marshalling")
	}
	if err := json.Unmarshal(bytes, to); err != nil {
		return errors.Wrap(err, "Deep copy error while unmarshalling")
	}
	return nil
}
