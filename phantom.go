package cdnfinder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

// rawResource holds each resource accessed by the webpage
type rawResource struct {
	Count    int          `json:"count"`
	Bytes    int          `json:"bytes"`
	IsBase   bool         `json:"isbase"`
	Hostname string       `json:"hostname"`
	Headers  *http.Header `json:"headers"`
}

// rawDiscovery Parses the stdout from phantomjs process
type rawDiscovery struct {
	BasePageHost string                 `json:"basepagehost"`
	Resources    map[string]rawResource `json:"resources"`
	Error        string                 `json:"error"`
}

// discoverResources loads the url in phantomjs and fetches the resources on the page
func discoverResources(url string, timeout time.Duration, verbose bool) (*rawDiscovery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, phantomjsbin, resourcefinderjs, url)
	if verbose {
		log.Println(cmd)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Println(string(out))
	}
	res := &rawDiscovery{}
	err = json.Unmarshal(out, res)
	if verbose {
		log.Println(res)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("Phantomjs had an error")
	}
	return res, err
}
