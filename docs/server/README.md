# ELH: Server Doc

## Helper function usage examples

- HTTP server handler
  ```go
  package main
  
  import (
    "log"
    "net/http"
    elh "github.com/Supraboy981322/ELH"
  )

  func main() {
    elh.Logger = func(str string) { log.Print(str) }
    http.HandleFunc("/", elh.HttpServer)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      log.Fatal("server failed:  %v", err)
    }
  }
  ```
  This function automatically handles http GET requests to the server, so your project can just serve files without having to write the singular function to do it yourself.

---

For the function signatures see [this doc](https://github.com/Supraboy981322/ELH/tree/master/docs/server/funcs.md)
