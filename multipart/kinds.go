package multipart

import "reflect"

var basicKinds = append([]reflect.Kind{
	reflect.Bool,
	reflect.String,
}, numericKinds...)

var intKinds = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
}

var floatKinds = []reflect.Kind{
	reflect.Float32,
	reflect.Float64,
}

var numericKinds = append(intKinds, floatKinds...)

func isNumericKind(k reflect.Kind) bool {
	for _, i := range numericKinds {
		if k == i {
			return true
		}
	}

	return false
}

func isBasicKind(k reflect.Kind) bool {
	for _, i := range basicKinds {
		if k == i {
			return true
		}
	}

	return false
}

func isIntKind(k reflect.Kind) bool {
	for _, i := range intKinds {
		if k == i {
			return true
		}
	}
	return false
}

func isFloatKind(k reflect.Kind) bool {
	for _, i := range floatKinds {
		if k == i {
			return true
		}
	}
	return false
}
