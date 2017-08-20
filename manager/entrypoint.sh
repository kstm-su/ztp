#!/bin/sh
yarn add $NODE_MODULE
(cd webpack && npm run build)
node index.js
