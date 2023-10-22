package problems

import (
	"fmt"
	"strings"
)

type EntityNotFoundError struct {
	Type string
	Path []string
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s/%s", e.Type, strings.Join(e.Path, "/"))
}
