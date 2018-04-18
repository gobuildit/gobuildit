build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags "-X main.buildSHA=$(shell cat VERSION)" -o bot ./...

docker: version
	docker build -t gcr.io/enocom-dev/readinggopher:$(shell cat VERSION) .

.PHONY: version
version:
	git rev-parse HEAD > VERSION
