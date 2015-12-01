package response

import (
  "bytes"
  "fmt"
  "http/types"
)

type HttpResponse struct {
  Protocol types.HttpProtocol
  Status   types.HttpStatus
  Headers  []string
  Body     []byte
}

func (r *HttpResponse) ToBytes() []byte {
  // A proper response looks like the following:
  //
  //  <HTTP PROTOCOL> <STATUS>\r\n
  //  <HEADER 1>\r\n
  //  ...
  //  <HEADER N>\r\n
  //  \r\n
  //  <BODY BYTES>\r\n
  //  \r\n
  //
  var buffer bytes.Buffer // using a byte buffer for efficient concatenation
  buffer.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Protocol, r.Status.StatusCode, r.Status.StatusMessage))

  for _, header := range r.Headers {
    buffer.WriteString(fmt.Sprintf("%s\r\n", header))
  }

  buffer.WriteString("\r\n") // blank line between headers and body

  buffer.Write(r.Body)
  buffer.WriteString("\r\n\r\n")

  fmt.Println(buffer.String())
  return buffer.Bytes()
}

func NewCloseResponse(protocol types.HttpProtocol, status types.HttpStatus) HttpResponse {
  headers := []string{"Connection: close"}
  return HttpResponse{protocol, status, headers, nil}
}

func NewHttpResponse(protocol types.HttpProtocol, status types.HttpStatus, headers []string, body []byte) HttpResponse {
  return HttpResponse{protocol, status, headers, body}
}
