#!/bin/bash
if [ -z "$1" ]
then
    TAG=$1
else
    TAG="testing"
fi

VERSION=`curl -fsSLI -o /dev/null -w %{url_effective} https://github.com/mattkasun/timetrace-gui/releases/latest | awk -F/ '{print $8}'`
echo $VERSION
VER=`curl -fsSLI -o /dev/null -w %{url_effective} https://github.com/dominikbraun/timetrace/releases/latest | awk -F/ '{print $8}'`
echo $VER
echo "docker build -t combined:v0.1 -f Dockerfile.combined . --build-arg VERSION=${VERSION} --build-arg VER=${VER}"
docker build -t combined:${TAG} -f Dockerfile.combined . --build-arg VERSION=${VERSION} --build-arg VER=${VER}
