package utils

import (
	"errors"
	"regexp"
	"strings"
)

const variableNamingRegexPattern = `([a-zA-Z0-9\-\_]+)`
const descriptionNamingRegexPattern = `([A-Z]+[a-zA-Z0-9\-\_ ]*)`
const variableDeclarationRegexPattern = `^< *` + variableNamingRegexPattern + ` *\| *` + descriptionNamingRegexPattern + ` *>$`
const variableReferenceRegexPattern = `\$<` + variableNamingRegexPattern + `>`

// IsValidVariableDeclaration receives a string and returns a boolean value indicating
// whether or not the string is a valid variable declaration expression
func IsValidVariableDeclaration(expression string) bool {
	matched, _ := regexp.Match(variableDeclarationRegexPattern, []byte(expression))
	return matched
}

// DoesExpressionContainsVariableReference receives a string and returns a boolean value indicating
// whether or not the string contains a variable reference
func DoesExpressionContainsVariableReference(expression string) bool {
	matched, _ := regexp.Match(variableReferenceRegexPattern, []byte(expression))
	return matched
}

// ExtractVariableNameFromVariableDeclaration receives a string and returns the string representing
// the variable name in the declaration expression. It returns an error if the expression is not valid
func ExtractVariableNameFromVariableDeclaration(expression string) (string, error) {
	if !IsValidVariableDeclaration(expression) {
		return "", errors.New("Invalid variable declaration:" + expression)
	}
	rx := regexp.MustCompile(variableDeclarationRegexPattern)
	return rx.FindStringSubmatch(expression)[1], nil
}

// EvaluateVariablesInExpression receives an expression and a replacementMap and returns the expression with all
// its variables replaced. It returns an error if the expression does not contain a variable reference or if a variable
// reference is not in the replacement map
func EvaluateVariablesInExpression(expression string, replacementMap map[string]string) (string, error) {
	if !DoesExpressionContainsVariableReference(expression) {
		return expression, nil
	}
	rx := regexp.MustCompile(variableReferenceRegexPattern)
	matches := rx.FindAllStringSubmatch(expression, -1)
	for _, match := range matches {
		if _, ok := replacementMap[match[1]]; !ok {
			return "", errors.New("Variable: " + match[0] + " could not be evaluated")
		}
		expression = strings.ReplaceAll(expression, match[0], replacementMap[match[1]])
	}
	return expression, nil
}
