# Grace pkg

Package with single function for create base context which will be canceled on signals:
`SIGINT`, `SIGTERM`, `SIGHUP`.

## Example

```go
package main

import (
    "context"

    "github.com/chapsuk/grace"
)

func main() {
    ctx := grace.ShutdownContext(context.Background())
    <-ctx.Done()
    // do graceful shutdown after context was canceled
}
```
