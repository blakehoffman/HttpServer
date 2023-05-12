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

	_ = parse_http_status_line(buffer, HttpParseResult[requestLine]{})
}