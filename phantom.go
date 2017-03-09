package cdnfinder

import (
	"context"
	"encoding/json"
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
}

// discoverResources loads the url in phantomjs and fetches the resources on the page
func discoverResources(url string, timeout time.Duration) (*rawDiscovery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, phantomjsbin, resourcefinderjs, url)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	//log.Println(string(out))
	res := &rawDiscovery{}
	err = json.Unmarshal(out, res)
	return res, err
}
