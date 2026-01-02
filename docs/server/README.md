# ELH: Server Doc

## Http handler helper function

- HTTP server handler
  ```go
  package main
  
  import (
    "log"
    "net/http"
    elh "github.com/Supraboy981322/ELH"
  )

  func main() {
    //optional: configure the server 
    server := elh.ServerOpts{
      Registry: elh.DefaultRegistry(),
      Log: elh.Logger{
        Func: func(str string) { log.Print(str) },
      },
    }

    //create the http handler
    http.HandleFunc("/", server.HttpHandler)

    //start the server 
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      log.Fatal("server failed:  %v", err)
    }
  }
  ```
  This function handles http GET requests to the server, so your project can just serve files (including rendering `.elh` files) without having to write the singular function to do it yourself.

---

## Custom registry and/or runner

- To create a custom runner (eg: a script/system command or a function), you have to create a registry with a custom runner.

You can define a registry like so:
  ```go
  registry := map[string]elh.Runner {
    "notify_access": &elh.ExternalRunner {
      CmdName: "curl",
      Args: []string{"example.com", "-d", "someone accessed the dashboard"},
      Timeout: 5 * time.Second,
      Env: os.Environ(),
    },
  }
  ```

See the [custom runners doc](custom_runner.md) for more information about custom runners and registries.

---

For the function signatures see [this doc](https://github.com/Supraboy981322/ELH/tree/master/docs/server/funcs.md)
