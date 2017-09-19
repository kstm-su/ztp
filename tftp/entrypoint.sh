#!/bin/sh

in.tftpd -L -v &
rsyslogd
echo > var/log/messages
tail -f /var/log/messages
