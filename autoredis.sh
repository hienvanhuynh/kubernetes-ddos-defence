kubectl -n kube-system apply -f redis/deployment/redis-config.yaml
kubectl -n kube-system apply -f redis/deployment/redis.yaml