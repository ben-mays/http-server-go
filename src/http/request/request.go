package request

import (
  "bufio"
  "fmt"
  "http/types"
  "net"
  "os"
  "strings"
)

type HttpRequest struct {
  Method   types.HttpMethod
  URI      string
  Protocol types.HttpProtocol
  Headers  []string
  Body     []byte
}

func FromConn(conn net.Conn) *HttpRequest {
  reader := bufio.NewReader(conn)
  scanner := bufio.NewScanner(reader)

  // parse the request line; e.g. GET /path/to/file/index.html HTTP/1.0
  scanner.Scan()                               // populate the internal buffer
  requestStr := scanner.Text()                 // read the scanner buffer as a string
  requestArr := strings.Split(requestStr, " ") // split it into a []string

  if len(requestArr) != 3 {
    panic("Invalid Request!")
  }

  method := requestArr[0]
  uri := requestArr[1]
  protocol := requestArr[2]

  // Scan for headers until we hit a blank line, this signals the body is following.
  headers := make([]string, 0)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
    if line == "" {
      break
    }
    headers = append(headers, line)
  }

  body := make([]byte, 0)

  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "Error parsing request", err)
  }

  return &HttpRequest{types.NewHttpMethod(method), uri, types.NewHttpProtocol(protocol), headers, body}
}
