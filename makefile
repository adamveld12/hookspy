.PHONY: api_dev bindata clean cli client db

clean:
	rm -rf ./main ./main.exe ./1 ./client/build ./bindata_assetfs.go
	docker rm -f hookspy-db

db:
	docker run -d -p 8080:8080 -p 28015:28015 -p 29015:29105 --name hookspy-db rethinkdb

setup: client/node_modules
	curl https://glide.sh/get | sh
	glide install

client:
	cd ./client && npm run start

api_dev:
	HOOKSPY_ADDR=:3001 HOOKSPY_DEBUG=true HOOKSPY_DB=localhost:28015 go run *.go

cli: bindata api
	docker build -t adamveld12/hookspy .

client/node_modules:
	cd ./client && npm install

client/build: client/node_modules
	cd ./client && npm run build

bindata: client/build
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go install github.com/jteeuwen/go-bindata/...
	go install github.com/elazarl/go-bindata-assetfs/...
	go-bindata-assetfs client/build/...

api:
	GOOS=linux CGO_ENABLED=0 go build -o api .
