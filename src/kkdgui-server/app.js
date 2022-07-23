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

const command = "kubectl get ccnp -o json";

const getDetectorsCommand = "kubectl get deployment --selector=module=detector --output=jsonpath={.items..metadata.name}";

const command2 = 'curl https://kubernetes.default.svc/apis/cilium.io/v2/ciliumclusterwidenetworkpolicies \
--cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
--header "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"';

const k8sNetworkApi = kc.makeApiClient(k8s.NetworkingV1Api);


const addDetectorCommand = 'cat detetorDeployment.yaml | sed "s/{{DETECTOR_NAME}}/$DETECTOR_NAME/g" | sed "s/{{IMAGE_NAME}}/$IMAGE_NAME/g" | kubectl apply -f -';

app.get('/apinode/addDetector', async (req, res) => {

  const name = req.query.name;
  const imageName = req.query.imageName;
  
  console.log('imageName', imageName);
  console.log('name', name);

  await exec(`export DETECTOR_NAME=name1`);
  await exec(`export IMAGE_NAME=image2`);
  
  const { stdout, stderr } = await exec('echo $IMAGE_NAME');
  console.log('echo', stdout);
  
  await exec(`cat detetorDeployment.yaml | sed "s/{{DETECTOR_NAME}}/${name}/g" | sed "s/{{IMAGE_NAME}}/${imageName}/g" | kubectl apply -f -`);
  

  res.status(200).send();
});

app.get('/apinode/deletePolicy', async (req, res) => {

  const name = req.query.name;
  console.log('name', name);

  const { stdout, stderr } = await exec(`kubectl delete ccnp ${name}`);
  console.log('echo', stdout);

  res.status(200).send();
});



app.get('/apinode/getAllPolicies', async (req, res) => {
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

app.get('/apinode/getDetectors', async (req, res) => {
  const { stdout, stderr } = await exec(getDetectorsCommand);
  console.log('getDetectors', stdout);

  // let object = {stdout};
  // console.log('object', object);
  // const fixOutput = fixJson(object.stdout);

  
  const list = stdout.split(" ");

  res.json({nameList: list});
})

app.get('/', (req, res) => {
  console.log('run');
  res.send('Hello World!')
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
