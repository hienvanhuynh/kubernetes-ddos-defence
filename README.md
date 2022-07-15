# KDD
 KDD is a software run on top of Cilium that helps mitigate network attacks.
<hr>

## Installation
### Prepare a registry
First you need a registry where you can push built images to and pull images from. If you already have, ignore this step.
If you want a simple way to setup a registry, take a look at [how to setup docker insecure registry](#setup-a-docker-insecure-registry)
### Build and push images to registry
Set your registry to environment variable
> $ export REGISTRY="\<YOUR-REGISTRY-URL\>"

> $ ./autobuild.sh
### Deploy Cilium
> $ helm upgrade --install cilium cilium/cilium --version 1.11.6 \\ \
    --namespace kube-system \\ \
    --set kubeProxyReplacement=strict \\ \
    --set k8sServiceHost=\<CONTROL-PLANE-IP\> \\ \
    --set k8sServicePort=6443 \\ \
    --set tunnel=disabled \\ \
    --set autoDirectNodeRoutes=true \\ \
    --set loadBalancer.mode=dsr \\ \
    --set loadBalancer.algorithm=maglev \\ \
    --set nativeRoutingCIDR="10.0.0.0/8" \\ \
    --set hubble.enabled=true \\ \
    --set hubble.relay.enabled=true \\ \
    --set hubble.ui.enabled=true

If you want to expose hubble-ui to NodePort
> $ kubectl -n kube-system patch service hubble-ui -p '{"spec":{"type":"NodePort", "ports":[{"port":80, "nodePort":30100}]}}'
### Deploy
> $ helm upgrade --install kdd ./helm/kdd --set imageRegistry="\<YOUR-REGISTRY-URL\>"
### Install Prometheus
> $ kubectl create namespace prometheus
> $ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

> $ helm upgrade -i prometheus prometheus-community/prometheus \\ \
    --namespace prometheus \\ \
    --set alertmanager.enabled=false \\ \
    --set server.persistentVolume.storageClass="local-storage" \\ \
    --set kubeStateMetrics.enabled=false \\ \
    --set nodeExporter.enabled=false \\ \
    --set pushgateway.enabled=false \\ \
    --set server.persistentVolume.size=1.5Gi

> $ kubectl --namespace=prometheus port-forward deploy/prometheus-server 9090

> $ mkdir /home/\<CLUSTER_NODE_NAME\>/log

> $ cat <<EOF | kubectl apply -f - \
apiVersion: v1 \
kind: PersistentVolume \
metadata: \
  name: prometheus-server \
spec: \
  capacity: \
    storage: 1.5Gi \
  accessModes: \
    - ReadWriteOnce \
  persistentVolumeReclaimPolicy: Retain \
  storageClassName: local-storage \
  claimRef: \
    name: prometheus-server \
    namespace: prometheus \
  local: \
    path: /home/\<CLUSTER_NODE_NAME\>/log/ \
  nodeAffinity: \
    required: \
      nodeSelectorTerms: \
       - matchExpressions: \
          - key: kubernetes.io/hostname \
            operator: In \
            values: \
             - \<CLUSTER_NODE_NAME\> \
EOF

## Supplementary
### Setup a docker insecure registry
> $ docker run -d -p 5000:5000 --restart=always --name registry registry:2

Or you can check for more details here:
https://docs.docker.com/registry/deploying/
### Build up a cluster using kubeadm
Prepare some nodes to run the cluster
Create the config file for kubelet: 
> $ sudo nano /etc/default/kubelet

Add this line to it
> KUBELET_EXTRA_ARGS=--node-ip=\<YOUR-NODE-IP\>

Disable firewall and turn off swap for kubelet to be able to run
> $ sudo ufw disable

> $ sudo swapoff -a; sudo sed -i '/swap/d' /etc/fstab

Check if port 6443 is ready for use:
> $ nc 127.0.0.1 6443

Install container runtime (containerd) by follow links here: \
https://docs.docker.com/engine/install/ubuntu/ \
https://docs.docker.com/engine/install/linux-postinstall/

Edit containerd config (/etc/containerd/config.toml), for kubelet to be able to run, by comment out this line
<pre>
disabled_plugins = ["cri"]
</pre>
> $ sudo systemctl restart containerd

Do above steps for every node

Edit containerd config (/etc/containerd/config.toml), for pulling image from insecure local registry:
<pre>
version = 2
[plugins."io.containerd.grpc.v1.cri".registry]
  config_path = "/etc/containerd/certs.d"
</pre>

> $ sudo mkdir -p /etc/containerd/certs.d/\<INSECURE-REGISTRY-IP\>:\<INSECURE-REGISTRY-PORT\>

> $ sudo nano /etc/containerd/certs.d/\<INSECURE-REGISTRY-IP\>:\<INSECURE-REGISTRY-PORT\>/hosts.toml

<pre>
server = "http://&ltINSECURE-REGISTRY-IP&gt:&ltINSECURE-REGISTRY-PORT&gt"

[host."http://&ltINSECURE-REGISTRY-IP&gt:&ltINSECURE-REGISTRY-PORT&gt"]
  capabilities = ["pull", "resolve", "push"]
  skip_verify = true
  plain-http = true
</pre>
> $ sudo systemctl restart containerd

Do above steps for every node

Install kubeadm, kubelet, and kubectl
> $ sudo apt-get update

> $ sudo apt-get install -y apt-transport-https ca-certificates curl

> $ sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

> $ echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

> $ sudo apt-get update

> $ sudo apt-get install -y kubelet kubeadm kubectl

As my experience, no need to configure cgroup driver

Init control plane node
> $ sudo kubeadm init --skip-phases=addon/kube-proxy --ignore-preflight-errors=all --apiserver-advertise-address=\<CONTROL-PLANE-NODE-IP\>

> $ mkdir -p \$HOME/.kube

> $ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config

> $ sudo chown \$(id -u)\:\$(id -g) \$HOME/.kube/config


Then join in the worker nodes, running the following on each node as root:

> $ kubeadm join \<CONTROL-PLANE-NODE-IP\>:6443 --token \<TOKEN-SHOWED-ON-INIT-COMMAND-LOG\> \\ \
	--discovery-token-ca-cert-hash sha256:\<HASH-SHOWED-ON-INIT-COMMAND-LOG\>

If control plane is in NotReady status, you can wait for some minutes. \
Then the cluster is ready.
## Team Members
* [Huynh Van Hien](https://github.com/hvhq)
* [Nguyen Hoang Thai Duong](https://github.com/somethingintheway)
