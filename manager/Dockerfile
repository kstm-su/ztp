FROM node:8.4-alpine

ADD ./package.json ./package.json
ADD ./webpack/package.json ./webpack/package.json
RUN yarn install && cd webpack && yarn install
ADD . .
RUN cd webpack && npm run build && cd .. && rm -rf webpack

ENV PORT=8080
EXPOSE 8080

CMD ["./entrypoint.sh"]
