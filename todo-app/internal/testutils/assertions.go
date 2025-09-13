package testutils

import (
	"testing"
)

// AssertNoError проверяет, что ошибка равна nil
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError проверяет, что ошибка НЕ равна nil
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", message)
	}
}

// AssertEqual проверяет равенство двух значений
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v but got %v", message, expected, actual)
	}
}

// AssertNotNil проверяет, что значение не равно nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s: expected non-nil value", message)
	}
}

// AssertTrue проверяет, что значение равно true
func AssertTrue(t *testing.T, value bool, message string) {
	t.Helper()
	if !value {
		t.Fatalf("%s: expected true but got false", message)
	}
}

// AssertFalse проверяет, что значение равно false
func AssertFalse(t *testing.T, value bool, message string) {
	t.Helper()
	if value {
		t.Fatalf("%s: expected false but got true", message)
	}
}

// AssertNotEqual проверяет, что два значения не равны
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected == actual {
		t.Fatalf("%s: expected values to be different, but both are %v", message, expected)
	}
}

// AssertGreaterThan проверяет, что первое значение больше второго
func AssertGreaterThan(t *testing.T, actual, expected int, message string) {
	t.Helper()
	if actual <= expected {
		t.Fatalf("%s: expected %d to be greater than %d", message, actual, expected)
	}
}

// AssertContains проверяет, что слайс содержит элемент
func AssertContains(t *testing.T, slice []string, element string, message string) {
	t.Helper()
	for _, item := range slice {
		if item == element {
			return // Found it
		}
	}
	t.Fatalf("%s: slice %v does not contain element %s", message, slice, element)
}
