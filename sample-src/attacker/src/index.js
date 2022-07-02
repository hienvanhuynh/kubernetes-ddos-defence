serviceurl="myservice.default.svc.cluster.local:5050"
const { exec } = require("child_process");

milisecondsPerRequest=1
numberOfRequestsSent=0
prevTimeInSeconds=Math.floor(Date.now() /  milisecondsPerRequest)

while (true) {
  currentTimeInSeconds=Math.floor(Date.now() / milisecondsPerRequest)
  prevTimeInSeconds = currentTimeInSeconds;
  exec(`curl ${serviceurl}`)
  numberOfRequestsSent += 1
}