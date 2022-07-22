#!/bin/bash

if [ $# != 1 ]; then
    echo "e.g.: bash $0 v1.0"
    exit 1
fi

docker_tag=$1

docker_repository="polarismesh"

echo "docker repository : ${docker_repository}/polaris-server, tag : ${docker_tag}"

bash build.sh ${docker_tag}

if [ $? != 0 ]; then
    echo "build polaris-server failed"
    exit 1
fi

docker build --network=host -t ${docker_repository}/polaris-server:${docker_tag} ./

docker push ${docker_repository}/polaris-server:${docker_tag}

pre_release=`echo ${docker_tag}|egrep "(alpha|beta|rc|[T|t]est)"|wc -l`
if [ ${pre_release} == 0 ]; then
  docker tag ${docker_repository}/polaris-server:${docker_tag} ${docker_repository}/polaris-server:latest
  docker push ${docker_repository}/polaris-server:latest
fi