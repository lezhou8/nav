package main

import "testing"

func TestIsDirAccessible(t *testing.T) {
	if !isDirAccessible("/tmp") {
		t.Error("Expected /tmp to be accessible")
	}
	if isDirAccessible("/tmp/foobar") {
		t.Error("Expected /tmp/foobar to be inaccessible")
	}
}
