# kubernetes-ddos-defence
Build images and push to registry
`Change environment variable $REGISTRY in .env`
`Then run`
``./imageregistry.sh``
Deploy all needed component
`When in one node of the cluster, run`
``./autodeploy.sh``
Wait for some minutes and deploy the attacker
``./attackerdeploy.sh``
Now we can check for existing cnp (or it can be delayed upto 3s)
``kubectl get cnp``