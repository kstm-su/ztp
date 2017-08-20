FROM node:8.4-alpine

ADD . .
RUN yarn install

ENV PORT=8080
EXPOSE 8080

CMD ["./entrypoint.sh"]
