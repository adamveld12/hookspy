.PHONY: api_dev bindata clean cli client db

clean:
	rm -rf ./main ./main.exe ./1 ./client/build
	docker rm -f hookspy-db

db:
	docker run -d -p 8080:8080 -p 28015:28015 -p 29015:29105 --name hookspy-db rethinkdb

client:
	cd ./client && npm run start

api_dev:
	HOOKSPY_ADDR=:3001 HOOKSPY_DEBUG=true HOOKSPY_DB=localhost:28015 go run *.go

cli: bindata api
	docker build -t adamveld12/hookspy .

client/build:
	npm run build

bindata: client/build
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go install github.com/elazarl/go-bindata-assetfs/...
	go-bindata-assetfs client/build/...

api:
	GOOS=linux CGO_ENABLED=0 go build -o api .
