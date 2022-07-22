import * as React from 'react';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { Button, Grid, Stack, Chip } from '@mui/material';
import PolicyIcon from '@mui/icons-material/Policy';

import SyntaxHighlighter from 'react-syntax-highlighter';
import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Box, Container } from '@mui/system';
import OverflowTip from './component/overflowTip';
import { Client1_13 } from 'kubernetes-client';

//const client = new Client1_13({ version: '1.13' });

const k8s = require('@kubernetes/client-node');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();

const k8sApi = kc.makeApiClient(k8s.CoreV1Api);
// const k8sNetworkingApi = kc.makeApiClient(k8s.NetworkingV1Api);

k8sApi.listNamespacedPod('kube-system').then((res) => {
  console.log('pod', res.body);
});

// k8sNetworkingApi.listNetworkPolicyForAllNamespaces().then((res => {
//   console.log('policies', res);
// }))

const policy1 = {
  "apiVersion": "v1",
  "items": [
      {
          "apiVersion": "cilium.io/v2",
          "kind": "CiliumClusterwideNetworkPolicy",
          "metadata": {
              "creationTimestamp": "2022-07-21T08:50:23Z",
              "generation": 1,
              "name": "blacklist-rule0581",
              "resourceVersion": "1842683",
              "uid": "38cd48a2-0b86-4fd8-b44f-d6aa55b7faab"
          },
          "spec": {
              "endpointSelector": {
                  "matchLabels": {
                      "app": "web",
                      "io.cilium.k8s.namespace.labels.kubernetes.io/metadata.name": "default",
                      "io.cilium.k8s.policy.cluster": "default",
                      "io.cilium.k8s.policy.serviceaccount": "default",
                      "io.kubernetes.pod.namespace": "default"
                  }
              },
              "ingress": [
                  {
                      "fromEntities": [
                          "all"
                      ]
                  }
              ],
              "ingressDeny": [
                  {
                      "fromCIDR": [
                          "192.168.56.1/32"
                      ]
                  }
              ]
          }
      }
  ],
  "kind": "List",
  "metadata": {
      "resourceVersion": ""
  }
};

const getAllPolicies = async () => {
  // const policies = await client.apis.networking.k8s.io.v1.networkpolicies.get();
  const policies = [policy1];
  return policies;
}


const convertJSONtoList = (json) =>
{
  let list = [];

  for (let key in json) {
    if (json.hasOwnProperty(key)) {
        list.push(key + ": " + json[key])
    }
  }

  return list;
}

const getDataFromPolicy = (policy) => {
  let policyName = policy["items"][0].metadata.name;

  let creationTimestamp = policy["items"][0].metadata.creationTimestamp;

  let labelObject = policy["items"][0].spec.endpointSelector.matchLabels;

  let toLabels = convertJSONtoList(labelObject);

  let ingressDeny = policy["items"][0].spec.ingressDeny;

  let fromLabels = [];


  ingressDeny.forEach(endpoint => {

    if (endpoint.hasOwnProperty("fromCIDR"))
    {
      fromLabels.push(endpoint.fromCIDR)
    } else {
      endpoint.fromEndpoints.forEach(element => {
        fromLabels.push(convertJSONtoList(element));
      });
    }
    
  });


  console.log(fromLabels);

  return {
    "name":policyName,
    "fromLabels":fromLabels,
    "toLabels": toLabels
  }
}

export default function PolicyListAccordion() {

  const [policies, setPolices] = React.useState([policy1]);

  React.useEffect(() => {
    const policiesFetched = getAllPolicies();
    console.log('policies', policiesFetched);
    //setPolices(policiesFetched);
  }, []);

  //let policy = getDataFromPolicy(policy1);

  return (
    <Box sx={{ borderColor: 'grey.500', padding: 1 }}>
      <Typography variant="h5" gutterBottom component="div">
        Policies
      </Typography>


      {policies.map((policyRaw) => {
        let policy = getDataFromPolicy(policyRaw);
        return (
        <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel1a-content"
          id="panel1a-header"
        >
            <Typography alignItems={'center'} sx={{ width: '70%', flexShrink: 0 }}>
                <PolicyIcon/> {policy.name}
            </Typography>

            {/* <Typography sx={{ color: 'text.secondary', fontSize: '8' }}>2.5h ago</Typography> */}
            
        </AccordionSummary>

        <AccordionDetails alignItems={'center'}>
          
          <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
            flexDirection="column"
            border={2}
            borderColor="red"
            marginBottom={3}
            padding={1}
          >
            
            <Typography variant="h6" color="red">
              Endpoints Denied
            </Typography>
            
            {policy.fromLabels.map((endpoint) => (
              <Box
                display="flex"
                flexDirection="row"
                alignItems="start"
                borderColor="red"
                border={1}
                margin={1}
                padding={1}
              >
                {endpoint.map((label) => (
                  <OverflowTip value={label} mainColor="red"
                  hoverColor="#347aeb"/>
                ))}
                
              </Box>
            ))}
             
            
          </Box>

          <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
            flexDirection="column"
            border={2}
            borderColor="green"
            marginBottom={3}
            padding={1}
          >
            
            <Typography variant="h6" color="green">
              To Endpoints
            </Typography>

            <Box
                display="flex"
                flexDirection="row"
                flexWrap="wrap"
                margin={1}
                padding={1}
                maxWidth="100%"
              >
                {policy.toLabels.map((label) => (
                // <Chip label={label} color="primary" sx={{overflow:"hidden", maxWidth:"30%", marginBottom:"3px", marginLeft:"2px"}}/>
                <OverflowTip value={label} mainColor="#0f52bd"
                  hoverColor="#347aeb"/>
              ))}
                
              </Box>
             
            
          </Box>   
          

          <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
          >
            <Button variant="contained" color="error">
              Deactivate
            </Button> 
          </Box>        


        </AccordionDetails>
      </Accordion>
      )

      })}
      
    </Box>
  );
}
