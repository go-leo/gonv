// Package gonv provides utility functions for the gonv type conversion library.
// This file contains helper functions used internally by other conversion functions.
package gonv

import (
	"reflect"
)

// trimZeroDecimal removes trailing zeros and decimal points from a numeric string.
// For example, "10.00" becomes "10", "5.10" becomes "5.1", and "3.0" becomes "3".
//
// Example:
//
//	result := trimZeroDecimal("10.00") // returns "10"
//	result := trimZeroDecimal("5.10")  // returns "5.1"
//	result := trimZeroDecimal("3.14")  // returns "3.14"
func trimZeroDecimal(s string) string {
	var foundZero bool
	// Process the string from right to left
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			// If we've found zeros before the decimal point, remove the decimal point
			if foundZero {
				return s[:i-1]
			}
		case '0':
			// Mark that we've found a zero
			foundZero = true
		default:
			// For any other character, return the string as is
			return s
		}
	}
	return s
}

// indirectValue dereferences pointers in a reflect.Value until it reaches a non-pointer value or nil.
// This function is useful when working with reflected values that might be pointers to the actual data.
//
// Example:
//
//	var x int = 42
//	var px *int = &x
//	var ppx **int = &px
//	indirectValue(reflect.ValueOf(ppx)) would return reflect.ValueOf(x)
func indirectValue(v reflect.Value) reflect.Value {
	// Keep dereferencing while the value is a pointer and not nil
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v
}
