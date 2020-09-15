package rotatify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRotatify(t *testing.T) {
	r := New()
	r.UpdateProxies([]string{
		"http://xxx.example.com:8000",
		"http://yyy.example.com:9000",
	})

	r.RotateInterval = 500 * time.Millisecond
	go r.StartRotateProxies()
	defer r.StopRotateProxies()

	i := r.proxyIndex

	time.Sleep(1 * time.Second)

	assert.NotEqual(t, i, r.proxyIndex)
}
