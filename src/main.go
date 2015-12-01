package main

import (
  "fmt"
  "http/request"
  "http/response"
  "http/server"
  "http/types"
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

func handler() server.Handler {
  return func(r request.HttpRequest) (response.HttpResponse, error) {
    fmt.Println("Received request")

    headers := []string{"Connection: keepalive", "Content-Size: 15"}
    body := []byte("Here is some data!")
    response := response.NewHttpResponse(r.Protocol, types.NewHttpStatus(200), headers, body)
    //sendFile(conn, request)
    fmt.Println("Finished creating response!")
    return response, nil
  }
}

func main() {
  server := server.NewServer()
  server.AddRoute(types.NewHttpMethod("GET"), "*", handler())
  server.Run(8080)
}
