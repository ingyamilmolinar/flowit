package io

import (
	"fmt"

	"github.com/pkg/errors"
)

type ConsoleWriter struct{}

func NewConsoleWriter() ConsoleWriter {
	return ConsoleWriter{}
}

func (w ConsoleWriter) Write(s string) error {
	return Println(s)
}

// Print receives a list of anything and writes them to standard output. Spaces are added between arguments.
// Returns an error in case of failure
func Print(a ...interface{}) error {
	if _, err := fmt.Print(a...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Println receives a list of anything and writes them to standard output appending a new line.
// Spaces are added between arguments. Returns an error in case of failure
func Println(a ...interface{}) error {
	if _, err := fmt.Println(a...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Printf receives a formatter string and list of anything and writes them to standard output
// following the specified format. Returns an error in case of failure
func Printf(format string, a ...interface{}) error {
	if _, err := fmt.Printf(format, a...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Printfln receives a formatter string and list of anything and writes them to standard output appending a new line
// following the specified format. Returns an error in case of failure
func Printfln(format string, a ...interface{}) error {
	if _, err := fmt.Printf(format+"\n", a...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func PrintSlice(s ...string) error {
	for _, e := range s {
		err := Print(e + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
