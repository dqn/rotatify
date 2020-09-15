# rotatify

Proxy-Rotatable HTTP(S) client.

## Installation

```bash
$ go get github.com/dqn/rotatify
```

## Usage

`rotatify.Rotatify` includes all `http.Client` properties and methods. Supports HTTP and SOCKS5 protocol for proxy.

```go
package main

import (
  "fmt"
  "io/ioutil"
  "time"

  "github.com/dqn/rotatify"
)

func printIP(r *rotatify.Rotatify) {
  resp, _ := r.Get("https://ifconfig.me")
  defer resp.Body.Close()
  b, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(b))
}

func main() {
  r := rotatify.New()

  r.UpdateProxies([]string{
    "http://XXX.XXX.XXX.XXX:8080",
    "socks5://YYY.YYY.YYY.YYY:7070",
  })
  r.RotateInterval = 10 * time.Second

  go r.StartRotateProxies()
  defer r.StopRotateProxies()

  printIP(r) // => XXX.XXX.XXX.XXX
  time.Sleep(10 * time.Second)
  printIP(r) // => YYY.YYY.YYY.YYY
}

```

## License

MIT
