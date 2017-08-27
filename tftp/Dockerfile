FROM alpine:3.6
RUN apk add --no-cache tftp-hpa
ENV TFTP_ROOT=/var/tftpboot
EXPOSE 69/udp
ENTRYPOINT ["in.tftpd"]
CMD ["-L", "-v", "-s", "/"]
