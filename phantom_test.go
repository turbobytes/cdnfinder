package cdnfinder

import (
	"testing"
	"time"
)

func TestPhantom(t *testing.T) {
	Init()
	_, err := discoverResources("http://www.sajalkayan.com/", time.Minute, false)
	if err != nil {
		t.Error(err)
	}
}

func TestPhantomTimeout(t *testing.T) {
	Init()
	_, err := discoverResources("http://blackhole.webpagetest.org/", time.Second, false)
	if err == nil {
		t.Errorf("Expected it to timeout")
	}
}
