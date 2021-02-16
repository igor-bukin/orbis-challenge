// nolint
package config

import (
	"errors"
	"fmt"
	"reflect"
)

// Clone clones src to dst, break all references
func Clone(dst, src interface{}) error {
	rvd := indirect(reflect.ValueOf(dst))
	rvs := indirect(reflect.ValueOf(src))

	return clone(rvd, rvs)
}

// clone recursion helper
func clone(rvd, rvs reflect.Value) error {
	if !rvs.IsValid() {
		return nil
	}

	if !rvd.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	if !rvd.CanSet() {
		return fmt.Errorf("can't set dst: %s, from src: %s", rvd, rvs)
	}

	// skip nil
	if rvs.Kind() == reflect.Ptr && rvs.IsNil() || rvs.Kind() == reflect.Slice && rvs.IsNil() {
		return nil
	}

	switch rvs.Kind() {
	case reflect.Struct:
		l := rvs.NumField()
		structDst := reflect.New(rvs.Type())

		for i := 0; i < l; i++ {
			fieldValue := indirect(rvs.Field(i))
			fieldType := indirectType(rvs.Type())
			name := fieldType.Field(i).Name
			dst := indirect(structDst).FieldByName(name)

			if dst.IsValid() && dst.CanSet() {
				if err := clone(dst, fieldValue); err != nil {
					return err
				}
			}
		}

		rvs = structDst

	case reflect.Slice:
		l := rvs.Len()
		sliceDst := reflect.MakeSlice(rvs.Type(), l, l)

		for i := 0; i < l; i++ {
			fieldValue := indirect(rvs.Index(i))
			dst := sliceDst.Index(i)
			if dst.IsValid() && dst.CanSet() {
				if err := clone(dst, fieldValue); err != nil {
					return err
				}
			}
		}

		rvs = sliceDst
	}

	set(rvd, rvs)

	return nil
}

// set properly sets value from to
func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
