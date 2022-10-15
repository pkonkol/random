#!/usr/bin/bash
eval $(minikube docker-env --shell='bash')
docker build master/ -t 'scraper-master'
docker build slave/ -t 'scraper-slave'
# scale down and back up to reload locally built images
minikube kubectl -- scale --replicas=0 deployment scraper-master
minikube kubectl -- scale --replicas=0 deployment scraper-slave
minikube kubectl -- apply -f k8s/
