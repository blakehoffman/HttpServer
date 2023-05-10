package http

type HttpRequestType int

const (
	HttpGet HttpRequestType = iota
	HttpPost
	HttpPut
	HttpDelete
	HttpHead
)