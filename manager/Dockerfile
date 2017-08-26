FROM node:8.4-alpine

ADD . .
RUN yarn install && \
	(cd webpack && yarn install && npm run build) && rm -rf webpack

ENV PORT=8080
EXPOSE 8080

CMD ["./entrypoint.sh"]
