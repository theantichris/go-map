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
	ErrNilType      = errors.New("map result type cannot be nil")
)

// ParallelMap an array of something into another thing using a go routine.
// Example: Map([]int{1,2,3}, func(num int) int { return num+1 })
// Results: []int{2,3,4}
func ParallelMap(source, transform interface{}, T reflect.Type) (interface{}, error) {
	sourceValue := reflect.ValueOf(source)

	if err := validate(sourceValue, transform, T); err != nil {
		return nil, err
	}

	result := reflect.MakeSlice(reflect.SliceOf(T), sourceValue.Len(), sourceValue.Cap())

	return result, nil
}

func validate(sourceValue reflect.Value, transform interface{}, T reflect.Type) error {
	if err := isSliceOrArray(sourceValue); err != nil {
		return err
	}

	if err := isFunc(transform); err != nil {
		return err
	}

	if err := isType(T); err != nil {
		return err
	}

	return nil
}

func isSliceOrArray(sourceValue reflect.Value) error {
	kind := sourceValue.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return ErrNotArray
	}

	return nil
}

func isFunc(transform interface{}) error {
	if transform == nil {
		return ErrNilTransform
	}

	transformValue := reflect.ValueOf(transform)
	if transformValue.Kind() != reflect.Func {
		return ErrNotFunction
	}

	return nil
}

func isType(T reflect.Type) error {
	if T == nil {
		return ErrNilType
	}

	return nil
}
