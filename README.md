# guruguru

Proxy-Rotatable HTTP(S) client.

## Installation

```bash
$ go get github.com/dqn/guruguru
```

## Usage

`guruguru.Guruguru` includes all `http.Client` properties and methods. Supports HTTP or SOCKS5 protocol for proxy.

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
  println(resp)
  defer resp.Body.Close()
  b, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(b))
}

func main() {
  g := guruguru.New()
  g.UpdateProxies([]string{
    "http://http.example.com:8080",
    "socks5://socks5.example.com:7070",
  })

  go g.StartRotateProxies()
  defer g.StopRotateProxies()

  g.RotateInterval = 10 * time.Second

  printIP(g) // => XXX.XXX.XXX.XXX
  time.Sleep(10 * time.Second)
  printIP(g) // => YYY.YYY.YYY.YYY
}

```

## License

MIT
