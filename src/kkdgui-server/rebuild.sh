docker build -t nguyenduongag/kkd-gui-server .
docker push nguyenduongag/kkd-gui-server
kubectl rollout restart deployment/kkd-gui-server-deployment
kubectl get pod
