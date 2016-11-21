FROM alpine

ENV HOOKSPY_ADDR=:80
ENV HOOKSPY_DEBUG=true
ENV HOOKSPY_DB=db:28015

WORKDIR /app/
COPY api /app

EXPOSE 80

ENTRYPOINT /app/api
