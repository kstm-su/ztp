FROM golang:1.8-alpine
RUN apk add --update --no-cache git musl-dev gcc ca-certificates
ADD . $GOPATH/src/github.com/kstm-su/ool-specialist-2017/moby
WORKDIR $GOPATH/src/github.com/kstm-su/ool-specialist-2017/moby
RUN go get -v ./... && go build --ldflags '-s -w -linkmode external -extldflags -static' -o /moby

FROM alpine:3.6
RUN apk add --update --no-cache ca-certificates
COPY --from=0 /moby /moby
CMD ["/moby"]
