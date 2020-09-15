package rotatify

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func getIP(g *Rotatify) (ip string, err error) {
	resp, err := g.Get("https://httpbin.org/ip")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ip = fmt.Sprintf("%s", b)
	return
}

func TestRotatify(t *testing.T) {
	g := New()
	g.UpdateProxies([]string{
		"http://95.179.165.80:8080",
		"http://8.12.22.124:8080",
	})
	go g.StartRotateProxies()
	defer g.StopRotateProxies()

	g.RotateInterval = 1 * time.Second

	if ip, err := getIP(g); err != nil {
		t.Fatal(err)
	} else {
		println(ip)
	}

	time.Sleep(1 * time.Second)

	if ip, err := getIP(g); err != nil {
		t.Fatal(err)
	} else {
		println(ip)
	}
}
