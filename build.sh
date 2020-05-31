#!/bin/bash

# Go related variables.
GOBASE=$(pwd)
GOPATH=$(go env GOPATH)
GOHOSTOS=$(go env GOOS)
GOHOSTARCH=$(go env GOARCH)

BIN_DIR=${GOHOSTOS}_${GOHOSTARCH}

mkdir -p $BIN_DIR
go get -u github.com/mariiatuzovska/user-service
go install github.com/mariiatuzovska/user-service
cp -a -v ${GOPATH}/bin/user-service ${BIN_DIR}/user-service
echo '{
    "DBContext": {
        "Shema": "postgres",
        "User": "postgres",
        "Password": "postgres",
        "Host": "127.0.0.1",
        "Port": "5432"
    },
    "APIContext": {
        "Host": "127.0.0.1",
        "Port": "8080"
    }
}' > ${BIN_DIR}/user-configuration.json
echo 'FROM golang:1.13
COPY user-configuration.json user-configuration.json
COPY user-service app
EXPOSE 8080
CMD ["./app", "start"]' > ${BIN_DIR}/Dockerfile
