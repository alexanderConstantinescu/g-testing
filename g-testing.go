package g

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type T interface {
}

type TTF struct {
	TName                string
	TF                   T
	TFInput              []T
	TFOutput             T
	TFPipeOutput         T
	TFOverride           T
	TFOverrideAssignment T
	TCleanup             T
}

type TTO struct {
	TName    string
	TOMethod string
	TOInput  []T
	TOCanary T
	TOExpect T
}

func TPipef(expected string, t *testing.T) func(format string, a ...interface{}) (int, error) {
	return func(format string, a ...interface{}) (int, error) {
		s1 := fmt.Sprintf(format, a...)
		s2 := fmt.Sprintf("%v", expected)
		if strings.Compare(s1, s2) != 0 {
			t.Fatalf("expected string: \"%s\", got: \"%s\"", s2, s1)
		}
		return 0, nil
	}
}

func TPipeln(expected string, t *testing.T) func(a ...interface{}) (int, error) {
	return func(a ...interface{}) (int, error) {
		s1 := fmt.Sprintln(a...)
		s2 := fmt.Sprintln(expected)
		if strings.Compare(s1, s2) != 0 {
			t.Fatalf("expected string: \"%s\", got: \"%s\"", s2, s1)
		}
		return 0, nil
	}
}

func TestStructMethod(tt []TTO, t *testing.T) {
	for _, tc := range tt {
		// prepare comparison objects for field comparison
		fields := reflect.TypeOf(tc.TOExpect)
		expectObj := reflect.ValueOf(tc.TOExpect)
		canaryObj := reflect.ValueOf(tc.TOCanary)

		// prepare input args for reflection
		inputs := make([]reflect.Value, len(tc.TOInput))
		for i := range tc.TOInput {
			inputs[i] = reflect.ValueOf(tc.TOInput[i])
		}

		// Call method on test object with appropriate input args
		canaryObj.MethodByName(tc.TOMethod).Call(inputs)
		for i := 0; i < fields.NumField(); i++ {
			if !reflect.DeepEqual(expectObj.Field(i).Interface(), reflect.ValueOf(tc.TOCanary).Field(i).Interface()) {
				t.Fatalf("test: \"%s\" failed! Expected: %v, got: %v for field: %v", tc.TName, expectObj.Field(i).Interface(), reflect.ValueOf(tc.TOCanary).Field(i).Interface(), fields.Field(i).Name)
			}
		}
	}
}

func TestPackageMethod(tt []TTF, t *testing.T) {
	for _, tc := range tt {
		test := tc.TFOverride
		fmt.Println(reflect.ValueOf(&test).CanAddr())
		inputs := make([]reflect.Value, len(tc.TFInput))
		for i := range tc.TFInput {
			inputs[i] = reflect.ValueOf(tc.TFInput[i])
		}
		for _, v := range reflect.ValueOf(tc.TF).Call(inputs) {
			t1 := fmt.Sprintf("%v", reflect.ValueOf(tc.TFOutput).Interface())
			t2 := fmt.Sprintf("%v", reflect.ValueOf(v).Interface())
			if t1 != t2 {
				t.Fatalf("test: \"%s\" failed! Expected: '%v', got: '%v'", tc.TName, t1, t2)
			}
		}
		if tc.TCleanup != nil {
			os.Remove(tc.TCleanup.(string))
		}
	}
}
