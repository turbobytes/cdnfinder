package cdnfinder

import (
	"net/http"
	"testing"
)

func TestHeader(t *testing.T) {
	cases := [][]string{
		{"SerVeR", "cLoudflarE-nginx", "Cloudflare"},
		{"powered-by-chinacache", "meh", "ChinaCache"},
		{"x-edge-Location", "whatever", "OnApp"},
		{"x-amz-cF-id", "whatever", "Amazon Cloudfront"},
		{"Via", "something.bitgravity.com:3826", "Bitgravity"},
		{"Via", "foo.somethingelse.com", ""},
		{"X-CDN-Provider", "whatever", ""},
		{"X-CDN-Provider", "SkyparkCDN", "Skypark"},
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
