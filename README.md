# guruguru

Proxy-Rotatable HTTP(S) client.

## Installation

```bash
$ go get github.com/dqn/guruguru
```

## Usage

`guruguru.Guruguru` includes all `http.Client` properties and methods. Supports HTTP and SOCKS5 protocol for proxy.

```go
package main

import (
  "fmt"
  "io/ioutil"
  "time"

  "github.com/dqn/guruguru"
)

func printIP(g *guruguru.Guruguru) {
  resp, _ := g.Get("https://ifconfig.me")
  defer resp.Body.Close()
  b, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(b))
}

func main() {
  g := guruguru.New()

  g.UpdateProxies([]string{
    "http://XXX.XXX.XXX.XXX:8080",
    "socks5://YYY.YYY.YYY.YYY:7070",
  })
  g.RotateInterval = 10 * time.Second

  go g.StartRotateProxies()
  defer g.StopRotateProxies()

  printIP(g) // => XXX.XXX.XXX.XXX
  time.Sleep(10 * time.Second)
  printIP(g) // => YYY.YYY.YYY.YYY
}

```

## License

MIT
