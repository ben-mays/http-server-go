package http

import (
  "bytes"
  "fmt"
)

type HttpResponse struct {
  Protocol HttpProtocol
  Status   HttpStatus
  // TODO: Use map[string]string
  Headers []string
  Body    []byte
}

func (r *HttpResponse) ToBytes() []byte {
  // A proper response looks like the following:
  //
  //  <HTTP PROTOCOL> <STATUS>\r\n
  //  <HEADER 1>\r\n
  //  ...
  //  <HEADER N>\r\n
  //  \r\n
  //  <BODY BYTES>
  var buffer bytes.Buffer // using a byte buffer for efficient concatenation
  buffer.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Protocol, r.Status.StatusCode, r.Status.StatusMessage))

  for _, header := range r.Headers {
    buffer.WriteString(fmt.Sprintf("%s\r\n", header))
  }

  buffer.WriteString("\r\n") // blank line between headers and body

  if r.Body != nil || len(r.Body) != 0 {
    buffer.Write(r.Body)
  }

  return buffer.Bytes()
}

func NewCloseResponse(protocol HttpProtocol, status HttpStatus) HttpResponse {
  headers := []string{"Connection: close"}
  return HttpResponse{protocol, status, headers, nil}
}

func NewHttpResponse(protocol HttpProtocol, status HttpStatus, headers []string, body []byte) HttpResponse {
  return HttpResponse{protocol, status, headers, body}
}
