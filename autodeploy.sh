kubectl apply -f nodeserver/nodeserver.yaml
kubectl apply -f normal-user/deployment/normaluser.yaml
kubectl -n kube-system apply -f scraper/deployment/scraper.yaml
kubectl -n kube-system apply -f ddos-detection/deployment/ddosdetection.yaml
kubectl -n kube-system apply -f executor/deployment/executor.yaml
#kubectl apply -f kddui/deployment/kddui.yaml

