package cdnfinder

import (
	"log"
	"sort"
	"sync"
	"time"
)

// Hdr holds the name/value pair for http headers in output
// Need this to maintain api compatibility
type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FullResource is description of each individual resource
type FullResource struct {
	Count       int      `json:"count"`
	Bytes       int      `json:"bytes"`
	IsBase      bool     `json:"isbase"`
	Hostname    string   `json:"hostname"`
	Headers     []Hdr    `json:"headers"`
	CNAMES      []string `json:"cnames"`
	CDN         *string  `json:"cdn"`
	HeaderGuess *string  `json:"headerguess"`
}

// FullOutput is the result of FullFinder
type FullOutput struct {
	BaseCDN   *string        `json:"basecdn"`
	AssetCDN  *string        `json:"assetcdn"`
	Resources []FullResource `json:"everything"`
}

// fullSort is intermediary type to make FullOutput.Resources sortable
type fullSort []FullResource

// Len satisfies sort.Sort interface
func (a fullSort) Len() int { return len(a) }

// Swap satisfies sort.Sort interface
func (a fullSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less satisfies sort.Sort interface
func (a fullSort) Less(i, j int) bool { return a[i].Count < a[j].Count }

type hostcdn struct {
	cdn    string
	cnames []string
}

// parseraw populates CDN information to raw output from phantomjs
func parseraw(raw *rawDiscovery, server string, verbose bool) *FullOutput {
	out := &FullOutput{
		Resources: make([]FullResource, 0),
	}
	hostmap := make(map[string]hostcdn)
	var wg sync.WaitGroup
	mut := &sync.Mutex{}
	for k := range raw.Resources {
		wg.Add(1)
		go func(k, server string) {
			cdn, cnames, _ := HostnametoCDN(k, server)
			mut.Lock()
			hostmap[k] = hostcdn{cdn, cnames}
			mut.Unlock()
			wg.Done()
		}(k, server)
	}
	wg.Wait()
	//Populate resources
	for k, v := range raw.Resources {
		hm := hostmap[k]
		cdn := hm.cdn
		cnames := hm.cnames
		res := FullResource{
			Count:    v.Count,
			Bytes:    v.Bytes,
			IsBase:   v.IsBase,
			Headers:  make([]Hdr, 0),
			Hostname: k,
			CNAMES:   cnames,
		}

		for key := range *v.Headers {
			res.Headers = append(res.Headers, Hdr{key, v.Headers.Get(key)})
		}
		if cdn != "" {
			res.CDN = &cdn
		}
		//Header Guess
		res.HeaderGuess = headerguess(v.Headers)
		if res.CDN == nil {
			res.CDN = res.HeaderGuess
		}
		if k == raw.BasePageHost && res.CDN != nil {
			out.BaseCDN = res.CDN
		}
		out.Resources = append(out.Resources, res)
	}
	sort.Sort(sort.Reverse(fullSort(out.Resources)))
	//Most popular hostname by count decides AssetCDN
	if len(out.Resources) > 0 {
		out.AssetCDN = out.Resources[0].CDN
	}
	return out
}

// FullFinder detects the CDN(s) used by a url by loading it in browser
func FullFinder(url, server string, timeout time.Duration, verbose bool) (*FullOutput, error) {
	raw, err := discoverResources(url, timeout, verbose)
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Println(raw)
	}
	return parseraw(raw, server, verbose), nil
}
