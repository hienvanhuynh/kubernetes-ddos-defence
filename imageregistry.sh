cd normaluser/deployment
pwd
docker-compose build
docker push $REGISTRY/normal-user
cd ../../nodeserver
pwd
docker-compose build
docker push $REGISTRY/nodeserver
cd ../attacker/deployment
pwd
docker-compose build
docker push $REGISTRY/attacker
cd ../../controller/deployment
pwd
docker-compose build
docker push $REGISTRY/controller
cd ../../ddos-detection/deployment
pwd
docker-compose build
docker push $REGISTRY/ddos-detection
cd ../