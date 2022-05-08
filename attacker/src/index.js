serviceurl="myservice.default.svc.cluster.local:5050"
const { exec } = require("child_process");

/*const express = require('express');
const app = express();

const cors = require('cors')


app.use(express.json());
app.use(express.urlencoded({
    extended: true
  }));
app.use(cors())
*/
//setInterval(pingServer, 1000)

//function pingServer() {
    //exec(`hping3 -i u20 -S -p 31000 -c 1000000 ${}`)
//}
milisecondsPerRequest=1
numberOfRequestsSent=0
prevTimeInSeconds=Math.floor(Date.now() /  milisecondsPerRequest)

while (true) {
  currentTimeInSeconds=Math.floor(Date.now() / milisecondsPerRequest)
  prevTimeInSeconds = currentTimeInSeconds;
  exec(`curl ${serviceurl}`)
  numberOfRequestsSent += 1
  //console.log('sent', numberOfRequestsSent, 'request')
  //userWantToRequest=true
  //if (currentTimeInSeconds > prevTimeInSeconds) {
  //}
}