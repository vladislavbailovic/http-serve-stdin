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

func TestGetDefaultHeaders(t *testing.T) {
	headers := getDefaultHeaders()

	if headers["content-type"] != "text/plain; charset=utf-8" {
		t.Log(headers)
		t.Fatalf("expected text as default content type, got %s", headers["content-type"])
	}
}

func TestParseHeadersReturnsHeadersMap(t *testing.T) {
	test := []string{
		"cOntEnt-type: application/json",
		"Server: in2http",
		"whatevers: Has:Some:Colons",
	}
	expected := map[string]string{
		"content-type": "application/json",
		"server":       "in2http",
		"whatevers":    "Has:Some:Colons",
	}

	headers := getParsedHeaders(test)
	for name, value := range expected {
		actual, ok := headers[name]
		if !ok {
			t.Fatalf("expected header {%s} to be parsed, it wasn't", name)
		}

		if value != actual {
			t.Fatalf("expected header {%s} to be {%s} - got {%s} instead", name, value, actual)
		}
	}
}

func TestGetHeadersReturnsHeadersMap(t *testing.T) {
	test := []string{
		"Server: in2http",
		"whatevers: Has:Some:Colons",
	}
	expected := map[string]string{
		"content-type": getDefaultHeaders()["content-type"],
		"server":       "in2http",
		"whatevers":    "Has:Some:Colons",
	}

	headers := getHeaders(test)
	for name, value := range expected {
		actual, ok := headers[name]
		if !ok {
			t.Fatalf("expected header {%s} to be parsed, it wasn't", name)
		}

		if value != actual {
			t.Fatalf("expected header {%s} to be {%s} - got {%s} instead", name, value, actual)
		}
	}
}

func TestGetParsedFlagsReturnDefaultsWithNoFlagsPassed(t *testing.T) {
	port, headers := getParsedFlags([]string{})
	if port != 8080 {
		t.Fatalf("expected port to be {%d} by default, got {%d}", 8080, port)
	}

	if len(headers) != 0 {
		t.Fatalf("expected headers to be empty, but they aren't: {%s}", strings.Join(headers, "|"))
	}
}

func TestGetParsedFlagsReturnPortIfSet(t *testing.T) {
	port1, _ := getParsedFlags([]string{"-p", "666"})
	if port1 != 666 {
		t.Fatalf("expected port to be {%d} by default, got {%d}", 666, port1)
	}
	port2, _ := getParsedFlags([]string{"--port", "667"})
	if port2 != 667 {
		t.Fatalf("expected port to be {%d} by default, got {%d}", 667, port2)
	}
}

func TestGetParsedFlagsReturnHeadersIfSet(t *testing.T) {
	_, h1 := getParsedFlags([]string{"-h", "content-type: text/html"})
	if h1[0] != "content-type: text/html" {
		t.Fatalf("expected port to be {%s} by default, got {%s}", "content-type: text/html", h1)
	}
	_, h2 := getParsedFlags([]string{"--header", "content-type: application/json"})
	if h2[0] != "content-type: application/json" {
		t.Fatalf("expected port to be {%s} by default, got {%s}", "content-type: application/json", h2)
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
	handler := http.HandlerFunc(getStdinHandler([]string{}, expectedReader))
	handler.ServeHTTP(recorder, req)

	if actual := recorder.Body.String(); actual != expected {
		t.Log(recorder)
		t.Fatalf("expected server to respond with {%s}, got {%s} instead", expected, actual)
	}
}

func TestStdinIsServedWithDefaultHeaders(t *testing.T) {
	expected := ""
	expectedReader := strings.NewReader(expected)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getStdinHandler([]string{}, expectedReader))
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Log(recorder)
		t.Fatalf("expected success requesting stdin, got %d instead", recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != getDefaultHeaders()["content-type"] {
		t.Log(recorder.Header())
		t.Fatalf("expected {%s} content type, got {%s} instead", getDefaultHeaders()["content-type"], contentType)
	}
}

func TestStdinIsServedWithCustomHeaders(t *testing.T) {
	expected := ""
	expectedReader := strings.NewReader(expected)
	headers := []string{
		"content-type: text/html",
		"server: in2http",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getStdinHandler(headers, expectedReader))
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Log(recorder)
		t.Fatalf("expected success requesting stdin, got %d instead", recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "text/html" {
		t.Log(recorder.Header())
		t.Fatalf("expected {%s} content type, got {%s} instead", "text/html", contentType)
	}

	if server := recorder.Header().Get("Server"); server != "in2http" {
		t.Log(recorder.Header())
		t.Fatalf("expected {%s} server, got {%s} instead", "in2http", server)
	}
}
