FROM alpine:3.6
RUN apk add --no-cache rsyslog tftp-hpa syslinux
ADD entrypoint.sh /entrypoint.sh
ENV TFTP_ROOT=/var/tftpboot
EXPOSE 69/udp
ENTRYPOINT ["/entrypoint.sh"]
