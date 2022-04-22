CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o qos-server .
if [[ $? == 0 ]]; then echo "successfully build..."; else echo "build failed!!! exiting..."; exit 1; fi
image="qos-server:`date +%Y%m%d%H%M%S`"
docker build -t ${image} .
