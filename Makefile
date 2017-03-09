
docker:
	mkdir -p bin #Create bin directory, ignored in .gitignore
	go generate #Regenerate assets
	go run cmd/cdnfindercli/main.go --phantomjsbin="bin/phantomjs"  --host www.cdnplanet.com #Preload phantomjs binary locally
	go build -o bin/cdnfinderserver cmd/cdnfinderserver/main.go
	go build -o bin/cdnfindercli cmd/cdnfindercli/main.go
	docker build -t turbobytes/cdnfinder .

test:
	go generate
	go test