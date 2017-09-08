FROM alpine:3.6
RUN apk add --no-cache tftp-hpa syslinux
ADD ./pxelinux.cfg /usr/share/syslinux/pxelinux.cfg
ENV TFTP_ROOT=/var/tftpboot
EXPOSE 69/udp
ENTRYPOINT ["in.tftpd"]
CMD ["-L", "-v", "-s", "/"]
