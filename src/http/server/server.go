// The HTTP package consist of a base of primitives (HttpRequest, HttpResponse, HttpStatus, HttpMethod, HttpProtocol) and a server implementation that handles interfacing with the raw TCP connections.
// To use the server, define a set of handlers and routes using the AddRoute method. Each handler accepts a HttpRequest instance and returns a HttpResponse instance.
package server

import (
  "errors"
  "fmt"
  . "http"
  "net"
)

// Handler is a function type that can be implemented by the consuming library to be bound to a route.
type Handler func(HttpRequest) (HttpResponse, error)

// Private to prevent instantiation outside of the package. Can be instantiated using NewServer.
type server struct {

  // TODO: Improve the routing granularity, using a trie.
  routes map[HttpMethod]map[string]Handler
}

// Returns a reference to a new server instance with no routes specified.
func NewServer() *server {
  s := server{}
  s.routes = make(map[HttpMethod]map[string]Handler)
  return &s
}

// AddRoute adds a new route mapping the given Handler to a given HttpMethod and string path. "*" is considered a default route and must be defined for each HttpMethod.
func (s *server) AddRoute(method HttpMethod, path string, handler Handler) {
  if _, exist := s.routes[method]; !exist {
    s.routes[method] = make(map[string]Handler)
  }
  s.routes[method][path] = handler
  fmt.Println(s.routes)
}

// Route takes a HttpRequest and looks up a matching route in the server's routing mechanism. If no route is found, the router will look for a default route (specified by "*") for the given HttpMethod.
func (s *server) Route(request HttpRequest) (Handler, error) {
  _, hasMethod := s.routes[request.Method]
  if hasMethod {
    if handler, exist := s.routes[request.Method][request.URI]; exist {
      //fmt.Println("Found route for request.")
      // route exist, return the handler
      return handler, nil
    }

    fmt.Println("No route found, trying default route.")

    // does a default handler exist? e.g. "*"
    if handler, exist := s.routes[request.Method]["*"]; exist {
      //fmt.Println("Default route found, using.")
      return handler, nil
    }
  }

  fmt.Println("No routes found for request.", request)
  return nil, errors.New("No routes for request")
}

// closeConn will 'try' to send a HttpResponse to inform the client on the impending close and then close the connection.
// The client should honor the `Connection: close` header, in some cases it will attempt to re-open the connection.
func (s *server) closeConn(tcpConn net.Conn, protocol HttpProtocol, status HttpStatus) {
  defer tcpConn.Close()
  response := NewCloseResponse(protocol, status)
  tcpConn.Write(response.ToBytes())
}

// handleError will accept an error, connection and request. If the error is not nil, it will invoke closeConn on the server with a 500 Server Internal Error.
func (s *server) handleError(err error, conn net.Conn) {
  if err != nil {
    // assume a server internal error for now
    s.closeConn(conn, Protocol11, Status500)
  }
}

// process accepts a net.Conn and parses a HttpRequest, routes the request to a handler, executes the handler, writes a response and then closes the connection.
// Improvements could be made here:
//   * Use a chann primitive to keep the tcp connection alive; e.g HTTP pipelining
func (s *server) process(tcpConn net.Conn) {
  var request *HttpRequest
  var response HttpResponse
  var handler Handler
  var err error

  request, err = RequestFromConn(tcpConn)
  if err == nil {
    handler, err = s.Route(*request)
    if err == nil {
      response, err = handler(*request)
      if err == nil {
        _, err = tcpConn.Write(response.ToBytes())
      }
    }
  }

  if err != nil {
    s.handleError(err, tcpConn)
  }
}

// Run starts listening for new TCP connections on the port given. It launches a isolated goroutine to process each incoming socket connection.
func (s *server) Run(port int) {
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

  if err != nil {
    panic("Error binding!")
  }

  for {
    tcpConn, err := listener.Accept()
    if err != nil {
      fmt.Errorf("Error accepting new connection! %v", err)
    } else {
      go s.process(tcpConn)
    }
  }
}
