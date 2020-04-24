package mapping

import (
	"errors"
	"reflect"
)

// Errors
var (
	ErrNotArray     = errors.New("source value is not an array or slice")
	ErrNilTransform = errors.New("transform cannot be nil")
	ErrNotFunction  = errors.New("transform must be a function")
)

// ParallelMap an array of something into another thing using a go routine.
// Example: Map([]int{1,2,3}, func(num int) int { return num+1 })
// Results: []int{2,3,4}
func ParallelMap(source interface{}, transform interface{}) (interface{}, error) {
	if err := isSliceOrArray(source); err != nil {
		return nil, err
	}

	if err := isFunc(transform); err != nil {
		return nil, err
	}

	return nil, nil
}

func isSliceOrArray(source interface{}) error {
	sourceValue := reflect.ValueOf(source)
	kind := sourceValue.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return ErrNotArray
	}

	return nil
}

func isFunc(transform interface{}) (err error) {
	if transform == nil {
		return ErrNilTransform
	}

	transformValue := reflect.ValueOf(transform)
	if transformValue.Kind() != reflect.Func {
		return ErrNotFunction
	}

	return nil
}
