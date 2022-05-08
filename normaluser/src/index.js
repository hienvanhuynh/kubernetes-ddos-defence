serviceurl="myservice.default.svc.cluster.local:5050"
const { exec } = require("child_process");

milisecondsPerRequest=1000
numberOfRequestsSent=0
prevTimeInSeconds=Math.floor(Date.now()/milisecondsPerRequest)

while (true) {
  currentTimeInSeconds=Math.floor(Date.now()/milisecondsPerRequest)
  r=Math.floor(Math.random() * 100)
  userWantToRequest=! (r % 3)
  if (currentTimeInSeconds > prevTimeInSeconds) {
    if (userWantToRequest) {
      prevTimeInSeconds = currentTimeInSeconds;
      exec(`curl ${serviceurl}`)
      numberOfRequestsSent += 1
      console.log('sent', numberOfRequestsSent, 'request')
    }
  }
}