docker build -t nguyenduongag/kkd-gui .
docker push nguyenduongag/kkd-gui
kubectl rollout restart deployment/kkd-gui-deployment
kubectl get pod