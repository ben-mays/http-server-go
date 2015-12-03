package main

import (
  "bufio"
  "fmt"
  "http"
  HttpServer "http/server"
  "io"
  "os"
)

func main() {
  rootDir := "/tmp/files"
  var sendFileHandler HttpServer.Handler
  sendFileHandler = func(request http.HttpRequest) (http.HttpResponse, error) {
    // get the file
    file := request.URI
    if file == "/" {
      file = "/index.html"
    }
    fmt.Printf("Sending file %s%s\n", rootDir, file)
    f, ferr := os.Open(fmt.Sprintf("%s%s", rootDir, file))
    defer f.Close()
    if ferr != nil {
      return http.NewCloseResponse(http.Protocol11, http.Status200), ferr
    }

    // wrap the file with a reader
    chunksize := 10240 // 10kB chunks
    reader := bufio.NewReader(f)
    body := make([]byte, 0)
    buffer := make([]byte, chunksize)

    for {
      n, err := reader.Read(buffer) // n is the number of bytes read, can be less than chunksize
      if err == io.EOF {
        break
      }
      // TODO: Send chunk encoded responses
      body = append(body, buffer[:n]...) // append buffer[0-n) to the body buffer.
    }

    // TODO: Figure out proper content-length. CURL is reporting excess bytes in the body, should we count the CLRF bytes here?
    headers := []string{"Connection: close", fmt.Sprintf("Content-Length: %d", len(body)+2)}
    return http.NewHttpResponse(http.Protocol11, http.Status200, headers, body), nil
  }

  server := HttpServer.NewServer()
  server.AddRoute(http.GET, "*", sendFileHandler)
  server.Run(8080)
}
