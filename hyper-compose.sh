version: '2'
services:
  web:
    size: S1
    environment:
      - "HOOKSPY_DEBUG=false"
    image: adamveld12/hookspy
    links:
      - db:db
    depends_on:
      - db
  db:
    size: S2
    image: rethinkdb
    volumes:
      - dbstorage:/data
  lb:
    size: S2
    image: 'dockercloud/haproxy:latest'
    fip: 209.177.93.123
    links:
      - web
    depends_on:
      - web
    ports:
      - '80:80'
      - '443:443'
