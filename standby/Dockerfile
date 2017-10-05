FROM golang:1.8
WORKDIR /usr/local/go/src/github.com/kstm-su/ztp/standby
ADD . ./
RUN go get ./... && go build --ldflags '-s -w -linkmode external -extldflags -static' -o /agent

FROM scratch
COPY --from=0 /agent /agent
CMD ["/agent"]
