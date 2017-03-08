package cdnfinder

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"
)

type RawResource struct {
	Count    int          `json: "count"`
	Bytes    int          `json: "bytes"`
	IsBase   bool         `json: "isbase"`
	Hostname string       `json: "hostname"`
	Headers  *http.Header `json: "headers"`
}

type RawDiscovery struct {
	BasePageHost string                 `json: "basepagehost"`
	Resources    map[string]RawResource `json: "resources"`
}

func discoverResources(url string, timeout time.Duration) (*RawDiscovery, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, phantomjsbin, resourcefinderjs, url)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	//log.Println(string(out))
	res := &RawDiscovery{}
	err = json.Unmarshal(out, res)
	return res, err
}
