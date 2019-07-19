package rules

import (
	"reflect"
)

type RuleRequired struct {
	*Rule
}

func NewRequiredRule(base *Rule) ValidateRule {
	const errMessage = "Required validation err"
	const ruleName = "required"

	return &RuleRequired{
		Rule: base.init(ruleName, errMessage, 0),
	}
}

func (r *RuleRequired) Validate(field reflect.Value) (vr ValidateRule) {

	if !field.IsValid() {
		goto err
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) == 0 {
			goto err
		}
	case reflect.Ptr:
		if field.IsNil() {
			goto err
		}
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if field.IsNil() || field.Len() == 0 {
			goto err
		}
	case reflect.Array:
		if field.Len() == 0 {
			goto err
		}
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		if field.Int() == 0 {
			goto err
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if field.Float() == 0 {
			goto err
		}
	}

	return r
err:
	return r.validationFailed()
}
