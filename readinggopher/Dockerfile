FROM golang AS build-env
WORKDIR /go/src/readinggopher
COPY . .
RUN make

FROM alpine
WORKDIR /readinggopher
COPY --from=build-env /go/src/readinggopher/bot .
RUN apk --update add openssl ca-certificates
ENTRYPOINT ["./bot"]
