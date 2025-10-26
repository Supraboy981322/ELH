# ELH (Embed Languages in HTML)

>[!WARNING]
>HIGHLY EXPERIMENTAL AND EARLY DEVELOPMENT!

Freely embed various programming languages into HTML; like PHP, but without PHP.

---

### known working languages:
- Java (tested with Java 23)
- Python (tested with Python 3)
- Bash
- JavaScript (tested with Bun) 

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
- [more examples](https://github.com/Supraboy981322/ELH/tree/master/examples)
