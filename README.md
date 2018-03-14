Violations Receiver
-------------------

Content Security Policy **vio**lations **r**eceiver

Currently it uses Elasticsearch as a storage, but other dbs could be easily implemented

# How to start

Listen on `:8080` and use `127.0.0.1:9200` as the Elastic server storage

```
$: VIOR_PORT=8080 \
   VIOR_ELASTIC_URL=http://127.0.0.1:9200 \
   go run cmd/vior-http/main.go
```