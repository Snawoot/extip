# extip

Go package which retrieves external IP address using STUN servers

---

:heart: :heart: :heart:

You can say thanks to the author by donations to these wallets:

- ETH: `0xB71250010e8beC90C5f9ddF408251eBA9dD7320e`
- BTC:
  - Legacy: `1N89PRvG1CSsUk9sxKwBwudN6TjTPQ1N8a`
  - Segwit: `bc1qc0hcyxc000qf0ketv4r44ld7dlgmmu73rtlntw`

---

## Usage

Example:

```go
package main

import (
    "fmt"
    "github.com/Snawoot/extip"
    "context"
    "os"
    "time"
)

func main() {
    ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)
    ip, err := extip.QueryMultipleServers(ctx, nil, 2, false)
    if err != nil {
        switch res := err.(type) {
        case extip.InconclusiveResult:
            fmt.Fprintf(os.Stderr, "Inconclusive result:\n")
            fmt.Fprintf(os.Stderr, "Required quorum = %v\n", res.Quorum)
            for k, v := range res.Results {
                fmt.Fprintf(os.Stderr, "Server %s responded: %s\n", k, v)
            }
            for k, v := range res.Errors {
                fmt.Fprintf(os.Stderr, "Server %s failed: %v\n", k, v)
            }
        default:
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
        return
    }
    fmt.Println(ip)
}
```

See GoDoc and sources for more details.
