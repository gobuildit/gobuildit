FROM golang:alpine AS build-env

WORKDIR /go/src/readinggopher
COPY . .

RUN go get -d -v ./...
RUN go build -v -ldflags "-X main.buildSHA=$(git rev-parse HEAD)" -o app ./...

FROM alpine
WORKDIR /readinggopher
COPY --from=build-env /go/readinggopher/app .
RUN apk --update add openssl ca-certificates
ENTRYPOINT ["./app"]
