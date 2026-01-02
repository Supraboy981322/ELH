# ELH: Server - custom runners and registries

ELH exposes all the helper functions used internally to write runners as part of the module

---

## Creating registry with a custom runner

To create a custom runner, you need to have a registry. Custom runners typically point to a function. By defining a function, you override the default ELH logic (which, for most runners, simply runs an os command). Using a custom registry overrides the default registry, meaning anything not defined in your registry won't be available to any functions using that registry. However, there is a helper function to use the default registry for a specific runner (`UseDefault(name string) *ExternalRunner`).

Below is a basic registry that uses the default lua runner, and an os command for Bash, and a custom markdown runner
  ```go
  registry := map[string]elh.Runner {
    "lua": elh.UseDefault("lua"),
    "bash": &elh.ExternalRunner {
      CmdName: "bash",
      Args: []string{}, //passed to the command when executing
      Timeout: 5 * time.Second,
      Env: os.Environ(),
    },
    "md": &elh.ExternalRunner {
      Func: customMarkdownRunner(),
      Timeout: 5 * time.Second,
    },
  }
  ```

---

## Interacting with ELH

ELH exposes structs which're passed to runner functions.

The a struct of type `RunOpts` is directly passed to a runner when it's called by ELH, `RunOpts`.
  ```go
  type RunOpts struct {
    Code string       //the extracted code from the requested ELH file
    TmpDir string     //the path to a temporary directory, if needed
    TmpFile *os.File  //pointer to a temporary file, if needed
    Req *http.Request //the request that was made to the server
    Wr http.ResponseWriter //to respond to the client before the function returns
    Lang string       //the language tag used in the elh openning element ('<$')
    Registry map[string]Runner //the registry the function was called from
    Api API    //used in the default ELH runner ('<$elh'), exposed incase needed
    Log Logger //to log messages before the runner returns
  }
  ```

The `Logger` struct type (`RunOpts.Log`) contains 2 values, `Func` (type `func (string)`) and `Runner` (a struct)

The `Logger` struct type is defined below:
  ```go
  type Logger struct {
    Func func(string)
    Runner struct {
      LogStderr bool
      LogStderrToStdout bool
    }
  }
  ```
The value of `Func` in the `Logger` struct is the same value defined as `ServerOpts.Log.Func` (defined when configuring the ELH server), it's just passed to the `RunOpts` struct for easier access to it. If it's not defined in the server configuration, the value will be `nil`.

Below is an example of a runner that simply logs the "code" (input) and returns it as the result:
  ```go
  func echoRunner(opts elh.RunOpts) (string, string, error) {
    input := opts.Code

    //recommended: check if it's nil
    if opts.Log.Func != nil {
      opts.Log.Func(input)
    } else {
      fmt.Println(input)
    }

    return input, "", nil
  }
  ```

---

## Helper functions

### `ToArgs(src string, by string)`

A single-purpose string parser that converts "code" to slice of strings (`[]string`) for parsing as arguments
  
  (takes src (the "code") and a string to split by)
  ```go
  ToArgs(src string, by string) []string
  ```

### `(opts *RunOpts) Unwrap()`

Returns the basic `RunOpts` as individual values.

(returns `string`, `string`, `*os.File`, `*http.Request`, and `http.ResponseWriter`)

Below is an example
  ```go
  func someRunner(opts elh.RunOpts) (string, string, error) {
    code, tmpDir, tmpFile, request, respWriter := opts.Unwrap()

    ip := request.RemoteAddr
    respWriter.Write([]byte("hello, "+ip))

    print(tmpFile.Name())
    print(tmpDir)

    return code,  "", nil
  }
  ```
