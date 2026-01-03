# ELH: docs - langs - Bash

The ELH default registry has an entry for the [GNU Bash Unix shell](<https://en.wikipedia.org/wiki/Bash_(Unix_shell)>). The Bash runner for ELH just calls the Bash command on the machine running the server.

### Functions available to the Bash runner

When ELH starts the Bash runner, it generates a temporary file containing functions that Bash sources before running your code. 

>![NOTE]
>Just like logging to the server's log on every other runner, you must enabled the `elh.Server.Log.LogStderr` option when configuring the server.

- Logging to the server
  ```bash
  elhlog "the number %d is neat" "5652733"
  ```
  Function source
  ```bash
  elhlog() (
    set -euo pipefail
    printf "${@}" 1>&2
    printf "\n" 1>&2
  )
  ```
