# dumbo
A blob with a mutex over http

# why?
Most databases perform queries on data on behalf of clients. Queries are resource hungry. By the server only providing the blob - the burden will be shifted to the client and thus save server resources. Clients will run queries on in-memory data, which is likely to be fast, saving network hops.

The server is agnostic on blob contents. It just provides it over http.

# api
- GET /read

Server returns compressed blob to client.

- POST /write

Expects a compressed blob as request body and header `Content-Encoding:gzip`. The blob is persisted to disk. As an experminent we keep this file open for reads and writes for the duration of the server lifetime.

# usage
```bash
# Start server with blob
$ go run cmd/server/main.go testdata/nice.json
# client read
$ curl -s 0.0.0.0:7979/read | gzip -d | jq
# client write
$ cat testdata/false.gz | curl -H Content-Encoding:gzip --data-binary @- 0.0.0.0:7979/write
```

# test
```bash
$ go test -v -race
```

# todos
- [x] compress data over wire
- [ ] compare hashed blob before persist to disk
- [ ] send hash as Etag to read requests
- [ ] k8s operator
- [ ] investigate graceful exits and mutex
- [ ] encryption
- [ ] timeout on read lock contention
- [ ] relax mutex on reads
- [ ] optional basic auth for writes
- [ ] maybe make filewr a package?

# license
MIT
