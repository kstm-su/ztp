#!/bin/sh

in.tftpd -L -v -s $TFTP_ROOT &
rsyslogd
echo > var/log/messages
tail -f /var/log/messages
