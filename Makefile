dev:
	cd .. && air -c gophp-dev-ui/all.air.toml

fe-build:
	cd frontend && npm run build

api-build:
	go build -o tmp/main

build: fe-build api-build