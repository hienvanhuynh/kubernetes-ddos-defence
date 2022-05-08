kubectl apply -f nodeserver/nodeserver.yaml
kubectl apply -f normal-user/deployment/normaluser.yaml
kubectl -n kube-system apply -f controller/deployment/controller.yaml
kubectl -n kube-system apply -f ddos-detection/deployment/ddosdetection.yaml
