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

## Dig usage example

```go
package main

import (
    "github.com/chapsuk/grace"
    "go.uber.org/dig"
)

func main() {
    c := dig.New()

    c.Provide(grace.NewShutdownContext)
    c.Invoke(func(p grace.ContextParams) {
        // nodes: {
        //     context.Context[name="grace_context"] -> deps: [], ctor: func() grace.ContextResult
        // }
        // values: {
        //     context.Context[name="grace_context"] => context.Background.WithCancel
        // }
        <-p.Context.Done()
    })
}
```
