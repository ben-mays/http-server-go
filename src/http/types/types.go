package types

// This a collection of 'enums' to be used by the server.

// Currently only support GET.
type HttpMethod string

const (
  GET HttpMethod = "GET"
)

// Currently only implement HTTP/1.1
type HttpProtocol string

const (
  Protocol11 HttpProtocol = "HTTP/1.1"
  Protocol2  HttpProtocol = "HTTP/2.0"
)

type HttpStatusCode string

const (
  Status200 HttpStatusCode = "200"
  Status400 HttpStatusCode = "400"
  Status500 HttpStatusCode = "500"
)

type HttpStatusMessage string

const (
  Message200 HttpStatusMessage = "OK"
  Message400 HttpStatusMessage = "Bad request"
  Message500 HttpStatusMessage = "Server internal error"
)

type HttpStatus struct {
  StatusCode    HttpStatusCode
  StatusMessage HttpStatusMessage
}

func NewHttpProtocol(protocol string) HttpProtocol {
  switch protocol {
  case "HTTP/1.1":
    return Protocol11
  case "HTTP/2.0":
    return Protocol2
  default:
    panic("Invalid protocol")
  }
}

func NewHttpMethod(method string) HttpMethod {
  switch method {
  case "GET":
    return GET
  default:
    panic("Invalid method")
  }
}

func NewHttpStatus(status int) HttpStatus {
  switch status {
  case 200:
    return HttpStatus{Status200, Message200}
  case 400:
    return HttpStatus{Status400, Message400}
  case 500:
    return HttpStatus{Status500, Message500}
  default:
    panic("Invalid status")
  }
}
