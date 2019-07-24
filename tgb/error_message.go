package tgb

import "fmt"

type ErrMessage struct {
	FieldName string
	RuleName  string
	Message   string
}

func (em *ErrMessage) Error() string {
	return fmt.Sprintf("Validation error -> Field: \"%s\", Rule:\"%s\", Message: \"%v\"", em.FieldName, em.RuleName, em.Message)
}

type ErrMessagesList []*ErrMessage

func (eml ErrMessagesList) Error() string {
	var result string
	for _, v := range eml {
		if v == nil {
			continue
		}

		result = fmt.Sprintf("%s,%s", result, v.Error())
	}

	return result
}
