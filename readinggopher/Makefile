build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags "-X main.buildSHA=$(cat VERSION)" -o bot ./...

docker: version
	docker build -t gcr.io/enocom-dev/readinggopher:$(shell cat VERSION) .

version:
	git rev-parse HEAD > VERSION
