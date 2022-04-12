package core

import (
	"fmt"
)

type sliceFlag[T fmt.Stringer] struct {
	values         []T
	fromStringFunc func(string) (T, error)
}

func NewSliceFlag[T fmt.Stringer](fromStringFunc func(string) (T, error)) sliceFlag[T] {
	return sliceFlag[T]{
		fromStringFunc: fromStringFunc,
		values:         []T{},
	}
}

func (s *sliceFlag[T]) String() string {
	var stringRepr string
	for index := range s.values {
		// last entry, dont append ,
		if index == len(s.values)-1 {
			stringRepr += fmt.Sprintf("%s", s.values[index])
			break
		}
		stringRepr += fmt.Sprintf("%s, ", s.values[index])
	}
	return stringRepr
}

func (s *sliceFlag[T]) Set(value string) error {
	valFromString, err := s.fromStringFunc(value)
	if err != nil {
		return err
	}
	s.values = append(s.values, valFromString)
	return nil
}

func (s *sliceFlag[T]) GetValues() (objs []T) {
	return s.values
}
