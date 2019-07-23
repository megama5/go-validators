package rule

import (
	"fmt"
	"log"
	"reflect"
)

type ValidateRule interface {
	Validate(in reflect.Value) ValidateRule
	GetErrorMessage() string
	IsSuccessful() bool
}

type Rule struct {
	ValidateRule
	RuleName             string
	FieldName            string
	ErrorMessage         string
	Restrictions         []string
	MinRestrictionsCount int
	IsValidationFailed   bool
}

func NewRule(fieldName string, restrictions []string) *Rule {
	return &Rule{
		FieldName:    fieldName,
		Restrictions: restrictions,
	}
}

func (r *Rule) ValidationFailed() ValidateRule {
	r.validationFailed()
	return r
}

func (r *Rule) Init(ruleName, errMessage string, minRestrictionsCount int) *Rule {
	r.RuleName = ruleName
	r.ErrorMessage = errMessage
	r.IsValidationFailed = false
	r.MinRestrictionsCount = minRestrictionsCount

	return r
}

func (r *Rule) validationFailed() {
	r.IsValidationFailed = true
}

func (r *Rule) IsSuccessful() bool {
	return !r.IsValidationFailed
}

func (r *Rule) GetErrorMessage() string {
	return fmt.Sprintf("Field '%v' failed to pass %v validation rule ", r.FieldName, r.RuleName)
}

func (r *Rule) Log(message string) {
	log.Printf("Tag-Validator[%v] - %v", r.RuleName, message)
}
