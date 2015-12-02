package main

import (
  "fmt"
  "http"
  HttpServer "http/server"
)

// func createChunkResponse(data []byte, request *request.HttpRequest) response.HttpResponse {
//   // while not eof, create a response
//   headers := make([]string, 0)
//   response := response.NewHttpResponse(request.Protocol, types.NewHttpStatus(200), headers, data)
//   // create content length
//   return response
// }

// func sendFile(request *request.HttpRequest) response.HttpResponse {
//   // get the file
//   // file := request.URI
//   // rdr = reader.open(file, 'r')
//   // while(!rdr.end()) {
//   //   buffer := rdr.read(10000) // read 10kB per response
//   //   response := createChunkResponse(buffer)
//   //   conn.send(response)
//   // }
// }

func main() {
  var echoHandler HttpServer.Handler
  echoHandler = func(r http.HttpRequest) (http.HttpResponse, error) {
    bodySize := len(r.Body) + 4 // Even an empty body takes up a word
    contentLength := fmt.Sprintf("Content-Length: %d", bodySize)
    headers := []string{"Connection: close", "Content-Type: plain/text", contentLength}
    response := http.NewHttpResponse(r.Protocol, http.Status200, headers, r.Body)
    return response, nil
  }
  echoServer := HttpServer.NewServer()
  echoServer.AddRoute(http.POST, "*", echoHandler)
  echoServer.Run(8080)
}
