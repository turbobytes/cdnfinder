package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/turbobytes/cdnfinder"
)

var server = flag.String("server", "8.8.8.8:53", "dns server for resolution")
var full = flag.String("full", "", "URL for full finder")
var hostname = flag.String("host", "", "hostname for single hostname finder")

func init() {
	flag.Parse()
	cdnfinder.Init()
}

func main() {
	if *full != "" {
		out, err := cdnfinder.FullFinder(*full, *server, time.Minute)
		if err != nil {
			log.Fatal(err)
		}
		d, err := json.MarshalIndent(out, " ", " ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(d))
	} else if *hostname != "" {
		c, _, err := cdnfinder.HostnametoCDN(*hostname, *server)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(c)
	} else {
		log.Fatal("full or host needs to be specified")
	}

	//
}
