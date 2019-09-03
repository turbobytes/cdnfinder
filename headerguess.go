package cdnfinder

import (
	"net/http"
	"strings"
)

//TODO: define this somehow in json...
func headerguessStr(hdr *http.Header) string {
	//Cloudflare advertises a custom Server header
	if strings.ToLower(hdr.Get("Server")) == "cloudflare-nginx" {
		return "Cloudflare"
	}
	//China cache sends a Powered-By-Chinacache header
	if hdr.Get("powered-by-chinacache") != "" {
		return "ChinaCache"
	}
	//OnApp edge servers use X-Edge-Location to indicate the location
	if hdr.Get("x-edge-location") != "" {
		return "OnApp"
	}
	//CloudFront adds in some custom tracking id
	if hdr.Get("x-amz-cf-id") != "" {
		return "Amazon Cloudfront"
	}
	//Bitgravity adds edge hostname to Via header
	if strings.Contains(strings.ToLower(hdr.Get("via")), "bitgravity.com") {
		return "Bitgravity"
	}
	//Skypark sends a X header with their brand name
	if hdr.Get("X-CDN-Provider") == "SkyparkCDN" {
		return "Skypark"
	}
	//BaishanCloud uses BC prefix in X-Ser header
	if strings.HasPrefix(hdr.Get("X-Ser"), "BC") {
		return "BaishanCloud"
	}
	return ""
}

func headerguess(hdr *http.Header) *string {
	cdn := headerguessStr(hdr)
	if cdn != "" {
		return &cdn
	}
	return nil
}
