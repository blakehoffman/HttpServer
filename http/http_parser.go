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

func parse_http_status_line(buffer []byte, parseResult HttpParseResult[requestLine]) (HttpParseResult[requestLine], *byte) {
	var httpRequestLine requestLine
	var dataStringBuilder strings.Builder
	var lastByte byte
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
		currentByte := buffer[i]
		/* the buffer is initialized with zeros. If we find a zero, return because we
		   know we've hit the end of the message within the buffer since we are hitting
		   the default values
		*/
		if currentByte == 0{
			return HttpParseResult[requestLine]{
				result: httpRequestLine,
				lastByte: 0,
			}, nil
		}
		
		if currentByte == ' '{
			write_string_builder_to_request_line(parseLocation, dataStringBuilder, &httpRequestLine)
			parseLocation++
			dataStringBuilder.Reset()
		} else if currentByte == '\n' && lastByte == '\r'{
			lastByte = currentByte
			write_string_builder_to_request_line(parseLocation, dataStringBuilder, &httpRequestLine)

			return HttpParseResult[requestLine]{
				result: httpRequestLine,
				lastByte: lastByte,
				completed: true,
			}, &buffer[i]
		} else{
			dataStringBuilder.WriteByte(currentByte)
		}

		lastByte = currentByte
	}

	write_string_builder_to_request_line(parseLocation, dataStringBuilder, &httpRequestLine)
	
	return HttpParseResult[requestLine]{
		result: httpRequestLine,
		lastByte: lastByte,
	}, nil
}

func write_string_builder_to_request_line(parseLocation int, stringBuilder strings.Builder, requestLinePointer *requestLine) {
	if parseLocation == InVerb{
		requestLinePointer.verb = stringBuilder.String()
	} else if parseLocation == InUrl{
		requestLinePointer.url = stringBuilder.String()
	} else if parseLocation == InVersion{
		requestLinePointer.version = stringBuilder.String()
	}
}