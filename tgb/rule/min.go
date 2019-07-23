package rule

import (
	"fmt"
	"reflect"
	"strconv"
)

type RuleMin struct {
	*Rule
	min int64
}

func NewMinRule(base *Rule) ValidateRule {
	const errMessage = "Min validation err"
	const ruleName = "min"

	return &RuleMin{
		Rule: base.init(ruleName, errMessage, 1),
	}
}

func (r *RuleMin) Validate(field reflect.Value) (vr ValidateRule) {

	if !field.IsValid() || !r.isRestrictionValid() {
		goto err
	}

	switch field.Kind() {
	case reflect.String:
		if int64(len(field.String())) < r.min {
			goto err
		}
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if field.IsNil() || int64(len(field.String())) < r.min {
			goto err
		}
	case reflect.Array:
		if int64(field.Len()) < r.min {
			goto err
		}
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if field.Int() < r.min {
			goto err
		}
	}

	return r
err:
	return r.validationFailed()
}

func (r *RuleMin) isRestrictionValid() bool {
	if len(r.Restrictions) < r.MinRestrictionsCount {
		return false
	}

	val, err := strconv.Atoi(r.Restrictions[0])
	if err != nil {
		fmt.Printf("Wrong restriction parameter in %v rule", r.RuleName)
		return false
	}

	r.min = int64(val)

	return r.min >= 0
}
