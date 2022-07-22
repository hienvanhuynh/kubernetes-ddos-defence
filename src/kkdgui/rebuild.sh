docker build -t nguyenduongag/kkd-gui .
docker rmi -f 1271e75dcfa1
docker push nguyenduongag/kkd-gui
kubectl rollout restart deployment/kkd-gui-deployment
