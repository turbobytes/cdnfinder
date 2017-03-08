package cdnfinder

import (
	"net/http"
	"testing"
)

func TestHeader(t *testing.T) {
	cases := [][]string{
		[]string{"SerVeR", "cLoudflarE-nginx", "Cloudflare"},
		[]string{"powered-by-chinacache", "meh", "ChinaCache"},
		[]string{"x-edge-Location", "whatever", "OnApp"},
		[]string{"x-amz-cF-id", "whatever", "Amazon Cloudfront"},
		[]string{"Via", "something.bitgravity.com:3826", "Bitgravity"},
		[]string{"Via", "foo.somethingelse.com", ""},
		[]string{"X-CDN-Provider", "whatever", ""},
		[]string{"X-CDN-Provider", "SkyparkCDN", "Skypark"},
	}
	for _, tcase := range cases {
		hdr := make(http.Header)
		hdr.Set(tcase[0], tcase[1])
		cdn := headerguessStr(&hdr)
		if cdn != tcase[2] {
			t.Errorf("Header key: %s, Value: %s should have returned %s. Got: %s", tcase[0], tcase[1], tcase[2], cdn)
		}
	}
}
