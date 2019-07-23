package tgb

import (
	"github.com/megama5/go-validators/tgb/rule"
	"reflect"
	"strconv"
	"testing"
)

func TestRequiredRule(t *testing.T) {

	type Foo struct {
		Name           string                 `validate:"required"`
		Age            int32                  `validate:"required"`
		Info           map[string]interface{} `validate:"required"`
		AdditionalInfo map[string]interface{}
		About          string
	}

	type TCase struct {
		TestObject      interface{}
		TargetResult    bool
		CaseDescription string
	}

	cases := []*TCase{
		{
			TestObject: &Foo{
				Name:           "name",
				Age:            11,
				Info:           map[string]interface{}{"Addr": "addr"},
				AdditionalInfo: map[string]interface{}{"Addr1": "addr1"},
				About:          "about",
			},
			TargetResult:    true,
			CaseDescription: "All fields are populated",
		},
		{
			TestObject:      nil,
			TargetResult:    false,
			CaseDescription: "Nil object",
		},
		{
			TestObject: &Foo{
				Name:           "",
				Age:            0,
				Info:           nil,
				AdditionalInfo: nil,
				About:          "",
			},
			TargetResult:    false,
			CaseDescription: "Some fields will rise validation err",
		},
	}

	var rules = map[string]func(base *rule.Rule) rule.ValidateRule{
		"required":     rule.NewRequiredRule,
		"min":          rule.NewMinRule,
		"max":          rule.NewMaxRule,
		"array-unique": rule.NewArrayUniqueRule,
	}

	validator := NewValidator(rules)

	for _, v := range cases {
		func(c *TCase) {
			t.Run("Sub test case", func(t *testing.T) {
				t.Logf("obj: %+v\n", v.TestObject)
				result, errList := validator.Validate(v.TestObject, nil)
				t.Logf("\nResult: %v\n", result)
				if errList != nil {
					t.Logf("Error list:%v", errList.Error())
				}
				if result != v.TargetResult {
					t.Fatalf("Failed to pass case!\nDescription:%v\nTargetResult:%v -> Actual Result:%v\nErrList:%v",
						v.CaseDescription, v.TargetResult, result, errList)
				}
				t.Logf("Successfully passed test: %v", v.CaseDescription)
			})
		}(v)
	}
}

func TestAddCustomRule(t *testing.T) {

	type Foo struct {
		Name string `validate:"len:10"`
		Age  int32  `validate:"positive"`
	}

	type TCase struct {
		TestObject      interface{}
		TargetResult    bool
		CaseDescription string
	}

	cases := []*TCase{
		{
			TestObject: &Foo{
				Name: "name",
				Age:  1,
			},
			TargetResult:    true,
			CaseDescription: "All fields are populated",
		},
		{
			TestObject:      nil,
			TargetResult:    false,
			CaseDescription: "Nil object",
		},
		{
			TestObject: &Foo{
				Name: "",
				Age:  -1,
			},
			TargetResult:    false,
			CaseDescription: "Int32 field will rise validation err",
		},
		{
			TestObject: &Foo{
				Name: "ddddddddddddddddddddddd",
				Age:  1,
			},
			TargetResult:    false,
			CaseDescription: "String field will rise validation err",
		},
	}

	var rules = map[string]func(base *rule.Rule) rule.ValidateRule{
		"len":      NewLenRule,
		"positive": NewPositiveRule,
	}

	validator := NewValidator(rules)

	for _, v := range cases {
		func(c *TCase) {
			t.Run("Sub test case", func(t *testing.T) {
				t.Logf("obj: %+v\n", v.TestObject)
				result, errList := validator.Validate(v.TestObject, nil)
				t.Logf("\nResult: %v\n", result)
				if errList != nil {
					t.Logf("Error list:%v", errList.Error())
				}
				if result != v.TargetResult {
					t.Fatalf("Failed to pass case!\nDescription:%v\nTargetResult:%v -> Actual Result:%v\nErrList:%v",
						v.CaseDescription, v.TargetResult, result, errList)
				}
				t.Logf("Successfully passed test: %v", v.CaseDescription)
			})
		}(v)
	}
}

type Positive struct {
	*rule.Rule
}

func NewPositiveRule(base *rule.Rule) rule.ValidateRule {
	const errMessage = "Positive validation err"
	const ruleName = "positive"

	return &Positive{
		Rule: base.Init(ruleName, errMessage, 0),
	}
}

func (p *Positive) Validate(field reflect.Value) (vr rule.ValidateRule) {

	if !field.IsValid() {
		return p.ValidationFailed()
	}

	switch field.Kind() {
	case reflect.String:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Int:
		if field.Int() < 0 {
			return p.ValidationFailed()
		}
	case reflect.Int8:
		if field.Int() < 0 {
			return p.ValidationFailed()
		}
	case reflect.Int32:
		if field.Int() < 0 {
			return p.ValidationFailed()
		}
	case reflect.Int64:
		if field.Int() < 0 {
			return p.ValidationFailed()
		}
	case reflect.Float32:
		if field.Float() < 0 {
			return p.ValidationFailed()
		}
	case reflect.Float64:
		if field.Float() < 0 {
			return p.ValidationFailed()
		}
	}

	return p
}

type Len struct {
	*rule.Rule
	len int
}

func NewLenRule(base *rule.Rule) rule.ValidateRule {
	const errMessage = "Len validation err"
	const ruleName = "len"

	return &Len{
		Rule: base.Init(ruleName, errMessage, 1),
	}
}

func (l *Len) Validate(field reflect.Value) (vr rule.ValidateRule) {

	if !field.IsValid() || !l.isRestrictionValid() {
		return l.ValidationFailed()
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > l.len {
			return l.ValidationFailed()
		}
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		if field.IsNil() || len(field.String()) > l.len {
			return l.ValidationFailed()
		}
	case reflect.Array:
		if field.Len() > l.len {
			return l.ValidationFailed()
		}
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
	}

	return l
}

func (l *Len) isRestrictionValid() bool {
	if len(l.Restrictions) < l.MinRestrictionsCount {
		return false
	}

	val, err := strconv.Atoi(l.Restrictions[0])
	if err != nil {
		return false
	}

	l.len = val

	return l.len >= 0
}
