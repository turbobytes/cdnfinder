[![Go Report Card](https://goreportcard.com/badge/github.com/turbobytes/cdnfinder)](https://goreportcard.com/report/github.com/turbobytes/cdnfinder)
[![GoDoc](https://godoc.org/github.com/turbobytes/cdnfinder?status.svg)](https://godoc.org/github.com/turbobytes/cdnfinder)
# cdnfinder

Previously known as [cdnfinder.js](https://github.com/sajal/cdnfinder.js).

Webapp and cli-tool to detect CDN usage of websites. This is the backend for CDNPlanet's [CDN Finder tool.](http://www.cdnplanet.com/tools/cdnfinder/)

- Test single hostname or full webpage
- Automatically downloads the compatible phantomjs executable

## Install


### TODO: Test on darwin, linux/386, windows


### TODO: Binary releases


### Docker image

[turbobytes/cdnfinder](https://hub.docker.com/r/turbobytes/cdnfinder/)

cli: `docker run -it turbobytes/cdnfinder cdnfindercli --phantomjsbin="/bin/phantomjs"  --full http://www.cdnplanet.com/`

server: `docker run -it turbobytes/cdnfinder cdnfinderserver --phantomjsbin="/bin/phantomjs"`

The `--phantomjsbin="/bin/phantomjs"` portion is important to avoid re-downloading phantomjs each time you launch a container. I will get rid of it in the future using environment variables.

### TODO: Install from source

## Usage

TODO

### cdnfindercli

````
Usage of cdnfindercli:
  -full string
    	URL for full finder
  -host string
    	hostname for single hostname finder
  -phantomjsbin string
    	path to phantomjs, if blank tmp dir is used
  -server string
    	dns server for resolution (default "8.8.8.8:53")
````

Either `-full` or `-host` must be provided

### cdnfinderserver

````
Usage of cdnfinderserver:
  -phantomjsbin string
    	path to phantomjs, if blank tmp dir is used
  -server string
    	dns server for resolution (default "8.8.8.8:53")

````

The server listens on port 1337 on all interfaces.

#### Server API

Single host: `curl -X POST -d '{"hostname": "st.cdnplanet.com"}' -H "Content-Type: application/json" http://127.0.0.1:1337/hostname/`

Full site: `curl -X POST -d '{"url": "http://www.turbobytes.com"}' -H "Content-Type: application/json" http://127.0.0.1:1337/`

## CDN mappings

CNAME mappings are located in [assets/cnamechain.json](assets/cnamechain.json). It is a list of pair of strings where first item is part of the hostname to be matched and the second is the name of the CDN.

To update the list..

1. Fork this repo
2. Make your changes to `assets/cnamechain.json`
3. Run `make test`

If all passes, send a pull request. If the nature of the change requires changes in the tests then please do so. Bonus points for expanding on the tests

CDN header detection logic is currently located in [headerguess.go](headerguess.go). If you have some ideas on how to express as json, I would like to hear about it in issues.
