package rule

import (
	"fmt"
	"reflect"
	"strconv"
)

type Max struct {
	*Rule
	max int64
}

func NewMaxRule(base *Rule) ValidateRule {
	const errMessage = "Max validation err"
	const ruleName = "max"

	return &Max{
		Rule: base.Init(ruleName, errMessage, 1),
	}
}

func (r *Max) Validate(field reflect.Value) (vr ValidateRule) {

	if !field.IsValid() || !r.isRestrictionValid() {
		goto err
	}

	switch field.Kind() {
	case reflect.String:
		if int64(len(field.String())) > r.max {
			goto err
		}
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if field.IsNil() || int64(len(field.String())) > r.max {
			goto err
		}
	case reflect.Array:
		if int64(field.Len()) > r.max {
			goto err
		}
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if field.Int() > r.max {
			goto err
		}
	}

	return r
err:
	return r.ValidationFailed()
}

func (r *Max) isRestrictionValid() bool {
	if len(r.Restrictions) < r.MinRestrictionsCount {
		return false
	}

	val, err := strconv.Atoi(r.Restrictions[0])
	if err != nil {
		fmt.Printf("Wrong restriction parameter in %v rule", r.RuleName)
		return false
	}

	r.max = int64(val)

	return r.max >= 1
}
