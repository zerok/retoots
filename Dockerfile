FROM golang:1.24.5-alpine as builder
RUN mkdir /src && apk add --no-cache make
COPY . /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod \
    make

FROM alpine:3.22
COPY --from=builder /src/bin/retoots /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/retoots"]
