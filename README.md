# ELH (Embed Languages in HTML)

>[!WARNING]
>HIGHLY EXPERIMENTAL AND EARLY DEVELOPMENT!

Freely embed various programming languages into HTML; like PHP, but without being bound to one language.

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
  - bun
- Java
  - java (tested with jdk-23)
- Bash
  - bash
- Python
  - python3
- Markdown
  - marked
- Lua
  - lua
- Brainfuck
  - beef
- Go
  - go
- R
  - r-base
- Ruby
  - ruby
- PHP
  - php (tested with php8.4-cli)
- Perl
  - perl
- Basic
  - bwbasic

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

~~-------------~~  ***`languages`***  ~~-------------~~
- [x] Go
- [x] Python
- [x] Bash
- [x] Java
- [x] R
- [x] JavaScript
- [x] Markdown
- [x] Ruby
- [x] PHP
- [x] Perl
- [x] Lua
- [x] Brainfuck
- [ ] Basic
- [ ] Dart
