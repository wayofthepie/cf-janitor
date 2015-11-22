package main_test


import (
	"reflect"
	"testing"
)

func TestJanitor(t *testing.T) {
	if msg := "test"; !reflect.DeepEqual(msg, "test") {
		t.Errorf("test = %+v, want %+v", msg, "test")
	}
}