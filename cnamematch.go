package cdnfinder

import (
	"strings"
)

//Given a hostname, detect cname
func HostnametoCDN(hostname, server string) (string, []string, error) {
	if !strings.HasSuffix(hostname, ".") {
		//Make FQDN
		hostname = hostname + "."
	}
	cnames, err := detectcnames(hostname, server)
	if err != nil {
		return "", cnames, err
	}
	for _, cname := range cnames {
		for _, cdn := range cdnmatches {
			if strings.Contains(cname, cdn[0]) {
				return cdn[1], cnames, nil
			}
		}
	}
	return "", cnames, nil
}
