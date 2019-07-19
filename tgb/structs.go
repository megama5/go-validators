package tgb

import "fmt"

type ErrMessage struct {
	FieldName string
	RuleName  string
	Message   string
}

func (em ErrMessage) String() string {
	//return fmt.Sprintf("Field \"%s\" rised error on rule \"%s\" with message: \"%v\"", em.FieldName, em.RuleName, em.Message)
	return fmt.Sprintf("Validation error -> Field: \"%s\", Rule:\"%s\", Message: \"%v\"", em.FieldName, em.RuleName, em.Message)
}

type ErrMessagesList []*ErrMessage

func (eml ErrMessagesList) String() string {
	var result string
	for _, v := range eml {
		if v == nil {
			continue
		}

		result = fmt.Sprintf("%s,%s", result, v.String())
	}

	return result
}
