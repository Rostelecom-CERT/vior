[![Build Status](https://travis-ci.org/Rostelecom-CERT/vior.svg?branch=master)](https://travis-ci.org/Rostelecom-CERT/vior) [![](https://godoc.org/github.com/Rostelecom-CERT/vior?status.svg)](http://godoc.org/github.com/Rostelecom-CERT/vior)

Violations Receiver
-------------------

Content Security Policy **vio**lations **r**eceiver.

Currently it uses Elasticsearch as a storage, but other dbs could be easily implemented.

# How to start

## Docker

```
sudo docker-compose up -d
```

Do not forget to specify volume for the Elasticsearch data if you want to persist the data.

## Development version

Only Go 1.9+ is supported.

Listen on `:8080` and use `127.0.0.1:9200` as the Elastic server storage.

```
$: VIOR_PORT=8080 \
   VIOR_ELASTIC_URL=http://127.0.0.1:9200 \
   go run cmd/vior-http/main.go
```
