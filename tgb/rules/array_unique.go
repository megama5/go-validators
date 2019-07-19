package rules

import (
	"fmt"
	"reflect"
	"strconv"
)

type ArrayUnique struct {
	*Rule
	uniqueFieldName string
	canBeEmpty      bool
}

func NewArrayUniqueRule(base *Rule) ValidateRule {
	const errMessage = "Array element unique validation err"
	const ruleName = "array-unique"

	return &ArrayUnique{
		Rule: base.init(ruleName, errMessage, 2),
	}
}

func (r *ArrayUnique) Validate(field reflect.Value) (vr ValidateRule) {

	values := make(map[interface{}]interface{})

	if !field.IsValid() || !r.isRestrictionValid() {
		goto err
	}

	switch field.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:

		switch {
		case field.IsNil():
			r.Log(fmt.Sprintf("Nill pointer value on Uniquer Array rule"))
			goto err
		case field.Len() == 0 && !r.canBeEmpty:
			r.Log(fmt.Sprintf("Empty array recived"))
			goto err
		}

		errFlag := false
		for i := 0; i < field.Len(); i++ {

			switch v := field.Index(i); v.Kind() {
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32,
				reflect.Uint64, reflect.Float32, reflect.Float64:

				element := v.Interface()
				if _, ok := values[element]; ok {
					r.Log(fmt.Sprintf("Duplicate value on index: %v with value: %v", i, element))
					errFlag = true
				} else {
					values[element] = nil
				}

			case reflect.Ptr, reflect.Struct:

				element := v.Interface()
				reflectedElementValue := reflect.ValueOf(element)
				if reflectedElementValue.Type().Kind() == reflect.Ptr {
					reflectedElementValue = reflectedElementValue.Elem()
				}

				if !reflectedElementValue.IsValid() {
					r.Log("Couldn't convert element into struct")
					goto err
				}

				elementFieldValue := reflectedElementValue.FieldByName(r.uniqueFieldName)
				if !elementFieldValue.IsValid() {
					r.Log(fmt.Sprintf("Can't extract element field value by field name: %v", r.uniqueFieldName))
					errFlag = true
					continue
				}

				stringElementFieldValue := elementFieldValue.String()
				if len(stringElementFieldValue) == 0 {
					r.Log(fmt.Sprintf("Can't convert array element field value into string or empty string value"))
					errFlag = true
					continue
				}

				if _, ok := values[stringElementFieldValue]; ok {
					r.Log(fmt.Sprintf("Duplicate value on index: %v with value: %+v\n", i, element))
					errFlag = true
				} else {
					values[stringElementFieldValue] = nil
				}

			}
		}

		if errFlag {
			goto err
		}
	default:
		r.Log(fmt.Sprintf("Unsupported field kind. Kind: %v", field.Kind()))
		goto err
	}

	return r
err:
	return r.validationFailed()
}

func (r *ArrayUnique) isRestrictionValid() bool {
	if l := len(r.Restrictions); l < r.MinRestrictionsCount {
		r.Log(fmt.Sprintf("%v require %v param but got %v", r.RuleName, r.MinRestrictionsCount, l))
		return false
	}

	r.uniqueFieldName = r.Restrictions[0]

	canBeEmpty, err := strconv.ParseBool(r.Restrictions[1])
	if err != nil {
		r.Log(fmt.Sprintf("Unsupported bool representation in third parametr. Val: %v", r.Restrictions[2]))
	}
	r.canBeEmpty = canBeEmpty

	return true
}
