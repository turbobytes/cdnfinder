package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nytimes/gziphandler"
	"github.com/turbobytes/cdnfinder"
)

var server = flag.String("server", "8.8.8.8:53", "dns server for resolution")
var failure = []byte(`{"status":"FAILURE"}`)

func init() {
	flag.Parse()
	cdnfinder.Init()
}

type SingleHostReq struct {
	Hostname string `json:"hostname"`
}

type FullPageReq struct {
	Url string `json:"url"`
}

func handleFullPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(failure)
		return
	}
	defer r.Body.Close()
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failure)
		return
	}
	log.Println(string(d))
	req := &FullPageReq{}
	err = json.Unmarshal(d, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(failure)
		return
	}
	out, err := cdnfinder.FullFinder(req.Url, *server, time.Minute)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failure)
		return
	}
	//No resource detected = some booboo
	if len(out.Resources) == 0 {
		w.Write(failure)
		return
	}
	d, err = json.MarshalIndent(out, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failure)
		return
	}
	w.Write(d)
}

func handleSingleHost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(string(d))
	req := &SingleHostReq{}
	err = json.Unmarshal(d, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c, _, err := cdnfinder.HostnametoCDN(req.Hostname, *server)
	if err != nil {
		log.Println(err) //Informational..
	}
	w.Write([]byte(c))
}

func main() {
	http.HandleFunc("/hostname/", handleSingleHost)
	http.HandleFunc("/", handleFullPage)

	s := &http.Server{
		Addr:           ":1337",
		Handler:        gziphandler.GzipHandler(http.DefaultServeMux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on :1337")
	log.Fatal(s.ListenAndServe())
}
