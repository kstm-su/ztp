FROM alpine:3.6
RUN apk add --no-cache rsyslog tftp-hpa syslinux
ADD entrypoint.sh /entrypoint.sh
ENV INTFTPD_PATH=/var/tftpboot
EXPOSE 69/udp
ENTRYPOINT ["/entrypoint.sh"]
