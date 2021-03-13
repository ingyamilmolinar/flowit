package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"
)

func Sha256(o interface{}) ([]byte, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(fmt.Sprintf("%+v", o))); err != nil {
		return nil, errors.WithStack(err)
	}
	return hash.Sum(nil), nil
}
