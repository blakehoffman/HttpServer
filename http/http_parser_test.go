package http

import (
	"testing"
)

func Test_Get_Status_Line_Parse_Location_With_Zero_Last_Byte_Returns_InVerb(t *testing.T){
	requestLine := requestLine{}
	
	actual := get_status_line_parse_location(0, requestLine)

	if actual != InVerb{
		t.Fatalf("get_stats_line_parse_location(0, requestLine) = %d, want %d", actual, InVerb)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Space_Last_Byte_Returns_InUrl(t *testing.T){
	requestLine := requestLine{}
	requestLine.verb = "get"

	actual := get_status_line_parse_location(' ', requestLine)

	if actual != InUrl{
		t.Fatalf("get_stats_line_parse_location(' ', requestLine) = %d, want %d", actual, InUrl)
	}
}

func Test_Get_Status_Line_Parse_Location_With_Space_Last_Byte_Returns_InVersion(t *testing.T){
	requestLine := requestLine{
		verb: "get",
		url: "test.com",
	}
	
	actual := get_status_line_parse_location(' ', requestLine)

	if actual != InVersion{
		t.Fatalf("get_stats_line_parse_location(' ', requestLine) = %d, want %d", actual, InVersion)
	}
}