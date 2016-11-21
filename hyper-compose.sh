version: '2'
services:
  api:
    environment:
      HOOKSPY_DEBUG: false
    image: adamveld12/hookspy
    fip: 209.177.93.123
    links:
      - db:db
    depends_on:
      - db
    ports:
      - "80:80"
  db:
    image: rethinkdb
