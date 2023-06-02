package http

import (
	"strings"
)

const MaxHttpStatusLineLength = 1024

const (
	InVerb int = iota
	InUrl
	InVersion
	AtEnd
)

type httpHeader struct {
	name  string
	value string
}

type requestLine struct {
	verb    string
	url     string
	version string
}

type httpParseHeadersResult struct {
	httpParseResult[[]httpHeader]
	atNewHeader bool
	atHeaderValue bool
	headerParseIndex int
}

type httpParseResult[T any] struct {
	result    T
	lastByte  byte
	completed bool
}

func get_status_line_parse_location(lastByte byte, requestLine requestLine) int {
	if lastByte == 0 {
		return InVerb
	}

	if lastByte == ' ' {
		if requestLine.verb != "" && requestLine.url == "" {
			return InUrl
		} else {
			return InVersion
		}
	}

	if requestLine.version != "" {
		return InVersion
	} else if requestLine.url != "" {
		return InUrl
	} else {
		return InVerb
	}
}

func parse_http_headers(buffer []byte, parseResult httpParseHeadersResult) (httpParseHeadersResult, *byte) {
	var dataStringBuilder strings.Builder
	lastByte := parseResult.lastByte
	write_http_partial_header_to_string_builder(parseResult, dataStringBuilder)

	for i := 0; i <= len(buffer); i++ {
		currentByte := buffer[i]

		// if current byte is zero, we know we've hit the end of data in the buffer
		if currentByte == 0{
			return httpParseHeadersResult {
				httpParseResult: parseResult.httpParseResult,
				atNewHeader: parseResult.atNewHeader,
				atHeaderValue: parseResult.atHeaderValue,
				headerParseIndex: parseResult.headerParseIndex,
			}, nil
		} else if parseResult.atNewHeader && parseResult.lastByte == '\r' && currentByte == '\n'{
			write_string_builder_to_http_header(&parseResult.result[parseResult.headerParseIndex], parseResult.atHeaderValue, dataStringBuilder)
			lastByte = currentByte

			return httpParseHeadersResult { 
				httpParseResult: parseResult.httpParseResult,
			}, &buffer[i]
		}
	}

	return httpParseHeadersResult{}, nil
} 

func parse_http_status_line(buffer []byte, parseResult httpParseResult[requestLine]) (httpParseResult[requestLine], *byte) {
	var dataStringBuilder strings.Builder
	lastByte := parseResult.lastByte
	parseLocation := get_status_line_parse_location(parseResult.lastByte, parseResult.result)

	// if part of the string was written from the last call to this method, make sure to concatenate it
	if parseLocation == InVerb && parseResult.result.verb != "" {
		dataStringBuilder.WriteString(parseResult.result.verb)
	} else if parseLocation == InUrl && parseResult.result.url != "" {
		dataStringBuilder.WriteString(parseResult.result.url)
	} else if parseLocation == InVersion && parseResult.result.version != "" {
		dataStringBuilder.WriteString(parseResult.result.version)
	}

	for i := 0; i <= len(buffer); i++ {
		currentByte := buffer[i]
		/* the buffer is initialized with zeros. If we find a zero, return because we
		   know we've hit the end of the message within the buffer since we are hitting
		   the default values
		*/
		if currentByte == 0 {
			return parseResult, nil
		} else if currentByte == ' ' {
			write_string_builder_to_request_line(parseLocation, dataStringBuilder, &parseResult.result)
			parseLocation++
			dataStringBuilder.Reset()
		} else if currentByte == '\n' && lastByte == '\r' {
			lastByte = currentByte
			write_string_builder_to_request_line(parseLocation, dataStringBuilder, &parseResult.result)

			return httpParseResult[requestLine]{
				result:    parseResult.result,
				lastByte:  lastByte,
				completed: true,
			}, &buffer[i]
		} else if currentByte != '\n' && currentByte != '\r' {
			dataStringBuilder.WriteByte(currentByte)
		}

		lastByte = currentByte
	}

	write_string_builder_to_request_line(parseLocation, dataStringBuilder, &parseResult.result)

	return httpParseResult[requestLine]{
		result:   parseResult.result,
		lastByte: lastByte,
	}, nil
}

func write_http_partial_header_to_string_builder(parseResult httpParseHeadersResult, dataStringBuilder strings.Builder){
	if !parseResult.atNewHeader && !parseResult.atHeaderValue {
		dataStringBuilder.WriteString(parseResult.httpParseResult.result[parseResult.headerParseIndex].name)
	} else if !parseResult.atNewHeader && parseResult.atHeaderValue {
		dataStringBuilder.WriteString(parseResult.httpParseResult.result[parseResult.headerParseIndex].value)
	}
}

func write_string_builder_to_http_header(httpHeader *httpHeader, atHeaderValue bool, stringBuilder strings.Builder) {
	if atHeaderValue {
		httpHeader.value = stringBuilder.String()
	} else {
		httpHeader.name = stringBuilder.String()
	}
}

func write_string_builder_to_request_line(parseLocation int, stringBuilder strings.Builder, requestLinePointer *requestLine) {
	if parseLocation == InVerb {
		requestLinePointer.verb = stringBuilder.String()
	} else if parseLocation == InUrl {
		requestLinePointer.url = stringBuilder.String()
	} else if parseLocation == InVersion {
		requestLinePointer.version = stringBuilder.String()
	}
}
