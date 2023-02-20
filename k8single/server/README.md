# K8Single Server

## Local environment

### Generate swagger document

```shell
# $PRJDIR/k8single/server
$ swag init -d cmd,api --exclude api/docs -o api/docs
```

### Run server

```shell
# $PRJDIR/k8single/server
$ go run .
```

http://localhost:8080