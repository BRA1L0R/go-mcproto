package tagutils_test

import (
	"reflect"
	"testing"

	"github.com/BRA1L0R/go-mcproto/packets/serialization/tagutils"
)

func TestCheckDependency(t *testing.T) {
	type DependsOn struct {
		Condition  bool
		Dependency string `depends_on:"Condition"`
	}

	trueDependency := DependsOn{Condition: true, Dependency: "Hello World"}
	inter := reflect.ValueOf(trueDependency)
	interValue := inter.Type().Field(1)
	if !tagutils.CheckDependency(inter, interValue) {
		t.Error("dependency should have been true but function returned a false")
	}

	falseDependency := DependsOn{Condition: false, Dependency: "Hello World"}
	inter = reflect.ValueOf(falseDependency)
	interValue = inter.Type().Field(1)
	if tagutils.CheckDependency(inter, interValue) {
		t.Error("dependency should have been true but function returned a false")
	}
}
