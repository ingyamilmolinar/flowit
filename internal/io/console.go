package io

import (
	"fmt"

	"github.com/pkg/errors"
)

// ConsoleWriter exposes useful methods to interact with the console
type ConsoleWriter struct{}

// NewConsoleWriter returns a new ConsoleWriter instance
func NewConsoleWriter() ConsoleWriter {
	return ConsoleWriter{}
}

// Write receives a string and prints it to the console including a new line
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
