package request

import (
	"errors"
	"fmt"
)

var ErrTypeInvalid = errors.New("invalid type")

// Validator is an object that can be validated.
type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid() (problems Problems)
}

type Problems struct {
	List map[string]string
}

//type Problems map[string]string

func IsValid[T Validator](v T) (Problems, error) {
	problems := v.Valid()
	if len(problems.List) > 0 {
		err := fmt.Errorf("%w: type %T has %d problems", ErrTypeInvalid, v, len(problems.List))
		return problems, err
	}
	return problems, nil
}

func (p Problems) String() string {
	str := ""
	for k, v := range p.List {
		str += fmt.Sprintf(" %s: %s;", k, v)
	}
	return str
}
