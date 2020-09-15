package rotatify

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Rotatify struct {
	http.Client
	RotateInterval      time.Duration
	Proxies             []url.URL
	proxyIndex          int
	mux                 *sync.Mutex
	stopRotateProxiesCh chan struct{}
}

func New() *Rotatify {
	return &Rotatify{
		RotateInterval:      10 * time.Minute,
		Proxies:             make([]url.URL, 0),
		proxyIndex:          0,
		mux:                 &sync.Mutex{},
		stopRotateProxiesCh: make(chan struct{}),
	}
}

func (r *Rotatify) UpdateProxies(rawURLs []string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.proxyIndex = 0
	r.Proxies = make([]url.URL, 0, len(rawURLs))

	for _, rawURL := range rawURLs {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}

		r.Proxies = append(r.Proxies, *u)
	}

	r.rotateProxy()

	return nil
}

func (r *Rotatify) rotateProxy() {
	if len(r.Proxies) == 0 {
		r.Transport = nil
		return
	}

	r.Transport = &http.Transport{Proxy: http.ProxyURL(&r.Proxies[r.proxyIndex])}
	r.proxyIndex = (r.proxyIndex + 1) % len(r.Proxies)
}

func (r *Rotatify) StartRotateProxies() {
	t := time.NewTicker(r.RotateInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			r.mux.Lock()
			r.rotateProxy()
			r.mux.Unlock()
		case <-r.stopRotateProxiesCh:
			return
		}
	}
}

func (r *Rotatify) StopRotateProxies() {
	r.stopRotateProxiesCh <- struct{}{}
}
