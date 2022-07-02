#Make sure a registry is specified
test -n "$REGISTRY" || eval 'echo "Need to specify a REGISTRY enviroment variable where to push image to" ; exit'

pushd scraper/deployment
docker-compose build
docker push $REGISTRY/scraper
popd

pushd ddos-detector/deployment
docker-compose build
docker push $REGISTRY/ddos-detector
popd

pushd executor/deployment
docker-compose build
docker push $REGISTRY/executor
popd

pushd kddui/deployment
docker-compose build
docker push $REGISTRY/kddui
popd

pushd nodeserver
docker-compose build
docker push $REGISTRY/nodeserver
popd

pushd normaluser/deployment
docker-compose build
docker push $REGISTRY/normal-user
popd

pushd python-attacker/deployment
docker-compose build
docker push $REGISTRY/python-attacker
popd

pushd attacker/deployment
docker-compose build
docker push $REGISTRY/attacker
popd