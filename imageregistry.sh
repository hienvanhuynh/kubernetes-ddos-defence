REGISTRY='192.168.56.1:5000'
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
cd ../../scraper/deployment
pwd
docker-compose build
docker push $REGISTRY/scraper
cd ../../ddos-detection/deployment
pwd
docker-compose build
docker push $REGISTRY/ddos-detection
cd ../../executor/deployment
pwd
docker-compose build
docker push $REGISTRY/executor
cd ../../kddui/deployment
pwd
docker-compose build
docker push $REGISTRY/kddui
cd ../../python-attacker/deployment
pwd
docker-compose build
docker push $REGISTRY/python-attacker
cd ../..