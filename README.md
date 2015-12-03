# http-server-go
HTTP server in Go

# Quick Start

Simply define a handler that accepts the HttpRequest type and returns a HttpResponse, and wire up a route to it.

```
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
```

# Design

The goal was to create a Server API that was abstracted away from the application logic. 

The core components:
 * HttpRequest
 * HttpResponse
 * Server
 * Handler

 Both the HttpRequest and HttpResponse are structs containing things like the headers, body, protocol, etc.

 ```
 type HttpRequest struct {
  Method   HttpMethod
  URI      string
  Protocol HttpProtocol
  Headers []string
  Body    []byte
}

type HttpResponse struct {
  Protocol HttpProtocol
  Status   HttpStatus
  Headers  []string
  Body     []byte
}
```

The server API has two public methods:

* `Run(port int)`
* `AddRoute(method HttpMethod, path string, handler Handler)`

Routes can be added after the server has started.

Defining a handler is simple; it's just a reference to a function with the header `func(r http.HttpRequest) (http.HttpResponse, error)`.

# Improvements

Create a mechanism for adding middleware, perhaps adding a generic map to requests for storing data.
