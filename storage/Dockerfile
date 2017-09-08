FROM golang:1.8-alpine
RUN apk add --update --no-cache git musl-dev gcc ca-certificates syslinux
WORKDIR $GOPATH/src/github.com/kstm-su/ztp/storage
ADD . .
RUN go get -v ./... && go build --ldflags '-s -w -linkmode external -extldflags -static' -o /storage && \
	mv pxelinux.cfg /pxelinux.cfg

FROM docker:latest
COPY --from=0 /storage /storage
COPY --from=0 /usr/share/syslinux /usr/share/syslinux
COPY --from=0 /pxelinux.cfg /usr/share/syslinux/pxelinux.cfg
CMD ["/storage"]
