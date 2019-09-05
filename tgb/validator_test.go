package tgb

import "testing"

type TCase struct {
	TestObject      interface{}
	TargetResult    bool
	CaseDescription string
}

func TestRequiredRule(t *testing.T) {

	type Foo struct {
		Name           string                 `validate:"required"`
		Age            int32                  `validate:"required"`
		Info           map[string]interface{} `validate:"required"`
		AdditionalInfo map[string]interface{}
		About          string
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

	for _, v := range cases {
		func(c *TCase) {
			t.Run("Sub test case", func(t *testing.T) {
				t.Logf("obj: %+v\n", v.TestObject)
				result, errList := Validate(v.TestObject, nil)
				t.Logf("\nResult: %v\nerrList:%v", result, errList.String())
				if result != v.TargetResult {
					t.Fatalf("Failed to pass case!\nDescription:%v\nTargetResult:%v -> Actual Result:%v\nErrList:%v",
						v.CaseDescription, v.TargetResult, result, errList)
				}
				t.Logf("Successfully passed test: %v", v.CaseDescription)
			})
		}(v)
	}
}
