FROM golang:alpine AS build-env

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine
WORKDIR /app
COPY --from=build-env /go/bin/app .
ENTRYPOINT ./app
