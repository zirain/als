#!/bin/sh

# build
GOOS=linux GOARCH=amd64 go build -o envoy-als-server ../cmd/main.go
chmod 777 envoy-als-server
sha256sum envoy-als-server

# clean pod
kubectl delete -f k8s.yaml -nistio-system

# build image
docker rmi envoy-als-server:latest
docker build -t envoy-als-server:latest .
rm -rf envoy-als-server

# run pod
kubectl apply -f k8s.yaml -nistio-system