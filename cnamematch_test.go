package cdnfinder

import (
	"fmt"
	"net"
	"testing"

	"github.com/miekg/dns"
)

// Returns an available UDP port from kernel
func getfreeport() (int, error) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.LocalAddr().(*net.UDPAddr).Port, nil
}

type tstcase struct {
	cnames []string
	cdn    string
}

func TestCNAME(t *testing.T) {
	Init()
	cases := make(map[string]tstcase)
	cases["tb.foo.pulse."] = tstcase{[]string{"b.c.d.e.", "something.clients.turbobytes.net.", "somethingelse.example.com."}, "TurboBytes"}
	cases["nobody.foo.pulse."] = tstcase{[]string{"b.c.d.e.", "something.clients.someoneelse.net.", "somethingelse.example.com."}, ""}
	cases["tbcdn.foo.pulse."] = tstcase{[]string{"b.c.d.e.", "something.clients.someoneelse.net.", "somethingelse.turbobytes-cdn.com."}, "TurboBytes"}
	cases["something.clients.turbobytes.net."] = tstcase{[]string{}, "TurboBytes"} //Direct cdn hostname
	// Start mock server
	port, err := getfreeport()
	if err != nil {
		t.Fatal(err)
	}
	mock := fmt.Sprintf("127.0.0.1:%d", port)
	server := &dns.Server{Addr: mock, Net: "udp"}
	wait := make(chan struct{})
	go func() {
		//t.Log("serving")
		close(wait) // Signal start of goroutine
		//Fail if any errors creating mock server
		err := server.ListenAndServe()
		t.Fatal(err)
	}()
	// Wait for server goroutine start
	<-wait
	//Setup handlers
	//Always responds 1.1.1.1 and only to qtype A
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = false
		if len(r.Question) > 0 {
			if r.Question[0].Qtype == dns.TypeA {
				//Only include answer for type A
				tcase, ok := cases[r.Question[0].Name]
				prev := r.Question[0].Name
				if ok {
					for _, cname := range tcase.cnames {
						cRec := &dns.CNAME{
							Hdr: dns.RR_Header{
								Name:   prev,
								Rrtype: dns.TypeCNAME,
								Class:  dns.ClassINET,
								Ttl:    10,
							},
							Target: cname,
						}
						prev = cname
						m.Answer = append(m.Answer, cRec)
					}

				}
				aRec := &dns.A{
					Hdr: dns.RR_Header{
						Name:   prev,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    10,
					},
					A: net.ParseIP("1.1.1.1").To4(),
				}
				m.Answer = append(m.Answer, aRec)
			}
		}
		//t.Log(m)
		w.WriteMsg(m)
	})
	for k, v := range cases {
		cdn, chain, err := HostnametoCDN(k, mock)
		if cdn != v.cdn {
			t.Errorf("Expected %s got: %s. Q = %s, chain = %v", v.cdn, cdn, k, chain)
		}
		if err != nil {
			t.Error(err)
		}
	}
	//Real world test cases
	cases_real := map[string]string{
		"cloudfront.cdnplanet.com":                 "Amazon Cloudfront",
		"hwnd.cdnplanet.com":                       "Highwinds",
		"ec.cdnplanet.com":                         "EdgeCast",
		"cdnw.cdnplanet.com":                       "CDNetworks",
		"cachefly.cdnplanet.com":                   "Cachefly",
		"lvl3.cdnplanet.com":                       "Level3",
		"netdna.cdnplanet.com":                     "MaxCDN",
		"lw.cdnplanet.com":                         "LeaseWeb CDN",
		"bg.cdnplanet.com":                         "Tata communications",
		"cdn77.cdnplanet.com":                      "CDN77",
		"llnw.cdnplanet.com":                       "Limelight",
		"azion.turbobytes.net":                     "Azion",
		"medianova.cdnplanet.com":                  "Medianova",
		"cnc.cdnplanet.com.picnorecord.qtlglb.com": "QUANTIL/ChinaNetCenter",
		"cloudflare.cdnplanet.com":                 "Cloudflare",
		"ss.cdnplanet.com":                         "SwiftServe",
		"dogfood.turbobytes.com":                   "TurboBytes",
		"ecgs.cdnplanet.com":                       "EdgeCast",
		"ecgs-large.cdnplanet.com":                 "EdgeCast",
		"fastly-global.cdnplanet.com":              "Fastly",
		"hwndssl.turbobytes.net":                   "Highwinds",
		"ecssl.turbobytes.net":                     "EdgeCast",
		"akamai.turbobytes.net":                    "Akamai",
		"ss-sec.cdnplanet.com":                     "SwiftServe",
		"rum.trbbts14.net":                         "EdgeCast",
		"ecl.cdnplanet.com":                        "EdgeCast",
	}
	for k, v := range cases_real {
		cdn, chain, err := HostnametoCDN(k, "8.8.8.8:53")
		if cdn != v {
			t.Errorf("Expected %s got: %s. Q = %s, chain = %v", v, cdn, k, chain)
		}
		if err != nil {
			t.Error(err)
		}
	}
}
