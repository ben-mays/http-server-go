package main

import (
  "fmt"
  "http/request"
  "http/response"
  "http/types"
  "net"
)

func createChunkResponse(data []byte, request *request.HttpRequest) response.HttpResponse {
  // while not eof, create a response
  headers := make([]string, 0)
  response := response.NewHttpResponse(request.Protocol, types.NewHttpStatus(200), headers, data)
  // create content length
  return response
}

func sendFile(conn net.Conn, request *request.HttpRequest) {
  // get the file
  // file := request.URI
  // rdr = reader.open(file, 'r')
  // while(!rdr.end()) {
  //   buffer := rdr.read(10000) // read 10kB per response
  //   response := createChunkResponse(buffer)
  //   conn.send(response)
  // }
}

func handleConnection(conn net.Conn) {
  fmt.Println("Received request")
  request := request.FromConn(conn)
  response := response.NewCloseResponse(types.Protocol11, types.NewHttpStatus(200))
  response.Write(conn)
  //sendFile(conn, request)
  fmt.Println("Finished sending!")
  conn.Close()
}

func main() {
  listener, err := net.Listen("tcp", ":8080")

  if err != nil {
    fmt.Errorf("Error binding! %v", err)
  }

  for {
    tcpConn, err := listener.Accept()

    if err != nil {
      // handle error
    }

    go handleConnection(tcpConn)
  }
}
