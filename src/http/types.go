// Collection of useful const types. Using NewXY will coerce the primitive type to the const typed value.
package http

import "errors"

type HttpMethod string

const (
  GET           HttpMethod = "GET"
  POST          HttpMethod = "POST"
  OPTION        HttpMethod = "OPTION"
  InvalidMethod HttpMethod = "INVALID"
)

type HttpProtocol string

const (
  Protocol11      HttpProtocol = "HTTP/1.1"
  Protocol2       HttpProtocol = "HTTP/2.0"
  InvalidProtocol HttpProtocol = "INVALID"
)

type HttpStatusCode string

const (
  StatusCode200     HttpStatusCode = "200"
  StatusCode400     HttpStatusCode = "400"
  StatusCode500     HttpStatusCode = "500"
  InvalidStatusCode                = "INVALID"
)

type HttpStatusMessage string

const (
  Message200           HttpStatusMessage = "OK"
  Message400           HttpStatusMessage = "Bad request"
  Message500           HttpStatusMessage = "Server internal error"
  InvalidStatusMessage                   = "INVALID"
)

type HttpStatus struct {
  StatusCode    HttpStatusCode
  StatusMessage HttpStatusMessage
}

var Status200 = HttpStatus{StatusCode200, Message200}
var Status400 = HttpStatus{StatusCode400, Message400}
var Status500 = HttpStatus{StatusCode500, Message500}
var InvalidStatus = HttpStatus{InvalidStatusCode, InvalidStatusMessage}

func NewHttpProtocol(protocol string) (HttpProtocol, error) {
  switch protocol {
  case "HTTP/1.1":
    return Protocol11, nil
  case "HTTP/2.0":
    return Protocol2, nil
  default:
    return InvalidProtocol, errors.New("Invalid protocol")
  }
}

func NewHttpMethod(method string) (HttpMethod, error) {
  switch method {
  case "GET":
    return GET, nil
  case "POST":
    return POST, nil
  case "OPTION":
    return OPTION, nil
  default:
    return InvalidMethod, errors.New("Invalid method")
  }
}

func NewHttpStatus(status int) (HttpStatus, error) {
  switch status {
  case 200:
    return Status200, nil
  case 400:
    return Status400, nil
  case 500:
    return Status500, nil
  default:
    return InvalidStatus, errors.New("Invalid status")
  }
}
