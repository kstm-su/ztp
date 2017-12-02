FROM node:8.4-alpine

ADD . /manager
WORKDIR /manager
RUN yarn install && cd webpack && yarn install && yarn run build && cd .. && rm -rf webpack

ENV PORT=8080
EXPOSE 8080

CMD ["/manager/entrypoint.sh"]
