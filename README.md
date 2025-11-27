<h1 align="center">ELH (Embed Languages in HTML)</h1>

>[!WARNING]
>HIGHLY EXPERIMENTAL AND EARLY DEVELOPMENT!

Freely embed various programming languages into HTML; like PHP, but without being bound to one language.

---

### Import module

- The repo was created with capital letters, so it's recommended to import with an alias of `elh`:
  ```go
  import (
      elh "github.com/Supraboy981322/ELH"
  )
  ```

---

### Function signatures
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

---

### Usage examples

- list files in directory with Bash
  ```html
  <!DOCTYPE html>
  <body>
    <h1>files</h1>
    <ul>
      <$bash
        for file in $(ls); do
          printf "<li>$file</li>\n"
        done
      $>
    </ul>
  </body>
  ```
- [more examples](https://github.com/Supraboy981322/ELH/tree/master/docs/examples)

---

### Dependencies
>[!NOTE]
>You only need to install the dependencies for the languages you're using

- JS/TS
  - [bun](https://bun.sh)
- Java
  - [java](https://java.com) (tested with jdk-23)
- Bash
  - [bash](https://gnu.org/software/bash/)
- Python
  - [python3](https://python.org/)
- Markdown
  - [marked](https://github.com/markedjs/marked)
- Lua
  - [lua](https://lua.org)
- Brainfuck
  - beef
- Go
  - [go](https://go.dev)
- R
  - r-base
- Ruby
  - [ruby](https://ruby-lang.org)
- PHP
  - [php](https://php.net) (tested with php8.4-cli)
- Perl
  - [perl](https://perl.org)
- Basic
  - bwbasic
- VimScript
  - [VIMc](https://github.com/Supraboy981322/vimc/)

---

### TODO

>[!NOTE]
>This project will reach v1.0.0 when I finish the TODO list; until then, there will be no version numbering
>  (as I add more to the list as I think of things that it needs for completion) 

~~-----------~~  ***`functionality`***  ~~-----------~~
- [x] HTTP server
- [x] Parsing `.elh` for language tags
- [x] Basic definitions for language tags
- [x] Executing code in language tags
- [x] Injecting formatting requirements for various languages 
- [x] Mime types for files other-than `.elh`
- [ ] Passing headers and params
- [ ] imports
- [ ] Fix indentation bug
- [ ] user-defined languages (?)
- [x] module (?)

~~-------------~~  ***`languages`***  ~~-------------~~
- [x] Go
  - [x] Works
  - [x] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Python
  - [x] Works
  - [x] Imports
  - [ ] Headers
  - [ ] Request params
  - [ ] Fixed all known bugs
    - [ ] whitespace requirements
- [x] Bash
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Java
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] R
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] JavaScript
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Markdown
  - [x] Works
  - [ ] Fixed all known bugs
    - [ ] whitespace requirements
- [x] Ruby
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] PHP
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Perl
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Lua
  - [x] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [x] Fixed all known bugs
- [x] Brainfuck
  - [x] Works
- [x] Basic
  - [x] Works
  - [x] Fixed all known bugs
- [x] VimScript
  - [x] Works
  - [x] Imports
  - [x] Headers
  - [ ] Request params
  - [ ] Fixed all known bugs
    - [ ] Strange behavior when printing integers
- [ ] Dart
  - [ ] Works
  - [ ] Imports
  - [ ] Headers
  - [ ] Request params
  - [ ] Fixed all known bugs
