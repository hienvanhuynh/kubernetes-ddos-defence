serviceurl="myservice.default.svc.cluster.local:5050"
const { exec } = require("child_process");

milisecondsPerRequest=1000
numberOfRequestsSent=0
prevTimeInSeconds=Math.floor(Date.now()/milisecondsPerRequest)

while (true) {
  currentTimeInSeconds=Math.floor(Date.now()/milisecondsPerRequest)
  r=Math.floor(Math.random() * 100)
  userWantToRequest=(r % 5) > 1
  if (currentTimeInSeconds > prevTimeInSeconds) {
    prevTimeInSeconds = currentTimeInSeconds;
    if (userWantToRequest) {
      exec(`curl ${serviceurl}`)
      numberOfRequestsSent += 1
      console.log('sent', numberOfRequestsSent, 'request')
    }
  }
}