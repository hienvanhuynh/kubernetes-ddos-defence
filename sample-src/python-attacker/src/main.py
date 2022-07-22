import subprocess

command = "python3 /home/server/ihulk/src/ihulk.py http://web.default.svc.cluster.local safe"
process = subprocess.Popen(command.split(), stdout=subprocess.PIPE)
output, error = process.communicate()

k=0
while True:
   k=k+1
   if k>1000000000:
       k=0 