FROM alpine:latest

ADD ./_output/release /usr/local/bin/

ARG ENV
ENV ENV=$ENV

ARG DB
ENV DB=$DB

ARG API_GATEWAY
ENV API_GATEWAY=${API_GATEWAY}

ENV PORT=8080

ENTRYPOINT neuron-agent