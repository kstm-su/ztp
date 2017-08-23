FROM node:8.4-alpine

ADD . .
RUN npm install && (cd webpack && npm install) && rm -rf webpack

ENV PORT=8080
EXPOSE 8080

CMD ["./entrypoint.sh"]
