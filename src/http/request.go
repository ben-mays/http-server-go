package http

import (
  "bufio"
  "errors"
  "net"
  "strconv"
  "strings"
)

type HttpRequest struct {
  Method   HttpMethod
  URI      string
  Protocol HttpProtocol
  // TODO: Use map[string][string] instead of []string for headers
  Headers []string
  Body    []byte
}

// parse the request line; e.g. GET /path/to/file/index.html HTTP/1.0
func parseStatusLine(scanner *bufio.Scanner) (HttpMethod, string, HttpProtocol, error) {
  scanner.Scan()                               // populate the internal buffer
  requestStr := scanner.Text()                 // read the scanner buffer as a string
  requestArr := strings.Split(requestStr, " ") // split it into a []string

  if len(requestArr) < 3 {
    return InvalidMethod, "", InvalidProtocol, errors.New("Invalid request!")
  }

  method, methodParsingErr := NewHttpMethod(requestArr[0])
  protocol, protocolParsingErr := NewHttpProtocol(requestArr[2])

  if methodParsingErr != nil {
    return InvalidMethod, "", InvalidProtocol, methodParsingErr
  }

  if protocolParsingErr != nil {
    return InvalidMethod, "", InvalidProtocol, protocolParsingErr
  }

  return method, requestArr[1], protocol, nil
}

// Parses header text out of a bufio.Scanner, should be called immediately after parseStatusLine.
func parseHeaders(scanner *bufio.Scanner, method HttpMethod) []string {
  headers := make([]string, 0)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    line := scanner.Text()

    // Break if we hit a empty line.
    if len(line) == 0 {
      break
    }

    headers = append(headers, line)
  }
  return headers
}

// Returns the value of the content length header or an error if it doesn't exist.
func getContentLength(headers []string) (int64, error) {
  contentLengthStr := "Content-Length:"
  for _, header := range headers {
    if strings.Contains(header, contentLengthStr) {
      sArr := strings.Split(header, " ")
      contentLength, _ := strconv.ParseInt(sArr[1], 0, 64)
      return contentLength, nil
    }
  }
  return -1, errors.New("No content-length header")
}

func parseBody(scanner *bufio.Scanner, contentLength int64) []byte {
  bytes := make([]byte, contentLength)
  scanner.Split(bufio.ScanBytes)
  var i int64
  for i = 0; i < contentLength; i++ {
    scanner.Scan()
    bytes[i] = scanner.Bytes()[0]
  }
  return bytes
}

// This is only public for testing, any implementation should use the conn interface function, RequestFromConn
func RequestFromReader(reader *bufio.Reader) (*HttpRequest, error) {
  scanner := bufio.NewScanner(reader)

  method, uri, protocol, err := parseStatusLine(scanner)
  if err != nil {
    return nil, err
  }

  // Scan for headers until we hit a blank line, this signals the body is following.
  headers := parseHeaders(scanner, method)

  var body []byte
  if method == POST {
    contentLength, err := getContentLength(headers)
    if err != nil {
      return nil, errors.New("Content-Length header not set for POST request.")
    }
    body = parseBody(scanner, contentLength)
  }

  request := HttpRequest{method, uri, protocol, headers, body}
  return &request, nil
}

// Returns a HttpRequest instance parsed from the given net.Conn instance.
func RequestFromConn(conn net.Conn) (*HttpRequest, error) {
  reader := bufio.NewReader(conn)
  return RequestFromReader(reader)
}
