FROM golang:1.8-alpine
RUN apk add --update --no-cache git musl-dev gcc ca-certificates
WORKDIR $GOPATH/src/github.com/kstm-su/ztp/builder
ADD . .
RUN go get -v ./... && go build --ldflags '-s -w -linkmode external -extldflags -static' -o /build

FROM docker:latest
COPY --from=0 /build /build
COPY ./pxelinux.cfg /usr/share/syslinux/pxelinux.cfg
CMD ["/build"]
