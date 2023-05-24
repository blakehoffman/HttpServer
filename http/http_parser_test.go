package http

import (
	"testing"
)

func Test_Get_Status_Line_Parse_Location_With_Character_Last_Byte_And_Non_Empty_Url_Returns_InUrl(t *testing.T) {
	requestLine := requestLine{
		verb: "GET",
		url:  "test",
	}

	actual := get_status_line_parse_location('t', requestLine)

	if actual != InUrl {
		t.Fatalf("get_status_line_parse_location('a', requestLine) = %d, want %d", actual, InUrl)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Character_Last_Byte_And_Non_Empty_Verb_Returns_InVerb(t *testing.T) {
	requestLine := requestLine{
		verb: "test",
	}

	actual := get_status_line_parse_location('t', requestLine)

	if actual != InVerb {
		t.Fatalf("get_status_line_parse_location('a', requestLine) = %d, want %d", actual, InVerb)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Character_Last_Byte_And_Non_Empty_Version_Returns_InVersion(t *testing.T) {
	requestLine := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "test",
	}

	actual := get_status_line_parse_location('t', requestLine)

	if actual != InVersion {
		t.Fatalf("get_status_line_parse_location('a', requestLine) = %d, want %d", actual, InVersion)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Zero_Last_Byte_Returns_InVerb(t *testing.T) {
	requestLine := requestLine{}

	actual := get_status_line_parse_location(0, requestLine)

	if actual != InVerb {
		t.Fatalf("get_status_line_parse_location(0, requestLine) = %d, want %d", actual, InVerb)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Space_Last_Byte_Returns_InUrl(t *testing.T) {
	requestLine := requestLine{
		verb: "GET",
	}

	actual := get_status_line_parse_location(' ', requestLine)

	if actual != InUrl {
		t.Fatalf("get_status_line_parse_location(' ', requestLine) = %d, want %d", actual, InUrl)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Space_Last_Byte_Returns_InVersion(t *testing.T) {
	requestLine := requestLine{
		verb: "get",
		url:  "test.com",
	}

	actual := get_status_line_parse_location(' ', requestLine)

	if actual != InVersion {
		t.Fatalf("get_status_line_parse_location(' ', requestLine) = %d, want %d", actual, InVersion)
	}
}

func Test_Parse_Http_Status_Line_With_Valid_Line_Returns_Status_Line(t *testing.T) {
	buffer := []byte{'G', 'E', 'T', ' ', 't', 'e', 's', 't', '.', 'h', 't', 'm', 'l', ' ', 'H', 'T', 'T', 'P', '/', '1', '.', '1', '\r', '\n'}

	expected := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "HTTP/1.1",
	}

	actual, _ := parse_http_status_line(buffer, httpParseResult[requestLine]{})

	if actual.result != expected {
		t.Fatalf("parse_http_status_line(), got %v want %v", actual.result, expected)
	}
}

func Test_Parse_Http_Status_Line_With_Space_As_Last_Byte_Between_Returns_Status_Line(t *testing.T) {
	buffer := []byte{'t', 'e', 's', 't', '.', 'h', 't', 'm', 'l', ' ', 'H', 'T', 'T', 'P', '/', '1', '.', '1', '\r', '\n'}
	parseResult := httpParseResult[requestLine]{
		lastByte: ' ',
		result: requestLine{
			verb: "GET",
		},
	}

	expected := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "HTTP/1.1",
	}

	actual, _ := parse_http_status_line(buffer, parseResult)

	if actual.result != expected {
		t.Fatalf("parse_http_status_line() with space as last byte, got %v want %v", actual.result, expected)
	}
}

func Test_Parse_Http_Status_Line_With_Partial_Url_Returns_Status_Line(t *testing.T) {
	buffer := []byte{'s', 't', '.', 'h', 't', 'm', 'l', ' ', 'H', 'T', 'T', 'P', '/', '1', '.', '1', '\r', '\n'}
	parseResult := httpParseResult[requestLine]{
		lastByte: 'e',
		result: requestLine{
			verb: "GET",
			url:  "te",
		},
	}

	expected := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "HTTP/1.1",
	}

	actual, _ := parse_http_status_line(buffer, parseResult)

	if actual.result != expected {
		t.Fatalf("parse_http_status_line() with partial url, got %v want %v", actual.result, expected)
	}
}

func Test_Parse_Http_Status_Line_With_Partial_Verb_Returns_Status_Line(t *testing.T) {
	buffer := []byte{'T', ' ', 't', 'e', 's', 't', '.', 'h', 't', 'm', 'l', ' ', 'H', 'T', 'T', 'P', '/', '1', '.', '1', '\r', '\n'}
	parseResult := httpParseResult[requestLine]{
		lastByte: 'E',
		result: requestLine{
			verb: "GE",
		},
	}

	expected := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "HTTP/1.1",
	}

	actual, _ := parse_http_status_line(buffer, parseResult)

	if actual.result != expected {
		t.Fatalf("parse_http_status_line() with partial verb, got %v want %v", actual.result, expected)
	}
}

func Test_Parse_Http_Status_Line_With_Partial_Version_Returns_Status_Line(t *testing.T) {
	buffer := []byte{'T', 'P', '/', '1', '.', '1', '\r', '\n'}
	parseResult := httpParseResult[requestLine]{
		lastByte: 'T',
		result: requestLine{
			verb:    "GET",
			url:     "test.html",
			version: "HT",
		},
	}

	expected := requestLine{
		verb:    "GET",
		url:     "test.html",
		version: "HTTP/1.1",
	}

	actual, _ := parse_http_status_line(buffer, parseResult)

	if actual.result != expected {
		t.Fatalf("parse_http_status_line() with partial version, got %v want %v", actual.result, expected)
	}
}
