package rule

import (
	"reflect"
)

type Required struct {
	*Rule
}

func NewRequiredRule(base *Rule) ValidateRule {
	const errMessage = "Required validation err"
	const ruleName = "required"

	return &Required{
		Rule: base.Init(ruleName, errMessage, 0),
	}
}

func (r *Required) Validate(field reflect.Value) (vr ValidateRule) {

	if !field.IsValid() {
		return r.ValidationFailed()
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) == 0 {
			return r.ValidationFailed()
		}
	case reflect.Ptr:
		if field.IsNil() {
			return r.ValidationFailed()
		}
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if field.IsNil() || field.Len() == 0 {
			return r.ValidationFailed()
		}
	case reflect.Array:
		if field.Len() == 0 {
			return r.ValidationFailed()
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
			return r.ValidationFailed()
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if field.Float() == 0 {
			return r.ValidationFailed()
		}
	}

	return r
}
