# http-server-go
HTTP server in Go

Currently only supports HTTP/1.1 (possibly HTTP/1.0?)

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

```
  rootDir := "/tmp/files"
  fileTypes := map[string]string{
    "jpg":  "image/jpeg",
    "jpeg": "image/jpeg",
    "png":  "image/png",
    "html": "text/html",
    "mp4":  "video/mp4",
  }

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

    // detect content type by parsing the file extension
    contentType := "application/octet-bytes"
    if strings.Contains(file, ".") {
      fileTypeArr := strings.Split(file, ".")
      fileType := fileTypeArr[len(fileTypeArr)-1]
      if val, exist := fileTypes[fileType]; exist {
        contentType = val
      }
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

    headers := []string{"Connection: close", fmt.Sprintf("Content-Length: %d", len(body)), fmt.Sprintf("Content-Type: %s", contentType)}
    return http.NewHttpResponse(http.Protocol11, http.Status200, headers, body), nil
  }

  server := HttpServer.NewServer()
  server.AddRoute(http.GET, "*", sendFileHandler)
  server.Run(8080)
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

* Create a mechanism for adding middleware, perhaps adding a generic map to requests for storing data.
* Add CGI support
** For TB: Create dynamic (use cookies) web app in another language and served via CGI

CGI : Accepts input from ENV vars, forks a binary and returns STDOUT as response