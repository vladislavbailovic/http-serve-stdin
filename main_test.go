package main

import (
	"strings"
	"testing"
)

func TestGetStdinHappyPath(t *testing.T) {
	expected := "Hello there"
	expectedReader := strings.NewReader(expected)

	actual := getStdin(expectedReader)
	if actual != expected {
		t.Log(actual)
		t.Fatalf("expected {%s} but got {%s} from stdin instead", expected, actual)
	}
}

func TestGetStdinReturnsEmptyStringIfNeeded(t *testing.T) {
	expected := ""
	expectedReader := strings.NewReader(expected)

	actual := getStdin(expectedReader)
	if actual != expected {
		t.Log(actual)
		t.Fatalf("expected {%s} but got {%s} from stdin instead", expected, actual)
	}
}
