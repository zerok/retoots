FROM golang:1.23.3-alpine as builder
RUN mkdir /src && apk add --no-cache make
COPY . /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod \
    make

FROM alpine:3.20
COPY --from=builder /src/bin/retoots /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/retoots"]
