package http

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
	Result T
	LastByte byte
	Completed bool
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
	parseLocation := get_status_line_parse_location(parseResult.LastByte, parseResult.Result)

	for i := 0; i <= len(buffer); i++ {
		/* the buffer is initialized with zeros. If we find a zero, return because we
		   know we've hit the end of the message within the buffer since we are hitting
		   the default values
		*/
		if buffer[i] == 0{
			return HttpParseResult[requestLine]{
				Result: httpRequestLine,
				LastByte: 0,
			}
		}


	}

	return HttpParseResult[requestLine]{
		Result: httpRequestLine,
		LastByte: lastByte,
	}
}