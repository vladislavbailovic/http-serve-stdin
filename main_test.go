package main

import (
	"net/http"
	"net/http/httptest"
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

func TestStdinHandlerResponse(t *testing.T) {
	expected := "Hello there"
	expectedReader := strings.NewReader(expected)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getStdinHandler(expectedReader))
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Log(recorder)
		t.Fatalf("expected success requesting stdin, got %d instead", recorder.Code)
	}

	if actual := recorder.Body.String(); actual != expected {
		t.Log(recorder)
		t.Fatalf("expected server to respond with {%s}, got {%s} instead", expected, actual)
	}
}
