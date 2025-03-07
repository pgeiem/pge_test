package engine

import (
	"testing"
)

type NestedStruct struct {
	NestedField1 string
	NestedField2 int
}

type TestStruct struct {
	Field1      string
	Field2      int
	Field3      bool
	NestedField *NestedStruct
}

func TestOnlyOneFieldSet(t *testing.T) {
	tests := []struct {
		name     string
		input    TestStruct
		expected bool
	}{
		{"No fields set", TestStruct{}, true},
		{"One field set", TestStruct{Field1: "value"}, true},
		{"Two fields set", TestStruct{Field1: "value", Field2: 1}, false},
		{"All fields set", TestStruct{Field1: "value", Field2: 1, Field3: true}, false},
		{"Nested field set", TestStruct{NestedField: &NestedStruct{NestedField1: "nested"}}, true},
		{"Field and nested field set", TestStruct{Field1: "value", NestedField: &NestedStruct{NestedField1: "nested"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOnlyOneFieldSet(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
