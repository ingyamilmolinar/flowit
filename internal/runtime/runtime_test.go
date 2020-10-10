package runtime_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	r "github.com/yamil-rivera/flowit/internal/runtime"
)

type executor struct{}

func (e executor) Execute(command string) (string, error) {
	if command == "FAIL" {
		return command, errors.New("Command failed")
	}
	return command, nil
}

var _ = Describe("Runtime", func() {

	var e executor
	service := r.NewService(e)

	Context("Executing commands", func() {

		It("should execute commands successfully", func() {

			validCommands := []string{
				"COMMAND_1",
			}
			out, errIdx, err := service.Execute(validCommands, nil, 0)
			Expect(err).To(BeNil())
			Expect(errIdx).To(BeZero())
			Expect(out).To(BeEquivalentTo([]string{"COMMAND_1"}))

		})

		It("should return an error if a command fails", func() {

			invalidCommands := []string{
				"COMMAND_1",
				"FAIL",
			}
			out, errIdx, err := service.Execute(invalidCommands, nil, 0)
			Expect(err).ToNot(BeNil())
			Expect(errIdx).To(Equal(1))
			Expect(out).To(BeEquivalentTo([]string{"COMMAND_1", "FAIL"}))

		})

		It("should evaluate variables successfully", func() {

			commands := []string{
				"$<variable-1>",
				"$<variable-2>",
			}
			variables := map[string]interface{}{
				"variable-1": "value1",
				"variable-2": "value2",
			}
			out, errIdx, err := service.Execute(commands, variables, 0)
			Expect(err).To(BeNil())
			Expect(errIdx).To(BeZero())
			Expect(out).To(BeEquivalentTo([]string{"value1", "value2"}))

		})

		It("should return an error if a variable does not exist", func() {

			commands := []string{
				"$<variable-1> $<variable-2>",
			}
			variables := map[string]interface{}{
				"variable-1": "value",
			}
			out, errIdx, err := service.Execute(commands, variables, 0)
			Expect(err).ToNot(BeNil())
			Expect(errIdx).To(BeZero())
			Expect(len(out)).To(Equal(0))

		})

	})

})
