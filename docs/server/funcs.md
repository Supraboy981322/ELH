# ELH Function signatures

- Render from source
  (Returns the result, and an error)
  ```go
  Render(src string, r *http.Request, wr http.ResponseWriter) (string, error)
  ```

- Render with a custom registry
  (Returns the result, and an error)
  ```go
  RenderWithRegistry(src string, registry map[string]Runner, r *http.Request, wr *http.ResponseWriter) (string, error)
  ```

- Construct the ELH data API struct for a custom runner
  (returns type `elh.API`)
  ```go
  BldApiStruct(req *http.Request, wr http.ResponseWriter) API
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
  The resulting runner looks something like this:
  ```go
  cmd: &elh.ExternalRunner{
    Cmdname: cmd,
    Args: args,
    Timeout: 5 * time.Second,
    Env: os.Environ()
  }
  ```

- Read and render a file from the path/name
  (Returns the result if it's an ELH file, otherwise returns unchanged, and an error) 
  ```go
  RenderFile(file string, r *http.Request) ([]byte, error)
  ```

- Get the default registry
  (returns a registry, which is of type `map[string]elh.Runner`)
  ```go
  DefaultRegistry() map[string]Runner
  ```

- Use the default registry entry for a runner
  (returns type `*elh.ExternalRunner`)
  ```go
  UseDefault(name string) *ExternalRunner
  ```

- Run code (expects HTML to be pre-stripped) 
  (returns the both stdout and stderr as strings, and error)
  ```go
  (r *ExternalRunner) Run(code string, tmp *os.File) (string, string, error)
  ```

- Convert "code" to slice of args (`[]string`) for parsing in a custom runner (takes src of type string and a string to split by
  ```go
  ToArgs(src string, by string) []string
  ```

- http handler helper
  ```go
  (server *ServerOpts) HttpHandler(w http.ResponseWriter, r *http.Request)
  ```

---

>[!WARNING]
>The following functions need to be updated, functionality not guaranteed to be 100% up-to-date with other functions

- Serve a file with `"net/http"` and auto detect ELH files
  (Returns the relative filepath, for logging, and an error)
  ```go
  Serve(w http.ResponseWriter, r *http.Request) (string, error)
  ```

- Serve elh file with a custom registry
  (returns the relative filepath, for logging, and an error)
  ```go
  ServeWithRegistry(w http.ResponseWriter, r *http.Request, registry map[string]Runner) (string, error)
 ```
