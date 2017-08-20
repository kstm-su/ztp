#!/bin/sh
npm install $NODE_MODULE
(cd webpack && npm run build)
node index.js
