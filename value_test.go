package errors

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	err := New("root error")

	key := "testKey"

	// Test Value with no key-value pairs
	if Value(err, key) != nil {
		t.Error("expected nil value for root error")
	}

	// Test Value with a key-value pair
	expectedValue := "testValue"
	err = With(err, KV(key, "testValue"))
	if v := Value(err, key); v != expectedValue {
		t.Errorf("expected value %v for key %v, got %v", expectedValue, key, v)
	}

	// Test Value changing the value for the same key
	expectedValue = "anotherValue"
	err = With(err, KV(key, expectedValue))
	if v := Value(err, key); v != expectedValue {
		t.Errorf("expected value %v for key %v, got %v", expectedValue, key, v)
	}
}

func TestValueT(t *testing.T) {
	err := With(
		New("root error"),
		KV("int", 42),
		KV("string", "testValue"),
		KV("bool", true),
		KV("slice", []string{"a", "b", "c"}),
		KV("map", map[string]int{"key1": 1, "key2": 2}),
		KV("nil", nil),
		KV("empty", ""),
		KV("struct", struct{}{}),
	)

	if v := ValueT[int](err, "int"); v != 42 {
		t.Errorf("expected int value 42, got %d", v)
	}
	if v := ValueT[string](err, "string"); v != "testValue" {
		t.Errorf("expected string value 'testValue', got %s", v)
	}
	if v := ValueT[bool](err, "bool"); !v {
		t.Errorf("expected bool value true, got %t", v)
	}
	if v := ValueT[[]string](err, "slice"); !reflect.DeepEqual(v, []string{"a", "b", "c"}) {
		t.Errorf("expected slice value ['a', 'b', 'c'], got %v", v)
	}
	if v := ValueT[map[string]int](err, "map"); !reflect.DeepEqual(v, map[string]int{"key1": 1, "key2": 2}) {
		t.Errorf("expected map value {'key1': 1, 'key2': 2}, got %v", v)
	}
	if v := ValueT[any](err, "nil"); v != nil {
		t.Error("expected nil value for key 'nil'")
	}
	if v := ValueT[string](err, "empty"); v != "" {
		t.Error("expected empty string value for key 'empty'")
	}
	if v := ValueT[struct{}](err, "struct"); !reflect.DeepEqual(v, struct{}{}) {
		t.Error("expected empty struct value for key 'struct'")
	}
}

func TestValues(t *testing.T) {
	err := With(
		New("root error"),
		KV("key", "value1"),
		KV("key", "value2"),
		KV("key", "value3"),
	)
	values := Values(err, "key")
	if !reflect.DeepEqual(values, []any{"value3", "value2", "value1"}) {
		t.Errorf("expected values ['value3', 'value2', 'value1'], got %v", values)
	}
}

func TestValuesT(t *testing.T) {
	emprtyErr := Error{err: New("empty error")}
	if values := ValuesT[string](emprtyErr, "key"); values != nil {
		t.Error("expected nil values for empty error")
	}

	rootErr := New("root error")
	// Test ValuesT with no key-value pairs
	values := ValuesT[string](rootErr, "key")
	if values != nil {
		t.Error("expected nil values for root error")
	}

	// Test ValuesT with a key-value pair
	err := With(
		rootErr,
		KV("key", "value1"),
		KV("key", "value2"),
		KV("key", nil), // This should be ignored
		KV("key", "value3"),
	)
	values = ValuesT[string](err, "key")
	if !reflect.DeepEqual(values, []string{"value3", "value2", "value1"}) {
		t.Errorf("expected values ['value3', 'value2', 'value1'], got %v", values)
	}

	// Test ValuesT with mixed types
	err = With(
		rootErr,
		KV("key", "stringValue"),
		KV("key", 42),
		KV("key", true),
		KV("key", 123),
	)
	values2 := ValuesT[int](err, "key")
	if !reflect.DeepEqual(values2, []int{123, 42}) {
		t.Errorf("expected values [123, 42], got %v", values2)
	}
}

func TestValuesMapOf(t *testing.T) {
	emprtyErr := Error{err: New("empty error")}
	if valuesMap := ValuesMapOf(emprtyErr, "key"); len(valuesMap) != 0 {
		t.Error("expected empty map for empty error")
	}

	rootErr := New("root error")
	// Test ValuesMapOf with no key-value pairs
	valuesMap := ValuesMapOf(rootErr, "key")
	if len(valuesMap) != 0 {
		t.Errorf("expected empty map for root error: %v", valuesMap)
	}

	// Test ValuesMapOf with a key-value pair
	err := With(
		rootErr,
		KV("key1", "value1.1"),
		KV("key1", "value1.2"),
		KV("key2", "value2"),
		KV("key3", "value3"),
	)
	valuesMap = ValuesMapOf(err, "key")
	expectedMap := map[any][]any{
		"key1": {"value1.2", "value1.1"},
		"key2": {"value2"},
		"key3": {"value3"},
	}
	if !reflect.DeepEqual(expectedMap, valuesMap) {
		t.Errorf("expected values map %v, got %v", expectedMap, valuesMap)
	}
}

func TestValueMap(t *testing.T) {
	var err error = Error{err: New("empty error")}
	if valueMap := ValueMap(err); len(valueMap) != 0 {
		t.Error("expected empty map for empty error")
	}

	err = New("root error")
	err = With(
		err,
		KV("key1", "value1-wont-be-included"),
		KV("key1", "value1"),
		KV("key2", "value2"),
		KV("key3", "value3"),
	)

	valueMap := ValueMap(err)
	if len(valueMap) != 3 {
		t.Errorf("expected map with 3 key-value pairs, got %d", len(valueMap))
	}
	expectedMap := map[any]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	for k, v := range expectedMap {
		if value, exists := valueMap[k]; !exists || value != v {
			t.Errorf("expected key %v with value %v, got %v", k, v, valueMap[k])
		}
	}
}

func TestValueMapOf(t *testing.T) {
	emprtyErr := Error{err: New("empty error")}
	if valuesMap := ValueMapOf(emprtyErr, "key"); len(valuesMap) != 0 {
		t.Error("expected empty map for empty error")
	}

	rootErr := New("root error")
	// Test ValuesMapOf with no key-value pairs
	valuesMap := ValueMapOf(rootErr, "key")
	if len(valuesMap) != 0 {
		t.Errorf("expected empty map for root error: %v", valuesMap)
	}

	type anotherKeyType int
	// Test ValuesMapOf with a key-value pair
	err := With(
		rootErr,
		KV("key1", "value1-wont-be-included"),
		KV("key1", "value1"),
		KV("key2", "value2"),
		KV(anotherKeyType(0), "value3"),
	)
	valuesMap = ValueMapOf(err, "key")
	expectedMap := map[any]any{
		"key1": "value1",
		"key2": "value2",
	}
	if !reflect.DeepEqual(expectedMap, valuesMap) {
		t.Errorf("expected values map %v, got %v", expectedMap, valuesMap)
	}
}

func TestValueAllSlice(t *testing.T) {
	emptyErr := Error{err: New("empty error")}
	if values := ValueAllSlice(emptyErr); len(values) != 0 {
		t.Error("expected empty slice for empty error")
	}

	rootErr := New("root error")
	// Test ValueAllSlice with no key-value pairs
	values := ValueAllSlice(rootErr)
	if len(values) != 0 {
		t.Errorf("expected empty slice for root error, got %v", values)
	}

	// Test ValueAllSlice with a key-value pair
	err := With(
		rootErr,
		KV("key1", "value1"),
		KV("key2", "value2"),
	)
	values = ValueAllSlice(err)
	if len(values) != 2 {
		t.Errorf("expected slice with 2 key-value pairs, got %d", len(values))
	}
	// Check if the values are in the expected order
	// Note: The order is the reverse of how they were added
	expectedValues := []KeyValuer{
		KeyValue{"key2", "value2"},
		KeyValue{"key1", "value1"},
	}
	for i, v := range expectedValues {
		if values[i].Key() != v.Key() || values[i].Value() != v.Value() {
			t.Errorf("expected value %v, got %v", v, values[i])
		}
	}
}
