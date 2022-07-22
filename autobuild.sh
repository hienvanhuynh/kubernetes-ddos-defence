#Make sure a registry is specified
test -n "$REGISTRY" || eval 'echo "Need to specify a REGISTRY enviroment variable where to push image to" ; exit'

pushd src/scraper/deployment
docker-compose build
docker push $REGISTRY/kdd-scraper
popd

pushd src/ddos-detector/deployment
docker-compose build
docker push $REGISTRY/kdd-ddos-detector
popd

pushd src/executor/deployment
docker-compose build
docker push $REGISTRY/kdd-executor
popd

# Deprecated
#pushd src/kddui/deployment
#docker-compose build
#docker push $REGISTRY/kdd-ui
#popd

[[ "$SAMPLE_IMAGE" == "yes" ]] || exit

pushd sample-src/nodeserver
docker-compose build
docker push $REGISTRY/nodeserver
popd

pushd sample-src/normaluser/deployment
docker-compose build
docker push $REGISTRY/normal-user
popd

pushd sample-src/python-attacker/deployment
docker-compose build
docker push $REGISTRY/python-attacker
popd

pushd sample-src/attacker/deployment
docker-compose build
docker push $REGISTRY/attacker
popd