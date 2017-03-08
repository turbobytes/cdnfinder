package cdnfinder

import (
	"sort"
	"sync"
	"time"
)

//Need this to mantain api compatibility
type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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

type FullOutput struct {
	BaseCDN   *string        `json:"basecdn"`
	AssetCDN  *string        `json:"assetcdn"`
	Resources []FullResource `json:"everything"`
}

// to make FullOutput.Resources sortable
type FullSort []FullResource

func (a FullSort) Len() int           { return len(a) }
func (a FullSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FullSort) Less(i, j int) bool { return a[i].Count < a[j].Count }

type hostcdn struct {
	cdn    string
	cnames []string
}

func parseraw(raw *RawDiscovery, server string) *FullOutput {
	out := &FullOutput{
		Resources: make([]FullResource, 0),
	}
	hostmap := make(map[string]hostcdn)
	var wg sync.WaitGroup
	mut := &sync.Mutex{}
	for k, _ := range raw.Resources {
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

		for key, _ := range *v.Headers {
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
	sort.Sort(sort.Reverse(FullSort(out.Resources)))
	//Most popular hostname by count desides AssetCDN
	if len(out.Resources) > 0 {
		out.AssetCDN = out.Resources[0].CDN
	}
	return out
}

func FullFinder(url, server string, timeout time.Duration) (*FullOutput, error) {
	raw, err := discoverResources(url, timeout)
	if err != nil {
		return nil, err
	}
	return parseraw(raw, server), nil
}
