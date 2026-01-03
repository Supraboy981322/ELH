# ELH: docs - langs - JS

Thanks to the [v8go module](github.com/rogchap/v8go), ELH as built-in runner for JS, meaning no external binaries required.

In order to integrate v8go into the ELH codebase, some functions (like `console.log`) had to be substituted with injected functions.

### Function substitutes

- `console`

  Functions within the `console` object in JS had to be substituted with new functions within ELH's injected `elh` object for JS. (I tried to use a custom object named `console` so you don't have to remember not to use the `console` object, but it just wouldn't work, strangely, it didn't even return an error)

  The injected `elh` object contains the following functions

  Printing to stdout on the server (not to the HTML document, unlike other runners). It should be noted that this, as with any logging within runners, you must enable the `elh.Server.Log.LogStderr` option when configuring your server, otherwise it's discarded along with all other stderr output from runners.
  ```js
  elh.log(...);
  ```

  Printing an error to stdout on the server. (same `elh.Server.Log.LogStderr` note from above also applies here)
  ```js
  elh.error(...);
  ```

  Writing to the HTML document
  ```js
  elh.print(...);
  ```

---

### Basic example

```elh
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>foo</title>
  </head>
  <body>
    <$js
      elh.log("writing 10 times to the document...");
      for (i = 1; i <= 10; i++) {
        elh.print("<p>"+i+"</p>);
      }
      elh.log("wrote 10 times to the document");
    $>
  </body>
</html>
```
