package tgb

import (
	"errors"
	"fmt"
	. "github.com/megama5/go-validators/tgb/rules"
	"log"
	"reflect"
	"strings"
)

type Reflection struct {
	TypeOf  reflect.Type
	ValueOf reflect.Value
}

var rulesHolder = map[string]func(base *Rule) ValidateRule{
	"required":     NewRequiredRule,
	"min":          NewMinRule,
	"max":          NewMaxRule,
	"array-unique": NewArrayUniqueRule,
}

func Validate(in interface{}, requiredData map[string]interface{}) (bool, ErrMessagesList) {
	a, b := validate(in, "")

	err := required(requiredData)
	if err != nil {
		b = append(b, err)
	}

	return a, b
}

func required(data map[string]interface{}) *ErrMessage {
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

func validate(in interface{}, parentName string) (bool, ErrMessagesList) {
	// will holds field name of nested structure for the recursion call
	var nestedFieldNames []string

	// holder for the rules found in type description
	calculatedRules := make(map[string]map[string][]string, 0)

	// indicate of any error occurred
	errIndicator := false

	// holds all error messages
	errMessages := make(ErrMessagesList, 0)

	// if passed nil the validation is failed
	if in == nil {
		return false, errMessages
	}

	reflectObject, err := prepareReflection(in)
	if err != nil {
		log.Printf("Validator err: %v", err)
		return false, errMessages
	}

	prepareRules(reflectObject, calculatedRules, &nestedFieldNames)
	applyRules(calculatedRules, reflectObject, parentName, &errIndicator, &errMessages)
	applyRuleToNestedStructures(nestedFieldNames, reflectObject, parentName, &errIndicator, &errMessages)

	return !errIndicator, errMessages
}

func applyRuleToNestedStructures(
	nestedFieldNames []string,
	reflectObject *Reflection,
	parentName string,
	errIndicator *bool,
	errMessages *ErrMessagesList,
) {
	parentName = fmt.Sprintf("%s%s", prepareParentName(parentName), reflectObject.TypeOf.Name())
	// checking nested structures
	for _, v := range nestedFieldNames {
		field := reflectObject.ValueOf.FieldByName(v)
		if field.IsValid() && ((field.Kind() == reflect.Ptr && !field.IsNil()) || field.Kind() == reflect.Struct) {
			res, list := validate(field.Interface(), parentName)
			*errIndicator = *errIndicator && res
			*errMessages = append(*errMessages, list...)
		}
	}
}

func applyRules(
	calculatedRules map[string]map[string][]string,
	reflectObject *Reflection,
	parentName string,
	errIndicator *bool,
	errMessages *ErrMessagesList,
) {
	parentName = prepareParentName(parentName)
	// now we have calculatedRules for the structure related to the
	// fields and we need to get a field and validate it
	for fieldName, rules := range calculatedRules {

		//fmt.Printf("Field name under validation: %v\n", fieldName)
		field := reflectObject.ValueOf.FieldByName(fieldName)
		for ruleName, params := range rules {

			//fmt.Printf("Field Name:%v | Rule name: %v | params:%v\n", fieldName, ruleName, params)
			if rule, ruleFound := rulesHolder[ruleName]; ruleFound {

				//fmt.Printf("Rule found: %v\n", ruleName)
				result := rule(NewRule(fieldName, params)).Validate(field)

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

func prepareParentName(parentName string) string {
	if parentName != "" {
		parentName = fmt.Sprintf("%v.", parentName)
	}
	return parentName
}

func prepareReflection(in interface{}) (*Reflection, error) {
	typeOf, valueOf := reflect.TypeOf(in), reflect.ValueOf(in)
	if typeOf.Kind() == reflect.Ptr {
		// pointer received
		typeOf = typeOf.Elem()

		// can be pointer to nil value
		if valueOf.IsNil() {
			return nil, errors.New("Passed nil pointer interface into the validator ")
		}
		valueOf = valueOf.Elem()
	}
	return &Reflection{
		TypeOf:  typeOf,
		ValueOf: valueOf,
	}, nil
}

func prepareRules(
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

				// now we adding found rules into storage by field name as key
				// one field can have many rules
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
