const cors = require('cors');
const express = require('express')
const app = express()
const port = 3005
app.use(cors({origin: "*",}));

const util = require('util');
const exec = util.promisify(require('child_process').exec);

const k8s = require('@kubernetes/client-node');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();

const k8sApi = kc.makeApiClient(k8s.CoreV1Api);

const command = "kubectl get ccnp -o json"

const command2 = 'curl https://kubernetes.default.svc/apis/cilium.io/v2/ciliumclusterwidenetworkpolicies \
--cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
--header "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"';

const k8sNetworkApi = kc.makeApiClient(k8s.NetworkingV1Api);


// const fixJson = (string) => {
//   let newString = string.replace('\n','').replace("\\",'');
//   return JSON.parse(newString);
// }

// k8sNetworkApi.listNetworkPolicyForAllNamespaces().then((res) => {
//     console.log(res.body);
// });

// k8sApi.listNamespacedPod('kube-system').then((res) => {
//     console.log(res.body);
// });

app.get('/getAllPolicies', async (req, res) => {
  const { stdout, stderr } = await exec(command);
  console.log('stfgegdout', stdout);

  // let object = {stdout};
  // console.log('object', object);
  // const fixOutput = fixJson(object.stdout);
  console.log('typeeee', typeof(stdout))

  const fixOutput = JSON.parse(stdout);
  console.log(fixOutput);

  res.json(fixOutput);
})

app.get('/', (req, res) => {
  console.log('run');
  res.send('Hello World!')
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
