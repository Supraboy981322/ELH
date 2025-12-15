# ELH Function signatures

- Render from source
  (Returns the result, and an error)
  ```go
  Render(src string, r *http.Request) (string, error)
  ```

- Render with a custom registry
  (Returns the result, and an error)
  ```go
  RenderWithRegistry(src string, registry map[string]Runner, r *http.Request) (string, error)
  ```

- Make registry
  (returns a registry, which is of type `map[string]elh.Runner`)
  ```go
  MkReg(caller string, cmd string, args []string, timeout int, env []string) map[string]Runner
  ```

- Make a registry using the defaults
  (returns a registry, which is of type `map[string]elh.Runner`)
  ```go
  MkRegDefaults(cmd string, args []string) map[string]Runner
  ```

- Serve a file with `"net/http"` and auto detect ELH files
  (Returns the relative filepath, for logging, and an error)
  ```go
  Serve(w http.ResponseWriter, r *http.Request) (string, error)
  ```

- Read and render a file from the path/name
  (Returns the result if it's an ELH file, otherwise returns unchanged, and an error) 
  ```go
  RenderFile(file string, r *http.Request) ([]byte, error)
  ```

- Get the default registry
  (returns a registry, which is of type `map[string]elh.Runner`)
  ```go
  DefaultRegistry()
  ```

- Run code (expects HTML to be pre-stripped) 
  (returns the both stdout and stderr as strings, and error)
  ```go
  (r *ExternalRunner) Run(code string, tmp *os.File) (string, string, error)
  ```

- http handler helper
  ```go
  HttpServer(w http.ResponseWriter, r *http.Request)
  ```
