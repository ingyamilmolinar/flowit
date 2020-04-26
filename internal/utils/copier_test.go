package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	Describe("Copying a raw struct into a mirrored struct", func() {

		Context("Copying a property with a nil pointer", func() {

			It("should set the zero value equivalent", func() {

				source := struct {
					Bool   *bool
					Int    *int32
					Float  *float32
					Rune   *rune
					String *string
				}{}

				var target struct {
					Bool   bool
					Int    int32
					Float  float32
					Rune   rune
					String string
				}

				err := DeepCopy(source, &target)
				Expect(err).To(BeNil())
				Expect(target.Bool).To(BeZero())
				Expect(target.Int).To(BeZero())
				Expect(target.Float).To(BeZero())
				Expect(target.Rune).To(BeZero())
				Expect(target.String).To(BeZero())

			})
		})
	})

	Describe("Copying a struct into a mirrored raw struct", func() {

		Context("Copying a zero value property", func() {

			It("should set the pointer to the zero value", func() {

				source := struct {
					Bool   bool
					Int    int32
					Float  float32
					Rune   rune
					String string
				}{}

				var target struct {
					Bool   *bool
					Int    *int32
					Float  *float32
					Rune   *rune
					String *string
				}

				err := DeepCopy(source, &target)
				Expect(err).To(BeNil())
				Expect(target.Bool).To(Not(BeNil()))
				Expect(target.Int).To(Not(BeNil()))
				Expect(target.Float).To(Not(BeNil()))
				Expect(target.Rune).To(Not(BeNil()))
				Expect(target.String).To(Not(BeNil()))

				Expect(*target.Bool).To(BeZero())
				Expect(*target.Int).To(BeZero())
				Expect(*target.Float).To(BeZero())
				Expect(*target.Rune).To(BeZero())
				Expect(*target.String).To(BeZero())

			})

		})
	})

})
