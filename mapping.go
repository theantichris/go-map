package mapping

import (
	"errors"
	"reflect"
	"sync"
)

// Errors
var (
	ErrNotArray     = errors.New("source value is not an array or slice")
	ErrNilTransform = errors.New("transform cannot be nil")
	ErrNotFunction  = errors.New("transform must be a function")
	ErrNilType      = errors.New("map result type cannot be nil")
)

// ConcurrentMap an array of something into another thing using a go routine.
// Example: Map([]int{1,2,3}, func(num int) int { return num+1 })
// Results: []int{2,3,4}
func ConcurrentMap(source, transform interface{}, T reflect.Type) (interface{}, error) {
	sourceValue := reflect.ValueOf(source)
	transformValue := reflect.ValueOf(transform)

	if err := validate(sourceValue, transform, T); err != nil {
		return nil, err
	}

	result := reflect.MakeSlice(reflect.SliceOf(T), sourceValue.Len(), sourceValue.Cap())

	wg := &sync.WaitGroup{}
	wg.Add(sourceValue.Len())

	for i := 0; i < sourceValue.Len(); i++ {
		go func(index int, entry reflect.Value) {
			// call transform and store the result
			transformResults := transformValue.Call([]reflect.Value{entry})

			// store the result in result container
			resultEntry := result.Index(index)
			if len(transformResults) > 0 {
				resultEntry.Set(transformResults[0])
			} else {
				resultEntry.Set(reflect.Zero(T))
			}

			wg.Done()
		}(i, sourceValue.Index(i))
	}

	wg.Wait()

	return result.Interface(), nil
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
