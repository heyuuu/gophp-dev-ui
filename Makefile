fe-build:
	cd frontend && npm run build

api-build:
	go build -o tmp/main

build: fe-build api-build