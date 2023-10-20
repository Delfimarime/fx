package problems

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

type JsonValidationError struct {
	Schema   string
	Problems []gojsonschema.ResultError
}

func (j *JsonValidationError) Error() string {
	return fmt.Sprintf(`JsonSchemaValidationError[name="%s"]`, j.Schema)
}
