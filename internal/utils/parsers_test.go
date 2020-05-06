package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	Describe("Validating variable definition matcher", func() {

		It("should only match a valid variable definition", func() {

			valid := IsValidVariableDeclaration("<my-var-with-spaces | Desc>")
			Expect(valid).To(BeTrue())
			valid = IsValidVariableDeclaration("<my-var-without-spaces|Desc>")
			Expect(valid).To(BeTrue())
			valid = IsValidVariableDeclaration("<my-var-without-spaces| Desc with spaces >")
			Expect(valid).To(BeTrue())
			valid = IsValidVariableDeclaration("< myVarWithCamelCase | Desc >")
			Expect(valid).To(BeTrue())
			valid = IsValidVariableDeclaration("< myVarWithCamelCaseAndNumbers2 | Desc >")
			Expect(valid).To(BeTrue())

			valid = IsValidVariableDeclaration("<my-var-without-closing-bracket")
			Expect(valid).To(BeFalse())
			valid = IsValidVariableDeclaration("{my-var-with-different-bracket}")
			Expect(valid).To(BeFalse())
			valid = IsValidVariableDeclaration("< my var with spaces | Desc >")
			Expect(valid).To(BeFalse())
			valid = IsValidVariableDeclaration("< my-var-without-desc >")
			Expect(valid).To(BeFalse())
			valid = IsValidVariableDeclaration("< my-var-without-desc | >")
			Expect(valid).To(BeFalse())
			valid = IsValidVariableDeclaration("< my-var | desc not starting with uppercase >")
			Expect(valid).To(BeFalse())

		})

	})

	Describe("Validating variable expression matcher", func() {

		It("should only match a valid variable reference", func() {

			valid := DoesExpressionContainsVariableReference("my var = $<my-var>")
			Expect(valid).To(BeTrue())
			valid = DoesExpressionContainsVariableReference("$<my-var-1> $<my-var-2>")
			Expect(valid).To(BeTrue())
			valid = DoesExpressionContainsVariableReference("$<myVarWithCamelCaseAndNumbers1>")
			Expect(valid).To(BeTrue())

			valid = DoesExpressionContainsVariableReference("my var = <my-var-without-$>")
			Expect(valid).To(BeFalse())
			valid = DoesExpressionContainsVariableReference("my var = ${my-var-with-different-brackets}")
			Expect(valid).To(BeFalse())
			valid = DoesExpressionContainsVariableReference("my var = $<my var with spaces>")
			Expect(valid).To(BeFalse())

		})

	})

	Describe("Validating variable extraction in variable definition", func() {

		It("should extract variable in variable definition", func() {

			variable, err := ExtractVariableNameFromVariableDeclaration("< my-var | Description >")
			Expect(variable).To(BeIdenticalTo("my-var"))
			Expect(err).To(BeNil())

			variable, err = ExtractVariableNameFromVariableDeclaration("<myVar| Description >")
			Expect(variable).To(BeIdenticalTo("myVar"))
			Expect(err).To(BeNil())

			variable, err = ExtractVariableNameFromVariableDeclaration("< myVar1 | Description >")
			Expect(variable).To(BeIdenticalTo("myVar1"))
			Expect(err).To(BeNil())

			variable, err = ExtractVariableNameFromVariableDeclaration("< myVar1,myVar2 | Description >")
			Expect(variable).To(BeZero())
			Expect(err).To(Not(BeNil()))

			variable, err = ExtractVariableNameFromVariableDeclaration("< myVar | description >")
			Expect(variable).To(BeZero())
			Expect(err).To(Not(BeNil()))

		})

	})

	Describe("Validating variable evaluation from map", func() {

		It("should evaluate every variable defined in map", func() {

			replacementMap := map[string]string{
				"my-var":   "value!",
				"my-var-2": "another value!",
			}
			expression, err := EvaluateVariablesInExpression("my var = $<my-var>", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = value!"))
			Expect(err).To(BeNil())

			expression, err = EvaluateVariablesInExpression("my var = $<my-var>; my var 2 = $<my-var-2>", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = value!; my var 2 = another value!"))
			Expect(err).To(BeNil())

			expression, err = EvaluateVariablesInExpression("my var = $<my-var>; my var 2 = $<my-var>", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = value!; my var 2 = value!"))
			Expect(err).To(BeNil())

			expression, err = EvaluateVariablesInExpression("my var = <my-var>; my var 2 = $<my-var>", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = <my-var>; my var 2 = value!"))
			Expect(err).To(BeNil())

			// TODO: Should we raise an issue on this case? It is most likely a user error
			expression, err = EvaluateVariablesInExpression("my var = $<my-var>; my var 2 = $< my-var >", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = value!; my var 2 = $< my-var >"))
			Expect(err).To(BeNil())

			expression, err = EvaluateVariablesInExpression("my var = no variables here!", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = no variables here!"))
			Expect(err).To(BeNil())

			expression, err = EvaluateVariablesInExpression("my var = $<my-var>; my var 2 = $<my-unknown-var>", replacementMap)
			Expect(expression).To(BeIdenticalTo(""))
			Expect(err).ToNot(BeNil())

			replacementMap = map[string]string{}
			expression, err = EvaluateVariablesInExpression("my var = no variables here!", replacementMap)
			Expect(expression).To(BeIdenticalTo("my var = no variables here!"))
			Expect(err).To(BeNil())

		})

	})

})
