import React, { useEffect, useState } from 'react';
import { Box, Container } from '@mui/system';
import {Line} from 'react-chartjs-2';
import { Chart as ChartJS } from 'chart.js/auto'
import { Chart }            from 'react-chartjs-2'
import { Typography } from '@mui/material';
import BasicSelect from './component/select';


// const jsonData = {
//   "status": "success",
//   "data": {
//     "resultType": "vector",
//     "result": [
//       {
//         "metric": {
//           "__name__": "cilium_drop_count_total",
//           "controller_revision_hash": "76d7bb66f9",
//           "direction": "EGRESS",
//           "instance": "192.168.1.100:9090",
//           "job": "kubernetes-pods",
//           "k8s_app": "cilium",
//           "namespace": "kube-system",
//           "pod": "cilium-w2g4n",
//           "pod_template_generation": "2",
//           "reason": "Unsupported L2 protocol"
//         },
//         "value": [
//           1658405173.774,
//           "142"
//         ]
//       },
//       {
//         "metric": {
//           "__name__": "cilium_drop_count_total",
//           "controller_revision_hash": "76d7bb66f9",
//           "direction": "EGRESS",
//           "instance": "192.168.1.100:9090",
//           "job": "kubernetes-pods",
//           "k8s_app": "cilium",
//           "namespace": "kube-system",
//           "pod": "cilium-w2g4n",
//           "pod_template_generation": "2",
//           "reason": "Unsupported L3 protocol"
//         },
//         "value": [
//           1658405173.774,
//           "94"
//         ]
//       },
//       {
//         "metric": {
//           "__name__": "cilium_drop_count_total",
//           "controller_revision_hash": "76d7bb66f9",
//           "direction": "EGRESS",
//           "instance": "192.168.1.101:9090",
//           "job": "kubernetes-pods",
//           "k8s_app": "cilium",
//           "namespace": "kube-system",
//           "pod": "cilium-bhp6s",
//           "pod_template_generation": "2",
//           "reason": "Unsupported L3 protocol"
//         },
//         "value": [
//           1658405173.774,
//           "22"
//         ]
//       },
//       {
//         "metric": {
//           "__name__": "cilium_drop_count_total",
//           "controller_revision_hash": "76d7bb66f9",
//           "direction": "INGRESS",
//           "instance": "192.168.1.100:9090",
//           "job": "kubernetes-pods",
//           "k8s_app": "cilium",
//           "namespace": "kube-system",
//           "pod": "cilium-w2g4n",
//           "pod_template_generation": "2",
//           "reason": "Stale or unroutable IP"
//         },
//         "value": [
//           1658405173.774,
//           "18"
//         ]
//       },
//       {
//         "metric": {
//           "__name__": "cilium_drop_count_total",
//           "controller_revision_hash": "76d7bb66f9",
//           "direction": "INGRESS",
//           "instance": "192.168.1.101:9090",
//           "job": "kubernetes-pods",
//           "k8s_app": "cilium",
//           "namespace": "kube-system",
//           "pod": "cilium-bhp6s",
//           "pod_template_generation": "2",
//           "reason": "Stale or unroutable IP"
//         },
//         "value": [
//           1658405173.774,
//           "15"
//         ]
//       }
//     ]
//   }
// };

const getCurrentTotal = (data) => {
  let sum = 0;

  //console.log('data calculate', data);
  data.data.result.forEach(source => {
    let n = parseInt(source.value[1]);
    sum += n;
  });

  //console.log('sum',sum);

  return sum;
}

const controlDataLength = (data, max) => {
  const currentLength = data.length;

  if (currentLength > max)
  {
    data.shift();
  }

  return data;
}

const updateData = async (labels, dataset, n) => {

  let newLabels = labels;

  let second = new Date().getSeconds();
  let minute = new Date().getMinutes();
  let hours = new Date().getHours();
  
  let drop_data = await fetch("/api/v1/query?query=cilium_drop_count_total");
  
  let forward_data = await fetch("/api/v1/query?query=cilium_forward_count_total");
  
  //console.log('drop raw', drop_data);

  const dropJson = await drop_data.json();
  const forwardJson = await forward_data.json();
  
  newLabels.push(`${hours}:${minute}:${second}`);
  
 
  let dropDataset = dataset[0];
  dropDataset.push(getCurrentTotal(dropJson));
   
  let forwardDataset = dataset[1];
  forwardDataset.push(getCurrentTotal(forwardJson));
   
  let newDataset = [dropDataset, forwardDataset]


  return [newLabels, newDataset]

}

const MAX_GRAPH_LENGTH = 40;

const calculateDataGraph = (data) => {
  //console.log('data', data);
  if (data.length < 2) {
   return [];
  }
  let newData = [];
  for (let i = 1; i < data.length; i++)
  {
    newData.push(data[i]-data[i-1]);
  }
  
  return newData.slice(-MAX_GRAPH_LENGTH);
}

const calculateLabel = (labels) => {
  const indexToRemove = 0;

  const result = [...labels.slice(0, indexToRemove), ...labels.slice(indexToRemove + 1)];
  
  
  //console.log('labels', labels);
  //console.log('labels result', result);
  
  return result.slice(-MAX_GRAPH_LENGTH);	
}


export default function FlowGraph() {

  // getCurrentTotal(jsonData);

  const [labels, setLabels] = useState([]);
  const [dataset, setDataset] = useState([[], []]);
  let n = 1;

  useEffect(() => {
    setInterval(async () => {
      let newLabels;
      let newDataset;
      

      let values = await updateData(labels, dataset, n);
      newLabels = values[0];
      newDataset = values[1];

      // console.log(newLabels);
      // console.log(newDataset);
      // console.log(n);

      setLabels(newLabels);
      setDataset(newDataset);
      n += 1;
    }, 2000);
  }, []);

  const data = {
    labels: calculateLabel(labels),
    datasets: [{
      label: 'Dropped packets',
      data: calculateDataGraph(dataset[0]),
      fill: true,
      borderColor: 'red',
      tension: 0.1
    },
    {
      label: 'Forwarded packets',
      data: calculateDataGraph(dataset[1]),
      fill: true,
      borderColor: 'green',
      tension: 0.1
    }
    ]
  };

  return (
    <Container>

      <Box 
        display="flex"
        justifyContent="space-between"
        flexDirection="row"
      >
        <Typography variant="h5">
          Forward/Drop packets count
        </Typography>

        {/* <BasicSelect
          defaultValue={'prometheus'}
          labelName="Namespace"
          optionList={['prometheus', 'kube-system']}
        /> */}

      </Box>
      

      <Container>
        <Line
          height={500}
          width={300}
          data={data}
          options={{
            maintainAspectRatio: false,
            elements: {
                point:{
                    radius: 0
                }
            }
          }}
        />
      </Container>

    </Container>
    
  );

}
