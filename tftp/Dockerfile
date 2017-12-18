FROM alpine:3.6
ADD https://github.com/kstm-su/ztp/releases/download/standby/standby.tar.gz /standby.tar.gz
RUN mkdir /default && tar xzvf /standby.tar.gz -C /default

FROM alpine:3.6
COPY --from=0 /default /default
RUN apk add --no-cache tftp-hpa syslinux \
	&& ln -s /usr/share/syslinux/pxelinux.0 /default/syslinux/pxelinux.0 \
	&& ln -s /usr/share/syslinux/ldlinux.c32 /default/syslinux/ldlinux.c32
ENV TFTP_ROOT=/var/tftpboot
EXPOSE 69/udp
ENTRYPOINT ["in.tftpd"]
CMD ["-L", "-v", "-s", "/"]
