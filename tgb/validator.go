package tgb

import (
	"errors"
	"fmt"
	"github.com/megama5/go-validators/tgb/rule"
	"log"
	"reflect"
	"strings"
)

type Reflection struct {
	TypeOf  reflect.Type
	ValueOf reflect.Value
}

type Validator struct {
	rules map[string]func(base *rule.Rule) rule.ValidateRule
}

func NewValidator(rules map[string]func(base *rule.Rule) rule.ValidateRule) *Validator {
	return &Validator{
		rules: rules,
	}
}

func (v *Validator) Validate(in interface{}, requiredData map[string]interface{}) (bool, error) {
	a, b := v.validate(in, "")

	err := v.required(requiredData)
	if err != nil {
		b = append(b.(ErrMessagesList), err.(*ErrMessage))
	}

	return a, b
}

func (v *Validator) required(data map[string]interface{}) error {
	for k, v := range data {
		switch tt := v.(type) {
		case string:
			if len(tt) == 0 {
				return &ErrMessage{
					FieldName: k,
					RuleName:  "required",
					Message:   "No field provided",
				}
			}
		}
	}

	return nil
}

func (v *Validator) validate(in interface{}, parentName string) (bool, error) {
	// will holds field name of nested structure for the recursion call
	var nestedFieldNames []string

	// holder for the rule found in type description
	calculatedRules := make(map[string]map[string][]string, 0)

	// indicate of any error occurred
	errIndicator := false

	// holds all error messages
	errMessages := make(ErrMessagesList, 0)

	// if passed nil the validation is failed
	if in == nil {
		return false, errMessages
	}

	reflectObject, err := v.prepareReflection(in)
	if err != nil {
		log.Printf("Validator err: %v", err)
		return false, errMessages
	}

	v.prepareRules(reflectObject, calculatedRules, &nestedFieldNames)
	v.applyRules(calculatedRules, reflectObject, parentName, &errIndicator, &errMessages)
	v.applyRuleToNestedStructures(nestedFieldNames, reflectObject, parentName, &errIndicator, &errMessages)

	return !errIndicator, errMessages
}

func (v *Validator) applyRuleToNestedStructures(
	nestedFieldNames []string,
	reflectObject *Reflection,
	parentName string,
	errIndicator *bool,
	errMessages *ErrMessagesList,
) {
	parentName = fmt.Sprintf("%s%s", v.prepareParentName(parentName), reflectObject.TypeOf.Name())
	// checking nested structures
	for _, name := range nestedFieldNames {
		field := reflectObject.ValueOf.FieldByName(name)
		if field.IsValid() && ((field.Kind() == reflect.Ptr && !field.IsNil()) || field.Kind() == reflect.Struct) {
			res, list := v.validate(field.Interface(), parentName)
			*errIndicator = *errIndicator && res
			*errMessages = append(*errMessages, list.(ErrMessagesList)...)
		}
	}
}

func (v *Validator) applyRules(
	calculatedRules map[string]map[string][]string,
	reflectObject *Reflection,
	parentName string,
	errIndicator *bool,
	errMessages *ErrMessagesList,
) {
	parentName = v.prepareParentName(parentName)
	// now we have calculatedRules for the structure related to the
	// fields and we need to get a field and validate it
	for fieldName, rules := range calculatedRules {

		field := reflectObject.ValueOf.FieldByName(fieldName)
		for ruleName, params := range rules {

			if ruleF, ruleFound := v.rules[ruleName]; ruleFound {

				result := ruleF(rule.NewRule(fieldName, params)).Validate(field)

				if !result.IsSuccessful() {
					*errIndicator = true
					*errMessages = append(*errMessages, &ErrMessage{
						FieldName: fmt.Sprintf("%s%s.%s", parentName, reflectObject.TypeOf.Name(), fieldName),
						RuleName:  ruleName,
						Message:   result.GetErrorMessage(),
					})
				}
			}
		}
	}
}

func (v *Validator) prepareParentName(parentName string) string {
	if parentName != "" {
		parentName = fmt.Sprintf("%v.", parentName)
	}
	return parentName
}

func (v *Validator) prepareReflection(in interface{}) (*Reflection, error) {
	typeOf := reflect.TypeOf(in)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	valueOf := reflect.ValueOf(in)
	if valueOf.IsNil() {
		return nil, errors.New("Passed nil pointer interface into the validator ")
	}

	return &Reflection{
		TypeOf:  typeOf,
		ValueOf: reflect.Indirect(valueOf),
	}, nil
}

func (v *Validator) prepareRules(
	reflectObject *Reflection,
	calculatedRules map[string]map[string][]string,
	nestedFieldNames *[]string,
) {

	for i := 0; i < reflectObject.TypeOf.NumField(); i++ {

		field := reflectObject.TypeOf.Field(i)
		if tagValue := field.Tag.Get(tagName); tagValue != "" {

			// splits tag for rule that will have from "ruleName:val,val,val..."
			rulesArr := strings.Split(tagValue, ruleSeparator)
			for _, v := range rulesArr {

				// splits the rule name from rule value
				kv := strings.Split(v, ruleValueSeparator)
				if len(kv) == 1 {
					// there can be rule without values but names and we just adds some value
					kv = append(kv, "1")
				}

				// now we adding found rule into storage by field name as key
				// one field can have many rule
				if fieldRules, exists := calculatedRules[field.Name]; exists {
					// using existed mep
					fieldRules[strings.TrimSpace(kv[0])] = strings.Split(strings.TrimSpace(kv[1]), valuesSeparator)
				} else {
					// create new map if there no map for given field name
					calculatedRules[field.Name] = map[string][]string{strings.TrimSpace(kv[0]): strings.Split(strings.TrimSpace(kv[1]), valuesSeparator)}
				}
			}
		}

		if k := field.Type.Kind(); k == reflect.Struct || k == reflect.Ptr {
			*nestedFieldNames = append(*nestedFieldNames, field.Name)
		}
	}
}
