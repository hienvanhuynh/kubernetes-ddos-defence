cd normaluser/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/normal-user
cd ../../nodeserver
pwd
docker-compose build
docker push 192.168.56.1:5000/nodeserver
cd ../attacker/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/attacker
cd ../../scraper/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/scraper
cd ../../ddos-detection/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/ddos-detection
cd ../../executor/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/executor
cd ../../kddui/deployment
pwd
docker-compose build
docker push 192.168.56.1:5000/kddui
cd ../..