const express = require('express');
const app = express();



const cors = require('cors')

app.use(express.json({limit: '50mb'}));
app.use(express.urlencoded({
    limit: '50mb',
    extended: true
  }));
app.use(cors())

meanT=0
standardDeviation=0
alpha = 0.08

maxAttackHostsRatio = 0.25
R = 0

app.get('*', (req, res) => {
    return res.send('\
    This is DDoS detection module\n\
    To get analysist result, send a post request with flows to serviceip:serviceport/newpatch\n')
})
app.post('/newpatch', (req, res) => {
    //lengthBefore = req.body.length
    filterOnlyRequestTraffic(req.body)
    
    //req.body is the input json
    T = req.body.length
    
    //if (T!= lengthBefore){
    //    console.log("Really reduce ", lengthBefore, T)
    //}

    //count number of each host
    hostStats = countHostsAppearance(req.body)
    numberOfHost = hostStats.size 

    //for first scrape, we don't know meanT
    if (meanT==0) { meanT = T}
    if (R == 0) {
        R = meanT/numberOfHost
        console.log("R: ", R)
    }
    console.log('meanT:', meanT, " stdDev:", standardDeviation)    
    //Phase 1
    newMeanT = meanT + alpha * (T - meanT)
    newStandardDeviation = Math.sqrt(alpha * (T - meanT) ** 2
                             + (1 - alpha)*(standardDeviation ** 2))
    meanT = newMeanT
    standardDeviation = newStandardDeviation
    threshold = meanT + 3 * standardDeviation

    console.log("meanT:", meanT, "stdDev:",standardDeviation,"threshold: ", threshold)
    console.log("T:", T, "threshold: ", threshold)
    if (T <= threshold) {
        //
        R = meanT / numberOfHost
        console.log("ok")
        return res.send('ok')
    }

    //Phase 2
    console.log('possible attack detected')
    //
    //m
    attackTrafficsRatio = (T - (1-maxAttackHostsRatio) * numberOfHost * R) / (maxAttackHostsRatio * numberOfHost * R)
    //mR
    minAttackTraffics = attackTrafficsRatio * R

    console.log(minAttackTraffics)
    suspectedHosts = getSuspectedHosts(hostStats, minAttackTraffics)

    
    res.send(suspectedHosts)
})

const PORT = 5060
app.listen(PORT, () => console.log(`Server is listening on port ${PORT}`))

function filterOnlyRequestTraffic(requestbody)
{
    for (var i=0; i < requestbody.length; i++) {
        if (requestbody[i].is_reply==true){
            requestbody.splice(i,1)
            i--;
        }
    }
}
function countHostsAppearance(requestbody)
{
    groupip = new Map([])
    for (i in requestbody) {
        host = requestbody[i]
        if (!groupip.has(host.IP.source)) {
            groupip.set(host.IP.source, 1)
        } else {
            groupip.set(host.IP.source, groupip.get(host.IP.source) + 1)
        }
    }

    return groupip
}

function getSuspectedHosts(hostStats, minAttackTraffics)
{
    suspectedHosts = []
    for (host of hostStats.keys()) {
        if (hostStats.get(host) >= minAttackTraffics) {
            suspectedHosts.push(host)
        }
    }

    return suspectedHosts
}

/*function countNumberOfHost(requestbody)
{
    groupip = []
    for (host in requestbody) {
        if (!groupip.includes(host.ip)){
            groupip.push(host.ip)
        }
    }

    return groupip.length
}*/
