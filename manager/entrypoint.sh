#!/bin/sh
if ! [ -d "node_modules/$NODE_MODULE" ]; then
	yarn add $NODE_MODULE
fi
node index.js
