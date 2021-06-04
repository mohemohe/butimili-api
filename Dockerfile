FROM golang:alpine as server
ADD . /go/src/github.com/mohemohe/butimili-api/
WORKDIR /go/src/github.com/mohemohe/butimili-api/
RUN \
    set -xe; \
    apk --no-cache add alpine-sdk; \
    go build -ldflags "\
      -X go/src/github.com/mohemohe/butimili-api/util.version=$(date '+%y.%m.%U%u%H%M') \
      -X go/src/github.com/mohemohe/butimili-api/util.hash=$(git rev-parse HEAD) \
    " -o /app

FROM alpine
RUN \
    set -xe; \
    apk --no-cache add ca-certificates
COPY --from=server /app /app

EXPOSE 8100
CMD ["/app"]