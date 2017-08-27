#!/bin/sh
rm -rf $SOCKET
if ! [ -d "node_modules/$NODE_MODULE" ]; then
	/usr/local/bin/yarn add $NODE_MODULE
fi
/usr/local/bin/node index.js
