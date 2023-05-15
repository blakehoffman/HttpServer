package http

import (
	"fmt"
	"net"
)

type HttpRequestType int

const (
	HttpGet HttpRequestType = iota
	HttpPost
	HttpPut
	HttpDelete
	HttpHead
)

func read_http_request(connection net.Conn){
	buffer := make([]byte, 1024)
	_, error := connection.Read(buffer)

	if error != nil{
		fmt.Println("Error reading buffer", error.Error())
	}

	requestLineHttpParseResult, bufferPointer := parse_http_status_line(buffer, HttpParseResult[requestLine]{})

	//continue reading from connection and parsing http status line if it didn't complete on first call to parse_http_status_line
	for completed := requestLineHttpParseResult.completed; !completed; completed = requestLineHttpParseResult.completed{
		connection.Read(buffer)
		requestLineHttpParseResult, bufferPointer = parse_http_status_line(buffer, requestLineHttpParseResult)
	}

	fmt.Println(requestLineHttpParseResult)
	fmt.Println(bufferPointer)
}