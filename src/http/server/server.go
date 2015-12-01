package server

import (
  "errors"
  "fmt"
  "http/request"
  "http/response"
  "http/types"
  "net"
)

type Handler func(request.HttpRequest) (response.HttpResponse, error)

type server struct {
  routes map[types.HttpMethod]map[string]Handler
}

func NewServer() *server {
  s := server{}
  s.routes = make(map[types.HttpMethod]map[string]Handler)
  return &s
}

func (s *server) AddRoute(method types.HttpMethod, path string, handler Handler) {
  if _, exist := s.routes[method]; !exist {
    s.routes[method] = make(map[string]Handler)
  }
  s.routes[method][path] = handler
  fmt.Println(s.routes)
}

func (s *server) route(request request.HttpRequest) (Handler, error) {
  _, hasMethod := s.routes[request.Method]
  if hasMethod {
    if handler, exist := s.routes[request.Method][request.URI]; exist {
      // route exist, return the handler
      return handler, nil
    }

    fmt.Println("No route found, trying default")

    // does a default handler exist? e.g. "*"
    if handler, exist := s.routes[request.Method]["*"]; exist {
      return handler, nil
    }
  }

  fmt.Println("No routes found for new request", request)
  return nil, errors.New("No routes for request")
}

func (s *server) sendCloseSignal(tcpConn net.Conn, request request.HttpRequest, statusCode int) {
  response := response.NewCloseResponse(request.Protocol, types.NewHttpStatus(statusCode))
  tcpConn.Write(response.ToBytes())
  tcpConn.Close()
}

func (s *server) process(tcpConn net.Conn) {
  request, err1 := request.FromConn(tcpConn)
  if err1 != nil {
    fmt.Println("Recovered in process, error in HttpRequest.FromConn", err1)
    s.sendCloseSignal(tcpConn, *request, 500)
    return
  }

  handler, err2 := s.route(*request)
  if err2 != nil {
    fmt.Println("Recovered in process, error in routing", err2)
    s.sendCloseSignal(tcpConn, *request, 400)
    return
  }

  response, err3 := handler(*request)
  if err3 != nil {
    fmt.Println("Recovered in process, error in handler execution", err3)
    s.sendCloseSignal(tcpConn, *request, 500)
    return
  }

  _, err4 := tcpConn.Write(response.ToBytes())
  if err4 != nil {
    fmt.Println("Recovered in process, error writing response to socket", err4)
    s.sendCloseSignal(tcpConn, *request, 500)
    return
  }

  // send a OK, Connection: Close
  s.sendCloseSignal(tcpConn, *request, 200)
}

func (s *server) Run(port int) {
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

  if err != nil {
    panic("Error binding!")
  }

  for {
    tcpConn, err := listener.Accept()

    if err != nil {
      fmt.Errorf("Error accepting new connection! %v", err)
    }

    go s.process(tcpConn)
  }
}
