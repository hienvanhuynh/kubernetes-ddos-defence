kubectl apply -f nodeserver/nodeserver.yaml
kubectl apply -f normaluser/deployment/normaluser.yaml
kubectl -n kube-system apply -f scraper/deployment/scraper.yaml
kubectl -n kube-system apply -f ddos-detection/deployment/ddosdetection.yaml
kubectl -n kube-system apply -f executor/deployment/executor.yaml
kubectl -n kube-system apply -f kddui/deployment/kddui.yaml
kubectl apply -f python-attacker/deployment/python-attacker.yaml

