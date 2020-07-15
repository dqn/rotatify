package guruguru

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Guruguru struct {
	http.Client
	RotateInterval      time.Duration
	Proxies             []url.URL
	proxyIndex          int
	mux                 *sync.Mutex
	stopRotateProxiesCh chan struct{}
}

func New() *Guruguru {
	return &Guruguru{
		RotateInterval:      10 * time.Minute,
		Proxies:             make([]url.URL, 0),
		proxyIndex:          0,
		mux:                 &sync.Mutex{},
		stopRotateProxiesCh: make(chan struct{}),
	}
}

func (g *Guruguru) UpdateProxies(rawURLs []string) error {
	g.mux.Lock()
	defer g.mux.Unlock()

	g.proxyIndex = 0
	g.Proxies = make([]url.URL, 0, len(rawURLs))

	for _, rawURL := range rawURLs {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}

		g.Proxies = append(g.Proxies, *u)
	}

	g.rotateProxy()

	return nil
}

func (g *Guruguru) rotateProxy() {
	if len(g.Proxies) == 0 {
		g.Transport = nil
		return
	}

	g.Transport = &http.Transport{Proxy: http.ProxyURL(&g.Proxies[g.proxyIndex])}
	g.proxyIndex = (g.proxyIndex + 1) % len(g.Proxies)
}

func (g *Guruguru) StartRotateProxies() {
	t := time.NewTicker(g.RotateInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			g.mux.Lock()
			g.rotateProxy()
			g.mux.Unlock()
		case <-g.stopRotateProxiesCh:
			return
		}
	}
}

func (g *Guruguru) StopRotateProxies() {
	g.stopRotateProxiesCh <- struct{}{}
}
