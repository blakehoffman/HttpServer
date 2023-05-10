package http

import (
	"strings"
)

const MaxHttpStatusLineLength = 1024

const (
	InVerb int = iota
	InUrl
	InVersion
)

type requestLine struct {
	verb    string
	url     string
	version string
}

type HttpParseResult[T any] struct {
	result T
	lastByte byte
	completed bool
}

func get_status_line_parse_location(lastByte byte, requestLine requestLine) int{
	if lastByte == 0{
		return InVerb;
	}

	if lastByte == ' '{
		if requestLine.verb != "" && requestLine.url ==  ""{
			return InUrl
		} else{
			return InVersion
		}
	}

	if requestLine.version != ""{
		return InVersion
	} else if requestLine.url != ""{
		return InUrl
	} else{
		return InVerb
	}
}

func parse_http_status_line(buffer []byte, parseResult HttpParseResult[requestLine]) HttpParseResult[requestLine] {
	var httpRequestLine requestLine
	var dataStringBuilder strings.Builder
	parseLocation := get_status_line_parse_location(parseResult.lastByte, parseResult.result)

	// if part of the string was written from the last call to this method, make sure to concatenate it
	if parseLocation == InVerb && parseResult.result.verb != ""{
		dataStringBuilder.WriteString(parseResult.result.verb)
	} else if parseLocation == InUrl && parseResult.result.url != ""{
		dataStringBuilder.WriteString(parseResult.result.url)
	} else if parseLocation == InVersion && parseResult.result.version != ""{
		dataStringBuilder.WriteString(parseResult.result.version)
	}

	for i := 0; i <= len(buffer); i++ {
		/* the buffer is initialized with zeros. If we find a zero, return because we
		   know we've hit the end of the message within the buffer since we are hitting
		   the default values
		*/
		if buffer[i] == 0{
			return HttpParseResult[requestLine]{
				result: httpRequestLine,
				lastByte: 0,
			}
		}

		
	}

	return HttpParseResult[requestLine]{
		Result: httpRequestLine,
		LastByte: lastByte,
	}
}