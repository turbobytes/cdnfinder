package cdnfinder

import (
	"time"

	"github.com/miekg/dns"
)

func detectcnamesretry(hostname, server string, retry bool) ([]string, error) {
	out := make([]string, 0)
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{hostname, dns.TypeA, dns.ClassINET}
	c := new(dns.Client)
	c.Timeout = time.Second * 5
	msg, _, err := c.Exchange(m1, server)
	if err != nil {
		return out, err
	}
	for _, ans := range msg.Answer {
		if c, ok := ans.(*dns.CNAME); ok {
			out = append(out, c.Target)
		}
	}
	if err != nil {
		if retry {
			//If fail retry once
			return detectcnamesretry(hostname, server, false)
		}
	}
	return out, err
}

func detectcnames(hostname, server string) ([]string, error) {
	return detectcnamesretry(hostname, server, true)
}
