package configs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configs", func() {

	Describe("Validating configuration", func() {

		Context("Validating a valid configuration", func() {

			It("should return a nil error", func() {

				var flowit Flowit

				version := "0.1"
				branchID := "master"
				branchName := "master"
				branchEternal := true
				branchProtected := true
				workflowBranch := branch{
					ID:        &branchID,
					Name:      &branchName,
					Eternal:   &branchEternal,
					Protected: &branchProtected,
				}
				workflow := workflow{
					Branches: []*branch{
						&workflowBranch,
					},
				}

				flowit.Version = &version
				flowit.Workflow = &workflow
				err := validateConfig(&flowit)
				Expect(err).To(BeNil())

			})

		})

	})
})
