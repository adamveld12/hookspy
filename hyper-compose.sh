version: '2'
services:
  web:
    restart: always
    size: S1
    fip: 209.177.93.123
    environment:
      - "HOOKSPY_DEBUG=false"
    image: adamveld12/hookspy
    ports:
      - '80:80'
    links:
      - db:db
    depends_on:
      - db
  db:
    restart: always
    size: S2
    image: rethinkdb
    volumes:
      - dbstorage:/data
