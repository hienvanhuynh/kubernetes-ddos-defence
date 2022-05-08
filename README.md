# kubernetes-ddos-defence
## Build images and push to registry<br />
### Change environment variable `$REGISTRY` in `.env` file<br />
### Then run<br />
>./imageregistry.sh<br />

## Deploy all needed component<br />
### When in one node of the cluster, run<br />
>./autodeploy.sh<br />


## Wait for some minutes and deploy the attacker<br />
>./attackerdeploy.sh<br />


## Deploy all needed component<br />
### When in one node of the cluster, run<br />
>./autodeploy.sh<br />

