#!/bin/bash

# Build the project
make build

#upload the project to docker.io

var1=$(docker load -i ./_output/xpkg/linux_amd64/*)

#getstring from Loaded image ID: to end
var1=${var1##*Loaded image ID: }

echo $var1

docker tag $var1 datavix/flprocess-provider
docker push datavix/flprocess-provider