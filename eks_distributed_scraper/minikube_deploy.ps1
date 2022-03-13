gsudo minikube -p minikube docker-env --shell powershell | Invoke-Expression
docker build master/ -t 'scraper-master'
docker build slave/ -t 'scraper-slave'
minikube kubectl -- scale --replicas=0 deployment scraper-master
minikube kubectl -- scale --replicas=0 deployment scraper-slave
minikube kubectl -- apply -f k8s/